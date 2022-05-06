package restaurant

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const printerTerminalCollectionName = "printerTerminals"

type PrinterTerminal struct {
	Code    string `json:"code" bson:"code"`
	Name1   string `json:"name1" bson:"name1" gorm:"name1"`
	Name2   string `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3   string `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4   string `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5   string `json:"name5,omitempty" bson:"name5,omitempty"`
	Address string `json:"address" bson:"address" `
	Type    int8   `json:"type" bson:"type"`
}

type PrinterTerminalInfo struct {
	models.DocIdentity `bson:"inline"`
	PrinterTerminal    `bson:"inline"`
}

func (PrinterTerminalInfo) CollectionName() string {
	return printerTerminalCollectionName
}

type PrinterTerminalData struct {
	models.ShopIdentity `bson:"inline"`
	PrinterTerminalInfo `bson:"inline"`
}

type PrinterTerminalDoc struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PrinterTerminalData `bson:"inline"`
	models.ActivityDoc  `bson:"inline"`
	models.LastUpdate   `bson:"inline"`
}

func (PrinterTerminalDoc) CollectionName() string {
	return printerTerminalCollectionName
}

//Extra

type PrinterTerminalItemGuid struct {
	Code string `json:"code" bson:"code" gorm:"code"`
}

func (PrinterTerminalItemGuid) CollectionName() string {
	return printerTerminalCollectionName
}

type PrinterTerminalActivity struct {
	PrinterTerminalData `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PrinterTerminalActivity) CollectionName() string {
	return printerTerminalCollectionName
}

type PrinterTerminalDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PrinterTerminalDeleteActivity) CollectionName() string {
	return printerTerminalCollectionName
}

type PrinterTerminalInfoResponse struct {
	Success bool                `json:"success"`
	Data    PrinterTerminalInfo `json:"data,omitempty"`
}

type PrinterTerminalPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []PrinterTerminalInfo         `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}

type PrinterTerminalLastActivityResponse struct {
	New    []PrinterTerminalActivity       `json:"new" `
	Remove []PrinterTerminalDeleteActivity `json:"remove"`
}

type PrinterTerminalFetchUpdateResponse struct {
	Success    bool                                `json:"success"`
	Data       PrinterTerminalLastActivityResponse `json:"data,omitempty"`
	Pagination models.PaginationDataResponse       `json:"pagination,omitempty"`
}
