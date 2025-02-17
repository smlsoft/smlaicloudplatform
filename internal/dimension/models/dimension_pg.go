package models

import (
	"smlaicloudplatform/internal/models"
	"time"
)

// ✅ **โครงสร้างหลัก Dimension**
type DimensionPg struct {
	ShopID     string            `json:"shopid" gorm:"column:shopid;primaryKey"`
	GuidFixed  string            `json:"guidfixed" gorm:"column:guidfixed;primaryKey"`
	Names      models.JSONB      `json:"names" gorm:"column:names;type:jsonb"`
	IsDisabled bool              `json:"isdisabled" gorm:"column:isdisabled"`
	CreatedAt  time.Time         `json:"createdat" gorm:"column:createdat;autoCreateTime"`
	CreatedBy  string            `json:"createdby" gorm:"column:createdby"`
	UpdatedAt  time.Time         `json:"updatedat" gorm:"column:updatedat;autoUpdateTime"`
	UpdatedBy  string            `json:"updatedby" gorm:"column:updatedby"`
	Items      []DimensionItemPg `json:"items" gorm:"-"`
}

func (DimensionPg) TableName() string {
	return "dimensions"
}
