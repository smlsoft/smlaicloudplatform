package models

import (
	"smlaicloudplatform/internal/models"
)

// ✅ **โครงสร้างรายการของ Dimension**
type DimensionItemPg struct {
	ShopID        string       `json:"shopid" gorm:"column:shopid;primaryKey"`
	GuidFixed     string       `json:"guidfixed" gorm:"column:guidfixed;primaryKey"`
	DimensionGuid string       `json:"dimension_guid" gorm:"column:dimension_guid;index"`
	Names         models.JSONB `json:"names" gorm:"column:names;type:jsonb"`
	IsDisabled    bool         `json:"isdisabled" gorm:"column:isdisabled"`
}

func (DimensionItemPg) TableName() string {
	return "dimension_items"
}
