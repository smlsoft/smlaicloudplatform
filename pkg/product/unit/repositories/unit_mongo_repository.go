package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/unit/models"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IUnitRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.UnitDoc) (string, error)
	CreateInBatch(docList []models.UnitDoc) error
	Update(shopID string, guid string, doc models.UnitDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.UnitInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.UnitDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.UnitItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.UnitDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.UnitInfo, mongopagination.PaginationData, error)
}

type UnitRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.UnitDoc]
	repositories.SearchRepository[models.UnitInfo]
	repositories.GuidRepository[models.UnitItemGuid]
}

func NewUnitRepository(pst microservice.IPersisterMongo) *UnitRepository {

	insRepo := &UnitRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.UnitDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.UnitInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.UnitItemGuid](pst)

	return insRepo
}
