package saleinvoice

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISaleinvoiceService interface {
	CreateSaleinvoice(shopID string, username string, trans models.Saleinvoice) (string, error)
	UpdateSaleinvoice(guid string, shopID string, username string, trans models.Saleinvoice) error
	DeleteSaleinvoice(guid string, shopID string, username string) error
	InfoSaleinvoice(guid string, shopID string) (models.SaleinvoiceInfo, error)
	SearchSaleinvoice(shopID string, q string, page int, limit int) ([]models.SaleinvoiceInfo, paginate.PaginationData, error)
	SearchItemsSaleinvoice(guid string, shopID string, q string, page int, limit int) ([]models.SaleinvoiceInfo, paginate.PaginationData, error)
}

type SaleinvoiceService struct {
	saleinvoiceRepository ISaleinvoiceRepository
	mqRepo                ISaleinvoiceMQRepository
}

func NewSaleinvoiceService(saleinvoiceRepository ISaleinvoiceRepository, mqRepo ISaleinvoiceMQRepository) SaleinvoiceService {

	return SaleinvoiceService{
		saleinvoiceRepository: saleinvoiceRepository,
		mqRepo:                mqRepo,
	}
}

func (svc SaleinvoiceService) CreateSaleinvoice(shopID string, username string, trans models.Saleinvoice) (string, error) {

	sumAmount := 0.0
	for i, transDetail := range *trans.Items {
		transDetail.LineNumber = i + 1
		sumAmount += transDetail.Price * transDetail.Qty
	}

	newGuidFixed := utils.NewGUID()

	transDoc := models.SaleinvoiceDoc{}
	transDoc.ShopID = shopID
	transDoc.GuidFixed = newGuidFixed
	transDoc.SumAmount = sumAmount
	transDoc.Saleinvoice = trans

	transDoc.CreatedBy = username
	transDoc.CreatedAt = time.Now()

	_, err := svc.saleinvoiceRepository.Create(transDoc)

	if err != nil {
		return "", err
	}

	err = svc.mqRepo.Create(transDoc.SaleinvoiceData)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc SaleinvoiceService) UpdateSaleinvoice(guid string, shopID string, username string, trans models.Saleinvoice) error {

	findDoc, err := svc.saleinvoiceRepository.FindByGuid(guid, shopID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("guid invalid")
	}

	sumAmount := 0.0
	for i, transDetail := range *trans.Items {
		transDetail.LineNumber = i + 1
		sumAmount += transDetail.Price * transDetail.Qty
	}

	findDoc.Saleinvoice = trans
	findDoc.SumAmount = sumAmount

	findDoc.UpdatedBy = username
	findDoc.UpdatedAt = time.Now()

	err = svc.saleinvoiceRepository.Update(guid, findDoc)

	if err != nil {
		return err
	}

	err = svc.mqRepo.Create(findDoc.SaleinvoiceData)
	if err != nil {
		return err
	}

	return nil
}

func (svc SaleinvoiceService) DeleteSaleinvoice(guid string, shopID string, username string) error {

	err := svc.saleinvoiceRepository.Delete(guid, shopID, username)
	if err != nil {
		return err
	}

	docIndentity := models.Identity{
		ShopID:    shopID,
		GuidFixed: guid,
	}

	err = svc.mqRepo.Delete(docIndentity)
	if err != nil {
		return err
	}

	return nil
}

func (svc SaleinvoiceService) InfoSaleinvoice(guid string, shopID string) (models.SaleinvoiceInfo, error) {
	trans, err := svc.saleinvoiceRepository.FindByGuid(guid, shopID)

	if err != nil {
		return models.SaleinvoiceInfo{}, err
	}

	return trans.SaleinvoiceInfo, nil
}

func (svc SaleinvoiceService) SearchSaleinvoice(shopID string, q string, page int, limit int) ([]models.SaleinvoiceInfo, paginate.PaginationData, error) {
	transList, pagination, err := svc.saleinvoiceRepository.FindPage(shopID, q, page, limit)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}

func (svc SaleinvoiceService) SearchItemsSaleinvoice(guid string, shopID string, q string, page int, limit int) ([]models.SaleinvoiceInfo, paginate.PaginationData, error) {
	transList, pagination, err := svc.saleinvoiceRepository.FindItemsByGuidPage(guid, shopID, q, page, limit)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}
