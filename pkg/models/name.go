package models

type Name struct {
	Name1 *string `json:"name1" bson:"name1" validate:"required"`
	Name2 *string `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3 *string `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4 *string `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5 *string `json:"name5,omitempty" bson:"name5,omitempty"`
}

type UnitName struct {
	UnitName1 *string `json:"unitname1" bson:"unitname1" gorm:"unitname1"`
	UnitName2 *string `json:"unitname2,omitempty" bson:"unitname2,omitempty" gorm:"unitname2,omitempty"`
	UnitName3 *string `json:"unitname3,omitempty" bson:"unitname3,omitempty" gorm:"unitname3,omitempty"`
	UnitName4 *string `json:"unitname4,omitempty" bson:"unitname4,omitempty" gorm:"unitname4,omitempty"`
	UnitName5 *string `json:"unitname5,omitempty" bson:"unitname5,omitempty" gorm:"unitname5,omitempty"`
}

type Description struct {
	Description1 *string `json:"description1,omitempty" bson:"description1,omitempty" gorm:"description1,omitempty"`
	Description2 *string `json:"description2,omitempty" bson:"description2,omitempty" gorm:"description2,omitempty"`
	Description3 *string `json:"description3,omitempty" bson:"description3,omitempty" gorm:"description3,omitempty"`
	Description4 *string `json:"description4,omitempty" bson:"description4,omitempty" gorm:"description4,omitempty"`
	Description5 *string `json:"description5,omitempty" bson:"description5,omitempty" gorm:"description5,omitempty"`
}
