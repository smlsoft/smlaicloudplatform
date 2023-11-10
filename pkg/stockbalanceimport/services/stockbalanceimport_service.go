package services

import (
	"context"
	"errors"
	"fmt"
	"io"

	micromodels "smlcloudplatform/internal/microservice/models"
	productbarcode_repo "smlcloudplatform/pkg/product/productbarcode/repositories"
	"smlcloudplatform/pkg/stockbalanceimport/models"
	"smlcloudplatform/pkg/stockbalanceimport/repositories"
	stockbalance_models "smlcloudplatform/pkg/transaction/stockbalance/models"
	stockbalance_services "smlcloudplatform/pkg/transaction/stockbalance/services"
	stockbalancedetail_models "smlcloudplatform/pkg/transaction/stockbalancedetail/models"
	stockbalancedetail_services "smlcloudplatform/pkg/transaction/stockbalancedetail/services"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type IStockBalanceImportService interface {
	List(shopID string, taskID string, pageable micromodels.Pageable) ([]models.StockBalanceImportDoc, models.PaginationData, error)
	Create(shopID string, req *models.StockBalanceImport) error
	Update(shopID string, guid string, doc models.StockBalanceImportRaw) error
	Delete(shopID string, guid string) error
	DeleteTask(shopID string, taskID string) error
	ImportFromFile(shopID string, fileUpload io.Reader) (string, error)
	SaveTask(shopID string, authUsername string, taskID string, headerDoc stockbalance_models.StockBalanceHeader) (string, error)
}

type StockBalanceImportService struct {
	deafultPartSize           int
	sizeID                    int
	cacheExpire               time.Duration
	chRepo                    repositories.IStockBalanceImportClickHouseRepository
	productBarcodeRepo        productbarcode_repo.IProductBarcodeRepository
	stockBalanceService       stockbalance_services.IStockBalanceHttpService
	stockBalanceDetailService stockbalancedetail_services.IStockBalanceDetailHttpService
	GenerateID                func(int) string
	GenerateGUID              func() string
}

func NewStockBalanceImportService(
	chRepo repositories.IStockBalanceImportClickHouseRepository,
	productBarcodeRepo productbarcode_repo.IProductBarcodeRepository,
	stockBalanceService stockbalance_services.IStockBalanceHttpService,
	stockBalanceDetailService stockbalancedetail_services.IStockBalanceDetailHttpService,
	GenerateID func(int) string,
	GenerateGUID func() string,
) *StockBalanceImportService {
	return &StockBalanceImportService{
		deafultPartSize:           100,
		sizeID:                    12,
		cacheExpire:               time.Minute * 60,
		chRepo:                    chRepo,
		productBarcodeRepo:        productBarcodeRepo,
		stockBalanceService:       stockBalanceService,
		stockBalanceDetailService: stockBalanceDetailService,
		GenerateID:                GenerateID,
		GenerateGUID:              GenerateGUID,
	}
}

func (svc *StockBalanceImportService) List(shopID string, taskID string, pageable micromodels.Pageable) ([]models.StockBalanceImportDoc, models.PaginationData, error) {
	return svc.chRepo.List(context.Background(), shopID, taskID, pageable)
}

func (svc *StockBalanceImportService) Create(shopID string, doc *models.StockBalanceImport) error {
	docData := models.StockBalanceImportDoc{}
	docData.ShopID = shopID
	docData.GUIDFixed = svc.GenerateGUID()
	docData.StockBalanceImport = *doc
	return svc.chRepo.Create(context.Background(), docData)
}

func (svc *StockBalanceImportService) ImportFromFile(shopID string, fileUpload io.Reader) (string, error) {

	f, err := excelize.OpenReader(fileUpload)
	if err != nil {
		return "", err
	}

	if len(f.GetSheetList()) == 0 {
		return "", errors.New("sheet not found")
	}

	sheetName := f.GetSheetList()[0]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", err
	}

	if len(rows) <= 1 {
		return "", errors.New("sheet is empty")
	}

	colIdxs := map[string]int{}

	cols := rows[0]
	for i, col := range cols {
		colIdxs[col] = i
	}

	expectColName := []string{
		"Barcode",
		"Name",
		"Unit Code",
		"Warehouse Code",
		"Shelf Code",
		"Qty",
		"Amount",
	}

	for _, colName := range expectColName {
		if _, ok := colIdxs[colName]; !ok {
			return "", fmt.Errorf("column %s not found", colName)
		}
	}

	prepareDataDoc := []models.StockBalanceImportDoc{}

	taskID := svc.GenerateID(svc.sizeID)

	for i, doc := range rows {

		if i == 0 {
			continue
		}

		tempData, err := svc.prepareData(shopID, taskID, float64(i), colIdxs, doc)
		if err != nil {
			return "", err
		}

		prepareDataDoc = append(prepareDataDoc, tempData)
	}

	err = svc.chRepo.CreateInBatch(context.Background(), prepareDataDoc)

	if err != nil {
		return "", err
	}

	return taskID, nil
}
func (svc *StockBalanceImportService) prepareData(shopID string, taskID string, rowNumber float64, colIdx map[string]int, doc []string) (models.StockBalanceImportDoc, error) {

	qty, err := strconv.ParseFloat(doc[colIdx["Qty"]], 64)

	if err != nil {
		return models.StockBalanceImportDoc{}, fmt.Errorf("qty in row %d invalid", int(rowNumber))
	}

	amount, err := strconv.ParseFloat(doc[colIdx["Amount"]], 64)

	if err != nil {
		return models.StockBalanceImportDoc{}, fmt.Errorf("amount in row %d invalid", int(rowNumber))
	}

	price := float64(0)

	if qty > 0 && amount > 0 {
		price = amount / qty
	}

	newGUID := svc.GenerateGUID()

	dataDoc := models.StockBalanceImportDoc{}

	dataDoc.GUIDFixed = newGUID
	dataDoc.ShopID = shopID
	dataDoc.TaskID = taskID
	dataDoc.RowNumber = rowNumber
	dataDoc.Barcode = doc[colIdx["Barcode"]]
	dataDoc.Name = doc[colIdx["Name"]]
	dataDoc.UnitCode = doc[colIdx["Unit Code"]]
	dataDoc.WarehouseCode = doc[colIdx["Warehouse Code"]]
	dataDoc.ShelfCode = doc[colIdx["Shelf Code"]]
	dataDoc.Qty = qty
	dataDoc.Price = price
	dataDoc.SumAmount = amount

	return dataDoc, nil
}

func (svc *StockBalanceImportService) Update(shopID string, guid string, doc models.StockBalanceImportRaw) error {
	return svc.chRepo.Update(context.Background(), shopID, guid, doc)
}

func (svc *StockBalanceImportService) Delete(shopID string, guid string) error {
	return svc.chRepo.DeleteByGUID(context.Background(), shopID, guid)
}

func (svc *StockBalanceImportService) DeleteTask(shopID string, taskID string) error {
	return svc.chRepo.DeleteByTaskID(context.Background(), shopID, taskID)
}

func (svc *StockBalanceImportService) SaveTask(shopID string, authUsername string, taskID string, headerDoc stockbalance_models.StockBalanceHeader) (string, error) {

	docs, err := svc.chRepo.All(context.Background(), shopID, taskID)

	if err != nil {
		return "", err
	}

	tempDetails := []stockbalancedetail_models.StockBalanceDetail{}

	barcodes := []string{}
	tempBarcodes := map[string]models.StockBalanceImportDoc{}
	for i, doc := range docs {
		barcodes = append(barcodes, doc.Barcode)
		tempBarcodes[docs[i].Barcode] = docs[i]
		if (i > 1 && i%5000 == 0) || i == len(docs)-1 {
			productList, err := svc.productBarcodeRepo.FindByBarcodes(context.Background(), shopID, barcodes)
			if err != nil {
				return "", err
			}

			for _, product := range productList {
				stockbalanceDetail := stockbalancedetail_models.StockBalanceDetail{}
				temp := tempBarcodes[product.Barcode]

				stockbalanceDetail.ItemCode = product.ItemCode
				stockbalanceDetail.Barcode = product.Barcode
				stockbalanceDetail.ItemNames = product.Names
				stockbalanceDetail.ItemType = product.ItemType
				stockbalanceDetail.TaxType = product.TaxType
				stockbalanceDetail.VatType = product.VatType
				stockbalanceDetail.DivideValue = product.DivideValue
				stockbalanceDetail.StandValue = product.StandValue
				stockbalanceDetail.VatCal = product.VatCal
				stockbalanceDetail.UnitCode = product.ItemUnitCode
				stockbalanceDetail.UnitNames = product.ItemUnitNames

				stockbalanceDetail.Qty = temp.Qty
				stockbalanceDetail.Price = temp.Price
				stockbalanceDetail.SumAmount = temp.SumAmount

				tempDetails = append(tempDetails, stockbalanceDetail)
			}

			barcodes = []string{}
			tempBarcodes = map[string]models.StockBalanceImportDoc{}
		}
	}

	tempTransaction := stockbalance_models.StockBalance{}

	tempTransaction.StockBalanceHeader = headerDoc

	_, docNo, err := svc.stockBalanceService.CreateStockBalance(shopID, authUsername, tempTransaction)
	if err != nil {
		return "", err
	}

	for i := range tempDetails {
		tempDetails[i].DocNo = docNo
	}

	err = svc.stockBalanceDetailService.CreateStockBalanceDetail(shopID, authUsername, tempDetails)
	if err != nil {
		return "", err
	}

	err = svc.DeleteTask(shopID, taskID)

	if err != nil {
		return "", err
	}

	return docNo, nil
}
