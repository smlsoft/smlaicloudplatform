package models

import inventoryModel "smlcloudplatform/pkg/product/inventory/models"

const INVENTORY_SEARCH_INDEXNAME string = "inventory-search-index"

type InventorySearch struct {
	inventoryModel.InventoryInfo
}

func (*InventorySearch) IndexName() string {
	return INVENTORY_SEARCH_INDEXNAME
}
