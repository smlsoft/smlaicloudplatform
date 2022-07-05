package models

import (
	common "smlcloudplatform/pkg/models"
	inventoryModel "smlcloudplatform/pkg/product/inventory/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryOptionGroup struct {
	ID                     primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	ShopID                 string                       `json:"shopid" bson:"shopid"`
	GuidFixed              string                       `json:"guidfixed" bson:"guidfixed"`
	OptionName1            string                       `json:"optionname1" bson:"optionname1"`
	ProductSelectOption1   bool                         `json:"productselectoption1" bson:"productselectoption1"`
	ProductSelectOption2   bool                         `json:"productselectoption2" bson:"productselectoption2"`
	ProductSelectOptionMin int                          `json:"productselectoptionmin" bson:"productselectoptionmin"`
	ProductSelectOptionMax int                          `json:"productselectoptionmax" bson:"productselectoptionmax"`
	Details                []InventoryOptionGroupDetail `json:"details" bson:"details"`
	common.ActivityDoc
}

func (*InventoryOptionGroup) CollectionName() string {
	return "inventoryOptionGroup"
}

type InventoryOptionGroupDetail struct {
	GuidFixed   string  `json:"guidfixed" bson:"guidfixed"`
	DetailName1 string  `json:"detailname1" bson:"detailname1"`
	Amount      float32 `json:"amount" bson:"amount"`
}

//swagger
type InventoryOptionGroupResponse struct {
	Success    bool                           `json:"success"`
	Data       []inventoryModel.InventoryInfo `json:"data,omitempty"`
	Pagination common.PaginationDataResponse  `json:"pagination,omitempty"`
}

type InventoryOptionGroupInfoResponse struct {
	Success bool                         `json:"success"`
	Data    inventoryModel.InventoryInfo `json:"data,omitempty"`
}
