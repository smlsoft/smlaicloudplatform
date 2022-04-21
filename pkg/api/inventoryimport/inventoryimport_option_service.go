package inventoryimport

import (
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
)

type IInventoryOptionMainImportService interface {
	CreateInBatch(shopID string, authUsername string, options []models.InventoryOptionMainImport) error
	Delete(shopID string, guidList []string) error
	ListInventory(shopID string, page int, limit int) ([]models.InventoryOptionMainImportInfo, paginate.PaginationData, error)
}

type InventoryOptionMainImportService struct {
	repo IInventoryOptionMainImportRepository
}

func NewInventoryOptionMainImportService(invImportOptionMainRepository IInventoryOptionMainImportRepository) InventoryOptionMainImportService {
	return InventoryOptionMainImportService{
		repo: invImportOptionMainRepository,
	}
}

func (svc InventoryOptionMainImportService) CreateInBatch(shopID string, authUsername string, options []models.InventoryOptionMainImport) error {

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
	//Clear old item
	svc.repo.DeleteInBatch(shopID, codeList)

	err := svc.repo.CreateInBatch(tempInvDataList)

	if err != nil {
		return err
	}

	return nil

}

func (svc InventoryOptionMainImportService) Delete(shopID string, guidList []string) error {

	err := svc.repo.DeleteInBatch(shopID, guidList)

	if err != nil {
		return err
	}

	return nil
}

func (svc InventoryOptionMainImportService) ListInventory(shopID string, page int, limit int) ([]models.InventoryOptionMainImportInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, page, limit)

	if err != nil {
		return []models.InventoryOptionMainImportInfo{}, pagination, err
	}

	return docList, pagination, nil
}
