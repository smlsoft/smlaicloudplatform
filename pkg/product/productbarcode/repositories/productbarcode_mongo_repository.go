package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IProductBarcodeRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.ProductBarcodeDoc) (string, error)
	CreateInBatch(docList []models.ProductBarcodeDoc) error
	Update(shopID string, guid string, doc models.ProductBarcodeDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ProductBarcodeDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ProductBarcodeItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.ProductBarcodeDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.ProductBarcodeInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.ProductBarcodeDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.ProductBarcodeActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.ProductBarcodeDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.ProductBarcodeActivity, error)
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
