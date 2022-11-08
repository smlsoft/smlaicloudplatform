package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const productcategoryCollectionName = "productCategories"

type ProductCategory struct {
	models.PartitionIdentity `bson:"inline"`
	ChildCount               int             `json:"childcount" bson:"childcount"`
	ParentGUID               string          `json:"parentguid" bson:"parentguid"`
	ParentGUIDAll            string          `json:"parentguidall" bson:"parentguidall"`
	ImageUri                 string          `json:"imageuri" bson:"imageuri"`
	Names                    *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	XSorts                   *[]models.XSort `json:"xsorts" bson:"xsorts" validate:"unique=Code,dive"`
	Barcodes                 *[]models.XSort `json:"barcodes" bson:"barcodes" validate:"unique=Code,dive"`
	UseImageOrColor          bool            `json:"useimageorcolor" bson:"useimageorcolor"`
	ColorSelect              string          `json:"colorselect" bson:"colorselect"`
	ColorSelectHex           string          `json:"colorselecthex" bson:"colorselecthex"`
}

type ProductCategoryInfo struct {
	models.DocIdentity `bson:"inline"`
	ProductCategory    `bson:"inline"`
}

func (ProductCategoryInfo) CollectionName() string {
	return productcategoryCollectionName
}

type ProductCategoryData struct {
	models.ShopIdentity `bson:"inline"`
	ProductCategoryInfo `bson:"inline"`
}

type ProductCategoryDoc struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductCategoryData `bson:"inline"`
	models.ActivityDoc  `bson:"inline"`
}

func (ProductCategoryDoc) CollectionName() string {
	return productcategoryCollectionName
}

type ProductCategoryItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (ProductCategoryItemGuid) CollectionName() string {
	return productcategoryCollectionName
}

type ProductCategoryActivity struct {
	ProductCategoryData `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (ProductCategoryActivity) CollectionName() string {
	return productcategoryCollectionName
}

type ProductCategoryDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (ProductCategoryDeleteActivity) CollectionName() string {
	return productcategoryCollectionName
}
