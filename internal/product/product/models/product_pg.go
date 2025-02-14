package models

import (
	"smlaicloudplatform/internal/models"
)

type ProductPg struct {
	ShopID                   string `json:"shopid" gorm:"column:shopid;primaryKey"`
	Code                     string `json:"code" gorm:"column:code;primaryKey"`
	GuidFixed                string `json:"guidfixed" gorm:"column:guidfixed"`
	models.PartitionIdentity `gorm:"embedded;"`
	Names                    models.JSONB  `json:"names"  gorm:"column:names;type:jsonb"`
	GroupGuid                string        `json:"groupguid" gorm:"column:groupguid"`
	UnitGuid                 string        `json:"unitguid" gorm:"column:unitguid"`
	ItemType                 int8          `json:"itemtype" gorm:"column:itemtype"`
	ManufacturerGUID         string        `json:"manufacturerguid" gorm:"column:manufacturerguid"`
	Dimensions               []DimensionPg `json:"dimensions" gorm:"many2many:product_dimensions;foreignKey:Code;joinForeignKey:ProductCode;References:Code;joinReferences:DimensionCode"`
}

func (ProductPg) TableName() string {
	return "productbarcode"
}

type DimensionPg struct {
	ShopID    string       `json:"shopid" gorm:"column:shopid;primaryKey"`
	GuidFixed string       `json:"guidfixed" gorm:"column:guidfixed"`
	Code      string       `json:"code" gorm:"column:code;primaryKey"`
	Names     models.JSONB `json:"names"  gorm:"column:names;type:jsonb"`
}

type ProductDimensionPg struct {
	ProductCode   string `json:"product_code" gorm:"column:product_code;primaryKey"`
	DimensionCode string `json:"dimension_code" gorm:"column:dimension_code;primaryKey"`
}
