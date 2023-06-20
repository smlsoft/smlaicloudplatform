package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const promotionCollectionName = "productPromotions"

type Promotion struct {
	models.PartitionIdentity `bson:"inline"`
	PromotionType            int8               `json:"promotiontype" bson:"promotiontype"`
	Code                     string             `json:"code" bson:"code"`
	Name                     string             `json:"name" bson:"name"`
	FromDate                 time.Time          `json:"fromdate" bson:"fromdate"`
	ToDate                   time.Time          `json:"todate" bson:"todate"`
	FromTime                 string             `json:"fromtime" bson:"fromtime"`
	ToTime                   string             `json:"totime" bson:"totime"`
	Barcode                  string             `json:"barcode" bson:"barcode"`
	IsMemberOnly             bool               `json:"ismemberonly" bson:"ismemberonly"`
	Remark                   string             `json:"remark" bson:"remark"`
	IsUseInMonday            bool               `json:"isuseinmonday" bson:"isuseinmonday"`
	IsUseInTuesday           bool               `json:"isuseintuesday" bson:"isuseintuesday"`
	IsUseInWednesday         bool               `json:"isuseinwednesday" bson:"isuseinwednesday"`
	IsUseInThursday          bool               `json:"isuseinthursday" bson:"isuseinthursday"`
	IsUseInFriday            bool               `json:"isuseinfriday" bson:"isuseinfriday"`
	IsUseInSaturday          bool               `json:"isuseinsaturday" bson:"isuseinsaturday"`
	IsUseInSunday            bool               `json:"isuseinsunday" bson:"isuseinsunday"`
	Details                  *[]PromotionDetail `json:"details" bson:"details"`
}

type PromotionDetail struct {
	DetailType int8    `json:"detailtype" bson:"detailtype"` // 0: discount, 1: buy x get y, 2: buy x get y with discount
	Minimum    float64 `json:"minimum" bson:"minimum"`
	Discount   float64 `json:"discount" bson:"discount"`
	Barcode    string  `json:"barcode" bson:"barcode"`
}

type PromotionInfo struct {
	models.DocIdentity `bson:"inline"`
	Promotion          `bson:"inline"`
}

func (PromotionInfo) CollectionName() string {
	return promotionCollectionName
}

type PromotionData struct {
	models.ShopIdentity `bson:"inline"`
	PromotionInfo       `bson:"inline"`
}

type PromotionDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PromotionData      `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (PromotionDoc) CollectionName() string {
	return promotionCollectionName
}

type PromotionItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (PromotionItemGuid) CollectionName() string {
	return promotionCollectionName
}

type PromotionActivity struct {
	PromotionData       `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PromotionActivity) CollectionName() string {
	return promotionCollectionName
}

type PromotionDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PromotionDeleteActivity) CollectionName() string {
	return promotionCollectionName
}
