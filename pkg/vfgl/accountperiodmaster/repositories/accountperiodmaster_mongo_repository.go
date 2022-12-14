package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/models"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IAccountPeriodMasterRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.AccountPeriodMasterDoc) (string, error)
	CreateInBatch(docList []models.AccountPeriodMasterDoc) error
	Update(shopID string, guid string, doc models.AccountPeriodMasterDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.AccountPeriodMasterInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.AccountPeriodMasterDoc, error)
	FindByDateRange(shopID string, startDate time.Time, endDate time.Time) (models.AccountPeriodMasterDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.AccountPeriodMasterItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.AccountPeriodMasterDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.AccountPeriodMasterInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.AccountPeriodMasterInfo, int, error)
}

type AccountPeriodMasterRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.AccountPeriodMasterDoc]
	repositories.SearchRepository[models.AccountPeriodMasterInfo]
	repositories.GuidRepository[models.AccountPeriodMasterItemGuid]
}

func NewAccountPeriodMasterRepository(pst microservice.IPersisterMongo) *AccountPeriodMasterRepository {

	insRepo := &AccountPeriodMasterRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.AccountPeriodMasterDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.AccountPeriodMasterInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.AccountPeriodMasterItemGuid](pst)

	return insRepo
}

func (repo AccountPeriodMasterRepository) FindByDateRange(shopID string, startDate time.Time, endDate time.Time) (models.AccountPeriodMasterDoc, error) {

	filterQuery := bson.D{
		bson.E{Key: "$or", Value: bson.A{
			bson.D{{"startdate", bson.D{{"$gte", startDate}}}},
			bson.D{{"enddate", bson.D{{"$gte", startDate}}}},
		}},
		bson.E{Key: "$or", Value: bson.A{
			bson.D{{"startdate", bson.D{{"$lte", endDate}}}},
			bson.D{{"enddate", bson.D{{"$lte", endDate}}}},
		}},
	}

	finDoc, err := repo.FindOne(shopID, filterQuery)

	if err != nil {
		return models.AccountPeriodMasterDoc{}, err
	}

	return finDoc, nil
}
