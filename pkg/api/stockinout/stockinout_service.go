package stockinout

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStockInOutService interface {
	CreateStockInOut(shopID string, username string, doc models.StockInOut) (string, error)
	UpdateStockInOut(guid string, shopID string, username string, doc models.StockInOut) error
	DeleteStockInOut(guid string, shopID string, username string) error
	InfoStockInOut(guid string, shopID string) (models.StockInOutInfo, error)
	SearchStockInOut(shopID string, q string, page int, limit int) ([]models.StockInOutInfo, paginate.PaginationData, error)
	SearchItemsStockInOut(guid string, shopID string, q string, page int, limit int) ([]models.StockInOutInfo, paginate.PaginationData, error)
}

type StockInOutService struct {
	repo   IStockInOutRepository
	mqRepo IStockInOutMQRepository
}

func NewStockInOutService(repo IStockInOutRepository, mqRepo IStockInOutMQRepository) StockInOutService {
	return StockInOutService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc StockInOutService) CreateStockInOut(shopID string, username string, doc models.StockInOut) (string, error) {

	sumAmount := 0.0
	for i, docDetail := range doc.Items {
		doc.Items[i].LineNumber = i + 1
		sumAmount += docDetail.Price * docDetail.Qty
	}

	newGuidFixed := utils.NewGUID()

	docData := models.StockInOutDoc{}

	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SumAmount = sumAmount
	docData.StockInOut = doc

	docData.CreatedBy = username
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	err = svc.mqRepo.Create(docData.StockInOutData)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc StockInOutService) UpdateStockInOut(guid string, shopID string, username string, doc models.StockInOut) error {

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

	findDoc.StockInOut = doc
	findDoc.SumAmount = sumAmount

	findDoc.UpdatedBy = username
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc StockInOutService) DeleteStockInOut(guid string, shopID string, username string) error {
	err := svc.repo.Delete(guid, shopID)
	if err != nil {
		return err
	}

	return nil
}

func (svc StockInOutService) InfoStockInOut(guid string, shopID string) (models.StockInOutInfo, error) {
	doc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return models.StockInOutInfo{}, err
	}

	return doc.StockInOutInfo, nil
}

func (svc StockInOutService) SearchStockInOut(shopID string, q string, page int, limit int) ([]models.StockInOutInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockInOutService) SearchItemsStockInOut(guid string, shopID string, q string, page int, limit int) ([]models.StockInOutInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindItemsByGuidPage(guid, shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
