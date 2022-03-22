package models

type Identity struct {
	ShopID    string `json:"shopID" bson:"shopID"`
	GuidFixed string `json:"guidFixed" bson:"guidFixed"`
}

type ShopIdentity struct {
	ShopID string `json:"shopID" bson:"shopID"`
}

type DocIdentity struct {
	GuidFixed string `json:"guidFixed" bson:"guidFixed"`
}
