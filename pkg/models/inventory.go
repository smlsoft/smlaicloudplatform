package models

import (
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
	Barcode      string  `json:"barcode" bson:"barcode" gorm:"barcode"`
	CategoryGuid string  `json:"categoryguid" bson:"categoryguid" gorm:"categoryguid"` // Guid กลุ่มสินค้า
	Price        float32 `json:"price" bson:"price" gorm:"price"`                      // ราคาพื้นฐาน (กรณีไม่มีตารางราคา และโปรโมชั่น)
	MemberPrice  float32 `json:"memberprice,omitempty" bson:"memberprice,omitempty" gorm:"memberprice,omitempty"`
	Recommended  bool    `json:"recommended,omitempty" bson:"recommended,omitempty" gorm:"recommended,omitempty"` // สินค้าแนะนำ
	Activated    bool    `json:"activated,omitempty" bson:"activated,omitempty" gorm:"activated,omitempty"`       // เปิดใช้งานอยู่

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

	UnitName1 string `json:"unitname1" bson:"unitname1" gorm:"unitname1" gorm:"unitname1"`
	UnitName2 string `json:"unitname2,omitempty" bson:"unitname2,omitempty" gorm:"unitname2,omitempty"`
	UnitName3 string `json:"unitname3,omitempty" bson:"unitname3,omitempty" gorm:"unitname3,omitempty"`
	UnitName4 string `json:"unitname4,omitempty" bson:"unitname4,omitempty" gorm:"unitname4,omitempty"`
	UnitName5 string `json:"unitname5,omitempty" bson:"unitname5,omitempty" gorm:"unitname5,omitempty"`

	ItemGuid        string  `json:"itemguid,omitempty" bson:"itemguid,omitempty" gorm:"itemguid,omitempty"`
	ItemCode        string  `json:"itemcode,omitempty" bson:"itemcode,omitempty" gorm:"itemcode,omitempty"`
	ItemUnitCode    string  `json:"itemunitcode,omitempty" bson:"itemunitcode,omitempty" gorm:"itemunitcode,omitempty"`
	ItemUnitStd     float64 `json:"itemunitstd,omitempty" bson:"itemunitstd,omitempty" gorm:"itemunitstd,omitempty"`
	ItemUnitDiv     float64 `json:"itemunitdiv,omitempty" bson:"itemunitdiv,omitempty" gorm:"itemunitdiv,omitempty"`
	ItemOptionCode  string  `json:"itemoptioncode,omitempty" bson:"itemoptioncode,omitempty" gorm:"itemoptioncode,omitempty"`
	ItemOptionOrder int8    `json:"itemoptionorder,omitempty" bson:"itemoptionorder,omitempty" gorm:"itemoptionorder,omitempty"`

	Options *[]Option         `json:"options,omitempty" bson:"options,omitempty" gorm:"many2many:inventoryoptions;foreignKey:GuidFixed;joinForeignKey:DocID;References:Code;joinReferences:OptID"`
	Images  *[]InventoryImage `json:"images,omitempty" bson:"images,omitempty" gorm:"images;foreignKey:DocID"`
	Tags    *[]InventoryTag   `json:"tags,omitempty" bson:"tags" gorm:"tags;foreignKey:DocID"`

	// WaitType         int             `json:"-" bson:"waitType"`                // ประเภทการรอ (สินค้าหมด)
	// WaitUntil        time.Time       `json:"-" bson:"waitUntil"`               // ระยะเวลาที่รอ
	// MultipleUnits    bool            `json:"-" bson:"multipleuUits" `          // สินค้าหลายหน่วยนับ
	// UnitStandardGuid string          `json:"-" bson:"unitStandardGuid" `       // หน่วยนับมาตรฐาน (นับสต๊อก)
	// UnitList         []InventoryUnit `json:"unitlist" bson:"unitList" `        // กรณีหลายหน่วยนับ ตารางหน่วบนับ
}

type InventoryOption struct {
	DocID string `bson:"-" gorm:"docid;primaryKey"`
	OptID string `bson:"-" gorm:"optid;primaryKey"`
}

func (InventoryOption) TableName() string {
	return "inventoryoptions"
}

type InventoryImage struct {
	DocID string `json:"-" bson:"-" gorm:"docid;primaryKey"`
	Url   string `json:"url" bson:"url" gorm:"url;primaryKey"`
}

type InventoryTag struct {
	DocID string `json:"-" bson:"-" gorm:"docid;primaryKey"`
	Name  string `json:"name" bson:"name" gorm:"name;primaryKey"`
}

type Option struct {
	Code       string    `json:"code" bson:"code" gorm:"code;primaryKey"`
	Required   bool      `json:"required,omitempty" bson:"required,omitempty" gorm:"required,omitempty"`
	ChoiceType int8      `json:"choicetype,omitempty" bson:"choicetype,omitempty" gorm:"choicetype,omitempty"`
	MaxSelect  int8      `json:"maxselect,omitempty" bson:"maxselect,omitempty" gorm:"maxselect,omitempty"`
	Name1      string    `json:"name1" bson:"name1" gorm:"name1"`
	Name2      string    `json:"name2,omitempty" bson:"name2,omitempty" gorm:"name2,omitempty"`
	Name3      string    `json:"name3,omitempty" bson:"name3,omitempty" gorm:"name3,omitempty"`
	Name4      string    `json:"name4,omitempty" bson:"name4,omitempty" gorm:"name4,omitempty"`
	Name5      string    `json:"name5,omitempty" bson:"name5,omitempty" gorm:"name5,omitempty"`
	Choices    *[]Choice `json:"choices" bson:"choices" gorm:"choices;foreignKey:OptCode"`
}

type Choice struct {
	OptCode     string  `bson:"-" gorm:"optcode;primaryKey" `
	Barcode     string  `json:"barcode" bson:"barcode" gorm:"barcode;primaryKey"`
	SuggestCode string  `json:"suggestcode,omitempty" bson:"suggestcode,omitempty" gorm:"suggestcode,omitempty"`
	Price       float64 `json:"price" bson:"price" gorm:"price"`
	Qty         float64 `json:"qty" bson:"qty" gorm:"qty"`
	QtyMax      float64 `json:"qtymax" bson:"qtymax" gorm:"qtymax"`
	Name1       string  `json:"name1" bson:"name1" gorm:"name1"`
	Name2       string  `json:"name2,omitempty" bson:"name2,omitempty" gorm:"name2,omitempty"`
	Name3       string  `json:"name3,omitempty" bson:"name3,omitempty" gorm:"name3,omitempty"`
	Name4       string  `json:"name4,omitempty" bson:"name4,omitempty" gorm:"name4,omitempty"`
	Name5       string  `json:"name5,omitempty" bson:"name5,omitempty" gorm:"name5,omitempty"`
	ItemUnit    string  `json:"itemunit,omitempty" bson:"itemunit" gorm:"itemunit,omitempty"`
	Selected    bool    `json:"selected,omitempty" bson:"selected,omitempty" gorm:"selected,omitempty"`
	Default     bool    `json:"default,omitempty" bson:"default,omitempty" gorm:"default,omitempty"`
}

type InventoryInfo struct {
	DocIdentity `bson:"inline" gorm:"embedded;"`
	Inventory   `bson:"inline" gorm:"embedded;"`
}

func (InventoryInfo) CollectionName() string {
	return inventoryCollectionName
}

type InventoryData struct {
	ShopIdentity  `bson:"inline" gorm:"embedded;"`
	InventoryInfo `bson:"inline" gorm:"embedded;"`
}

func (InventoryData) TableName() string {
	return inventoryTableName
}

type InventoryDoc struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	InventoryData `bson:"inline"`
	Activity      `bson:"inline"`
}

func (InventoryDoc) CollectionName() string {
	return inventoryCollectionName
}

type InventoryIndex struct {
	Index `bson:"inline"`
}

func (InventoryIndex) TableName() string {
	return inventoryIndexName
}

/* */
type InventoryOptionGroup struct {
	ID                     primitive.ObjectID          `json:"id" bson:"_id,omitempty"`
	ShopID                 string                      `json:"shopid" bson:"shopid"`
	GuidFixed              string                      `json:"guidfixed" bson:"guidfixed"`
	OptionName1            string                      `json:"optionname1" bson:"optionname1"`
	ProductSelectOption1   bool                        `json:"productselectoption1" bson:"productselectoption1"`
	ProductSelectOption2   bool                        `json:"productselectoption2" bson:"productselectoption2"`
	ProductSelectOptionMin int                         `json:"productselectoptionmin" bson:"productselectoptionmin"`
	ProductSelectOptionMax int                         `json:"productselectoptionmax" bson:"productselectoptionmax"`
	Details                []InventoryOptonGroupDetail `json:"details" bson:"details"`
	Activity
}

func (*InventoryOptionGroup) CollectionName() string {
	return "inventoryOptionGroup"
}

type InventoryOptonGroupDetail struct {
	GuidFixed   string  `json:"guidfixed" bson:"guidfixed"`
	DetailName1 string  `json:"detailname1" bson:"detailname1"`
	Amount      float32 `json:"amount" bson:"amount"`
}

type InventoryOptionMain struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID        string             `json:"shopid" bson:"shopid"`
	GuidFixed     string             `json:"guidfixed" bson:"guidfixed"`
	InventoryID   string             `json:"inventoryid" bson:"inventoryid"`
	OptionGroupID string             `json:"optiongroupid" bson:"optiongroupid"`
	Activity
}

func (*InventoryOptionMain) CollectionName() string {
	return "inventoryOption"
}

type InventoryPageResponse struct {
	Success    bool                   `json:"success"`
	Data       []InventoryInfo        `json:"data,omitempty"`
	Pagination PaginationDataResponse `json:"pagination,omitempty"`
}

type InventoryInfoResponse struct {
	Success bool          `json:"success"`
	Data    InventoryInfo `json:"data,omitempty"`
}

type InventoryOptionGroupResponse struct {
	Success    bool                   `json:"success"`
	Data       []InventoryInfo        `json:"data,omitempty"`
	Pagination PaginationDataResponse `json:"pagination,omitempty"`
}

type InventoryOptionGroupInfoResponse struct {
	Success bool          `json:"success"`
	Data    InventoryInfo `json:"data,omitempty"`
}
