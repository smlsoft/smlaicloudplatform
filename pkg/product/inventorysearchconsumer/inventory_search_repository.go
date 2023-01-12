package inventorysearchconsumer

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/inventorysearchconsumer/models"
)

type IInventorySearchRepository interface {
	UpSert(inventory *models.InventorySearch) error
	Delete(guidfixed string) error
}

type InventorySearchRepository struct {
	pst microservice.IPersisterOpenSearch
}

func NewInventorySearchRepository(pst microservice.IPersisterOpenSearch) *InventorySearchRepository {
	return &InventorySearchRepository{
		pst: pst,
	}
}

func (repo *InventorySearchRepository) UpSert(inventory *models.InventorySearch) error {
	err := repo.pst.CreateWithID(inventory.GuidFixed, inventory)

	if err != nil {
		return err
	}
	return nil
}

func (repo *InventorySearchRepository) Delete(guidfixed string) error {

	deleteItem := &models.InventorySearch{}
	err := repo.pst.Delete(guidfixed, deleteItem)

	if err != nil {
		return err
	}
	return nil
}
