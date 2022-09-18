package models

import (
	"os"
	common "smlcloudplatform/pkg/models"
	inventoryModel "smlcloudplatform/pkg/product/inventory/models"
)

const INVENTORY_SEARCH_INDEXNAME string = "inventory-search-index"

type InventorySearch struct {
	inventoryModel.InventoryInfo
	common.ShopIdentity `bson:"inline" gorm:"embedded;"`
}

func (*InventorySearch) IndexName() string {

	inventorySearchIndexName := os.Getenv("INVENTORY_SEARCH_INDEXNAME")
	if inventorySearchIndexName == "" {
		inventorySearchIndexName = INVENTORY_SEARCH_INDEXNAME
	}
	return inventorySearchIndexName
}
