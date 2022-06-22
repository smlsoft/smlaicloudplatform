package services

import (
	"smlcloudplatform/pkg/mastersync/repositories"
	"time"
)

type IMasterSyncService interface {
	GetStatus(shopID string) (map[string]time.Time, error)
}

type MasterSyncService struct {
	cacheRepo repositories.IMasterSyncCacheRepository
}

func NewMasterSyncService(cacheRepo repositories.IMasterSyncCacheRepository) MasterSyncService {

	return MasterSyncService{
		cacheRepo: cacheRepo,
	}
}

func (svc MasterSyncService) GetStatus(shopID string) (map[string]time.Time, error) {
	syncModules := []string{"all", "category", "member", "inventory", "kitchen", "shopprinter", "shoptable", "shopzone"}

	moduleStatus := map[string]time.Time{}

	for _, moduleName := range syncModules {
		lastTime, err := svc.cacheRepo.GetWithModule(shopID, moduleName)
		moduleStatus[moduleName] = lastTime
		if err != nil {
			return nil, err
		}
	}

	return moduleStatus, nil
}
