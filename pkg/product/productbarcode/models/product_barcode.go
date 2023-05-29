package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const productBarcodeCollectionName = "productBarcodes"

type ProductBarcodeBase struct {
	ItemCode  string          `json:"itemcode" bson:"itemcode"`
	Barcode   string          `json:"barcode" bson:"barcode" validate:"required,min=1"`
	GroupCode string          `json:"groupcode" bson:"groupcode"`
	GroupName *[]models.NameX `json:"groupnames" bson:"groupnames"`
	Names     *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	XSorts    *[]models.XSort `json:"xsorts" bson:"xsorts" validate:"unique=Code,dive"`

	ItemUnitCode    string           `json:"itemunitcode" bson:"itemunitcode"`
	ItemUnitNames   *[]models.NameX  `json:"itemunitnames" bson:"itemunitnames"`
	Prices          *[]ProductPrice  `json:"prices" bson:"prices"`
	ImageURI        string           `json:"imageuri" bson:"imageuri"`
	Options         *[]ProductOption `json:"options" bson:"options"`
	Images          *[]ProductImage  `json:"images" bson:"images"`
	UseImageOrColor bool             `json:"useimageorcolor" bson:"useimageorcolor"`
	ColorSelect     string           `json:"colorselect" bson:"colorselect"`
	ColorSelectHex  string           `json:"colorselecthex" bson:"colorselecthex"`

	Condition        bool    `json:"condition" bson:"condition"`
	DivideValue      float64 `json:"dividevalue" bson:"dividevalue"`
	StandValue       float64 `json:"standvalue" bson:"standvalue"`
	IsUseSubBarcodes bool    `json:"isusesubbarcodes" bson:"isusesubbarcodes"`

	ItemType    int8   `json:"itemtype" bson:"itemtype"`
	TaxType     int8   `json:"taxtype" bson:"taxtype"`
	VatType     int8   `json:"vattype" bson:"vattype"`
	IsSumPoint  bool   `json:"issumpoint" bson:"issumpoint"`
	MaxDiscount string `json:"maxdiscount" bson:"maxdiscount"`
	IsDividend  bool   `json:"isdividend" bson:"isdividend"`

	RefUnitNames   *[]models.NameX `json:"refunitnames" bson:"refunitnames"`
	StockBarcode   string          `json:"stockbarcode" bson:"stockbarcode"`
	Qty            float64         `json:"qty" bson:"qty"`
	RefDivideValue float64         `json:"refdividevalue" bson:"refdividevalue"`
	RefStandValue  float64         `json:"refstandvalue" bson:"refstandvalue"`
	VatCal         int             `json:"vatcal" bson:"vatcal"`
}

type RefProductBarcode struct {
	GuidFixed     string          `json:"guidfixed" bson:"guidfixed"`
	Names         *[]models.NameX `json:"names" bson:"names"`
	ItemUnitCode  string          `json:"itemunitcode" bson:"itemunitcode"`
	ItemUnitNames *[]models.NameX `json:"itemunitnames" bson:"itemunitnames"`
	Barcode       string          `json:"barcode" bson:"barcode" validate:"required,min=1"`
	Condition     bool            `json:"condition" bson:"condition"`
	DivideValue   float64         `json:"dividevalue" bson:"dividevalue"`
	StandValue    float64         `json:"standvalue" bson:"standvalue"`
	Qty           float64         `json:"qty" bson:"qty"`
}

type ProductBarcode struct {
	models.PartitionIdentity `bson:"inline"`
	ProductBarcodeBase       `bson:"inline"`
	RefBarcodes              *[]RefProductBarcode `json:"refbarcodes" bson:"refbarcodes"`
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

func (doc ProductBarcodeDoc) ToRefBarcode() RefProductBarcode {
	return RefProductBarcode{
		GuidFixed:     doc.GuidFixed,
		Names:         doc.Names,
		ItemUnitCode:  doc.ItemUnitCode,
		ItemUnitNames: doc.ItemUnitNames,
		Barcode:       doc.Barcode,
	}
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

type ProductBarcodeRequest struct {
	ProductBarcodeBase
	RefBarcodes []BarcodeRequest `json:"refbarcodes"`
}

type BarcodeRequest struct {
	Barcode     string  `json:"barcode" bson:"barcode" validate:"required,min=1"`
	Condition   bool    `json:"condition" bson:"condition"`
	DivideValue float64 `json:"dividevalue" bson:"dividevalue"`
	StandValue  float64 `json:"standvalue" bson:"standvalue"`
	Qty         float64 `json:"qty" bson:"qty"`
}

func (p ProductBarcodeRequest) ToProductBarcode() ProductBarcode {
	return ProductBarcode{
		ProductBarcodeBase: p.ProductBarcodeBase,
	}
}

type ProductBarcodeSearch struct {
	ICCode   string   `json:"iccode" ch:"iccode"`
	Barcode  string   `json:"barcode" ch:"barcode"`
	UnitCode string   `json:"unitcode" ch:"unitcode"`
	Price    string   `json:"price" ch:"price"`
	Names    []string `json:"names" ch:"names"`
}

func (ProductBarcodeSearch) TableName() string {
	return productBarcodeCollectionName
}
