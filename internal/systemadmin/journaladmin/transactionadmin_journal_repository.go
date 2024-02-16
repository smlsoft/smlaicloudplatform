package journal

import (
	"context"
	journalModels "smlcloudplatform/internal/vfgl/journal/models"
	"smlcloudplatform/pkg/microservice"

	"go.mongodb.org/mongo-driver/bson"
)

type IJournalTransactionAdminRepository interface {
	FindJournalTransactionDocByShopID(ctx context.Context, shopID string) ([]journalModels.JournalDoc, error)
}

type JournalTransactionAdminRepository struct {
	pst microservice.IPersisterMongo
}

func NewJournalTransactionAdminRepository(pst microservice.IPersisterMongo) IJournalTransactionAdminRepository {
	return &JournalTransactionAdminRepository{
		pst: pst,
	}
}

func (r JournalTransactionAdminRepository) FindJournalTransactionDocByShopID(ctx context.Context, shopID string) ([]journalModels.JournalDoc, error) {

	docList := []journalModels.JournalDoc{}

	err := r.pst.Find(ctx, &journalModels.JournalDoc{}, bson.M{"shopid": shopID}, &docList)
	if err != nil {
		return nil, err
	}

	return docList, nil
}
