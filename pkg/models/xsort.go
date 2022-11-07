package models

type XSort struct {
	Code   string `json:"code" bson:"code"`
	XOrder uint   `json:"xorder" bson:"xorder,omitempty" validate:"min=0,max=4294967295"`
}
