package models

const INVENTORY_SEARCH_INDEXNAME string = "inventory-search-index"

type InventorySearch struct {
	InventoryInfo
}

func (*InventorySearch) IndexName() string {
	return INVENTORY_SEARCH_INDEXNAME
}
