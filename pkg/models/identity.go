package models

type Identity struct {
	ShopID    string `json:"shopid" bson:"shopid" gorm:"shopid"`
	GuidFixed string `json:"guidfixed" bson:"guidfixed" gorm:"guidfixed;primaryKey"`
}

type ShopIdentity struct {
	ShopID string `json:"shopid" bson:"shopid" gorm:"column:shopid;primaryKey"`
}

type DocIdentity struct {
	GuidFixed string `json:"guidfixed" bson:"guidfixed" gorm:"guidfixed;primaryKey"`
}

type PartitionIdentity struct {
	ParID string `json:"parid"  gorm:"parid"`
}
