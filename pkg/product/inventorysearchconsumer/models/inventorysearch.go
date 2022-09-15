package models

import (
	common "smlcloudplatform/pkg/models"
	inventoryModel "smlcloudplatform/pkg/product/inventory/models"
)

const INVENTORY_SEARCH_INDEXNAME string = "inventory-search-index"

type InventorySearch struct {
	inventoryModel.InventoryInfo
	common.ShopIdentity `bson:"inline" gorm:"embedded;"`
}

func (*InventorySearch) IndexName() string {
	return INVENTORY_SEARCH_INDEXNAME
}
