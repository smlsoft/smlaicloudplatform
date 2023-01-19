package inventoryimport

import (
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/inventoryimport/models"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
)

type ICategoryImportService interface {
	CreateInBatch(shopID string, authUsername string, options []models.CategoryImport) error
	Delete(shopID string, guidList []string) error
	List(shopID string, pageable micromodels.Pageable) ([]models.CategoryImportInfo, mongopagination.PaginationData, error)
}

type CategoryImportService struct {
	repo ICategoryImportRepository
}

func NewCategoryImportService(repository ICategoryImportRepository) CategoryImportService {
	return CategoryImportService{
		repo: repository,
	}
}

func (svc CategoryImportService) CreateInBatch(shopID string, authUsername string, categories []models.CategoryImport) error {

	codeList := []string{}
	tempInvDataList := []models.CategoryImportDoc{}

	for _, category := range categories {
		codeList = append(codeList, category.GuidFixed)
		newGuid := utils.NewGUID()

		invDoc := models.CategoryImportDoc{}

		invDoc.GuidFixed = newGuid
		invDoc.ShopID = shopID
		invDoc.CategoryImport = category

		invDoc.CreatedBy = authUsername
		invDoc.CreatedAt = time.Now()

		tempInvDataList = append(tempInvDataList, invDoc)
	}
	//Clear old items
	svc.repo.DeleteInBatchCode(shopID, codeList)

	err := svc.repo.CreateInBatch(tempInvDataList)

	if err != nil {
		return err
	}

	return nil

}

func (svc CategoryImportService) Delete(shopID string, guidList []string) error {

	err := svc.repo.DeleteInBatch(shopID, guidList)

	if err != nil {
		return err
	}

	return nil
}

func (svc CategoryImportService) List(shopID string, pageable micromodels.Pageable) ([]models.CategoryImportInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, pageable)

	if err != nil {
		return []models.CategoryImportInfo{}, pagination, err
	}

	return docList, pagination, nil
}
