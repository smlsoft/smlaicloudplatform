package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IAccountPeriodMasterRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.AccountPeriodMasterDoc) (string, error)
	CreateInBatch(docList []models.AccountPeriodMasterDoc) error
	Update(shopID string, guid string, doc models.AccountPeriodMasterDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.AccountPeriodMasterInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.AccountPeriodMasterDoc, error)
	FindByDateRange(shopID string, startDate time.Time, endDate time.Time) (models.AccountPeriodMasterDoc, error)
	FindByPeriod(shopID string, period int) (models.AccountPeriodMasterDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.AccountPeriodMasterItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.AccountPeriodMasterDoc, error)
	FindAll(shopID string) ([]models.AccountPeriodMasterDoc, error)

	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.AccountPeriodMasterInfo, int, error)
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
	endDate = endDate.AddDate(0, 0, 1)

	filterQuery := bson.D{
		bson.E{Key: "$or", Value: bson.A{
			bson.D{{"startdate", bson.D{{"$gte", startDate}}}},
			bson.D{{"enddate", bson.D{{"$gte", startDate}}}},
		}},
		bson.E{Key: "$or", Value: bson.A{
			bson.D{{"startdate", bson.D{{"$lt", endDate}}}},
			bson.D{{"enddate", bson.D{{"$lt", endDate}}}},
		}},
	}

	finDoc, err := repo.FindOne(shopID, filterQuery)

	if err != nil {
		return models.AccountPeriodMasterDoc{}, err
	}

	return finDoc, nil
}

func (repo AccountPeriodMasterRepository) FindByPeriod(shopID string, period int) (models.AccountPeriodMasterDoc, error) {

	filterQuery := bson.D{
		bson.E{"period", period},
	}

	finDoc, err := repo.FindOne(shopID, filterQuery)

	if err != nil {
		return models.AccountPeriodMasterDoc{}, err
	}

	return finDoc, nil
}

func (repo AccountPeriodMasterRepository) FindAll(shopID string) ([]models.AccountPeriodMasterDoc, error) {

	filterQuery := bson.M{
		"shopid":     shopID,
		"isdisabled": false,
		"deletedat":  bson.M{"$exists": false},
	}

	findDocList := []models.AccountPeriodMasterDoc{}

	err := repo.pst.Find(models.AccountPeriodMasterDoc{}, filterQuery, &findDocList)

	if err != nil {
		return []models.AccountPeriodMasterDoc{}, err
	}

	return findDocList, nil
}
