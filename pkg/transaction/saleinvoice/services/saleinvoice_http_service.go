package services

import (
	"errors"
	micromodels "smlcloudplatform/internal/microservice/models"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISaleinvoiceService interface {
	CreateSaleinvoice(shopID string, username string, trans models.Saleinvoice) (string, error)
	UpdateSaleinvoice(shopID string, guid string, username string, trans models.Saleinvoice) error
	DeleteSaleinvoice(shopID string, guid string, username string) error
	InfoSaleinvoice(shopID string, guid string) (models.SaleinvoiceInfo, error)
	SearchSaleinvoice(shopID string, pageable micromodels.Pageable) ([]models.SaleinvoiceInfo, mongopagination.PaginationData, error)
	SearchItemsSaleinvoice(guid string, shopID string, pageable micromodels.Pageable) ([]models.SaleinvoiceInfo, mongopagination.PaginationData, error)
}

type SaleinvoiceService struct {
	saleinvoiceRepository repositories.ISaleinvoiceRepository
	mqRepo                repositories.ISaleinvoiceMQRepository
}

func NewSaleinvoiceService(saleinvoiceRepository repositories.ISaleinvoiceRepository, mqRepo repositories.ISaleinvoiceMQRepository) SaleinvoiceService {

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

func (svc SaleinvoiceService) UpdateSaleinvoice(shopID string, guid string, username string, trans models.Saleinvoice) error {

	findDoc, err := svc.saleinvoiceRepository.FindByGuid(shopID, guid)

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

	err = svc.saleinvoiceRepository.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	err = svc.mqRepo.Update(findDoc.SaleinvoiceData)
	if err != nil {
		return err
	}

	return nil
}

func (svc SaleinvoiceService) DeleteSaleinvoice(shopID string, guid string, username string) error {

	err := svc.saleinvoiceRepository.Delete(shopID, guid, username)
	if err != nil {
		return err
	}

	docIndentity := common.Identity{
		ShopID:    shopID,
		GuidFixed: guid,
	}

	err = svc.mqRepo.Delete(docIndentity)
	if err != nil {
		return err
	}

	return nil
}

func (svc SaleinvoiceService) InfoSaleinvoice(shopID string, guid string) (models.SaleinvoiceInfo, error) {
	trans, err := svc.saleinvoiceRepository.FindByGuid(shopID, guid)

	if err != nil {
		return models.SaleinvoiceInfo{}, err
	}

	return trans.SaleinvoiceInfo, nil
}

func (svc SaleinvoiceService) SearchSaleinvoice(shopID string, pageable micromodels.Pageable) ([]models.SaleinvoiceInfo, mongopagination.PaginationData, error) {
	transList, pagination, err := svc.saleinvoiceRepository.FindPage(shopID, pageable)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}

func (svc SaleinvoiceService) SearchItemsSaleinvoice(guid string, shopID string, pageable micromodels.Pageable) ([]models.SaleinvoiceInfo, mongopagination.PaginationData, error) {
	transList, pagination, err := svc.saleinvoiceRepository.FindItemsByGuidPage(guid, shopID, pageable)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}
