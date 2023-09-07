package models

import stockbalance_models "smlcloudplatform/pkg/transaction/stockbalance/models"

type StockBalanceImportTaskRequest struct {
	TotalItem int `json:"totalitem"`
}

type StockBalanceImportHeaderRequest struct {
	Header stockbalance_models.StockBalanceHeader `json:"header"`
}

type StockBalanceImportTask struct {
	TaskID    string `json:"taskid"`
	ChunkSize int    `json:"chunksize"`
	TotalItem int    `json:"totalitem"`
	// Header    stockbalance_models.StockBalanceHeader `json:"header"`
	Parts []StockBalanceImportPart `json:"parts"`
}

type StockBalanceImportPart struct {
	PartID     string `json:"partid"`
	PartNumber int    `json:"partnumber"`
}

type StockBalanceImportPartRequest struct {
	TaskID string `json:"taskid"`
	StockBalanceImportPart
}

type StockBalanceImportPartMeta struct {
	PartID     string     `json:"partid"`
	PartNumber int        `json:"partnumber"`
	Status     PartStatus `json:"status"`
}

type StockBalanceImportMeta struct {
	TaskID    string     `json:"taskid"`
	TotalItem int        `json:"totalitem"`
	Status    TaskStatus `json:"status"`
	// Header    stockbalance_models.StockBalanceHeader `json:"header"`
	Parts []StockBalanceImportPartMeta `json:"parts"`
}

type StockBalanceImportPartCache struct {
	TaskID string `json:"taskid"`
	StockBalanceImportPartMeta
	Detail []stockbalance_models.StockBalanceDetail `json:"body"`
}

type StockBalanceImportPartResponse struct {
	PartID     string                                   `json:"partid"`
	PartNumber int                                      `json:"partnumber"`
	Status     string                                   `json:"status"`
	Detail     []stockbalance_models.StockBalanceDetail `json:"body"`
}

type TaskStatus int8

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusProcessing
	TaskStatusDone
	TaskStatusError
	TaskStatusSaveSucceded
	TaskStatusSaveFailed
	TaskStatusNotFound
)
