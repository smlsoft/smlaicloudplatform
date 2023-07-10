package inventoryimport

import (
	"context"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/inventoryimport/models"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
)

type IInventoryImportService interface {
	CreateInBatch(shopID string, authUsername string, inventories []models.InventoryImport) error
	Delete(shopID string, guidList []string) error
	ListInventory(shopID string, pageable micromodels.Pageable) ([]models.InventoryImportInfo, mongopagination.PaginationData, error)
}

type InventoryImportService struct {
	invRepo        IInventoryImportRepository
	contextTimeout time.Duration
}

func NewInventoryImportService(inventoryRepo IInventoryImportRepository) InventoryImportService {

	contextTimeout := time.Duration(15) * time.Second

	return InventoryImportService{
		invRepo:        inventoryRepo,
		contextTimeout: contextTimeout,
	}
}

func (svc InventoryImportService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc InventoryImportService) CreateInBatch(shopID string, authUsername string, inventories []models.InventoryImport) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	codeList := []string{}

	tempInvDataList := []models.InventoryImportDoc{}

	for _, inventory := range inventories {
		codeList = append(codeList, inventory.ItemCode)

		newGuid := utils.NewGUID()
		invDoc := models.InventoryImportDoc{}

		invDoc.GuidFixed = newGuid
		invDoc.ShopID = shopID
		invDoc.InventoryImport = inventory

		invDoc.CreatedBy = authUsername
		invDoc.CreatedAt = time.Now()

		tempInvDataList = append(tempInvDataList, invDoc)
	}
	//Clear old items
	svc.invRepo.DeleteInBatchCode(ctx, shopID, codeList)

	err := svc.invRepo.CreateInBatch(ctx, tempInvDataList)

	if err != nil {
		return err
	}

	return nil

}

func (svc InventoryImportService) Delete(shopID string, guidList []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	err := svc.invRepo.DeleteInBatch(ctx, shopID, guidList)

	if err != nil {
		return err
	}

	return nil
}

func (svc InventoryImportService) ListInventory(shopID string, pageable micromodels.Pageable) ([]models.InventoryImportInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList, pagination, err := svc.invRepo.FindPage(ctx, shopID, pageable)

	if err != nil {
		return []models.InventoryImportInfo{}, pagination, err
	}

	return docList, pagination, nil
}
