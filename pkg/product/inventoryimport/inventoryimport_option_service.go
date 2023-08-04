package inventoryimport

import (
	"context"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/inventoryimport/models"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
)

type IInventoryOptionMainImportService interface {
	CreateInBatch(shopID string, authUsername string, options []models.InventoryOptionMainImport) error
	Delete(shopID string, guidList []string) error
	ListInventory(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionMainImportInfo, mongopagination.PaginationData, error)
}

type InventoryOptionMainImportService struct {
	repo           IInventoryOptionMainImportRepository
	contextTimeout time.Duration
}

func NewInventoryOptionMainImportService(invImportOptionMainRepository IInventoryOptionMainImportRepository) InventoryOptionMainImportService {

	contextTimeout := time.Duration(15) * time.Second

	return InventoryOptionMainImportService{
		repo:           invImportOptionMainRepository,
		contextTimeout: contextTimeout,
	}
}

func (svc InventoryOptionMainImportService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc InventoryOptionMainImportService) CreateInBatch(shopID string, authUsername string, options []models.InventoryOptionMainImport) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	codeList := []string{}
	tempInvDataList := []models.InventoryOptionMainImportDoc{}

	for _, opt := range options {

		codeList = append(codeList, opt.Code)

		newGuid := utils.NewGUID()

		invDoc := models.InventoryOptionMainImportDoc{}

		invDoc.GuidFixed = newGuid
		invDoc.ShopID = shopID
		invDoc.InventoryOptionMainImport = opt

		invDoc.CreatedBy = authUsername
		invDoc.CreatedAt = time.Now()

		tempInvDataList = append(tempInvDataList, invDoc)
	}
	//Clear old items
	svc.repo.DeleteInBatchCode(ctx, shopID, codeList)

	err := svc.repo.CreateInBatch(ctx, tempInvDataList)

	if err != nil {
		return err
	}

	return nil

}

func (svc InventoryOptionMainImportService) Delete(shopID string, guidList []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	err := svc.repo.DeleteInBatch(ctx, shopID, guidList)

	if err != nil {
		return err
	}

	return nil
}

func (svc InventoryOptionMainImportService) ListInventory(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionMainImportInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList, pagination, err := svc.repo.FindPage(ctx, shopID, pageable)

	if err != nil {
		return []models.InventoryOptionMainImportInfo{}, pagination, err
	}

	return docList, pagination, nil
}
