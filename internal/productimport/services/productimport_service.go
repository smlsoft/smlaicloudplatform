package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	common "smlcloudplatform/internal/models"
	product_models "smlcloudplatform/internal/product/productbarcode/models"
	productbarcode_repo "smlcloudplatform/internal/product/productbarcode/repositories"
	product_services "smlcloudplatform/internal/product/productbarcode/services"
	productunit_repo "smlcloudplatform/internal/product/unit/repositories"
	"smlcloudplatform/internal/productimport/models"
	"smlcloudplatform/internal/productimport/repositories"
	micromodels "smlcloudplatform/pkg/microservice/models"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type IProductImportService interface {
	List(shopID string, taskID string, pageable micromodels.Pageable) ([]models.ProductImportInfo, models.PaginationData, error)
	Create(shopID string, authUsername string, req *models.ProductImport) error
	Update(shopID string, guid string, doc models.ProductImportRaw) error
	Delete(shopID string, guid string) error
	DeleteTask(shopID string, taskID string) error
	ImportFromFile(shopID string, authUsername string, fileUpload io.Reader) (string, error)
	SaveTask(shopID string, authUsername string, taskID string, docHeader models.ProductImportHeader) error
	Verify(shopID string, taskID string) error
}

type ProductImportService struct {
	deafultPartSize       int
	sizeID                int
	cacheExpire           time.Duration
	chRepo                repositories.IProductImportClickHouseRepository
	productBarcodeRepo    productbarcode_repo.IProductBarcodeRepository
	productBarcodeService product_services.IProductBarcodeHttpService
	productUnitRepo       productunit_repo.IUnitRepository
	generateID            func(int) string
	generateGUID          func() string
	timeNow               func() time.Time
}

func NewProductImportService(
	chRepo repositories.IProductImportClickHouseRepository,
	productBarcodeRepo productbarcode_repo.IProductBarcodeRepository,
	productBarcodeService product_services.IProductBarcodeHttpService,
	productUnitRepo productunit_repo.IUnitRepository,
	generateID func(int) string,
	generateGUID func() string,
	timeNow func() time.Time,
) *ProductImportService {
	return &ProductImportService{
		deafultPartSize:       100,
		sizeID:                12,
		cacheExpire:           time.Minute * 60,
		chRepo:                chRepo,
		productBarcodeRepo:    productBarcodeRepo,
		productBarcodeService: productBarcodeService,
		productUnitRepo:       productUnitRepo,
		generateID:            generateID,
		generateGUID:          generateGUID,
		timeNow:               timeNow,
	}
}

func (svc *ProductImportService) List(shopID string, taskID string, pageable micromodels.Pageable) ([]models.ProductImportInfo, models.PaginationData, error) {
	findDocs, pagination, err := svc.chRepo.List(context.Background(), shopID, taskID, pageable)

	if err != nil {
		return []models.ProductImportInfo{}, models.PaginationData{}, err
	}

	results := []models.ProductImportInfo{}

	for _, doc := range findDocs {
		results = append(results, doc.ProductImportInfo)
	}

	return results, pagination, nil

}

func (svc *ProductImportService) Create(shopID string, authUsername string, doc *models.ProductImport) error {
	docData := models.ProductImportDoc{}
	docData.ShopID = shopID
	docData.GUIDFixed = svc.generateGUID()
	docData.ProductImport = *doc

	result, err := svc.chRepo.FindOne(context.Background(), shopID, doc.TaskID, []micromodels.KeyInt{
		{
			Key:   "rownumber",
			Value: -1,
		},
	})

	if err != nil {
		return err
	}

	if doc.RowNumber == 0 {
		docData.RowNumber = result.RowNumber + 1
	}

	docData.CreatedAt = svc.timeNow()
	docData.CreatedBy = authUsername

	return svc.chRepo.Create(context.Background(), docData)
}

func (svc *ProductImportService) ImportFromFile(shopID string, authUsername string, fileUpload io.Reader) (string, error) {

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
		"Price",
		"Price Member",
	}

	isNotfoundColumn := false
	columnNotFound := []string{}
	for _, colName := range expectColName {
		if _, ok := colIdxs[colName]; !ok {
			isNotfoundColumn = true
			columnNotFound = append(columnNotFound, colName)
		}
	}

	if isNotfoundColumn {
		return "", fmt.Errorf("column not found: %v", strings.Join(columnNotFound, ", "))
	}

	prepareDataDoc := []models.ProductImportDoc{}

	taskID := svc.generateID(svc.sizeID)

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

	createdAt := svc.timeNow()
	createdBy := authUsername

	for i := range prepareDataDoc {
		prepareDataDoc[i].CreatedAt = createdAt
		prepareDataDoc[i].CreatedBy = createdBy
	}

	err = svc.chRepo.CreateInBatch(context.Background(), prepareDataDoc)

	if err != nil {
		return "", err
	}

	return taskID, nil
}
func (svc *ProductImportService) prepareData(shopID string, taskID string, rowNumber float64, colIdx map[string]int, doc []string) (models.ProductImportDoc, error) {

	price, err := strconv.ParseFloat(doc[colIdx["Price"]], 64)
	if err != nil {
		return models.ProductImportDoc{}, fmt.Errorf("price in row %d invalid", int(rowNumber))
	}

	priceMember, err := strconv.ParseFloat(doc[colIdx["Price Member"]], 64)
	if err != nil {
		return models.ProductImportDoc{}, fmt.Errorf("price member in row %d invalid", int(rowNumber))
	}

	newGUID := svc.generateGUID()

	dataDoc := models.ProductImportDoc{}

	dataDoc.GUIDFixed = newGUID
	dataDoc.ShopID = shopID
	dataDoc.TaskID = taskID
	dataDoc.RowNumber = rowNumber
	dataDoc.Barcode = doc[colIdx["Barcode"]]
	dataDoc.Name = doc[colIdx["Name"]]
	dataDoc.UnitCode = doc[colIdx["Unit Code"]]

	dataDoc.Price = price
	dataDoc.PriceMember = priceMember

	return dataDoc, nil
}

func (svc *ProductImportService) Update(shopID string, guid string, doc models.ProductImportRaw) error {
	return svc.chRepo.Update(context.Background(), shopID, guid, doc)
}

func (svc *ProductImportService) Delete(shopID string, guid string) error {
	return svc.chRepo.DeleteByGUID(context.Background(), shopID, guid)
}

func (svc *ProductImportService) DeleteTask(shopID string, taskID string) error {
	return svc.chRepo.DeleteByTaskID(context.Background(), shopID, taskID)
}

func (svc ProductImportService) Verify(shopID string, taskID string) error {
	docs, err := svc.chRepo.All(context.Background(), shopID, taskID)

	if err != nil {
		return err
	}

	previousDuplicate := map[string]struct{}{}
	previousExist := map[string]struct{}{}
	previousUnitNotExist := map[string]struct{}{}

	itemDulpicated := map[string]struct{}{}
	itemExist := map[string]struct{}{}
	itemUnitNotExist := map[string]struct{}{}

	tempBarcodes := []string{}
	tempUnitCodes := []string{}

	itemDict := map[string]struct{}{}
	for i, doc := range docs {

		if doc.IsExist {
			previousExist[doc.Barcode] = struct{}{}
		}

		if _, ok := itemDict[doc.Barcode]; ok {
			itemDulpicated[doc.Barcode] = struct{}{}
		} else {
			itemDict[doc.Barcode] = struct{}{}
			tempBarcodes = append(tempBarcodes, doc.Barcode)

			if doc.IsDuplicate {
				previousDuplicate[doc.Barcode] = struct{}{}
			}
		}

		if doc.IsUnitNotExist {
			previousUnitNotExist[doc.UnitCode] = struct{}{}
		}

		tempUnitCodes = append(tempUnitCodes, doc.UnitCode)

		if (i > 1 && i%5000 == 0) || i == len(docs)-1 {
			productList, err := svc.productBarcodeRepo.FindByBarcodes(context.Background(), shopID, tempBarcodes)
			if err != nil {
				return err
			}

			for _, product := range productList {
				itemExist[product.Barcode] = struct{}{}
				delete(itemDulpicated, product.Barcode)
			}

			unitList, err := svc.productUnitRepo.FindByUnitCodes(context.Background(), shopID, tempUnitCodes)
			if err != nil {
				return err
			}

			tempDictUnitCodes := map[string]struct{}{}
			for _, unit := range unitList {
				tempDictUnitCodes[unit.UnitCode] = struct{}{}
			}

			for _, unitCode := range tempUnitCodes {
				if _, ok := tempDictUnitCodes[unitCode]; !ok {
					itemUnitNotExist[unitCode] = struct{}{}
				}
			}

			//Clear previous duplicate
			if err := svc.updateDuplicate(shopID, taskID, false, previousDuplicate); err != nil {
				return err
			}

			//Update duplicate
			if err := svc.updateDuplicate(shopID, taskID, true, itemDulpicated); err != nil {
				return err
			}

			//Clear previous exist
			if err := svc.updateExist(shopID, taskID, false, previousExist); err != nil {
				return err
			}

			// Update exist
			if err := svc.updateExist(shopID, taskID, true, itemExist); err != nil {
				return err
			}

			// clear previous unit not exist
			if err := svc.updateUnitExist(shopID, taskID, false, previousUnitNotExist); err != nil {
				return err
			}

			// Update unit not exist
			if err := svc.updateUnitExist(shopID, taskID, true, itemUnitNotExist); err != nil {
				return err
			}

			previousDuplicate = map[string]struct{}{}
			previousExist = map[string]struct{}{}
			itemDulpicated = map[string]struct{}{}
			itemExist = map[string]struct{}{}
			tempBarcodes = []string{}
		}
	}

	return nil

}

func (svc ProductImportService) updateDuplicate(shopID string, taskID string, isDuplicate bool, barcodes map[string]struct{}) error {

	if len(barcodes) == 0 {
		return nil
	}

	tempPreviousDuplicate := []string{}
	for barcode := range barcodes {
		tempPreviousDuplicate = append(tempPreviousDuplicate, barcode)
	}
	err := svc.chRepo.UpdateDuplicate(context.Background(), shopID, taskID, isDuplicate, tempPreviousDuplicate)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductImportService) updateExist(shopID string, taskID string, isExist bool, barcodes map[string]struct{}) error {
	if len(barcodes) == 0 {
		return nil
	}

	tempPreviousExist := []string{}
	for barcode := range barcodes {
		tempPreviousExist = append(tempPreviousExist, barcode)
	}
	err := svc.chRepo.UpdateExist(context.Background(), shopID, taskID, isExist, tempPreviousExist)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductImportService) updateUnitExist(shopID string, taskID string, isExist bool, unitCodes map[string]struct{}) error {
	if len(unitCodes) == 0 {
		return nil
	}

	tempPreviousExist := []string{}
	for barcode := range unitCodes {
		tempPreviousExist = append(tempPreviousExist, barcode)
	}
	err := svc.chRepo.UpdateUnitExist(context.Background(), shopID, taskID, isExist, tempPreviousExist)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductImportService) SaveTask(shopID string, authUsername string, taskID string, docHeader models.ProductImportHeader) error {

	err := svc.Verify(shopID, taskID)

	if err != nil {
		return err
	}

	countDuplicate, err := svc.chRepo.CountDuplicate(context.Background(), shopID, taskID, true)

	if err != nil {
		return errors.New("counting duplicate failed")
	}

	if countDuplicate > 0 {
		return errors.New("items barcode duplicate ")
	}

	countExist, err := svc.chRepo.CountExist(context.Background(), shopID, taskID, true)

	if err != nil {
		return errors.New("counting exist failed")
	}

	if countExist > 0 {
		return errors.New("items barcode exist")
	}

	countUnitNotExist, err := svc.chRepo.CountUnitExist(context.Background(), shopID, taskID, true)

	if err != nil {
		return errors.New("counting unit not exist failed")
	}

	if countUnitNotExist > 0 {
		return errors.New("items unit not exist")
	}

	docs, err := svc.chRepo.All(context.Background(), shopID, taskID)

	if err != nil {
		return err
	}

	dataDocs := svc.PrepareProductBarcodes(docHeader.LanguangeCode, docs)

	_, err = svc.productBarcodeService.SaveInBatch(shopID, authUsername, dataDocs)

	if err != nil {
		return err
	}

	err = svc.DeleteTask(shopID, taskID)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductImportService) PrepareProductBarcodes(langCode string, docs []models.ProductImportDoc) []product_models.ProductBarcode {

	dataDocs := []product_models.ProductBarcode{}

	for i := range docs {

		temp := product_models.ProductBarcode{}

		temp.Barcode = docs[i].Barcode
		temp.ItemUnitCode = docs[i].UnitCode

		productPrices := []product_models.ProductPrice{}

		productPrices = append(productPrices, product_models.ProductPrice{
			KeyNumber: 0,
			Price:     docs[i].Price,
		})

		productPrices = append(productPrices, product_models.ProductPrice{
			KeyNumber: 1,
			Price:     docs[i].PriceMember,
		})
		temp.Prices = &productPrices

		productNames := []common.NameX{}

		productNames = append(productNames, common.NameX{
			Code: &langCode,
			Name: &docs[i].Name,
		})

		temp.Names = &productNames

		dataDocs = append(dataDocs, temp)
	}

	return dataDocs
}
