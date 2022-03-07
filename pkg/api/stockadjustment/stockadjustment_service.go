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
	CreateStockAdjustment(merchantId string, username string, doc *models.StockAdjustment) (string, error)
	UpdateStockAdjustment(guid string, merchantId string, username string, doc models.StockAdjustment) error
	DeleteStockAdjustment(guid string, merchantId string, username string) error
	InfoStockAdjustment(guid string, merchantId string) (models.StockAdjustment, error)
	SearchStockAdjustment(merchantId string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error)
	SearchItemsStockAdjustment(guid string, merchantId string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error)
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

func (svc *StockAdjustmentService) CreateStockAdjustment(merchantId string, username string, doc *models.StockAdjustment) (string, error) {

	sumAmount := 0.0
	for i, docDetail := range doc.Items {
		doc.Items[i].LineNumber = i + 1
		sumAmount += docDetail.Price * docDetail.Qty
	}

	newGuidFixed := utils.NewGUID()
	doc.MerchantId = merchantId
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

func (svc *StockAdjustmentService) UpdateStockAdjustment(guid string, merchantId string, username string, doc models.StockAdjustment) error {

	findDoc, err := svc.repo.FindByGuid(guid, merchantId)

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

func (svc *StockAdjustmentService) DeleteStockAdjustment(guid string, merchantId string, username string) error {

	err := svc.repo.Delete(guid, merchantId)
	if err != nil {
		return err
	}

	return nil
}

func (svc *StockAdjustmentService) InfoStockAdjustment(guid string, merchantId string) (models.StockAdjustment, error) {
	doc, err := svc.repo.FindByGuid(guid, merchantId)

	if err != nil {
		return models.StockAdjustment{}, err
	}

	return doc, nil
}

func (svc *StockAdjustmentService) SearchStockAdjustment(merchantId string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(merchantId, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc *StockAdjustmentService) SearchItemsStockAdjustment(guid string, merchantId string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindItemsByGuidPage(guid, merchantId, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
