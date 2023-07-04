package repositories

import (
	"errors"
	"os"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/unit/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IUnitRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.UnitDoc) (string, error)
	CreateInBatch(docList []models.UnitDoc) error
	Update(shopID string, guid string, doc models.UnitDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.UnitInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.UnitDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.UnitDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.UnitItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.UnitDoc, error)

	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.UnitInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.UnitInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.UnitDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.UnitActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.UnitDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.UnitActivity, error)
	FindMasterInCodes(codes []string) ([]models.UnitInfo, error)

	Transaction(callback func() error) error
}

type UnitRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.UnitDoc]
	repositories.SearchRepository[models.UnitInfo]
	repositories.GuidRepository[models.UnitItemGuid]
	repositories.ActivityRepository[models.UnitActivity, models.UnitDeleteActivity]
}

func NewUnitRepository(pst microservice.IPersisterMongo) *UnitRepository {

	insRepo := &UnitRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.UnitDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.UnitInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.UnitItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.UnitActivity, models.UnitDeleteActivity](pst)

	return insRepo
}

func (repo UnitRepository) FindMasterInCodes(codes []string) ([]models.UnitInfo, error) {

	masterShopID := os.Getenv("MASTER_SHOP_ID")

	if len(masterShopID) == 0 {
		return []models.UnitInfo{}, errors.New("master shop id is empty")
	}

	docList := []models.UnitInfo{}

	filters := bson.M{
		"shopid": masterShopID,
		"unitcode": bson.M{
			"$in": codes,
		},
	}

	err := repo.pst.Find([]models.UnitInfo{}, filters, &docList)

	if err != nil {
		return nil, err
	}

	return docList, nil
}

func (repo UnitRepository) Transaction(callback func() error) error {
	return repo.pst.Transaction(callback)
}
