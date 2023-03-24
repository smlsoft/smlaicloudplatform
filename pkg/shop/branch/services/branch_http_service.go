package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/shop/branch/models"
	"smlcloudplatform/pkg/shop/branch/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IBranchHttpService interface {
	CreateBranch(shopID string, authUsername string, doc models.Branch) (string, error)
	UpdateBranch(shopID string, guid string, authUsername string, doc models.Branch) error
	DeleteBranch(shopID string, guid string, authUsername string) error
	DeleteBranchByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoBranch(shopID string, guid string) (models.BranchInfo, error)
	SearchBranch(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.BranchInfo, mongopagination.PaginationData, error)
	SearchBranchStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.BranchInfo, int, error)

	GetModuleName() string
}

type BranchHttpService struct {
	repo repositories.IBranchRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.BranchActivity, models.BranchDeleteActivity]
}

func NewBranchHttpService(repo repositories.IBranchRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *BranchHttpService {

	insSvc := &BranchHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.BranchActivity, models.BranchDeleteActivity](repo)

	return insSvc
}

func (svc BranchHttpService) CreateBranch(shopID string, authUsername string, doc models.Branch) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.BranchDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Branch = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc BranchHttpService) UpdateBranch(shopID string, guid string, authUsername string, doc models.Branch) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Branch = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc BranchHttpService) DeleteBranch(shopID string, guid string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc BranchHttpService) DeleteBranchByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc BranchHttpService) InfoBranch(shopID string, guid string) (models.BranchInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.BranchInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.BranchInfo{}, errors.New("document not found")
	}

	return findDoc.BranchInfo, nil

}

func (svc BranchHttpService) SearchBranch(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.BranchInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.BranchInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc BranchHttpService) SearchBranchStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.BranchInfo, int, error) {
	searchInFields := []string{
		"code",
	}

	selectFields := map[string]interface{}{
		"guidfixed": 1,
		"code":      1,
	}

	if langCode != "" {
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		selectFields["names"] = 1
	}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.BranchInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc BranchHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc BranchHttpService) GetModuleName() string {
	return "branch"
}
