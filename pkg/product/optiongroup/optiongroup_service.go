package optiongroup

import (
	"context"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/optiongroup/models"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/pkg/errors"
	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IOptionGroupService interface {
	CreateOptionGroup(shopID string, authUsername string, doc models.InventoryOptionGroup) (string, error)
	UpdateOptionGroup(shopID string, guid string, authUsername string, doc models.InventoryOptionGroup) error
	DeleteOptionGroup(shopID string, guid string, username string) error
	InfoOptionGroup(shopID string, guid string) (models.InventoryOptionGroup, error)
	SearchOptionGroup(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionGroup, mongopagination.PaginationData, error)
}

type OptionGroupService struct {
	repo           IOptionGroupRepository
	contextTimeout time.Duration
}

func NewOptionGroupService(optionGroupRepository IOptionGroupRepository) OptionGroupService {

	contextTimeout := time.Duration(15) * time.Second

	return OptionGroupService{
		repo:           optionGroupRepository,
		contextTimeout: contextTimeout,
	}
}

func (svc OptionGroupService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc OptionGroupService) CreateOptionGroup(shopID string, authUsername string, doc models.InventoryOptionGroup) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	newGuid := utils.NewGUID()
	doc.ShopID = shopID
	doc.GuidFixed = newGuid
	doc.CreatedBy = authUsername
	doc.CreatedAt = time.Now()

	for i := range doc.Details {
		doc.Details[i].GuidFixed = utils.NewGUID()
	}

	_, err := svc.repo.Create(ctx, doc)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc OptionGroupService) UpdateOptionGroup(shopID string, guid string, authUsername string, doc models.InventoryOptionGroup) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NewObjectID() {
		return errors.New("document not found")
	}

	findDoc.OptionName1 = doc.OptionName1
	findDoc.ProductSelectOption1 = doc.ProductSelectOption1
	findDoc.ProductSelectOption2 = doc.ProductSelectOption2
	findDoc.ProductSelectOptionMin = doc.ProductSelectOptionMin
	findDoc.ProductSelectOptionMax = doc.ProductSelectOptionMax
	findDoc.Details = doc.Details
	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc OptionGroupService) DeleteOptionGroup(shopID string, guid string, username string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	err := svc.repo.Delete(ctx, shopID, guid, username)

	if err != nil {
		return err
	}
	return nil
}

func (svc OptionGroupService) InfoOptionGroup(shopID string, guid string) (models.InventoryOptionGroup, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.InventoryOptionGroup{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.InventoryOptionGroup{}, errors.New("document not found")
	}

	return findDoc, nil

}

func (svc OptionGroupService) SearchOptionGroup(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionGroup, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList, pagination, err := svc.repo.FindPage(ctx, shopID, pageable)

	if err != nil {
		return []models.InventoryOptionGroup{}, pagination, err
	}

	return docList, pagination, nil
}
