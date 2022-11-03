package member

import (
	"errors"
	"fmt"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/member/models"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"sync"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMemberService interface {
	IsExistsGuid(shopID string, guidFixed string) (bool, error)
	//CreateIndex(doc models.MemberIndex) error
	CreateWithGuid(shopID string, username string, guid string, doc models.Member) (string, error)
	Create(shopID string, username string, doc models.Member) (string, error)
	Update(shopID string, guid string, username string, doc models.Member) error
	Delete(shopID string, guid string, username string) error
	Info(shopID string, guid string) (models.MemberInfo, error)
	Search(shopID string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error)

	LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, paginate.PaginationData, error)
	LastActivityOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) (common.LastActivity, error)
	GetModuleName() string
}

type MemberService struct {
	repo          IMemberRepository
	memberPgRepo  IMemberPGRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository
}

func NewMemberService(repo IMemberRepository, memberPgRepo IMemberPGRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) MemberService {
	return MemberService{
		repo:          repo,
		memberPgRepo:  memberPgRepo,
		syncCacheRepo: syncCacheRepo,
	}
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

// func (svc MemberService) CreateIndex(doc models.MemberIndex) error {

// 	err := svc.memberPgRepo.Create(doc)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

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

func (svc MemberService) Search(shopID string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc MemberService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, paginate.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.MemberDeleteActivity
	var pagination1 paginate.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.repo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.MemberActivity
	var pagination2 paginate.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.repo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return common.LastActivity{}, pagination1, err1
	}

	if err2 != nil {
		return common.LastActivity{}, pagination2, err2
	}

	lastActivity := common.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	pagination := pagination1

	if pagination.Total < pagination2.Total {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc MemberService) LastActivityOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) (common.LastActivity, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.MemberDeleteActivity
	var err1 error

	go func() {
		deleteDocList, err1 = svc.repo.FindDeletedOffset(shopID, lastUpdatedDate, skip, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.MemberActivity

	var err2 error

	go func() {
		createAndUpdateDocList, err2 = svc.repo.FindCreatedOrUpdatedOffset(shopID, lastUpdatedDate, skip, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return common.LastActivity{}, err1
	}

	if err2 != nil {
		return common.LastActivity{}, err2
	}

	lastActivity := common.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	return lastActivity, nil
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
