package purchase

import (
	"errors"
	"smlcloudplatform/pkg/transaction/purchase/models"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPurchaseService interface {
	CreatePurchase(shopID string, username string, doc models.Purchase) (string, error)
	UpdatePurchase(shopID string, guid string, username string, doc models.Purchase) error
	DeletePurchase(shopID string, guid string, username string) error
	InfoPurchase(shopID string, guid string) (models.PurchaseInfo, error)
	SearchPurchase(shopID string, q string, page int, limit int) ([]models.PurchaseInfo, mongopagination.PaginationData, error)
	SearchItemsPurchase(guid string, shopID string, q string, page int, limit int) ([]models.PurchaseInfo, mongopagination.PaginationData, error)
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

	for i, docDetail := range *doc.Items {
		docDetail.LineNumber = i + 1
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

func (svc PurchaseService) UpdatePurchase(shopID string, guid string, username string, doc models.Purchase) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("guid invalid")
	}

	findDoc.Purchase = doc

	sumAmount := 0.0
	for i, docDetail := range *findDoc.Items {
		docDetail.LineNumber = i + 1
		sumAmount += docDetail.Price * docDetail.Qty
	}

	findDoc.SumAmount = sumAmount

	findDoc.UpdatedBy = username
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc PurchaseService) DeletePurchase(shopID string, guid string, username string) error {

	err := svc.repo.Delete(shopID, guid, username)
	if err != nil {
		return err
	}

	return nil
}

func (svc PurchaseService) InfoPurchase(shopID string, guid string) (models.PurchaseInfo, error) {
	doc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.PurchaseInfo{}, err
	}

	return doc.PurchaseInfo, nil
}

func (svc PurchaseService) SearchPurchase(shopID string, q string, page int, limit int) ([]models.PurchaseInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc PurchaseService) SearchItemsPurchase(guid string, shopID string, q string, page int, limit int) ([]models.PurchaseInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repo.FindItemsByGuidPage(guid, shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
