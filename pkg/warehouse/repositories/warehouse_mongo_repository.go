package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/warehouse/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IWarehouseRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.WarehouseDoc) (string, error)
	CreateInBatch(docList []models.WarehouseDoc) error
	Update(shopID string, guid string, doc models.WarehouseDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.WarehouseInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.WarehouseDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.WarehouseItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.WarehouseDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.WarehouseInfo, mongopagination.PaginationData, error)
}

type WarehouseRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.WarehouseDoc]
	repositories.SearchRepository[models.WarehouseInfo]
	repositories.GuidRepository[models.WarehouseItemGuid]
}

func NewWarehouseRepository(pst microservice.IPersisterMongo) *WarehouseRepository {

	insRepo := &WarehouseRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.WarehouseDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.WarehouseInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.WarehouseItemGuid](pst)

	return insRepo
}
