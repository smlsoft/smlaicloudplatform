package inventoryimport

import (
	"context"
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
	repo           ICategoryImportRepository
	contextTimeout time.Duration
}

func NewCategoryImportService(repository ICategoryImportRepository) CategoryImportService {

	contextTimeout := time.Duration(15) * time.Second

	return CategoryImportService{
		repo:           repository,
		contextTimeout: contextTimeout,
	}
}

func (svc CategoryImportService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc CategoryImportService) CreateInBatch(shopID string, authUsername string, categories []models.CategoryImport) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

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
	svc.repo.DeleteInBatchCode(ctx, shopID, codeList)

	err := svc.repo.CreateInBatch(ctx, tempInvDataList)

	if err != nil {
		return err
	}

	return nil

}

func (svc CategoryImportService) Delete(shopID string, guidList []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	err := svc.repo.DeleteInBatch(ctx, shopID, guidList)

	if err != nil {
		return err
	}

	return nil
}

func (svc CategoryImportService) List(shopID string, pageable micromodels.Pageable) ([]models.CategoryImportInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList, pagination, err := svc.repo.FindPage(ctx, shopID, pageable)

	if err != nil {
		return []models.CategoryImportInfo{}, pagination, err
	}

	return docList, pagination, nil
}
