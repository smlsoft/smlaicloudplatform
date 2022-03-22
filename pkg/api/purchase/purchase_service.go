package purchase

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPurchaseService interface {
	CreatePurchase(shopID string, username string, doc models.Purchase) (string, error)
	UpdatePurchase(guid string, shopID string, username string, doc models.Purchase) error
	DeletePurchase(guid string, shopID string, username string) error
	InfoPurchase(guid string, shopID string) (models.PurchaseInfo, error)
	SearchPurchase(shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error)
	SearchItemsPurchase(guid string, shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error)
}

type PurchaseService struct {
	repo   IPurchaseRepository
	mqRepo IPurchaseMQRepository
}

func NewPurchaseService(repo IPurchaseRepository, mqRepo IPurchaseMQRepository) PurchaseService {

	return PurchaseService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc PurchaseService) CreatePurchase(shopID string, username string, doc models.Purchase) (string, error) {

	sumAmount := 0.0
	for i, docDetail := range doc.Items {
		doc.Items[i].LineNumber = i + 1
		sumAmount += docDetail.Price * docDetail.Qty
	}

	newGuidFixed := utils.NewGUID()

	docData := models.PurchaseDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SumAmount = sumAmount
	docData.Purchase = doc

	docData.CreatedBy = username
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	err = svc.mqRepo.Create(docData.PurchaseData)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc PurchaseService) UpdatePurchase(guid string, shopID string, username string, doc models.Purchase) error {

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

func (svc PurchaseService) DeletePurchase(guid string, shopID string, username string) error {

	err := svc.repo.Delete(guid, shopID)
	if err != nil {
		return err
	}

	return nil
}

func (svc PurchaseService) InfoPurchase(guid string, shopID string) (models.PurchaseInfo, error) {
	doc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return models.PurchaseInfo{}, err
	}

	return doc.PurchaseInfo, nil
}

func (svc PurchaseService) SearchPurchase(shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc PurchaseService) SearchItemsPurchase(guid string, shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindItemsByGuidPage(guid, shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
