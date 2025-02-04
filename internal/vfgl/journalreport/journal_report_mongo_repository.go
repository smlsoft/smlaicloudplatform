package journalreport

import (
	"context"
	"smlaicloudplatform/internal/vfgl/journalreport/models"
	"smlaicloudplatform/pkg/microservice"

	"go.mongodb.org/mongo-driver/bson"
)

type IJournalReportMongoRepository interface {
	FindCountDetailByDocs(ctx context.Context, shopID string, docs []string) ([]models.JournalSummary, error)
	FindCountImageByDocs(ctx context.Context, shopID string, docs []string) ([]models.JournalImageSummary, error)
}

type JournalMongoRepository struct {
	pst microservice.IPersisterMongo
}

func NewJournalMongoRepository(pst microservice.IPersisterMongo) *JournalMongoRepository {

	insRepo := &JournalMongoRepository{
		pst: pst,
	}

	return insRepo
}

func (repo *JournalMongoRepository) FindCountDetailByDocs(ctx context.Context, shopID string, docs []string) ([]models.JournalSummary, error) {

	matchQuery := bson.M{
		"shopid": shopID,
		"docno":  bson.M{"$in": docs},
		"vats":   bson.M{"$exists": true, "$type": "array"},
		"taxes":  bson.M{"$exists": true, "$type": "array"},
	}

	projectQuery := bson.M{
		"docno":    1,
		"countvat": bson.M{"$size": "$vats"},
		"counttax": bson.M{"$size": "$taxes"},
	}

	query := []interface{}{
		bson.M{"$match": matchQuery},
		bson.M{"$project": projectQuery},
	}

	docList := []models.JournalSummary{}
	err := repo.pst.Aggregate(ctx, &models.JournalSummary{}, query, &docList)

	if err != nil {
		return []models.JournalSummary{}, err
	}

	return docList, nil
}

func (repo *JournalMongoRepository) FindCountImageByDocs(ctx context.Context, shopID string, docs []string) ([]models.JournalImageSummary, error) {

	matchQuery := bson.M{
		"shopid":            shopID,
		"imagereferences":   bson.M{"$exists": true},
		"references.module": "GL",
		"references.docno":  bson.M{"$in": docs},
	}

	projectQuery := bson.M{
		"docno":      bson.M{"$first": "$references.docno"},
		"countimage": bson.M{"$size": "$imagereferences"},
	}

	query := []interface{}{
		bson.M{"$match": matchQuery},
		bson.M{"$project": projectQuery},
	}

	docList := []models.JournalImageSummary{}
	err := repo.pst.Aggregate(ctx, &models.JournalImageSummary{}, query, &docList)

	if err != nil {
		return []models.JournalImageSummary{}, err
	}

	return docList, nil
}
