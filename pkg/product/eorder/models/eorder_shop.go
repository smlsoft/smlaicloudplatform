package models

type EOrderShop struct {
	ShopID     string `json:"shopid" bson:"shopid"`
	Name1      string `json:"name1" bson:"name1"`
	TotalTable int    `json:"totaltable" bson:"totaltable"`
}
