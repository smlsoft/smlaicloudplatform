package models

import (
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
	ParID        string  `json:"parid" bson:"parid" gorm:"parid"`
	ItemSku      string  `json:"itemsku,omitempty" bson:"itemsku,omitempty" gorm:"itemsku,omitempty"`
	CategoryGuid string  `json:"categoryguid,omitempty" bson:"categoryguid" gorm:"categoryguid"` // Guid กลุ่มสินค้า
	Price        float32 `json:"price" bson:"price" gorm:"price"`                                // ราคาพื้นฐาน (กรณีไม่มีตารางราคา และโปรโมชั่น)
	MemberPrice  float32 `json:"memberprice,omitempty" bson:"memberprice,omitempty" gorm:"memberprice,omitempty"`
	Recommended  bool    `json:"recommended,omitempty" bson:"recommended,omitempty" gorm:"recommended,omitempty,type:bool,default:false"` // สินค้าแนะนำ
	Activated    bool    `json:"activated,omitempty" bson:"activated,omitempty" gorm:"activated,omitempty,type:bool,default:false"`       // เปิดใช้งานอยู่

	Name1        string `json:"name1" bson:"name1" gorm:"name1"` // ชื่อภาษาไทย
	Name2        string `json:"name2,omitempty" bson:"name2,omitempty" gorm:"name2,omitempty"`
	Name3        string `json:"name3,omitempty" bson:"name3,omitempty" gorm:"name3,omitempty"`
	Name4        string `json:"name4,omitempty" bson:"name4,omitempty" gorm:"name4,omitempty"`
	Name5        string `json:"name5,omitempty" bson:"name5,omitempty" gorm:"name5,omitempty"`
	Description1 string `json:"description1,omitempty" bson:"description1,omitempty" gorm:"description1,omitempty"` // รายละเอียดภาษาไทย
	Description2 string `json:"description2,omitempty" bson:"description2,omitempty" gorm:"description2,omitempty"`
	Description3 string `json:"description3,omitempty" bson:"description3,omitempty" gorm:"description3,omitempty"`
	Description4 string `json:"description4,omitempty" bson:"description4,omitempty" gorm:"description4,omitempty"`
	Description5 string `json:"description5,omitempty" bson:"description5,omitempty" gorm:"description5,omitempty"`

	// UnitName1 string `json:"unitname1" bson:"unitname1" gorm:"unitname1"`
	// UnitName2 string `json:"unitname2,omitempty" bson:"unitname2,omitempty" gorm:"unitname2,omitempty"`
	// UnitName3 string `json:"unitname3,omitempty" bson:"unitname3,omitempty" gorm:"unitname3,omitempty"`
	// UnitName4 string `json:"unitname4,omitempty" bson:"unitname4,omitempty" gorm:"unitname4,omitempty"`
	// UnitName5 string `json:"unitname5,omitempty" bson:"unitname5,omitempty" gorm:"unitname5,omitempty"`

	ItemGuid     string  `json:"itemguid,omitempty" bson:"itemguid,omitempty" gorm:"itemguid,omitempty"`
	ItemCode     string  `json:"itemcode" bson:"itemcode" gorm:"itemcode"`
	ItemUnitCode string  `json:"itemunitcode,omitempty" bson:"itemunitcode,omitempty" gorm:"itemunitcode,omitempty"`
	ItemUnitStd  float64 `json:"itemunitstd,omitempty" bson:"itemunitstd,omitempty" gorm:"itemunitstd,omitempty"`
	ItemUnitDiv  float64 `json:"itemunitdiv,omitempty" bson:"itemunitdiv,omitempty" gorm:"itemunitdiv,omitempty"`

	Options  *[]optionModel.Option   `json:"options,omitempty" bson:"options,omitempty" gorm:"many2many:inventoryoptions;foreignKey:GuidFixed;joinForeignKey:DocID;References:Code;joinReferences:OptID"`
	Images   *[]InventoryImage       `json:"images,omitempty" bson:"images,omitempty" gorm:"images;foreignKey:DocID"`
	Tags     *[]InventoryTag         `json:"tags,omitempty" bson:"tags" gorm:"tags;foreignKey:DocID"`
	Category *categoryModel.Category `json:"category,omitempty" bson:"category,omitempty"`
	XOrder   int8                    `json:"xorder" bson:"xorder" gorm:"xorder"`

	Barcodes *[]Barcode `json:"barcodes" bson:"barcodes" gorm:"barcodes"`
	UnitUses *[]UnitUse `json:"unituses" bson:"unituses" gorm:"units"`

	HaveSerialno bool `json:"haveserialno" bson:"haveserialno" gorm:"haveserialno,type:bool,default:false"`
	ItemVat      int8 `json:"itemvat" bson:"itemvat" gorm:"itemvat"`
	ItemType     int8 `json:"itemtype" bson:"itemtype" gorm:"itemtype"`
	HavePoint    bool `json:"havepoint" bson:"havepoint" gorm:"havepoint"`

	// WaitType         int             `json:"-" bson:"waitType"`                // ประเภทการรอ (สินค้าหมด)
	// WaitUntil        time.Time       `json:"-" bson:"waitUntil"`               // ระยะเวลาที่รอ
	// MultipleUnits    bool            `json:"-" bson:"multipleuUits" `          // สินค้าหลายหน่วยนับ
	// UnitStandardGuid string          `json:"-" bson:"unitStandardGuid" `       // หน่วยนับมาตรฐาน (นับสต๊อก)
	// UnitList         []InventoryUnit `json:"unitlist" bson:"unitList" `        // กรณีหลายหน่วยนับ ตารางหน่วบนับ
}

type Barcode struct {
	Barcode  string  `json:"barcode" bson:"barcode" gorm:"barcode"`
	UnitCode string  `json:"unitcode" bson:"unitcode" gorm:"unitcode"`
	UnitName string  `json:"unitname" bson:"unitname" gorm:"unitname"`
	Price    float64 `json:"price" bson:"price" gorm:"price"`
	Image    string  `json:"image" bson:"image" gorm:"image"`
}

type UnitUse struct {
	UnitCode string `json:"unitcode" bson:"unitcode" gorm:"unitcode"`
	// UnitName       string  `json:"unitname" bson:"unitname" gorm:"unitname"`
	common.Name    `bson:"inline"`
	ItemUnitSTD    float64 `json:"itemunitstd" bson:"itemunitstd" gorm:"itemunitstd"`
	ItemUnitDIV    float64 `json:"itemunitdiv" bson:"itemunitdiv" gorm:"itemunitdiv"`
	IsUnitCost     bool    `json:"isunitcost" bson:"isunitcost" gorm:"isunitcost,type:bool,default:false"`
	IsUnitStandard bool    `json:"isunitstandard" bson:"isunitstandard" gorm:"isunitstandard,type:bool,default:false"`
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
