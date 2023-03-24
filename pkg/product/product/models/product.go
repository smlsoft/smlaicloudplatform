package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const productCollectionName = "products"

type Product struct {
	models.PartitionIdentity `bson:"inline"`
	ItemCode                 string `json:"itemcode" bson:"itemcode" validate:"required,min=1,max=100"`

	Barcodes        *[]string       `json:"barcodes" bson:"barcodes"`
	Names           *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	MultiUnit       bool            `json:"multiunit" bson:"multiunit"`
	UseSerialNumber bool            `json:"useserialnumber" bson:"useserialnumber"`
	Units           *[]ProductUnit  `json:"units,omitempty" bson:"units,omitempty"`
	Unit            ProductUnit     `json:"unit" bson:"unit" validate:"required"`

	UnitCost          string          `json:"unitcost" bson:"unitcost"`
	UnitCostNames     *[]models.NameX `json:"unitcostnames" bson:"unitcostnames"`
	UnitStandard      string          `json:"unitstandard" bson:"unitstandard"`
	UnitStandardNames *[]models.NameX `json:"unitstandardnames" bson:"unitstandardnames"`
	ItemStockType     int8            `json:"itemstocktype" bson:"itemstocktype"`
	ItemType          int8            `json:"itemtype" bson:"itemtype"`
	VatType           int8            `json:"vattype" bson:"vattype"`
	IsSumPoint        bool            `json:"issumpoint" bson:"issumpoint"`
	Images            *[]ProductImage `json:"images" bson:"images"`
	Prices            *[]ProductPrice `json:"prices" bson:"prices"`
	GroupCode         string          `json:"groupcode" bson:"groupcode"`
	GroupName         *[]models.NameX `json:"groupname" bson:"groupname"`
	Branches          *[]uint16       `json:"branches" bson:"branches" validate:"unique"`
	IsMultipleBranch  bool            `json:"ismultiplebranch" bson:"ismultiplebranch"`
}

type ProductPrice struct {
	KeyNumber int     `json:"keynumber" bson:"keynumber"`
	Price     float64 `json:"price" bson:"price"`
}

type ProductUnit struct {
	UnitCode   string          `json:"unitcode" bson:"unitcode"`
	Names      *[]models.NameX `json:"names" bson:"names"`
	Divider    float64         `json:"divider" bson:"divider"`
	Stand      float64         `json:"stand" bson:"stand"`
	XOrder     int             `json:"xorder" bson:"xorder"`
	StockCount bool            `json:"stockcount" bson:"stockcount"`
}

type ProductImage struct {
	XOrder int    `json:"xorder" bson:"xorder"`
	URI    string `json:"uri" bson:"uri"`
}

type ProductInfo struct {
	models.DocIdentity `bson:"inline"`
	Product            `bson:"inline"`
}

func (ProductInfo) CollectionName() string {
	return productCollectionName
}

type ProductData struct {
	models.ShopIdentity `bson:"inline"`
	ProductInfo         `bson:"inline"`
}

type ProductDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductData        `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (ProductDoc) CollectionName() string {
	return productCollectionName
}

type ProductItemGuid struct {
	ItemCode string `json:"itemcode" bson:"itemcode"`
}

func (ProductItemGuid) CollectionName() string {
	return productCollectionName
}

type ProductActivity struct {
	ProductData         `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (ProductActivity) CollectionName() string {
	return productCollectionName
}

type ProductDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (ProductDeleteActivity) CollectionName() string {
	return productCollectionName
}
