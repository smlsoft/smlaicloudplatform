package inventorysearchconsumer

import (
	"encoding/json"
	"smlcloudplatform/pkg/product/inventorysearchconsumer/models"
)

type IInventorySearchConsumerService interface {
	Create(msg string) error
	Update(msg string) error
	Delete(msg string) error
}

type InventorySearchConsumerService struct {
	repo IInventorySearchRepository
}

func NewInventorySearchConsumerService(repo IInventorySearchRepository) *InventorySearchConsumerService {
	return &InventorySearchConsumerService{
		repo: repo,
	}
}

func (svc *InventorySearchConsumerService) Create(msg string) error {

	trans := models.InventorySearch{}
	err := json.Unmarshal([]byte(msg), &trans)

	if err != nil {
		return err
	}

	err = svc.repo.UpSert(&trans)
	if err != nil {
		return err
	}
	return nil
}

func (svc *InventorySearchConsumerService) Update(msg string) error {
	trans := models.InventorySearch{}
	err := json.Unmarshal([]byte(msg), &trans)

	if err != nil {
		return err
	}

	err = svc.repo.UpSert(&trans)
	if err != nil {
		return err
	}
	return nil
}

func (svc *InventorySearchConsumerService) Delete(msg string) error {
	trans := models.InventorySearch{}
	err := json.Unmarshal([]byte(msg), &trans)

	if err != nil {
		return err
	}

	err = svc.repo.Delete(trans.GuidFixed)
	if err != nil {
		return err
	}
	return nil
}
