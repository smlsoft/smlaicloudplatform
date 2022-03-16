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
	CreateStockAdjustment(shopId string, username string, doc *models.StockAdjustment) (string, error)
	UpdateStockAdjustment(guid string, shopId string, username string, doc models.StockAdjustment) error
	DeleteStockAdjustment(guid string, shopId string, username string) error
	InfoStockAdjustment(guid string, shopId string) (models.StockAdjustment, error)
	SearchStockAdjustment(shopId string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error)
	SearchItemsStockAdjustment(guid string, shopId string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error)
}

type StockAdjustmentService struct {
	repo   IStockAdjustmentRepository
	mqRepo IStockAdjustmentMQRepository
}

func NewStockAdjustmentService(repo IStockAdjustmentRepository, mqRepo IStockAdjustmentMQRepository) IStockAdjustmentService {

	return &StockAdjustmentService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc *StockAdjustmentService) CreateStockAdjustment(shopId string, username string, doc *models.StockAdjustment) (string, error) {

	sumAmount := 0.0
	for i, docDetail := range doc.Items {
		doc.Items[i].LineNumber = i + 1
		sumAmount += docDetail.Price * docDetail.Qty
	}

	newGuidFixed := utils.NewGUID()
	doc.ShopId = shopId
	doc.GuidFixed = newGuidFixed
	doc.SumAmount = sumAmount
	doc.Deleted = false
	doc.CreatedBy = username
	doc.CreatedAt = time.Now()

	idx, err := svc.repo.Create(*doc)

	if err != nil {
		return "", err
	}

	doc.Id = idx

	docReq := &models.StockAdjustmentRequest{}
	docReq.MapRequest(*doc)

	err = svc.mqRepo.Create(*docReq)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc *StockAdjustmentService) UpdateStockAdjustment(guid string, shopId string, username string, doc models.StockAdjustment) error {

	findDoc, err := svc.repo.FindByGuid(guid, shopId)

	if err != nil {
		return err
	}

	if findDoc.Id == primitive.NilObjectID {
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

func (svc *StockAdjustmentService) DeleteStockAdjustment(guid string, shopId string, username string) error {

	err := svc.repo.Delete(guid, shopId)
	if err != nil {
		return err
	}

	return nil
}

func (svc *StockAdjustmentService) InfoStockAdjustment(guid string, shopId string) (models.StockAdjustment, error) {
	doc, err := svc.repo.FindByGuid(guid, shopId)

	if err != nil {
		return models.StockAdjustment{}, err
	}

	return doc, nil
}

func (svc *StockAdjustmentService) SearchStockAdjustment(shopId string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopId, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc *StockAdjustmentService) SearchItemsStockAdjustment(guid string, shopId string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindItemsByGuidPage(guid, shopId, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
