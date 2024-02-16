package chartofaccountadmin

import (
	"context"
	chartOfAccountRepositories "smlcloudplatform/internal/vfgl/chartofaccount/repositories"
	"smlcloudplatform/pkg/microservice"
	"time"
)

type IChartOfAccountAdminService interface {
	ReSyncChartOfAccountDoc(shopID string) error
}

type ChartOfAccountAdminService struct {
	kafkaRepo       chartOfAccountRepositories.IChartOfAccountMQRepository
	mongoRepo       IChartOfAccountAdminRepository
	timeoutDuration time.Duration
}

func NewChartOfAccountAdminService(pst microservice.IPersisterMongo, kfProducer microservice.IProducer) IChartOfAccountAdminService {
	kafkaRepo := chartOfAccountRepositories.NewChartOfAccountMQRepository(kfProducer)
	mongoRepo := NewChartOfAccountAdminRepository(pst)

	return &ChartOfAccountAdminService{
		kafkaRepo:       kafkaRepo,
		mongoRepo:       mongoRepo,
		timeoutDuration: time.Duration(30) * time.Second,
	}
}
func (s *ChartOfAccountAdminService) ReSyncChartOfAccountDoc(shopID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), s.timeoutDuration)
	defer cancel()

	docs, err := s.mongoRepo.FindChartOfAccountDocByShopID(ctx, shopID)
	if err != nil {
		return err
	}

	err = s.kafkaRepo.CreateInBatch(docs)
	if err != nil {
		return err
	}
	return nil
}
