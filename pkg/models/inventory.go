package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryUnit struct {
	UnitGuid  string  `bson:"unitguid" json:"unitGuid"`   // Guid หน่วยนับ
	UnitName  string  `bson:"unitname" json:"unitName"`   // ชื่อหน่วยนับ
	Minuend   float32 `bson:"minuend" json:"minuend"`     // ตัวตั้ง
	Divisor   float32 `bson:"divisor" json:"divisor"`     // ตัวหาร
	Activated bool    `bson:"activated" json:"activated"` // เปิดใช้งานอยู่
}

type Inventory struct {
	Id               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ItemSku          string             `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	MerchantId       string             `json:"merchantId" bson:"merchantId"`               // รหัสร้าน
	GuidFixed        string             `json:"guidFixed,omitempty" bson:"guidFixed"`       // Guid สินค้า
	CategoryGuid     string             `json:"categoryGuid,omitempty" bson:"categoryGuid"` // Guid กลุ่มสินค้า
	LineNumber       int                `json:"lineNumber,omitempty" bson:"lineNumber"`     // บรรทัดที่ (เอาไว้เรียงลำดับ)
	Price            float32            `json:"price" bson:"price" `                        // ราคาพื้นฐาน (กรณีไม่มีตารางราคา และโปรโมชั่น)
	Recommended      bool               `json:"recommended" bson:"recommended" `            // สินค้าแนะนำ
	HaveImage        bool               `json:"haveImage" bson:"haveImage"`                 // มีรูปสินค้า
	Activated        bool               `json:"activated" bson:"activated"`                 // เปิดใช้งานอยู่
	WaitType         int                `json:"-" bson:"waitType"`                          // ประเภทการรอ (สินค้าหมด)
	WaitUntil        time.Time          `json:"-" bson:"waitUntil"`                         // ระยะเวลาที่รอ
	MultipleUnits    bool               `json:"-" bson:"multipleuUits" `                    // สินค้าหลายหน่วยนับ
	UnitStandardGuid string             `json:"-" bson:"unitStandardGuid" `                 // หน่วยนับมาตรฐาน (นับสต๊อก)
	UnitList         []InventoryUnit    `json:"unitList" bson:"unitList" `                  // กรณีหลายหน่วยนับ ตารางหน่วบนับ
	Name1            string             `json:"name1" bson:"name1"`                         // ชื่อภาษาไทย
	Description1     string             `json:"description1" bson:"description1"`           // รายละเอียดภาษาไทย
	Name2            string             `json:"name2,omitempty" bson:"name2,omitempty"`
	Description2     string             `json:"description2" bson:"description2,omitempty"`
	Name3            string             `json:"name3,omitempty" bson:"name3,omitempty"`
	Description3     string             `json:"description3" bson:"description3,omitempty"`
	Name4            string             `json:"name4,omitempty" bson:"name4,omitempty"`
	Description4     string             `json:"description4" bson:"description4,omitempty"`
	Name5            string             `json:"name5,omitempty" bson:"name5,omitempty"`
	Description5     string             `json:"description5" bson:"description5,omitempty"`
	Barcodes         []Barcode          `json:"barcodes,omitempty" bson:"barcodes,omitempty"`
	CreatedBy        string             `json:"-" bson:"createdBy"`
	CreatedAt        time.Time          `json:"-" bson:"createdAt"`
	UpdatedBy        string             `json:"-" bson:"updatedBy,omitempty"`
	UpdatedAt        time.Time          `json:"-" bson:"updatedAt,omitempty"`
	Deleted          bool               `json:"-" bson:"deleted"` // ลบแล้ว
}

func (*Inventory) CollectionName() string {
	return "inventory"
}

type Category struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MerchantId string             `json:"merchantId" bson:"merchantId"`
	GuidFixed  string             `json:"guidFixed" bson:"guidFixed"`
	LineNumber int                `json:"lineNumber" bson:"lineNumber"`
	Name1      string             `json:"name1" bson:"name1"`
	HaveImage  bool               `json:"haveImage" bson:"haveImage"`
	CreatedBy  string             `json:"-" bson:"createdBy"`
	CreatedAt  time.Time          `json:"-" bson:"createdAt"`
	UpdatedBy  string             `json:"-" bson:"updatedBy,omitempty"`
	UpdatedAt  time.Time          `json:"-" bson:"updatedAt,omitempty"`
	Deleted    bool               `json:"-" bson:"deleted"`
}

func (*Category) CollectionName() string {
	return "category"
}

type InventoryOptonGroupDetail struct {
	// Id          string  `json:"id" bson:"id,omitempty"`
	// MerchantId  string  `json:"merchantId" bson:"merchantId"`
	GuidFixed   string  `json:"guidFixed" bson:"guidFixed"`
	DetailName1 string  `json:"detailName1" bson:"detailName1"`
	Amount      float32 `json:"amount" bson:"amount"`
}

type InventoryOptionGroup struct {
	Id                     string                      `json:"id" bson:"id,omitempty"`
	MerchantId             string                      `json:"merchantId" bson:"merchantId"`
	GuidFixed              string                      `json:"guidFixed" bson:"guidFixed"`
	OptionName1            string                      `json:"optionName1" bson:"optionName1"`
	ProductSelectOption1   bool                        `json:"productSelectoPtion1" bson:"productSelectOption1"`
	ProductSelectOption2   bool                        `json:"productSelectoPtion2" bson:"productSelectOption2"`
	ProductSelectOptionMin int                         `json:"productSelectOptionMin" bson:"productSelectOptionMin"`
	ProductSelectOptionMax int                         `json:"productSelectOptionMax" bson:"productSelectOptionMax"`
	Details                []InventoryOptonGroupDetail `json:"details" bson:"details"`
	CreatedBy              string                      `json:"-" bson:"createdBy"`
	CreatedAt              time.Time                   `json:"-" bson:"createdAt"`
	UpdatedBy              string                      `json:"-" bson:"updatedBy,omitempty"`
	UpdatedAt              time.Time                   `json:"-" bson:"updatedAt,omitempty"`
	Deleted                bool                        `json:"-" bson:"deleted"`
}

func (*InventoryOptionGroup) CollectionName() string {
	return "inventoryOptionGroup"
}

type InventoryOption struct {
	Id            string `json:"id" bson:"id,omitempty"`
	MerchantId    string `json:"merchantId" bson:"merchantId"`
	OptionGroupId string `json:"optionGroupId" bson:"optionGroupId"`
}
