package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const productBarcodeCollectionName = "productBarcodes"

type ProductBarcode struct {
	models.PartitionIdentity `bson:"inline"`
	ItemCode                 string          `json:"itemcode" bson:"itemcode"`
	Barcode                  string          `json:"barcode" bson:"barcode" validate:"required,min=1"`
	GroupCode                string          `json:"groupcode" bson:"groupcode"`
	GroupName                *[]models.NameX `json:"groupname" bson:"groupname"`
	Names                    *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	XSorts                   *[]models.XSort `json:"xsorts" bson:"xsorts" validate:"unique=Code,dive"`

	ItemUnitCode     string           `json:"itemunitcode" bson:"itemunitcode"`
	ItemUnitNames    *[]models.NameX  `json:"itemunitnames" bson:"itemunitnames"`
	Prices           *[]ProductPrice  `json:"prices" bson:"prices"`
	ImageURI         string           `json:"imageuri" bson:"imageuri"`
	Options          *[]ProductOption `json:"options" bson:"options"`
	Images           *[]ProductImage  `json:"images" bson:"images"`
	UseImageOrColor  bool             `json:"useimageorcolor" bson:"useimageorcolor"`
	ColorSelect      string           `json:"colorselect" bson:"colorselect"`
	ColorSelectHex   string           `json:"colorselecthex" bson:"colorselecthex"`
	IsMultipleBranch bool             `json:"ismultiplebranch" bson:"ismultiplebranch"`
}

type ProductImage struct {
	XOrder int    `json:"xorder" bson:"xorder"`
	URI    string `json:"uri" bson:"uri"`
}

type ProductPrice struct {
	KeyNumber int     `json:"keynumber" bson:"keynumber"`
	Price     float64 `json:"price" bson:"price"`
}

type ProductBarcodeInfo struct {
	models.DocIdentity `bson:"inline"`
	ProductBarcode     `bson:"inline"`
}

func (ProductBarcodeInfo) CollectionName() string {
	return productBarcodeCollectionName
}

type ProductBarcodeData struct {
	models.ShopIdentity `bson:"inline"`
	ProductBarcodeInfo  `bson:"inline"`
}

type ProductBarcodeDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductBarcodeData `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (ProductBarcodeDoc) CollectionName() string {
	return productBarcodeCollectionName
}

type ProductBarcodeItemGuid struct {
	Barcode string `json:"barcode" bson:"barcode"`
}

func (ProductBarcodeItemGuid) CollectionName() string {
	return productBarcodeCollectionName
}

type ProductBarcodeActivity struct {
	ProductBarcodeData  `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (ProductBarcodeActivity) CollectionName() string {
	return productBarcodeCollectionName
}

type ProductBarcodeDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (ProductBarcodeDeleteActivity) CollectionName() string {
	return productBarcodeCollectionName
}
