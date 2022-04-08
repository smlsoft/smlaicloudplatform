package models

type Identity struct {
	ShopID    string `json:"shopid" bson:"shopid" gorm:"shopid"`
	GuidFixed string `json:"guidfixed" bson:"guidfixed" gorm:"guidfixed;primaryKey"`
}

type ShopIdentity struct {
	ShopID string `json:"shopid" bson:"shopid" gorm:"shopid"`
}

type DocIdentity struct {
	GuidFixed string `json:"guidfixed" bson:"guidfixed" gorm:"guidfixed;primaryKey"`
}
