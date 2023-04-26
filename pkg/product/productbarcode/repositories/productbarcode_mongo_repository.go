package repositories

import (
	"errors"
	"os"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IProductBarcodeRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.ProductBarcodeDoc) (string, error)
	CreateInBatch(docList []models.ProductBarcodeDoc) error
	Update(shopID string, guid string, doc models.ProductBarcodeDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ProductBarcodeDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.ProductBarcodeDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ProductBarcodeItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.ProductBarcodeDoc, error)
	FindByDocIndentityGuids(shopID string, indentityField string, indentityValues interface{}) ([]models.ProductBarcodeDoc, error)

	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductBarcodeDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductBarcodeActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeActivity, error)

	FindMasterInCodes(codes []string) ([]models.ProductBarcodeInfo, error)
	UpdateParentGuidByGuids(shopID string, parentGUID string, guids []string) error
	Transaction(fnc func() error) error
}

type ProductBarcodeRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ProductBarcodeDoc]
	repositories.SearchRepository[models.ProductBarcodeInfo]
	repositories.GuidRepository[models.ProductBarcodeItemGuid]
	repositories.ActivityRepository[models.ProductBarcodeActivity, models.ProductBarcodeDeleteActivity]
}

func NewProductBarcodeRepository(pst microservice.IPersisterMongo) *ProductBarcodeRepository {

	insRepo := &ProductBarcodeRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ProductBarcodeDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ProductBarcodeInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ProductBarcodeItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.ProductBarcodeActivity, models.ProductBarcodeDeleteActivity](pst)

	return insRepo
}

func (repo *ProductBarcodeRepository) UpdateParentGuidByGuids(shopID string, parentGUID string, guids []string) error {

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"guidfixed": bson.M{"$in": guids},
	}

	return repo.pst.Update(models.ProductBarcodeDoc{}, filters, bson.M{"$set": bson.M{"parentguid": parentGUID}})
}

func (repo *ProductBarcodeRepository) FindMasterInCodes(codes []string) ([]models.ProductBarcodeInfo, error) {

	masterShopID := os.Getenv("MASTER_SHOP_ID")

	if len(masterShopID) == 0 {
		return []models.ProductBarcodeInfo{}, errors.New("master shop id is empty")
	}

	docList := []models.ProductBarcodeInfo{}

	filters := bson.M{
		"shopid": masterShopID,
		"barcode": bson.M{
			"$in": codes,
		},
	}

	err := repo.pst.Find(models.ProductBarcodeInfo{}, filters, &docList)

	if err != nil {
		return nil, err
	}

	return docList, nil
}

func (repo *ProductBarcodeRepository) Transaction(fnc func() error) error {
	return repo.pst.Transaction(fnc)
}
