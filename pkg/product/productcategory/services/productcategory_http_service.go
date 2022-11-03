package services

import (
	"errors"
	"fmt"
	"smlcloudplatform/pkg/product/productcategory/models"
	"smlcloudplatform/pkg/product/productcategory/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"time"

	mastersync "smlcloudplatform/pkg/mastersync/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductCategoryHttpService interface {
	CreateProductCategory(shopID string, authUsername string, doc models.ProductCategory) (string, error)
	UpdateProductCategory(shopID string, guid string, authUsername string, doc models.ProductCategory) error
	DeleteProductCategory(shopID string, guid string, authUsername string) error
	InfoProductCategory(shopID string, guid string) (models.ProductCategoryInfo, error)
	SearchProductCategory(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ProductCategoryInfo, mongopagination.PaginationData, error)
	SearchProductCategoryStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.ProductCategoryInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ProductCategory) error

	GetModuleName() string
}

type ProductCategoryHttpService struct {
	repo          repositories.IProductCategoryRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository

	services.ActivityService[models.ProductCategoryActivity, models.ProductCategoryDeleteActivity]
}

func NewProductCategoryHttpService(repo repositories.IProductCategoryRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *ProductCategoryHttpService {

	insSvc := &ProductCategoryHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.ProductCategoryActivity, models.ProductCategoryDeleteActivity](repo)

	return insSvc
}

func (svc ProductCategoryHttpService) CreateProductCategory(shopID string, authUsername string, doc models.ProductCategory) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.ProductCategoryDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ProductCategory = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc ProductCategoryHttpService) UpdateProductCategory(shopID string, guid string, authUsername string, doc models.ProductCategory) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ProductCategory = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductCategoryHttpService) DeleteProductCategory(shopID string, guid string, authUsername string) error {

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

	return nil
}

func (svc ProductCategoryHttpService) InfoProductCategory(shopID string, guid string) (models.ProductCategoryInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ProductCategoryInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ProductCategoryInfo{}, errors.New("document not found")
	}

	return findDoc.ProductCategoryInfo, nil

}

func (svc ProductCategoryHttpService) SearchProductCategory(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ProductCategoryInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.ProductCategoryInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductCategoryHttpService) SearchProductCategoryStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.ProductCategoryInfo, int, error) {
	searchCols := []string{
		"guidfixed",
	}

	projectQuery := map[string]interface{}{
		"guidfixed":     1,
		"parentguid":    1,
		"parentguidall": 1,
		"imageuri":      1,
	}

	if langCode != "" {
		projectQuery["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		projectQuery["names"] = 1
	}

	docList, total, err := svc.repo.FindLimit(shopID, searchCols, q, skip, limit, sort, projectQuery)

	if err != nil {
		return []models.ProductCategoryInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ProductCategoryHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.ProductCategory) error {

	createDataList := []models.ProductCategoryDoc{}

	createdAt := time.Now()
	for _, doc := range dataList {

		newGuidFixed := utils.NewGUID()

		docData := models.ProductCategoryDoc{}
		docData.ShopID = shopID
		docData.GuidFixed = newGuidFixed
		docData.ProductCategory = doc

		docData.CreatedBy = authUsername
		docData.CreatedAt = createdAt

		createDataList = append(createDataList, docData)
	}

	if len(dataList) > 0 {
		err := svc.repo.CreateInBatch(createDataList)

		if err != nil {
			return err
		}

	}

	return nil
}

func (svc ProductCategoryHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc ProductCategoryHttpService) GetModuleName() string {
	return "productcategory"
}
