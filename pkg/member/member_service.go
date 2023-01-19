package member

import (
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

	LastActivity(shopID string, action string, lastUpdatedDate time.Time, pageable micromodels.Pageable) (common.LastActivity, mongopagination.PaginationData, error)
	GetModuleName() string
}

type MemberService struct {
	repo          IMemberRepository
	memberPgRepo  IMemberPGRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository

	services.ActivityService[models.MemberActivity, models.MemberDeleteActivity]
}

func NewMemberService(repo IMemberRepository, memberPgRepo IMemberPGRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) MemberService {
	insSvc := MemberService{
		repo:          repo,
		memberPgRepo:  memberPgRepo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.MemberActivity, models.MemberDeleteActivity](repo)
	return insSvc
}

func (svc MemberService) IsExistsGuid(shopID string, guidFixed string) (bool, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guidFixed)

	if err != nil {
		return false, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return false, nil
	}

	return true, nil

}

func (svc MemberService) CreateWithGuid(shopID string, username string, guid string, doc models.Member) (string, error) {
	dataDoc := models.MemberDoc{}

	newGuid := guid
	dataDoc.GuidFixed = newGuid
	dataDoc.ShopID = shopID
	dataDoc.CreatedBy = username
	dataDoc.CreatedAt = time.Now()

	dataDoc.LastUpdatedAt = time.Now()

	dataDoc.Member = doc

	_, err := svc.repo.Create(dataDoc)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc MemberService) Create(shopID string, username string, doc models.Member) (string, error) {
	dataDoc := models.MemberDoc{}

	newGuid := utils.NewGUID()
	dataDoc.GuidFixed = newGuid
	dataDoc.ShopID = shopID
	dataDoc.CreatedBy = username
	dataDoc.CreatedAt = time.Now()

	dataDoc.LastUpdatedAt = time.Now()

	dataDoc.Member = doc

	_, err := svc.repo.Create(dataDoc)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuid, nil
}

func (svc MemberService) Update(shopID string, guid string, username string, doc models.Member) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

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

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc MemberService) Delete(shopID string, guid string, username string) error {

	err := svc.repo.Delete(shopID, guid, username)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc MemberService) Info(shopID string, guid string) (models.MemberInfo, error) {
	doc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.MemberInfo{}, err
	}

	return doc.MemberInfo, nil
}

func (svc MemberService) Search(shopID string, pageable micromodels.Pageable) ([]models.MemberInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, pageable)

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
