package repositories

import (
	"context"
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
	Count(ctx context.Context, shopID string) (int, error)
	CountByRefBarcode(ctx context.Context, shopID string, refBarcode string) (int, error)
	CountByRefGuids(ctx context.Context, shopID string, GUIDs []string) (int, error)
	CountByBOM(ctx context.Context, shopID string, bomBarcode string) (int, error)
	CountByBOMGuids(ctx context.Context, shopID string, GUIDs []string) (int, error)
	CountByUnitCodes(ctx context.Context, shopID string, unitCodes []string) (int, error)
	CountByGroupCodes(ctx context.Context, shopID string, unitCodes []string) (int, error)
	CountByOrderTypes(ctx context.Context, shopID string, GUIDs []string) (int, error)
	CountByProductTypes(ctx context.Context, shopID string, GUIDs []string) (int, error)
	Create(ctx context.Context, doc models.ProductBarcodeDoc) (string, error)
	CreateInBatch(ctx context.Context, docList []models.ProductBarcodeDoc) error
	Update(ctx context.Context, shopID string, guid string, doc models.ProductBarcodeDoc) error
	DeleteByGuidfixed(ctx context.Context, shopID string, guid string, username string) error
	Delete(ctx context.Context, shopID string, username string, filters map[string]interface{}) error
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.ProductBarcodeDoc, error)
	FindByGuids(ctx context.Context, shopID string, guids []string) ([]models.ProductBarcodeDoc, error)
	FindPageFilter(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)

	FindInItemGuid(ctx context.Context, shopID string, columnName string, itemGuidList []string) ([]models.ProductBarcodeItemGuid, error)
	FindByDocIndentityGuid(ctx context.Context, shopID string, indentityField string, indentityValue interface{}) (models.ProductBarcodeDoc, error)
	FindByDocIndentityGuids(ctx context.Context, shopID string, indentityField string, indentityValues interface{}) ([]models.ProductBarcodeDoc, error)

	FindStep(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeInfo, int, error)

	FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductBarcodeDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductBarcodeActivity, mongopagination.PaginationData, error)
	FindDeletedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeDeleteActivity, error)
	FindCreatedOrUpdatedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeActivity, error)

	FindMasterInCodes(ctx context.Context, codes []string) ([]models.ProductBarcodeInfo, error)
	UpdateParentGuidByGuids(ctx context.Context, shopID string, parentGUID string, guids []string) error
	Transaction(ctx context.Context, fnc func(ctx context.Context) error) error
	FindByRefBarcode(ctx context.Context, shopID string, barcode string) ([]models.ProductBarcodeDoc, error)
	FindByBOMBarcode(ctx context.Context, shopID string, barcode string) ([]models.ProductBarcodeDoc, error)

	FindByBarcode(ctx context.Context, shopID string, barcode string) (models.ProductBarcodeDoc, error)
	FindByBarcodes(ctx context.Context, shopID string, barcodes []string) ([]models.ProductBarcodeInfo, error)
	FindPageByUnits(ctx context.Context, shopID string, unitCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	FindPageByGroups(ctx context.Context, shopID string, groupCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)

	UpdateRefBarcodeByGUID(ctx context.Context, shopID string, guid string, refBarcode models.RefProductBarcode) error
	UpdateAllProductTypeByGUID(ctx context.Context, shopID string, guid string, doc models.ProductType) error
	UpdateAllProductGroupByCode(ctx context.Context, shopID string, doc models.ProductGroup) error
	UpdateAllProductUnitByCode(ctx context.Context, shopID string, doc models.ProductUnit) error
	UpdateAllProductOrderTypeByGUID(ctx context.Context, shopID string, guid string, doc models.ProductOrderType) error
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

func (repo ProductBarcodeRepository) CountByRefBarcode(ctx context.Context, shopID string, refBarcode string) (int, error) {
	return repo.CountByKey(ctx, shopID, "refbarcodes.barcode", refBarcode)
}

func (repo ProductBarcodeRepository) CountByBOM(ctx context.Context, shopID string, bomBarcode string) (int, error) {
	return repo.CountByKey(ctx, shopID, "bom.barcode", bomBarcode)
}

func (repo ProductBarcodeRepository) CountByRefGuids(ctx context.Context, shopID string, GUIDs []string) (int, error) {
	return repo.CountByInKeys(ctx, shopID, "refbarcodes.guidfixed", GUIDs)
}

func (repo ProductBarcodeRepository) CountByBOMGuids(ctx context.Context, shopID string, GUIDs []string) (int, error) {
	return repo.CountByInKeys(ctx, shopID, "bom.guidfixed", GUIDs)
}

func (repo ProductBarcodeRepository) CountByUnitCodes(ctx context.Context, shopID string, unitCodes []string) (int, error) {
	return repo.CountByInKeys(ctx, shopID, "itemunitcode", unitCodes)
}

func (repo ProductBarcodeRepository) CountByGroupCodes(ctx context.Context, shopID string, unitCodes []string) (int, error) {
	return repo.CountByInKeys(ctx, shopID, "groupcode", unitCodes)
}

func (repo ProductBarcodeRepository) CountByOrderTypes(ctx context.Context, shopID string, GUIDs []string) (int, error) {
	return repo.CountByInKeys(ctx, shopID, "ordertypes.guidfixed", GUIDs)
}

func (repo ProductBarcodeRepository) CountByProductTypes(ctx context.Context, shopID string, GUIDs []string) (int, error) {
	return repo.CountByInKeys(ctx, shopID, "producttype.guidfixed", GUIDs)
}

func (repo ProductBarcodeRepository) UpdateParentGuidByGuids(ctx context.Context, shopID string, parentGUID string, guids []string) error {

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"guidfixed": bson.M{"$in": guids},
	}

	return repo.pst.Update(ctx, models.ProductBarcodeDoc{}, filters, bson.M{"$set": bson.M{"parentguid": parentGUID}})
}

func (repo ProductBarcodeRepository) FindMasterInCodes(ctx context.Context, codes []string) ([]models.ProductBarcodeInfo, error) {

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

	err := repo.pst.Find(ctx, models.ProductBarcodeInfo{}, filters, &docList)

	if err != nil {
		return nil, err
	}

	return docList, nil
}

func (repo ProductBarcodeRepository) FindByRefBarcode(ctx context.Context, shopID string, barcode string) ([]models.ProductBarcodeDoc, error) {

	docList := []models.ProductBarcodeDoc{}

	filters := bson.M{
		"shopid":              shopID,
		"refbarcodes.barcode": barcode,
		"itemtype":            bson.M{"$ne": 2},
		"deletedat":           bson.M{"$exists": false},
	}

	err := repo.pst.Find(ctx, models.ProductBarcodeDoc{}, filters, &docList)

	if err != nil {
		return nil, err
	}

	return docList, nil
}

func (repo ProductBarcodeRepository) FindByBOMBarcode(ctx context.Context, shopID string, barcode string) ([]models.ProductBarcodeDoc, error) {

	docList := []models.ProductBarcodeDoc{}

	filters := bson.M{
		"shopid":      shopID,
		"bom.barcode": barcode,
		"itemtype":    bson.M{"$ne": 2},
		"deletedat":   bson.M{"$exists": false},
	}

	err := repo.pst.Find(ctx, models.ProductBarcodeDoc{}, filters, &docList)

	if err != nil {
		return nil, err
	}

	return docList, nil
}

func (repo ProductBarcodeRepository) Transaction(ctx context.Context, fnc func(ctx context.Context) error) error {
	return repo.pst.Transaction(ctx, fnc)
}

func (repo ProductBarcodeRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func (repo ProductBarcodeRepository) FindByBarcode(ctx context.Context, shopID string, barcode string) (models.ProductBarcodeDoc, error) {

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"barcode":   barcode,
	}

	var result models.ProductBarcodeDoc
	err := repo.pst.Find(ctx, models.ProductBarcodeDoc{}, filters, &result)

	if err != nil {
		return models.ProductBarcodeDoc{}, err
	}

	return result, nil
}

func (repo ProductBarcodeRepository) FindByBarcodes(ctx context.Context, shopID string, barcodes []string) ([]models.ProductBarcodeInfo, error) {

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"barcode":   bson.M{"$in": barcodes},
	}

	var results []models.ProductBarcodeInfo
	err := repo.pst.Find(ctx, models.ProductBarcodeInfo{}, filters, &results)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo ProductBarcodeRepository) FindPageByUnits(ctx context.Context, shopID string, unitCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {

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
	pagination, err := repo.pst.FindPage(ctx, models.ProductBarcodeInfo{}, filters, pageable, &results)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func (repo ProductBarcodeRepository) FindPageByGroups(ctx context.Context, shopID string, groupCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {

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
	pagination, err := repo.pst.FindPage(ctx, models.ProductBarcodeInfo{}, filters, pageable, &results)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func (repo ProductBarcodeRepository) UpdateRefBarcodeByGUID(ctx context.Context, shopID string, guid string, refBarcode models.RefProductBarcode) error {

	filters := bson.M{
		"shopid":                shopID,
		"deletedat":             bson.M{"$exists": false},
		"refbarcodes.guidfixed": guid,
	}

	update := bson.M{
		"$set": bson.M{
			"refbarcodes.$.names":         refBarcode.Names,
			"refbarcodes.$.itemunitcode":  refBarcode.ItemUnitCode,
			"refbarcodes.$.itemunitnames": refBarcode.ItemUnitNames,
			"refbarcodes.$.condition":     refBarcode.Condition,
			"refbarcodes.$.dividevalue":   refBarcode.DivideValue,
			"refbarcodes.$.standvalue":    refBarcode.StandValue,
			"refbarcodes.$.qty":           refBarcode.Qty,
		},
	}

	return repo.pst.Update(ctx, models.ProductBarcodeDoc{}, filters, update)
}

func (repo ProductBarcodeRepository) UpdateAllProductTypeByGUID(ctx context.Context, shopID string, guid string, doc models.ProductType) error {
	filters := bson.M{
		"shopid":                shopID,
		"deletedat":             bson.M{"$exists": false},
		"producttype.guidfixed": guid,
	}

	update := bson.M{
		"$set": bson.M{
			"producttype.code":  doc.Code,
			"producttype.names": doc.Names,
		},
	}

	return repo.pst.Update(ctx, models.ProductBarcodeDoc{}, filters, update)
}

func (repo ProductBarcodeRepository) UpdateAllProductGroupByCode(ctx context.Context, shopID string, doc models.ProductGroup) error {
	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"groupcode": doc.Code,
	}

	update := bson.M{
		"$set": bson.M{
			"groupcode":  doc.Code,
			"groupnames": doc.Names,
		},
	}

	return repo.pst.Update(ctx, models.ProductBarcodeDoc{}, filters, update)
}

func (repo ProductBarcodeRepository) UpdateAllProductUnitByCode(ctx context.Context, shopID string, doc models.ProductUnit) error {
	filters := bson.M{
		"shopid":       shopID,
		"deletedat":    bson.M{"$exists": false},
		"itemunitcode": doc.UnitCode,
	}

	update := bson.M{
		"$set": bson.M{
			"itemunitcode":  doc.UnitCode,
			"itemunitnames": doc.Names,
		},
	}

	return repo.pst.Update(ctx, models.ProductBarcodeDoc{}, filters, update)
}

func (repo ProductBarcodeRepository) UpdateAllProductOrderTypeByGUID(ctx context.Context, shopID string, guid string, doc models.ProductOrderType) error {
	filters := bson.M{
		"shopid":               shopID,
		"deletedat":            bson.M{"$exists": false},
		"ordertypes.guidfixed": guid,
	}

	update := bson.M{
		"$set": bson.M{
			"ordertypes.$.code":  doc.Code,
			"ordertypes.$.names": doc.Names,
			"ordertypes.$.price": doc.Price,
		},
	}

	return repo.pst.Update(ctx, models.ProductBarcodeDoc{}, filters, update)
}
