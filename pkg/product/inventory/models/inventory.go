package models

import (
	"smlcloudplatform/pkg/models"
	common "smlcloudplatform/pkg/models"
	categoryModel "smlcloudplatform/pkg/product/category/models"
	optionModel "smlcloudplatform/pkg/product/option/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type InventoryUnit struct {
// 	UnitGuid  string  `bson:"unitguid" json:"unitguid"`   // Guid หน่วยนับ
// 	UnitName  string  `bson:"unitname" json:"unitname"`   // ชื่อหน่วยนับ
// 	Minuend   float32 `bson:"minuend" json:"minuend"`     // ตัวตั้ง
// 	Divisor   float32 `bson:"divisor" json:"divisor"`     // ตัวหาร
// 	Activated bool    `bson:"activated" json:"activated"` // เปิดใช้งานอยู่
// }

const inventoryCollectionName string = "inventories"
const inventoryTableName string = "inventories"
const inventoryIndexName string = "inventories_index"

type Inventory struct {
	ParID    string `json:"parid" bson:"parid" gorm:"parid"`
	ItemSku  string `json:"itemsku,omitempty" bson:"itemsku,omitempty" gorm:"itemsku,omitempty"`
	ItemGuid string `json:"itemguid,omitempty" bson:"itemguid,omitempty" gorm:"itemguid,omitempty"`
	ItemCode string `json:"itemcode" bson:"itemcode" gorm:"itemcode"`

	Barcode string `json:"barcode" bson:"barcode" gorm:"barcode"`

	UnitCode string `json:"unitcode" bson:"unitcode" gorm:"unitcode"`

	models.Name        `bson:"inline"`
	models.Description `bson:"inline"`
	ProductPrice       `bson:"inline"`
	CategoryGuid       string `json:"categoryguid,omitempty" bson:"categoryguid" gorm:"categoryguid"` // Guid กลุ่มสินค้า
	HaveSerialno       bool   `json:"haveserialno" bson:"haveserialno" gorm:"haveserialno,type:bool,default:false"`
	ItemVat            int8   `json:"itemvat" bson:"itemvat" gorm:"itemvat"`
	ItemType           int8   `json:"itemtype" bson:"itemtype" gorm:"itemtype"`
	HavePoint          bool   `json:"havepoint" bson:"havepoint" gorm:"havepoint"`
	XOrder             int8   `json:"xorder" bson:"xorder" gorm:"xorder"`

	IsStockProduct      bool   `json:"isstockproduct" bson:"isstockproduct" gorm:"isstockproduct"`
	StockProductGUIDRef string `json:"stockproductguidref" bson:"stockproductguidref" gorm:"stockproductguidref"`

	Activated   bool `json:"activated,omitempty" bson:"activated,omitempty" gorm:"activated,omitempty,type:bool,default:false"`       // เปิดใช้งานอยู่
	Recommended bool `json:"recommended,omitempty" bson:"recommended,omitempty" gorm:"recommended,omitempty,type:bool,default:false"` // สินค้าแนะนำ

	BarcodeDescriptionFromProduct bool `json:"barcodedescriptionfromproduct,omitempty" bson:"barcodedescriptionfromproduct,omitempty" gorm:"barcodedescriptionfromproduct,omitempty,type:bool,default:false"`

	Options  *[]optionModel.Option   `json:"options,omitempty" bson:"options,omitempty" gorm:"many2many:inventoryoptions;foreignKey:GuidFixed;joinForeignKey:DocID;References:Code;joinReferences:OptID"`
	Images   *[]InventoryImage       `json:"images,omitempty" bson:"images,omitempty" gorm:"images;foreignKey:DocID"`
	Tags     *[]InventoryTag         `json:"tags,omitempty" bson:"tags" gorm:"tags;foreignKey:DocID"`
	Category *categoryModel.Category `json:"category,omitempty" bson:"category,omitempty"`
	Barcodes *[]Barcode              `json:"barcodes" bson:"barcodes" gorm:"barcodes"`
	UnitUses *[]UnitUse              `json:"unituses" bson:"unituses" gorm:"units"`

	// WaitType         int             `json:"-" bson:"waitType"`                // ประเภทการรอ (สินค้าหมด)
	// WaitUntil        time.Time       `json:"-" bson:"waitUntil"`               // ระยะเวลาที่รอ
	// MultipleUnits    bool            `json:"-" bson:"multipleuUits" `          // สินค้าหลายหน่วยนับ
	// UnitStandardGuid string          `json:"-" bson:"unitStandardGuid" `       // หน่วยนับมาตรฐาน (นับสต๊อก)
	// UnitList         []InventoryUnit `json:"unitlist" bson:"unitList" `        // กรณีหลายหน่วยนับ ตารางหน่วบนับ
}

type ProductPrice struct {
	Price       float64 `json:"price" bson:"price" gorm:"price"` // ราคาพื้นฐาน (กรณีไม่มีตารางราคา และโปรโมชั่น)
	MemberPrice float32 `json:"memberprice,omitempty" bson:"memberprice,omitempty" gorm:"memberprice,omitempty"`
}

type Unit struct {
	UnitCode        string `json:"unitcode" bson:"unitcode" gorm:"unitcode"`
	models.UnitName `bson:"inline"`
}

type Barcode struct {
	Barcode            string `json:"barcode" bson:"barcode" gorm:"barcode"`
	UnitCode           string `json:"unitcode" bson:"unitcode" gorm:"unitcode"`
	Image              string `json:"image" bson:"image" gorm:"image"`
	IsPrimary          bool   `json:"isprimary" bson:"isprimary" gorm:"isprimary,type:bool,default:false"`
	ProductPrice       `bson:"inline"`
	models.UnitName    `bson:"inline"`
	models.Name        `bson:"inline"`
	models.Description `bson:"inline"`
}

type UnitUse struct {
	UnitCode           string `json:"unitcode" bson:"unitcode" gorm:"unitcode"`
	models.Description `bson:"inline"`
	models.UnitName    `bson:"inline"`
	ItemUnitSTD        float64 `json:"itemunitstd" bson:"itemunitstd" gorm:"itemunitstd"`
	ItemUnitDIV        float64 `json:"itemunitdiv" bson:"itemunitdiv" gorm:"itemunitdiv"`
	IsUnitCost         bool    `json:"isunitcost" bson:"isunitcost" gorm:"isunitcost,type:bool,default:false"`
	IsUnitStandard     bool    `json:"isunitstandard" bson:"isunitstandard" gorm:"isunitstandard,type:bool,default:false"`
	IsPrimary          bool    `json:"isprimary" bson:"isprimary" gorm:"isprimary,type:bool,default:false"`
}

type InventoryItemGuid struct {
	ItemGuid string `json:"itemguid,omitempty" bson:"itemguid,omitempty"`
}

func (InventoryItemGuid) CollectionName() string {
	return inventoryCollectionName
}

type InventoryImage struct {
	DocID string `json:"-" bson:"-" gorm:"docid;primaryKey"`
	Uri   string `json:"uri" bson:"uri" gorm:"uri;primaryKey"`
}

type InventoryTag struct {
	DocID string `json:"-" bson:"-" gorm:"docid;primaryKey"`
	Name  string `json:"name" bson:"name" gorm:"name;primaryKey"`
}

type InventoryInfo struct {
	common.DocIdentity `bson:"inline" gorm:"embedded;"`
	Inventory          `bson:"inline" gorm:"embedded;"`
	Unit               *Unit    `json:"unit,omitempty" bson:"unit,omitempty"`
	BarcodeDetail      *Barcode `json:"barcodedetail,omitempty" bson:"barcodedetail,omitempty"`
}

func (InventoryInfo) CollectionName() string {
	return inventoryCollectionName
}

type InventoryData struct {
	common.ShopIdentity `bson:"inline" gorm:"embedded;"`
	InventoryInfo       `bson:"inline" gorm:"embedded;"`
}

func (InventoryData) TableName() string {
	return inventoryTableName
}

type InventoryDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	InventoryData      `bson:"inline"`
	common.ActivityDoc `bson:"inline"`
	common.LastUpdate  `bson:"inline"`
}

func (InventoryDoc) CollectionName() string {
	return inventoryCollectionName
}

type InventoryActivity struct {
	InventoryData `bson:"inline"`
	CreatedAt     *time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
	UpdatedAt     *time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	DeletedAt     *time.Time `json:"deletedat,omitempty" bson:"deletedat,omitempty"`
}

func (InventoryActivity) CollectionName() string {
	return inventoryCollectionName
}

type InventoryDeleteActivity struct {
	common.Identity `bson:"inline"`
	CreatedAt       *time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
	UpdatedAt       *time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	DeletedAt       *time.Time `json:"deletedat,omitempty" bson:"deletedat,omitempty"`
}

func (InventoryDeleteActivity) CollectionName() string {
	return inventoryCollectionName
}

type InventoryIndex struct {
	common.Index `bson:"inline"`
}

func (InventoryIndex) TableName() string {
	return inventoryIndexName
}

// for swagger gen

type InventoryBulkImport struct {
	Created          []string `json:"created"`
	Updated          []string `json:"updated"`
	UpdateFailed     []string `json:"updateFailed"`
	PayloadDuplicate []string `json:"payloadDuplicate"`
}

type InventoryBulkReponse struct {
	Success bool `json:"success"`
	InventoryBulkImport
}

type InventoryPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []InventoryInfo               `json:"data,omitempty"`
	Pagination common.PaginationDataResponse `json:"pagination,omitempty"`
}

type InventoryInfoResponse struct {
	Success bool          `json:"success"`
	Data    InventoryInfo `json:"data,omitempty"`
}

type InventoryBulkInsertResponse struct {
	Success    bool     `json:"success"`
	Created    []string `json:"created"`
	Updated    []string `json:"updated"`
	Failed     []string `json:"updateFailed"`
	Duplicated []string `json:"payloadDuplicate"`
}

type InventoryLastActivityResponse struct {
	New    []InventoryActivity       `json:"new" `
	Remove []InventoryDeleteActivity `json:"remove"`
}

type InventoryFetchUpdateResponse struct {
	Success    bool                          `json:"success"`
	Data       InventoryLastActivityResponse `json:"data,omitempty"`
	Pagination common.PaginationDataResponse `json:"pagination,omitempty"`
}
