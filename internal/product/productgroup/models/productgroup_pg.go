package models

import (
	"smlaicloudplatform/internal/models"
	"time"
)

type ProductGroupPg struct {
	ShopID    string       `json:"shopid" gorm:"column:shopid;primaryKey"`
	GuidFixed string       `json:"guidfixed" bson:"guidfixed" gorm:"column:guidfixed;" `
	Code      string       `json:"code" gorm:"column:code;primaryKey"`
	Names     models.JSONB `json:"names" gorm:"column:names;type:jsonb"`
	CreatedBy string       `json:"createdby" gorm:"column:createdby"`
	CreatedAt time.Time    `json:"createdat" gorm:"column:createdat"`
	UpdatedBy string       `json:"updatedby" gorm:"column:updatedby"`
	UpdatedAt time.Time    `json:"updatedat" gorm:"column:updatedat"`
}

func (ProductGroupPg) TableName() string {
	return "product_groups"
}
