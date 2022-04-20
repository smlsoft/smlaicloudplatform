package models

import "time"

type Activity struct {
	CreatedBy string    `json:"-" bson:"createdby"`
	CreatedAt time.Time `json:"-" bson:"createdat"`
	UpdatedBy string    `json:"-" bson:"updatedby,omitempty"`
	UpdatedAt time.Time `json:"-" bson:"updatedat,omitempty"`
	DeletedBy string    `json:"-" bson:"deletedby,omitempty"`
	DeletedAt time.Time `json:"-" bson:"deletedat,omitempty"`
}

type LastActivity struct {
	New    interface{} `json:"new" `
	Remove interface{} `json:"remove"`
}

type LastUpdate struct {
	LastUpdatedAt time.Time `json:"lastupdatedat" bson:"lastupdatedat"`
}
