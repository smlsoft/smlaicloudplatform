package repositories

import (
	"context"
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/vfgl/journal/models"
	"smlaicloudplatform/pkg/microservice"
	micromodels "smlaicloudplatform/pkg/microservice/models"

	"github.com/smlsoft/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IJournalRepository interface {
	FindAll(ctx context.Context) ([]models.JournalDoc, error)
	Count(ctx context.Context, shopID string) (int, error)
	Create(ctx context.Context, doc models.JournalDoc) (string, error)
	CreateInBatch(ctx context.Context, docList []models.JournalDoc) error
	Update(ctx context.Context, shopID string, guid string, doc models.JournalDoc) error
	DeleteByGuidfixed(ctx context.Context, shopID string, guid string, username string) error
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.JournalInfo, mongopagination.PaginationData, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.JournalDoc, error)
	FindOne(ctx context.Context, shopID string, filters interface{}) (models.JournalDoc, error)
	FindFilter(ctx context.Context, shopID string, filters map[string]interface{}) ([]models.JournalDoc, error)
	IsAccountCodeUsed(ctx context.Context, shopID string, accountCode string) (bool, error)
	FindGUIDEmptyAll() []models.JournalDoc
	UpdateGuidEmpty(ctx context.Context, id string, guidfixed string) error
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

func (repo *JournalRepository) IsAccountCodeUsed(ctx context.Context, shopID string, accountCode string) (bool, error) {

	findDoc := models.JournalDoc{}

	filters := bson.M{
		"shopid":                    shopID,
		"journaldetail.accountcode": accountCode,
		"deletedat":                 bson.M{"$exists": false},
	}

	err := repo.pst.FindOne(ctx, models.JournalDoc{}, filters, &findDoc)

	if err != nil {
		return true, nil
	}

	return findDoc.ID != primitive.NilObjectID, nil

}

func (repo *JournalRepository) FindLastDocno(ctx context.Context, shopID string, docFormat string) (string, error) {

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

	err := repo.pst.Find(ctx, models.JournalDoc{}, filters, &findDocList, findOptions)

	if err != nil {
		return "", nil
	}

	if len(findDocList) < 1 {
		return "", nil
	}

	return findDocList[0].DocNo, nil

}

func (repo *JournalRepository) FindGUIDEmptyAll() ([]models.JournalDoc, error) {
	findDocList := []models.JournalDoc{}

	filters := bson.M{
		"guidfixed": "",
	}

	err := repo.pst.Find(context.Background(), models.JournalDoc{}, filters, &findDocList)

	if err != nil {
		return []models.JournalDoc{}, nil
	}

	return findDocList, nil
}

func (repo *JournalRepository) UpdateGuidEmpty(ctx context.Context, id primitive.ObjectID, guidfixed string) error {

	err := repo.pst.UpdateOne(ctx, models.JournalDoc{}, bson.M{"_id": id}, bson.M{"guidfixed": guidfixed})

	if err != nil {
		return err
	}

	return nil
}
