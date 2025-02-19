package models

import (
	dimension "smlaicloudplatform/internal/dimension/models"
	"smlaicloudplatform/internal/models"
	"time"
)

// ✅ **โครงสร้างหลักของ Product**
type ProductPg struct {
	ShopID                   string `json:"shopid" gorm:"column:shopid;primaryKey"`
	Code                     string `json:"code" gorm:"column:code;primaryKey"`
	GuidFixed                string `json:"guidfixed" gorm:"column:guidfixed"`
	models.PartitionIdentity `gorm:"embedded;"`
	Names                    models.JSONB            `json:"names"  gorm:"column:names;type:jsonb"`
	GroupGuid                *string                 `json:"groupguid,omitempty" gorm:"column:groupguid;default:null"`
	UnitGuid                 *string                 `json:"unitguid,omitempty" gorm:"column:unitguid;default:null"`
	ItemType                 int8                    `json:"itemtype" gorm:"column:itemtype;default:0"`
	ManufacturerGUID         *string                 `json:"manufacturerguid,omitempty" gorm:"column:manufacturerguid;default:null"`
	Dimensions               []dimension.DimensionPg `json:"dimensions" gorm:"-"`
	GroupCode                *string                 `json:"groupcode" gorm:"-"`
	GroupName                models.JSONB            `json:"groupname" gorm:"-"`
	ManufacturerCode         *string                 `json:"manufacturercode" gorm:"-"`
	ManufacturerName         []models.NameX          `json:"manufacturername" gorm:"-"`
	CreatedAt                time.Time               `json:"createdat" gorm:"column:createdat"`
	UpdatedAt                time.Time               `json:"updatedat" gorm:"column:updatedat"`
	CreatedBy                string                  `json:"createdby" gorm:"column:createdby"`
	UpdatedBy                string                  `json:"updatedby" gorm:"column:updatedby"`
}

func (ProductPg) TableName() string {
	return "products"
}
