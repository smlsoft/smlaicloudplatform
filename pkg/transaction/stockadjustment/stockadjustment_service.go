package stockadjustment

import (
	"errors"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/transaction/stockadjustment/models"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStockAdjustmentService interface {
	CreateStockAdjustment(shopID string, username string, doc models.StockAdjustment) (string, error)
	UpdateStockAdjustment(guid string, shopID string, username string, doc models.StockAdjustment) error
	DeleteStockAdjustment(guid string, shopID string, username string) error
	InfoStockAdjustment(guid string, shopID string) (models.StockAdjustmentInfo, error)
	SearchStockAdjustment(shopID string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error)
	SearchItemsStockAdjustment(guid string, shopID string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error)
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

func (svc StockAdjustmentService) CreateStockAdjustment(shopID string, username string, doc models.StockAdjustment) (string, error) {

	sumAmount := 0.0
	for i, docDetail := range *doc.Items {
		docDetail.LineNumber = i + 1
		sumAmount += docDetail.Price * docDetail.Qty
	}

	newGuidFixed := utils.NewGUID()

	docData := models.StockAdjustmentDoc{}

	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SumAmount = sumAmount
	docData.StockAdjustment = doc

	docData.CreatedBy = username
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	err = svc.mqRepo.Create(docData.StockAdjustmentData)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc StockAdjustmentService) UpdateStockAdjustment(guid string, shopID string, username string, doc models.StockAdjustment) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("guid invalid")
	}

	sumAmount := 0.0
	for i, docDetail := range *doc.Items {
		docDetail.LineNumber = i + 1
		sumAmount += docDetail.Price * docDetail.Qty
	}

	findDoc.StockAdjustment = doc
	findDoc.SumAmount = sumAmount

	findDoc.UpdatedBy = username
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc StockAdjustmentService) DeleteStockAdjustment(guid string, shopID string, username string) error {

	err := svc.repo.Delete(shopID, guid, username)
	if err != nil {
		return err
	}

	return nil
}

func (svc StockAdjustmentService) InfoStockAdjustment(guid string, shopID string) (models.StockAdjustmentInfo, error) {
	doc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.StockAdjustmentInfo{}, err
	}

	return doc.StockAdjustmentInfo, nil
}

func (svc StockAdjustmentService) SearchStockAdjustment(shopID string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, pageable)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockAdjustmentService) SearchItemsStockAdjustment(guid string, shopID string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repo.FindItemsByGuidPage(shopID, guid, pageable)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
