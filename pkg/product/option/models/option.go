package models

import (
	common "smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const inventoryOptionCollectionName string = "inventoryOptions"

type Option struct {
	Code           string    `json:"code" bson:"code" gorm:"code;primaryKey"`
	XOrder         int8      `json:"xorder" bson:"xorder" gorm:"xorder"`
	Required       bool      `json:"required" bson:"required" gorm:"required,type:bool,default:false"`
	ChoiceType     int8      `json:"choicetype" bson:"choicetype,omitempty" gorm:"choicetype,omitempty"`
	MaxSelect      int8      `json:"maxselect" bson:"maxselect,omitempty" gorm:"maxselect,omitempty"`
	Name1          string    `json:"name1" bson:"name1" gorm:"name1"`
	Name2          string    `json:"name2,omitempty" bson:"name2,omitempty" gorm:"name2,omitempty"`
	Name3          string    `json:"name3,omitempty" bson:"name3,omitempty" gorm:"name3,omitempty"`
	Name4          string    `json:"name4,omitempty" bson:"name4,omitempty" gorm:"name4,omitempty"`
	Name5          string    `json:"name5,omitempty" bson:"name5,omitempty" gorm:"name5,omitempty"`
	Choices        *[]Choice `json:"choices" bson:"choices" gorm:"choices;foreignKey:OptCode"`
	IsStockControl bool      `json:"isstockcontrol" bson:"isstockcontrol" gorm:"isstockcontrol"`

	OptionDetails []OptionDetail `json:"optiondetails" bson:"optiondetails"`
}

type OptionDetail struct {
	OptionDetailCode string         `json:"optiondetailcode" bson:"optiondetailcode"`
	Name1            string         `json:"name1" bson:"name1" gorm:"name1"`
	Name2            string         `json:"name2,omitempty" bson:"name2,omitempty" gorm:"name2,omitempty"`
	Name3            string         `json:"name3,omitempty" bson:"name3,omitempty" gorm:"name3,omitempty"`
	Name4            string         `json:"name4,omitempty" bson:"name4,omitempty" gorm:"name4,omitempty"`
	Name5            string         `json:"name5,omitempty" bson:"name5,omitempty" gorm:"name5,omitempty"`
	Image            string         `json:"image" bson:"image" gorm:"image"`
	ChoiceDetails    []IncudeChoice `json:"choicedetails" bson:"choicedetails"`
}

type Choice struct {
	OptCode     string  `json:"-" bson:"-" gorm:"optcode;primaryKey" `
	Code        string  `json:"code" bson:"code" gorm:"code;primaryKey"`
	SuggestCode string  `json:"suggestcode,omitempty" bson:"suggestcode,omitempty" gorm:"suggestcode,omitempty"`
	Barcode     string  `json:"barcode" bson:"barcode" gorm:"barcode;primaryKey"`
	Price       float64 `json:"price" bson:"price" gorm:"price"`
	Qty         float64 `json:"qty" bson:"qty" gorm:"qty"`
	QtyMax      float64 `json:"qtymax" bson:"qtymax" gorm:"qtymax"`
	Name1       string  `json:"name1" bson:"name1" gorm:"name1"`
	Name2       string  `json:"name2,omitempty" bson:"name2,omitempty" gorm:"name2,omitempty"`
	Name3       string  `json:"name3,omitempty" bson:"name3,omitempty" gorm:"name3,omitempty"`
	Name4       string  `json:"name4,omitempty" bson:"name4,omitempty" gorm:"name4,omitempty"`
	Name5       string  `json:"name5,omitempty" bson:"name5,omitempty" gorm:"name5,omitempty"`
	ItemUnit    string  `json:"itemunit,omitempty" bson:"itemunit" gorm:"itemunit,omitempty"`
	Selected    bool    `json:"selected" bson:"selected" gorm:"selected,type:bool,default:false"`
	Default     bool    `json:"default" bson:"default" gorm:"default,type:bool,default:false"`

	IncudeOptions []IncudeChoice `json:"choicedetails,omitempty" bson:"choicedetails,omitempty"`
}

type IncudeChoice struct {
	ChoiceCode string         `json:"choicecode" bson:"choicecode"`
	Details    []IncudeChoice `json:"choicedetails" bson:"choicedetails"`
}

type InventoryOptionMain struct {
	Option `bson:"inline" gorm:"embedded;"`
}

type InventoryOptionMainInfo struct {
	common.DocIdentity  `bson:"inline" gorm:"embedded;"`
	InventoryOptionMain `bson:"inline" gorm:"embedded;"`
}

func (InventoryOptionMainInfo) CollectionName() string {
	return inventoryOptionCollectionName
}

type InventoryOptionMainData struct {
	common.ShopIdentity     `bson:"inline" gorm:"embedded;"`
	InventoryOptionMainInfo `bson:"inline" gorm:"embedded;"`
}

type InventoryOptionMainDoc struct {
	ID                      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	InventoryOptionMainData `bson:"inline" gorm:"embedded;"`
	common.ActivityDoc      `bson:"inline" gorm:"embedded;"`
}

func (InventoryOptionMainDoc) CollectionName() string {
	return inventoryOptionCollectionName
}

type InventoryOption struct {
	DocID string `bson:"-" gorm:"docid;primaryKey"`
	OptID string `bson:"-" gorm:"optid;primaryKey"`
}

func (InventoryOption) TableName() string {
	return "inventoryoptions"
}

//swagger
type InventoryOptionPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []InventoryOptionMainInfo     `json:"data,omitempty"`
	Pagination common.PaginationDataResponse `json:"pagination,omitempty"`
}
