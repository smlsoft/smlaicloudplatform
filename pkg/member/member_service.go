package member

import (
	"context"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/member/models"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
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

	LastActivity(shopID string, action string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) (common.LastActivity, mongopagination.PaginationData, error)
	GetModuleName() string
}

type MemberService struct {
	repo          IMemberRepository
	memberPgRepo  IMemberPGRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository

	services.ActivityService[models.MemberActivity, models.MemberDeleteActivity]
	contextTimeout time.Duration
}

func NewMemberService(repo IMemberRepository, memberPgRepo IMemberPGRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) MemberService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := MemberService{
		repo:           repo,
		memberPgRepo:   memberPgRepo,
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

	findDoc.UpdatedBy = username
	findDoc.UpdatedAt = time.Now()
	findDoc.Member = doc

	findDoc.LastUpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

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
