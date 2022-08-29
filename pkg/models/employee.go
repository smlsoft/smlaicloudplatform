package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const employeeCollectionName string = "employees"

type Employee struct {
	Username       string `json:"username" bson:"username"`
	Password       string `json:"password" bson:"password"`
	Name           string `json:"name" bson:"name"`
	ProfilePicture string `json:"profilepicture" bson:"profilepicture"`
	Role           string `json:"role" bson:"role"`
}

type EmployeeInfo struct {
	DocIdentity `bson:"inline" gorm:"embedded;"`
	Employee    `bson:"inline" gorm:"embedded;"`
}

func (EmployeeInfo) CollectionName() string {
	return employeeCollectionName
}

type EmployeeData struct {
	ShopIdentity `bson:"inline" gorm:"embedded;"`
	EmployeeInfo `bson:"inline" gorm:"embedded;"`
}

type EmployeeDoc struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EmployeeData `bson:"inline"`
	ActivityDoc  `bson:"inline"`
}

func (EmployeeDoc) CollectionName() string {
	return employeeCollectionName
}

type EmployeeRequestLogin struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type EmployeeRequestUpdate struct {
	Username string  `json:"username" bson:"username"`
	Name     string  `json:"name" bson:"name"`
	Role     *string `json:"role" bson:"role"`
}

type EmployeeRequestPassword struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type EmployeePageResponse struct {
	Success    bool                   `json:"success"`
	Data       []EmployeeInfo         `json:"data,omitempty"`
	Pagination PaginationDataResponse `json:"pagination,omitempty"`
}
