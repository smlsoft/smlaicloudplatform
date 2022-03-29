package models

type Identity struct {
	ShopID    string `json:"shopID" bson:"shopID" gorm:"shop_id"`
	GuidFixed string `json:"guidFixed" bson:"guidFixed" gorm:"guid_fixed"`
}

type ShopIdentity struct {
	ShopID string `json:"shopID" bson:"shopID"`
}

type DocIdentity struct {
	GuidFixed string `json:"guidFixed" bson:"guidFixed"`
}
