package models

type Name struct {
	Name1 string  `json:"name1" bson:"name1" validate:"required,max=255"`
	Name2 *string `json:"name2,omitempty" bson:"name2" validate:"omitempty,max=255"`
	Name3 *string `json:"name3,omitempty" bson:"name3" validate:"omitempty,max=255"`
	Name4 *string `json:"name4,omitempty" bson:"name4" validate:"omitempty,max=255"`
	Name5 *string `json:"name5,omitempty" bson:"name5" validate:"omitempty,max=255"`
}

type UnitName struct {
	UnitName1 string  `json:"unitname1" bson:"unitname1" gorm:"unitname1" validate:"required,max=255"`
	UnitName2 *string `json:"unitname2,omitempty" bson:"unitname2,omitempty" gorm:"unitname2,omitempty" validate:"omitempty,max=255"`
	UnitName3 *string `json:"unitname3,omitempty" bson:"unitname3,omitempty" gorm:"unitname3,omitempty" validate:"omitempty,max=255"`
	UnitName4 *string `json:"unitname4,omitempty" bson:"unitname4,omitempty" gorm:"unitname4,omitempty" validate:"omitempty,max=255"`
	UnitName5 *string `json:"unitname5,omitempty" bson:"unitname5,omitempty" gorm:"unitname5,omitempty" validate:"omitempty,max=255"`
}

type Description struct {
	Description1 *string `json:"description1,omitempty" bson:"description1,omitempty" gorm:"description1,omitempty" validate:"omitempty,max=255"`
	Description2 *string `json:"description2,omitempty" bson:"description2,omitempty" gorm:"description2,omitempty" validate:"omitempty,max=255"`
	Description3 *string `json:"description3,omitempty" bson:"description3,omitempty" gorm:"description3,omitempty" validate:"omitempty,max=255"`
	Description4 *string `json:"description4,omitempty" bson:"description4,omitempty" gorm:"description4,omitempty" validate:"omitempty,max=255"`
	Description5 *string `json:"description5,omitempty" bson:"description5,omitempty" gorm:"description5,omitempty" validate:"omitempty,max=255"`
}

type NameX struct {
	Code   *string `json:"code" bson:"code" validate:"required,max=255"`
	Name   *string `json:"name" bson:"name" validate:"required,max=255"`
	IsAuto bool    `json:"isauto" bson:"isauto"`
}
