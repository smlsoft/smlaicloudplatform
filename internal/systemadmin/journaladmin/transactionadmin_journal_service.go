package journal

import (
	"context"
	journalRepositories "smlcloudplatform/internal/vfgl/journal/repositories"
	"smlcloudplatform/pkg/microservice"
	"time"
)

type IJournalTransactionAdminService interface {
	ReSyncJournalTransactionDoc(shopID string) error
}

type JournalTransactionAdminService struct {
	mongoRepo       IJournalTransactionAdminRepository
	kafkaRepo       journalRepositories.IJournalMqRepository
	timeoutDuration time.Duration
}

func NewJournalTransactionAdminService(pst microservice.IPersisterMongo, kfProducer microservice.IProducer) IJournalTransactionAdminService {

	mongorepo := NewJournalTransactionAdminRepository(pst)
	kafkarepo := journalRepositories.NewJournalMqRepository(kfProducer)

	return &JournalTransactionAdminService{
		mongoRepo:       mongorepo,
		kafkaRepo:       kafkarepo,
		timeoutDuration: time.Duration(30) * time.Second,
	}
}

func (s *JournalTransactionAdminService) ReSyncJournalTransactionDoc(shopID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), s.timeoutDuration)
	defer cancel()

	docs, err := s.mongoRepo.FindJournalTransactionDocByShopID(ctx, shopID)
	if err != nil {
		return err
	}

	err = s.kafkaRepo.CreateInBatch(docs)
	if err != nil {
		return err
	}

	return nil
}
