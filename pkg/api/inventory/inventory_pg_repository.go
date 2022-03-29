package inventory

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type IInventoryPGRepository interface {
	Count(shopID string, guid string) (int, error)
	Create(inventory models.InventoryIndex) error
	Delete(shopID string, guid string) error
	FindByGuid(shopID string, guid string) (models.InventoryIndex, error)
}

type InventoryPGRepository struct {
	pst microservice.IPersister
}

func NewInventoryPGRepository(pst microservice.IPersister) InventoryPGRepository {
	return InventoryPGRepository{
		pst: pst,
	}
}

func (repo InventoryPGRepository) Count(shopID string, guid string) (int, error) {
	count, err := repo.pst.Count(models.InventoryIndex{}, " shop_id = ? AND guid_fixed = ?", shopID, guid)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo InventoryPGRepository) Create(inventory models.InventoryIndex) error {
	err := repo.pst.Create(inventory)
	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryPGRepository) Delete(shopID string, guid string) error {
	tableName := models.InventoryIndex{}.TableName()
	err := repo.pst.Exec("DELETE FROM "+tableName+" WHERE shop_id = ? AND guid_fixed = ?", shopID, guid)
	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryPGRepository) FindByGuid(shopID string, guid string) (models.InventoryIndex, error) {
	inv := models.InventoryIndex{}
	_, err := repo.pst.Where(&inv, "  shop_id = ? AND guid_fixed = ?", shopID, guid)
	if err != nil {
		return inv, err
	}
	return inv, nil
}
