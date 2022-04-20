package inventoryimport

import (
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
)

type IInventoryImportService interface {
	CreateInBatch(shopID string, authUsername string, inventories []models.InventoryImport) error
	Delete(shopID string, guidList []string) error
	ListInventory(shopID string, page int, limit int) ([]models.InventoryImportInfo, paginate.PaginationData, error)
}

type InventoryImportService struct {
	invRepo IInventoryImportRepository
}

func NewInventoryImportService(inventoryRepo IInventoryImportRepository) InventoryImportService {
	return InventoryImportService{
		invRepo: inventoryRepo,
	}
}

func (svc InventoryImportService) CreateInBatch(shopID string, authUsername string, inventories []models.InventoryImport) error {

	tempInvDataList := []models.InventoryImportDoc{}

	for _, inventory := range inventories {
		newGuid := utils.NewGUID()

		invDoc := models.InventoryImportDoc{}

		invDoc.GuidFixed = newGuid
		invDoc.ShopID = shopID
		invDoc.InventoryImport = inventory

		invDoc.CreatedBy = authUsername
		invDoc.CreatedAt = time.Now()

		tempInvDataList = append(tempInvDataList, invDoc)
	}

	err := svc.invRepo.CreateInBatch(tempInvDataList)

	if err != nil {
		return err
	}

	return nil

}

func (svc InventoryImportService) Delete(shopID string, guidList []string) error {

	err := svc.invRepo.DeleteInBatch(shopID, guidList)

	if err != nil {
		return err
	}

	return nil
}

func (svc InventoryImportService) ListInventory(shopID string, page int, limit int) ([]models.InventoryImportInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.invRepo.FindPage(shopID, page, limit)

	if err != nil {
		return []models.InventoryImportInfo{}, pagination, err
	}

	return docList, pagination, nil
}
