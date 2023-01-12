package models

import (
	"time"

	memberModel "smlcloudplatform/pkg/member/models"
	common "smlcloudplatform/pkg/models"
	inventoryModel "smlcloudplatform/pkg/product/inventory/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const purchaseCollectionName = "purchases"
const purchaseIndexName = "purchases"

type Purchase struct {
	DocDate        *time.Time         `json:"docdate,omitempty" bson:"docdate,omitempty"`
	DocNo          string             `json:"docno,omitempty"  bson:"docno,omitempty"`
	Member         memberModel.Member `json:"member,omitempty"  bson:"member,omitempty"`
	Items          *[]PurchaseDetail  `json:"items" bson:"items" `
	TotalAmount    float64            `json:"totalamount" bson:"totalamount" `
	TaxRate        float64            `json:"taxrate" bson:"taxrate" `
	TaxAmount      float64            `json:"taxamount" bson:"taxamount" `
	TaxBaseAmount  float64            `json:"taxbaseamount" bson:"taxbaseamount" `
	DiscountAmount float64            `json:"discountamount" bson:"discountamount" `
	Payment        Payment            `json:"payment" bson:"payment"`
	SumAmount      float64            `json:"sumamount" bson:"sumamount" `
}

type PurchaseDetail struct {
	LineNumber                   int `json:"linenumber" bson:"linenumber"`
	inventoryModel.InventoryInfo `bson:"inline" gorm:"embedded;"`
	Price                        float64 `json:"price" bson:"price" `
	Qty                          float64 `json:"qty" bson:"qty" `
	DiscountAmount               float64 `json:"discountamount" bson:"discountamount"`
	DiscountText                 string  `json:"discounttext" bson:"discounttext"`
	Amount                       float64 `json:"amount" bson:"amount"`
}

type PurchaseInfo struct {
	common.DocIdentity `bson:"inline"`
	Purchase           `bson:"inline"`
}

func (PurchaseInfo) CollectionName() string {
	return purchaseCollectionName
}

type PurchaseData struct {
	common.ShopIdentity `bson:"inline"`
	PurchaseInfo        `bson:"inline"`
}

func (PurchaseData) IndexName() string {
	return purchaseIndexName
}

type PurchaseDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PurchaseData       `bson:"inline"`
	common.ActivityDoc `bson:"inline"`
}

func (PurchaseDoc) CollectionName() string {
	return purchaseCollectionName
}

type PurchaseListPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []PurchaseInfo                `json:"data,omitempty"`
	Pagination common.PaginationDataResponse `json:"pagination,omitempty"`
}
