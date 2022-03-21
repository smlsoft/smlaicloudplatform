package models

type Identity struct {
	ShopID    string `json:"shopID" bson:"shop_id"`
	GuidFixed string `json:"guidFixed" bson:"GuidFixed"`
}

type ShopIdentity struct {
	ShopID string `json:"shopID" bson:"shop_id"`
}

type DocIdentity struct {
	GuidFixed string `json:"guidFixed" bson:"GuidFixed"`
}
