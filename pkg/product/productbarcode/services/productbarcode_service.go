package services

import (
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/product/productbarcode/repositories"
	"sync"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
)

type IProductBarcodeService interface {
	LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, paginate.PaginationData, error)
}

type ProductBarcodeService struct {
	invRepo   repositories.IProductBarcodeRepository
	cacheRepo mastersync.IMasterSyncCacheRepository
}

func NewProductBarcodeService(productbarcodeRepo repositories.IProductBarcodeRepository, cacheRepo mastersync.IMasterSyncCacheRepository) ProductBarcodeService {
	return ProductBarcodeService{
		invRepo:   productbarcodeRepo,
		cacheRepo: cacheRepo,
	}
}

func (svc ProductBarcodeService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, paginate.PaginationData, error) {

	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.ProductBarcodeDeleteActivity
	var pagination1 paginate.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.invRepo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.ProductBarcodeActivity
	var pagination2 paginate.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.invRepo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return common.LastActivity{}, pagination1, err1
	}

	if err2 != nil {
		return common.LastActivity{}, pagination2, err2
	}

	lastActivity := common.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	pagination := pagination1

	if pagination.Total < pagination2.Total {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}
