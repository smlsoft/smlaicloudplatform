package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/journal/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IJournalRepository interface {
	Count(shopID string) (int, error)
	Create(category models.JournalDoc) (string, error)
	CreateInBatch(docList []models.JournalDoc) error
	Update(shopID string, guid string, category models.JournalDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.JournalInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.JournalDoc, error)
	FindOne(shopID string, filters interface{}) (models.JournalDoc, error)
	IsAccountCodeUsed(shopID string, accountCode string) (bool, error)

	// FindLastDocno(shopID string, docFormat string) (string, error)
}

type JournalRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.JournalDoc]
	repositories.SearchRepository[models.JournalInfo]
	repositories.GuidRepository[models.JournalItemGuid]
}

func NewJournalRepository(pst microservice.IPersisterMongo) JournalRepository {

	insRepo := JournalRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.JournalDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.JournalInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.JournalItemGuid](pst)

	return insRepo
}

func (repo *JournalRepository) IsAccountCodeUsed(shopID string, accountCode string) (bool, error) {

	findDoc := models.JournalDoc{}

	filters := bson.M{
		"shopid":                    shopID,
		"journaldetail.accountcode": accountCode,
		"deletedat":                 bson.M{"$exists": false},
	}

	err := repo.pst.FindOne(models.JournalDoc{}, filters, &findDoc)

	if err != nil {
		return true, nil
	}

	return findDoc.ID != primitive.NilObjectID, nil

}

func (repo *JournalRepository) FindLastDocno(shopID string, docFormat string) (string, error) {

	findDocList := []models.JournalDoc{}

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	if len(docFormat) < 1 {
		filters["$or"] = []interface{}{
			bson.M{"docformat": ""},
			bson.M{"docformat": bson.M{"$exists": false}},
		}
	} else {
		filters["docformat"] = docFormat
	}

	findOptions := options.Find()

	findOptions.SetSort(bson.M{"docformat": -1})
	findOptions.SetLimit(1)

	err := repo.pst.Find(models.JournalDoc{}, filters, &findDocList, findOptions)

	if err != nil {
		return "", nil
	}

	if len(findDocList) < 1 {
		return "", nil
	}

	return findDocList[0].DocNo, nil

}
