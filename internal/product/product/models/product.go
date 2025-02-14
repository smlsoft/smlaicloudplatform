package models

import (
	"smlaicloudplatform/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const productCollectionName = "products"

type Product struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string             `json:"code" bson:"code"`
	Names                    *[]models.NameX    `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	GroupCode                string             `json:"groupcode" bson:"groupcode"`
	GroupNames               *[]models.NameX    `json:"groupnames" bson:"groupnames"`
	ManufacturerGUID         string             `json:"manufacturerguid" bson:"manufacturerguid"`
	ManufacturerCode         string             `json:"manufacturercode" bson:"manufacturercode"`
	ManufacturerNames        *[]models.NameX    `json:"manufacturernames" bson:"manufacturernames"`
	Dimensions               []ProductDimension `json:"dimensions" bson:"dimensions"`
	VatType                  int8               `json:"vattype" bson:"vattype"`
	Barcodes                 []string           `json:"barcodes,omitempty"`
}

type JSONB []models.NameX

type ProductPg struct {
	ShopID                   string `json:"shopid" gorm:"column:shopid;primaryKey"`
	models.PartitionIdentity `gorm:"embedded;"`
	Barcode                  string  `json:"barcode" gorm:"column:barcode;primaryKey"`
	Names                    JSONB   `json:"names"  gorm:"column:names;type:jsonb" `
	UnitCode                 string  `json:"itemunitcode" gorm:"column:unitcode"`
	UnitNames                JSONB   `json:"itemunitnames" gorm:"column:unitnames;type:jsonb"`
	BalanceQty               float64 `json:"balanceqty" gorm:"column:balanceqty"`
	MainBarcodeRef           string  `json:"mainbarcoderef" gorm:"column:mainbarcoderef"`
	StandValue               float64 `json:"standvalue" gorm:"column:standvalue"`
	DivideValue              float64 `json:"dividevalue" gorm:"column:dividevalue"`
	BalanceAmount            float64 `json:"balanceamount" gorm:"column:balanceamount"`
	AverageCost              float64 `json:"averagecost" gorm:"column:averagecost"`
}

func (ProductPg) TableName() string {
	return "productbarcode"
}

type ProductDimension struct {
	models.DocIdentity `bson:"inline"`
	Names              *[]models.NameX      `json:"names" bson:"names"`
	IsDisabled         bool                 `json:"isdisabled" bson:"isdisabled"`
	Item               ProductDimensionItem `json:"item" bson:"item"`
}

type ProductDimensionItem struct {
	models.DocIdentity `bson:"inline"`
	Names              *[]models.NameX `json:"names" bson:"names"`
	IsDisabled         bool            `json:"isdisabled" bson:"isdisabled"`
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
	Code string `json:"code" bson:"code"`
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
