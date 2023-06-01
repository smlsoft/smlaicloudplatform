package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/optionpattern/models"
	"smlcloudplatform/pkg/repositories"

	"github.com/userplant/mongopagination"
)

type IOptionPatternRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.OptionPatternDoc) (string, error)
	CreateInBatch(docList []models.OptionPatternDoc) error
	Update(shopID string, guid string, doc models.OptionPatternDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.OptionPatternInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.OptionPatternDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.OptionPatternItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.OptionPatternDoc, error)
}

type OptionPatternRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.OptionPatternDoc]
	repositories.SearchRepository[models.OptionPatternInfo]
	repositories.GuidRepository[models.OptionPatternItemGuid]
}

func NewOptionPatternRepository(pst microservice.IPersisterMongo) *OptionPatternRepository {

	insRepo := &OptionPatternRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.OptionPatternDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.OptionPatternInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.OptionPatternItemGuid](pst)

	return insRepo
}
