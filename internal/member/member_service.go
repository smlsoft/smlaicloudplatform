package member

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	mastersync "smlcloudplatform/internal/mastersync/repositories"
	"smlcloudplatform/internal/member/models"
	common "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/services"
	"smlcloudplatform/internal/shop"
	"smlcloudplatform/internal/utils"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMemberService interface {
	IsExistsGuid(shopID string, guidFixed string) (bool, error)
	CreateWithGuid(shopID string, username string, guid string, doc models.Member) (string, error)
	Create(shopID string, username string, doc models.Member) (string, error)
	Update(shopID string, guid string, username string, doc models.Member) error
	Delete(shopID string, guid string, username string) error
	Info(shopID string, guid string) (models.MemberInfo, error)
	Search(shopID string, pageable micromodels.Pageable) ([]models.MemberInfo, mongopagination.PaginationData, error)

	AuthWithLine(lineAuth models.LineAuthRequest) (string, error)
	UpdateProfileWithLine(shopID string, lineUID string, doc models.Member) error
	LineProfileInfo(shopID string, lineUID string) (models.MemberInfo, error)

	LastActivity(shopID string, action string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) (common.LastActivity, mongopagination.PaginationData, error)
	GetModuleName() string
}

type MemberService struct {
	shopService   shop.IShopService
	repo          IMemberRepository
	memberPgRepo  IMemberPGRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository
	authService   *microservice.AuthService

	services.ActivityService[models.MemberActivity, models.MemberDeleteActivity]
	contextTimeout time.Duration
}

func NewMemberService(repo IMemberRepository, memberPgRepo IMemberPGRepository, shopService shop.IShopService, authService *microservice.AuthService, syncCacheRepo mastersync.IMasterSyncCacheRepository) MemberService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := MemberService{
		shopService:    shopService,
		repo:           repo,
		memberPgRepo:   memberPgRepo,
		authService:    authService,
		syncCacheRepo:  syncCacheRepo,
		contextTimeout: contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.MemberActivity, models.MemberDeleteActivity](repo)
	return insSvc
}

func (svc MemberService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc MemberService) IsExistsGuid(shopID string, guidFixed string) (bool, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guidFixed)

	if err != nil {
		return false, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return false, nil
	}

	return true, nil
}

func (svc MemberService) UpdateProfileWithLine(shopID string, lineUID string, doc models.Member) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByLineUID(ctx, shopID, lineUID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("guid invalid")
	}

	dataDoc := findDoc

	dataDoc.Member = doc

	dataDoc.LineUID = findDoc.LineUID
	dataDoc.UpdatedBy = lineUID
	dataDoc.UpdatedAt = time.Now()

	dataDoc.LastUpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, findDoc.GuidFixed, dataDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc MemberService) registerWithLine(shopID string, lineProfile models.LineProfile) (string, error) {

	idx, err := svc.Create(shopID, lineProfile.UserID, models.Member{
		LineUID:    lineProfile.UserID,
		Name:       lineProfile.DisplayName,
		PictureUrl: lineProfile.PictureUrl,
	})

	if err != nil {
		return "", err
	}

	shopInfo, err := svc.shopService.InfoShop(shopID)

	if err != nil {
		return "", err
	}

	if shopInfo.GuidFixed == "" {
		return "", errors.New("shop invalid")
	}

	return idx, nil
}

func (svc MemberService) AuthWithLine(lineAuth models.LineAuthRequest) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	_, lineProfile, err := svc.LineValidator(lineAuth.LineAccessToken)

	if err != nil {
		return "", err
	}

	findDoc, err := svc.repo.FindByLineUID(ctx, lineAuth.ShopID, lineProfile.UserID)

	if err != nil {
		return "", err
	}

	lineUID := ""
	memberName := ""

	if findDoc.GuidFixed == "" {
		_, err = svc.registerWithLine(lineAuth.ShopID, lineProfile)

		if err != nil {
			return "", err
		}

		svc.saveMasterSync(lineAuth.ShopID)

		lineUID = lineProfile.UserID
		memberName = lineProfile.DisplayName
	} else {
		lineUID = findDoc.Member.LineUID
		memberName = findDoc.Member.Name
	}

	userInfo := micromodels.UserInfo{
		Username: lineUID,
		Name:     memberName,
		ShopID:   lineAuth.ShopID,
		Role:     0,
	}

	tokenID, err := svc.authService.GenerateTokenWithRedisExpire(microservice.AUTHTYPE_BEARER, userInfo, time.Duration(24*30)*time.Hour)

	if err != nil {
		return "", err
	}

	return tokenID, nil

}

func (svc MemberService) LineValidator(lineToken string) (models.LineVerify, models.LineProfile, error) {
	lineVerify, err := svc.LineVerify(lineToken)

	if err != nil {
		return models.LineVerify{}, models.LineProfile{}, err
	}

	if lineVerify.ExpiresIn <= 0 {
		return models.LineVerify{}, models.LineProfile{}, errors.New("line token invalid")
	}

	if lineVerify.ClientID == "" {
		return models.LineVerify{}, models.LineProfile{}, errors.New("line token invalid")
	}

	lineProfile, err := svc.LineProfile(lineToken)

	if err != nil {
		return models.LineVerify{}, models.LineProfile{}, err
	}

	if lineProfile.UserID == "" {
		return models.LineVerify{}, models.LineProfile{}, errors.New("line token invalid")
	}

	return lineVerify, lineProfile, nil
}

func (svc MemberService) LineProfile(lineToken string) (models.LineProfile, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.line.me/v2/profile", nil)

	if err != nil {
		return models.LineProfile{}, err
	}

	req.Header.Add("Authorization", "Bearer "+lineToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return models.LineProfile{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return models.LineProfile{}, err
	}

	lineProfile := models.LineProfile{}

	err = json.Unmarshal(body, &lineProfile)
	if err != nil {
		return models.LineProfile{}, err
	}

	return lineProfile, nil
}

func (svc MemberService) LineVerify(lineToken string) (models.LineVerify, error) {

	req, err := http.NewRequest(http.MethodGet, "https://api.line.me/oauth2/v2.1/verify?access_token="+lineToken, nil)

	if err != nil {
		return models.LineVerify{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return models.LineVerify{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return models.LineVerify{}, err
	}

	lineVerify := models.LineVerify{}
	err = json.Unmarshal(body, &lineVerify)

	if err != nil {
		return models.LineVerify{}, err
	}

	return lineVerify, nil
}

func (svc MemberService) LineProfileInfo(shopID string, lineUID string) (models.MemberInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	doc, err := svc.repo.FindByLineUID(ctx, shopID, lineUID)

	if err != nil {
		return models.MemberInfo{}, err
	}

	return doc.MemberInfo, nil
}

func (svc MemberService) CreateWithGuid(shopID string, username string, guid string, doc models.Member) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	dataDoc := models.MemberDoc{}

	newGuid := guid
	dataDoc.GuidFixed = newGuid
	dataDoc.ShopID = shopID
	dataDoc.CreatedBy = username
	dataDoc.CreatedAt = time.Now()

	dataDoc.LastUpdatedAt = time.Now()

	dataDoc.Member = doc

	_, err := svc.repo.Create(ctx, dataDoc)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc MemberService) Create(shopID string, username string, doc models.Member) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	dataDoc := models.MemberDoc{}

	newGuid := utils.NewGUID()
	dataDoc.GuidFixed = newGuid
	dataDoc.ShopID = shopID
	dataDoc.CreatedBy = username
	dataDoc.CreatedAt = time.Now()

	dataDoc.LastUpdatedAt = time.Now()

	dataDoc.Member = doc

	_, err := svc.repo.Create(ctx, dataDoc)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuid, nil
}

func (svc MemberService) Update(shopID string, guid string, username string, doc models.Member) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("guid invalid")
	}

	dataDoc := findDoc

	dataDoc.Member = doc

	dataDoc.LineUID = findDoc.LineUID
	dataDoc.UpdatedBy = username
	dataDoc.UpdatedAt = time.Now()
	dataDoc.LastUpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, dataDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc MemberService) Delete(shopID string, guid string, username string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	err := svc.repo.Delete(ctx, shopID, guid, username)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc MemberService) Info(shopID string, guid string) (models.MemberInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	doc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.MemberInfo{}, err
	}

	return doc.MemberInfo, nil
}

func (svc MemberService) Search(shopID string, pageable micromodels.Pageable) ([]models.MemberInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList, pagination, err := svc.repo.FindPage(ctx, shopID, pageable)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc MemberService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc MemberService) GetModuleName() string {
	return "member"
}
