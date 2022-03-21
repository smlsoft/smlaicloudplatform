package stockadjustment

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStockAdjustmentService interface {
	CreateStockAdjustment(shopID string, username string, doc *models.StockAdjustment) (string, error)
	UpdateStockAdjustment(guid string, shopID string, username string, doc models.StockAdjustment) error
	DeleteStockAdjustment(guid string, shopID string, username string) error
	InfoStockAdjustment(guid string, shopID string) (models.StockAdjustment, error)
	SearchStockAdjustment(shopID string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error)
	SearchItemsStockAdjustment(guid string, shopID string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error)
}

type StockAdjustmentService struct {
	repo   IStockAdjustmentRepository
	mqRepo IStockAdjustmentMQRepository
}

func NewStockAdjustmentService(repo IStockAdjustmentRepository, mqRepo IStockAdjustmentMQRepository) StockAdjustmentService {

	return StockAdjustmentService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc StockAdjustmentService) CreateStockAdjustment(shopID string, username string, doc *models.StockAdjustment) (string, error) {

	sumAmount := 0.0
	for i, docDetail := range doc.Items {
		doc.Items[i].LineNumber = i + 1
		sumAmount += docDetail.Price * docDetail.Qty
	}

	newGuidFixed := utils.NewGUID()
	doc.ShopID = shopID
	doc.GuidFixed = newGuidFixed
	doc.SumAmount = sumAmount
	doc.Deleted = false
	doc.CreatedBy = username
	doc.CreatedAt = time.Now()

	idx, err := svc.repo.Create(*doc)

	if err != nil {
		return "", err
	}

	doc.ID = idx

	docReq := &models.StockAdjustmentRequest{}
	docReq.MapRequest(*doc)

	err = svc.mqRepo.Create(*docReq)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc StockAdjustmentService) UpdateStockAdjustment(guid string, shopID string, username string, doc models.StockAdjustment) error {

	findDoc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("guid invalid")
	}

	sumAmount := 0.0
	for i, docDetail := range doc.Items {
		findDoc.Items[i].LineNumber = i + 1
		sumAmount += docDetail.Price * docDetail.Qty
	}

	findDoc.Items = doc.Items
	findDoc.SumAmount = sumAmount
	findDoc.UpdatedBy = username
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc StockAdjustmentService) DeleteStockAdjustment(guid string, shopID string, username string) error {

	err := svc.repo.Delete(guid, shopID)
	if err != nil {
		return err
	}

	return nil
}

func (svc StockAdjustmentService) InfoStockAdjustment(guid string, shopID string) (models.StockAdjustment, error) {
	doc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return models.StockAdjustment{}, err
	}

	return doc, nil
}

func (svc StockAdjustmentService) SearchStockAdjustment(shopID string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockAdjustmentService) SearchItemsStockAdjustment(guid string, shopID string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindItemsByGuidPage(guid, shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
