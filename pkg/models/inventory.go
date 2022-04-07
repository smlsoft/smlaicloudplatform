package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type InventoryUnit struct {
// 	UnitGuid  string  `bson:"unitguid" json:"unitGuid"`   // Guid หน่วยนับ
// 	UnitName  string  `bson:"unitname" json:"unitName"`   // ชื่อหน่วยนับ
// 	Minuend   float32 `bson:"minuend" json:"minuend"`     // ตัวตั้ง
// 	Divisor   float32 `bson:"divisor" json:"divisor"`     // ตัวหาร
// 	Activated bool    `bson:"activated" json:"activated"` // เปิดใช้งานอยู่
// }

const inventoryCollectionName string = "inventories"
const inventoryTableName string = "inventories"
const inventoryIndexName string = "inventories_index"

type Inventory struct {
	ItemSku      string  `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	Barcode      string  `json:"barcode" bson:"barcode"`
	CategoryGuid string  `json:"categoryGuid" bson:"categoryGuid"` // Guid กลุ่มสินค้า
	Price        float32 `json:"price" bson:"price" `              // ราคาพื้นฐาน (กรณีไม่มีตารางราคา และโปรโมชั่น)
	MemberPrice  float32 `json:"memberPrice,omitempty" bson:"memberPrice,omitempty"`
	Recommended  bool    `json:"recommended,omitempty" bson:"recommended,omitempty" ` // สินค้าแนะนำ
	Activated    bool    `json:"activated,omitempty" bson:"activated,omitempty"`      // เปิดใช้งานอยู่

	Name1        string `json:"name1" bson:"name1"`                                   // ชื่อภาษาไทย
	Description1 string `json:"description1,omitempty" bson:"description1,omitempty"` // รายละเอียดภาษาไทย
	Name2        string `json:"name2,omitempty" bson:"name2,omitempty"`
	Description2 string `json:"description2,omitempty" bson:"description2,omitempty"`
	Name3        string `json:"name3,omitempty" bson:"name3,omitempty"`
	Description3 string `json:"description3,omitempty" bson:"description3,omitempty"`
	Name4        string `json:"name4,omitempty" bson:"name4,omitempty"`
	Description4 string `json:"description4,omitempty" bson:"description4,omitempty"`
	Name5        string `json:"name5,omitempty" bson:"name5,omitempty"`
	Description5 string `json:"description5,omitempty" bson:"description5,omitempty"`

	UnitName1 string            `json:"unitName1" bson:"unitName1" gorm:"unit_name1"`
	UnitName2 string            `json:"unitName2,omitempty" bson:"unitName2,omitempty" gorm:"unit_name1"`
	UnitName3 string            `json:"unitName3,omitempty" bson:"unitName3,omitempty" gorm:"unit_name2"`
	UnitName4 string            `json:"unitName4,omitempty" bson:"unitName4,omitempty" gorm:"unit_name3"`
	UnitName5 string            `json:"unitName5,omitempty" bson:"unitName5,omitempty" gorm:"unit_name4"`
	Options   *[]Option         `json:"options,omitempty" bson:"options,omitempty" gorm:"many2many:inventory_options;foreignKey:GuidFixed;joinForeignKey:DocID;References:Code;joinReferences:OptID"`
	Images    *[]InventoryImage `json:"images,omitempty" bson:"images,omitempty" gorm:"images;foreignKey:DocID"`
	Tags      *[]InventoryTag   `json:"tags,omitempty" bson:"tags" gorm:"tags;foreignKey:DocID"`

	// WaitType         int             `json:"-" bson:"waitType"`                // ประเภทการรอ (สินค้าหมด)
	// WaitUntil        time.Time       `json:"-" bson:"waitUntil"`               // ระยะเวลาที่รอ
	// MultipleUnits    bool            `json:"-" bson:"multipleuUits" `          // สินค้าหลายหน่วยนับ
	// UnitStandardGuid string          `json:"-" bson:"unitStandardGuid" `       // หน่วยนับมาตรฐาน (นับสต๊อก)
	// UnitList         []InventoryUnit `json:"unitList" bson:"unitList" `        // กรณีหลายหน่วยนับ ตารางหน่วบนับ
}

type InventoryOption struct {
	DocID string `bson:"-" gorm:"doc_id;primaryKey"`
	OptID string `bson:"-" gorm:"opt_id;primaryKey"`
}

func (InventoryOption) TableName() string {
	return "inventory_options"
}

type InventoryImage struct {
	DocID string `bson:"-" gorm:"doc_id;primaryKey"`
	Url   string `bson:"url" gorm:"url;primaryKey"`
}

type InventoryTag struct {
	DocID string `bson:"-" gorm:"doc_id;primaryKey"`
	Name  string `bson:"name" gorm:"name;primaryKey"`
}

type Option struct {
	Code       string    `json:"code" bson:"code" gorm:"primaryKey"`
	Required   bool      `json:"required" bson:"required"`
	SelectMode string    `json:"selectMode" bson:"selectMode"`
	MaxSelect  int       `json:"maxSelect" bson:"maxSelect"`
	Name1      string    `json:"name1" bson:"name1"`
	Name2      string    `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3      string    `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4      string    `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5      string    `json:"name5,omitempty" bson:"name5,omitempty"`
	Choices    *[]Choice `json:"choices" bson:"choices" gorm:"choices;foreignKey:OptCode"`
}

type Choice struct {
	OptCode     string  `bson:"-" gorm:"opt_code;primaryKey" `
	Barcode     string  `json:"barcode" bson:"barcode" gorm:"barcode;primaryKey"`
	SuggestCode string  `json:"suggestCode" bson:"suggestCode"`
	Price       float64 `json:"price" bson:"price"`
	Qty         int     `json:"qty" bson:"qty"`
	Name1       string  `json:"name1" bson:"name1"`
	Name2       string  `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3       string  `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4       string  `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5       string  `json:"name5,omitempty" bson:"name5,omitempty"`
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
	ShopID                 string                      `json:"shopID" bson:"shopID"`
	GuidFixed              string                      `json:"guidFixed" bson:"guidFixed"`
	OptionName1            string                      `json:"optionName1" bson:"optionName1"`
	ProductSelectOption1   bool                        `json:"productSelectoPtion1" bson:"productSelectOption1"`
	ProductSelectOption2   bool                        `json:"productSelectoPtion2" bson:"productSelectOption2"`
	ProductSelectOptionMin int                         `json:"productSelectOptionMin" bson:"productSelectOptionMin"`
	ProductSelectOptionMax int                         `json:"productSelectOptionMax" bson:"productSelectOptionMax"`
	Details                []InventoryOptonGroupDetail `json:"details" bson:"details"`
	Activity
}

func (*InventoryOptionGroup) CollectionName() string {
	return "inventoryOptionGroup"
}

type InventoryOptonGroupDetail struct {
	GuidFixed   string  `json:"guidFixed" bson:"guidFixed"`
	DetailName1 string  `json:"detailName1" bson:"detailName1"`
	Amount      float32 `json:"amount" bson:"amount"`
}

type InventoryOptionMain struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID        string             `json:"shopID" bson:"shopID"`
	GuidFixed     string             `json:"guidFixed" bson:"guidFixed"`
	InventoryID   string             `json:"inventoryID" bson:"inventoryID"`
	OptionGroupID string             `json:"optionGroupID" bson:"optionGroupID"`
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
