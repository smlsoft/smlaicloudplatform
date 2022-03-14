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
	CreatePurchase(shopId string, username string, doc *models.Purchase) (string, error)
	UpdatePurchase(guid string, shopId string, username string, doc models.Purchase) error
	DeletePurchase(guid string, shopId string, username string) error
	InfoPurchase(guid string, shopId string) (models.Purchase, error)
	SearchPurchase(shopId string, q string, page int, limit int) ([]models.Purchase, paginate.PaginationData, error)
	SearchItemsPurchase(guid string, shopId string, q string, page int, limit int) ([]models.Purchase, paginate.PaginationData, error)
}

type PurchaseService struct {
	repo   IPurchaseRepository
	mqRepo IPurchaseMQRepository
}

func NewPurchaseService(repo IPurchaseRepository, mqRepo IPurchaseMQRepository) IPurchaseService {

	return &PurchaseService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc *PurchaseService) CreatePurchase(shopId string, username string, doc *models.Purchase) (string, error) {

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

	docReq := &models.PurchaseRequest{}
	docReq.MapRequest(*doc)

	err = svc.mqRepo.Create(*docReq)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc *PurchaseService) UpdatePurchase(guid string, shopId string, username string, doc models.Purchase) error {

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

func (svc *PurchaseService) DeletePurchase(guid string, shopId string, username string) error {

	err := svc.repo.Delete(guid, shopId)
	if err != nil {
		return err
	}

	return nil
}

func (svc *PurchaseService) InfoPurchase(guid string, shopId string) (models.Purchase, error) {
	doc, err := svc.repo.FindByGuid(guid, shopId)

	if err != nil {
		return models.Purchase{}, err
	}

	return doc, nil
}

func (svc *PurchaseService) SearchPurchase(shopId string, q string, page int, limit int) ([]models.Purchase, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopId, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}

func (svc *PurchaseService) SearchItemsPurchase(guid string, shopId string, q string, page int, limit int) ([]models.Purchase, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindItemsByGuidPage(guid, shopId, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
