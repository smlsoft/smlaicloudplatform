package inventoryservice

import (
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
)

type IOptionGroupService interface {
	CreateOptionGroup(merchantId string, authUsername string, doc models.InventoryOptionGroup) (string, error)
	UpdateOptionGroup(guid string, merchantId string, authUsername string, doc models.InventoryOptionGroup) error
	DeleteOptionGroup(guid string, merchantId string) error
	InfoOptionGroup(guid string, merchantId string) (models.InventoryOptionGroup, error)
	SearchOptionGroup(merchantId string, q string, page int, limit int) ([]models.InventoryOptionGroup, paginate.PaginationData, error)
}

type OptionGroupService struct {
	repo IOptionGroupRepository
}

func NewOptionGroupService(optionGroupRepository IOptionGroupRepository) IOptionGroupService {
	return &OptionGroupService{
		repo: optionGroupRepository,
	}
}

func (svc *OptionGroupService) CreateOptionGroup(merchantId string, authUsername string, doc models.InventoryOptionGroup) (string, error) {

	newGuid := utils.NewGUID()
	doc.MerchantId = merchantId
	doc.GuidFixed = newGuid
	doc.CreatedBy = authUsername
	doc.CreatedAt = time.Now()
	doc.Deleted = false

	for i := range doc.Details {
		doc.Details[i].GuidFixed = utils.NewGUID()
	}

	_, err := svc.repo.Create(doc)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc *OptionGroupService) UpdateOptionGroup(guid string, merchantId string, authUsername string, doc models.InventoryOptionGroup) error {

	findDoc, err := svc.repo.FindByGuid(guid, merchantId)

	if err != nil {
		return err
	}

	findDoc.OptionName1 = doc.OptionName1
	findDoc.ProductSelectOption1 = doc.ProductSelectOption1
	findDoc.ProductSelectOption2 = doc.ProductSelectOption2
	findDoc.ProductSelectOptionMin = doc.ProductSelectOptionMin
	findDoc.ProductSelectOptionMax = doc.ProductSelectOptionMax
	findDoc.Details = doc.Details
	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc *OptionGroupService) DeleteOptionGroup(guid string, merchantId string) error {
	err := svc.repo.Delete(guid, merchantId)

	if err != nil {
		return err
	}
	return nil
}

func (svc *OptionGroupService) InfoOptionGroup(guid string, merchantId string) (models.InventoryOptionGroup, error) {

	findDoc, err := svc.repo.FindByGuid(guid, merchantId)

	if err != nil {
		return models.InventoryOptionGroup{}, err
	}

	return findDoc, nil

}

func (svc *OptionGroupService) SearchOptionGroup(merchantId string, q string, page int, limit int) ([]models.InventoryOptionGroup, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(merchantId, q, page, limit)

	if err != nil {
		return []models.InventoryOptionGroup{}, pagination, err
	}

	return docList, pagination, nil
}
