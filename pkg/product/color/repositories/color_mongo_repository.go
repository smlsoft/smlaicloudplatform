package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/color/models"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IColorRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.ColorDoc) (string, error)
	CreateInBatch(docList []models.ColorDoc) error
	Update(shopID string, guid string, doc models.ColorDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.ColorInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ColorDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ColorItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.ColorDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.ColorInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.ColorInfo, int, error)
}

type ColorRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ColorDoc]
	repositories.SearchRepository[models.ColorInfo]
	repositories.GuidRepository[models.ColorItemGuid]
}

func NewColorRepository(pst microservice.IPersisterMongo) *ColorRepository {

	insRepo := &ColorRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ColorDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ColorInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ColorItemGuid](pst)

	return insRepo
}
