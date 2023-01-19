package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/shopdesign/zonedesign/models"

	"github.com/userplant/mongopagination"
)

type IZoneDesignRepository interface {
	Count(shopID string) (int, error)
	Create(category models.ZoneDesignDoc) (string, error)
	CreateInBatch(docList []models.ZoneDesignDoc) error
	Update(shopID string, guid string, category models.ZoneDesignDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters interface{}) (models.ZoneDesignDoc, error)
	FindByGuid(shopID string, guid string) (models.ZoneDesignDoc, error)
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ZoneDesignInfo, mongopagination.PaginationData, error)
}

type ZoneDesignRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ZoneDesignDoc]
	repositories.SearchRepository[models.ZoneDesignInfo]
}

func NewZoneDesignRepository(pst microservice.IPersisterMongo) ZoneDesignRepository {
	insRepo := ZoneDesignRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ZoneDesignDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ZoneDesignInfo](pst)

	return insRepo

}
