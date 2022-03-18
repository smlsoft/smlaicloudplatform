package models

import "time"

type Activity struct {
	CreatedBy string    `json:"-" bson:"createdBy"`
	CreatedAt time.Time `json:"-" bson:"createdAt"`
	UpdatedBy string    `json:"-" bson:"updatedBy,omitempty"`
	UpdatedAt time.Time `json:"-" bson:"updatedAt,omitempty"`
	Deleted   bool      `json:"-" bson:"deleted"`
}
