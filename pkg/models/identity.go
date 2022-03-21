package models

type Identity struct {
	GuidFixed string `json:"guidFixed" bson:"GuidFixed"`
	ShopID    string `json:"shopID" bson:"shop_id"`
}
