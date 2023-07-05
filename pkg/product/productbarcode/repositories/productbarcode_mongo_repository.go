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
	CountByRefBarcode(shopID string, refBarcode string) (int, error)
	CountByRefGuids(shopID string, GUIDs []string) (int, error)
	CountByUnitCodes(shopID string, unitCodes []string) (int, error)
	CountByGroupCode(shopID string, unitCodes []string) (int, error)
	Create(doc models.ProductBarcodeDoc) (string, error)
	CreateInBatch(docList []models.ProductBarcodeDoc) error
	Update(shopID string, guid string, doc models.ProductBarcodeDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
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
	FindByRefBarcode(shopID string, barcode string) ([]models.ProductBarcodeDoc, error)

	FindByBarcodes(shopID string, barcodes []string) ([]models.ProductBarcodeInfo, error)
	FindPageByUnits(shopID string, unitCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	FindPageByGroups(shopID string, groupCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
}

type ProductBarcodeRepository struct {
	pst   microservice.IPersisterMongo
	cache microservice.ICacher
	repositories.CrudRepository[models.ProductBarcodeDoc]
	repositories.SearchRepository[models.ProductBarcodeInfo]
	repositories.GuidRepository[models.ProductBarcodeItemGuid]
	repositories.ActivityRepository[models.ProductBarcodeActivity, models.ProductBarcodeDeleteActivity]
}

func NewProductBarcodeRepository(pst microservice.IPersisterMongo, cache microservice.ICacher) *ProductBarcodeRepository {

	insRepo := &ProductBarcodeRepository{
		pst:   pst,
		cache: cache,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ProductBarcodeDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ProductBarcodeInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ProductBarcodeItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.ProductBarcodeActivity, models.ProductBarcodeDeleteActivity](pst)

	return insRepo
}

func (repo ProductBarcodeRepository) CountByRefBarcode(shopID string, refBarcode string) (int, error) {

	return repo.CountByKey(shopID, "refbarcodes.barcode", refBarcode)
}

func (repo ProductBarcodeRepository) CountByRefGuids(shopID string, GUIDs []string) (int, error) {

	return repo.CountByInKeys(shopID, "refbarcodes.guidfixed", GUIDs)
}

func (repo ProductBarcodeRepository) CountByUnitCodes(shopID string, unitCodes []string) (int, error) {

	return repo.CountByInKeys(shopID, "itemunitcode", unitCodes)
}

func (repo ProductBarcodeRepository) CountByGroupCode(shopID string, unitCodes []string) (int, error) {

	return repo.CountByInKeys(shopID, "groupcode", unitCodes)
}

func (repo ProductBarcodeRepository) UpdateParentGuidByGuids(shopID string, parentGUID string, guids []string) error {

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"guidfixed": bson.M{"$in": guids},
	}

	return repo.pst.Update(models.ProductBarcodeDoc{}, filters, bson.M{"$set": bson.M{"parentguid": parentGUID}})
}

func (repo ProductBarcodeRepository) FindMasterInCodes(codes []string) ([]models.ProductBarcodeInfo, error) {

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

func (repo ProductBarcodeRepository) FindByRefBarcode(shopID string, barcode string) ([]models.ProductBarcodeDoc, error) {

	docList := []models.ProductBarcodeDoc{}

	filters := bson.M{
		"shopid":              shopID,
		"refbarcodes.barcode": barcode,
		"itemtype":            bson.M{"$ne": 2},
		"deletedat":           bson.M{"$exists": false},
	}

	err := repo.pst.Find(models.ProductBarcodeDoc{}, filters, &docList)

	if err != nil {
		return nil, err
	}

	return docList, nil
}

func (repo ProductBarcodeRepository) Transaction(fnc func() error) error {
	return repo.pst.Transaction(fnc)
}

func (repo ProductBarcodeRepository) FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func (repo ProductBarcodeRepository) FindByBarcodes(shopID string, barcodes []string) ([]models.ProductBarcodeInfo, error) {

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"barcode":   bson.M{"$in": barcodes},
	}

	var results []models.ProductBarcodeInfo
	err := repo.pst.Find(models.ProductBarcodeInfo{}, filters, &results)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo ProductBarcodeRepository) FindPageByUnits(shopID string, unitCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {

	filters := bson.M{
		"shopid": shopID,
		"deletedat": bson.M{
			"$exists": false,
		},
		"itemunitcode": bson.M{
			"$in": unitCodes,
		},
	}

	results := []models.ProductBarcodeInfo{}
	pagination, err := repo.pst.FindPage(models.ProductBarcodeInfo{}, filters, pageable, &results)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func (repo ProductBarcodeRepository) FindPageByGroups(shopID string, groupCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {

	filters := bson.M{
		"shopid": shopID,
		"deletedat": bson.M{
			"$exists": false,
		},
		"groupcode": bson.M{
			"$in": groupCodes,
		},
	}

	results := []models.ProductBarcodeInfo{}
	pagination, err := repo.pst.FindPage(models.ProductBarcodeInfo{}, filters, pageable, &results)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}
