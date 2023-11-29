package services

import (
	"context"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	productbarcode_models "smlcloudplatform/pkg/product/productbarcode/models"
	productbarcode_repositories "smlcloudplatform/pkg/product/productbarcode/repositories"
	"smlcloudplatform/pkg/services"
	trancache "smlcloudplatform/pkg/transaction/repositories"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"strconv"
	"strings"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISaleInvoiceHttpService interface {
	CreateSaleInvoice(shopID string, authUsername string, doc models.SaleInvoice) (string, string, error)
	UpdateSaleInvoice(shopID string, guid string, authUsername string, doc models.SaleInvoice) error
	DeleteSaleInvoice(shopID string, guid string, authUsername string) error
	DeleteSaleInvoiceByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoSaleInvoice(shopID string, guid string) (models.SaleInvoiceInfo, error)
	InfoSaleInvoiceByCode(shopID string, code string) (models.SaleInvoiceInfo, error)
	SearchSaleInvoice(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceInfo, mongopagination.PaginationData, error)
	SearchSaleInvoiceStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.SaleInvoice) (common.BulkImport, error)
	GetLastPOSDocNo(shopID, posID, maxDocNo string) (string, error)
	Export(languageCode string, shopID string, languageHeader map[string]string) ([][]string, error)

	GetModuleName() string
}

const (
	MODULE_NAME = "SI"
)

type SaleInvoiceHttpService struct {
	repoMq             repositories.ISaleInvoiceMessageQueueRepository
	repo               repositories.ISaleInvoiceRepository
	repoCache          trancache.ICacheRepository
	productbarcodeRepo productbarcode_repositories.IProductBarcodeRepository
	cacheExpireDocNo   time.Duration
	syncCacheRepo      mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity]
	contextTimeout time.Duration
}

func NewSaleInvoiceHttpService(
	repo repositories.ISaleInvoiceRepository,
	productbarcodeRepo productbarcode_repositories.IProductBarcodeRepository,
	repoCache trancache.ICacheRepository,
	repoMq repositories.ISaleInvoiceMessageQueueRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,
) *SaleInvoiceHttpService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &SaleInvoiceHttpService{
		repo:               repo,
		productbarcodeRepo: productbarcodeRepo,
		repoMq:             repoMq,
		repoCache:          repoCache,
		syncCacheRepo:      syncCacheRepo,
		cacheExpireDocNo:   time.Hour * 24,
		contextTimeout:     contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity](repo)

	return insSvc
}

func (svc SaleInvoiceHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc SaleInvoiceHttpService) getDocNoPrefix(docDate time.Time) string {
	docDateStr := docDate.Format("20060102")
	return fmt.Sprintf("%s%s", MODULE_NAME, docDateStr)
}

func (svc SaleInvoiceHttpService) generateNewDocNo(ctx context.Context, shopID, prefixDocNo string, docNumber int) (string, int, error) {
	prevoiusDocNumber, err := svc.repoCache.Get(shopID, prefixDocNo)

	if prevoiusDocNumber == 0 || err != nil {
		lastDoc, err := svc.repo.FindLastDocNo(ctx, shopID, prefixDocNo)

		if err != nil {
			return "", 0, err
		}

		if len(lastDoc.DocNo) > 0 {
			rawNumber := strings.Replace(lastDoc.DocNo, prefixDocNo, "", -1)
			prevoiusDocNumber, err = strconv.Atoi(rawNumber)

			if err != nil {
				prevoiusDocNumber = 0
			}
		}

	}

	newDocNumber := prevoiusDocNumber + 1
	newDocNo := fmt.Sprintf("%s%05d", prefixDocNo, newDocNumber)

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", newDocNo)

	if err != nil {
		return "", 0, err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", 0, errors.New("DocNo is exists")
	}

	return newDocNo, newDocNumber, nil
}

func (svc SaleInvoiceHttpService) CreateSaleInvoice(shopID string, authUsername string, doc models.SaleInvoice) (string, string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	isGenerateDocNo := !doc.IsPOS

	prefixDocNo, docNo, newDocNumber := "", "", 0

	if isGenerateDocNo {
		docDate := doc.DocDatetime
		prefixDocNo = svc.getDocNoPrefix(docDate)

		tempNewDocNo, tempNewDocNumber, err := svc.generateNewDocNo(ctx, shopID, prefixDocNo, 1)

		if err != nil {
			return "", "", err
		}

		docNo = tempNewDocNo
		newDocNumber = tempNewDocNumber
	} else {
		if doc.DocNo == "" {
			return "", "", errors.New("docno is required")
		}

		docNo = doc.DocNo
	}

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", docNo)

	if err != nil {
		return "", "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", "", errors.New("DocNo is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.SaleInvoiceDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SaleInvoice = doc

	docData.DocNo = docNo
	if isGenerateDocNo && doc.TaxDocNo == "" {
		docData.TaxDocNo = docNo
	}

	err = svc.prepareSaleInvoiceDetail(ctx, shopID, docData.SaleInvoice.Details)
	if err != nil {
		return "", "", err
	}

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", "", err
	}

	go func() {
		svc.repoMq.Create(docData)
		svc.saveMasterSync(shopID)

		if isGenerateDocNo {
			svc.repoCache.Save(shopID, prefixDocNo, newDocNumber, svc.cacheExpireDocNo)
		}
	}()

	return newGuidFixed, docNo, nil
}

func (svc SaleInvoiceHttpService) prepareSaleInvoiceDetail(ctx context.Context, shopID string, details *[]models.SaleInvoiceDetail) error {
	var tempBarcodes []string
	for _, doc := range *details {
		tempBarcodes = append(tempBarcodes, doc.Barcode)
	}

	productBarcodes, err := svc.productbarcodeRepo.FindByBarcodes(ctx, shopID, tempBarcodes)
	if err != nil {
		return err
	}

	productBarcodeDict := map[string]productbarcode_models.ProductBarcodeInfo{}
	for _, doc := range productBarcodes {
		productBarcodeDict[doc.Barcode] = doc
	}

	for i := 0; i < len(*details); i++ {
		tempDetail := (*details)[i]
		tempProduct := productBarcodeDict[tempDetail.Barcode]
		if _, ok := productBarcodeDict[tempDetail.Barcode]; ok {
			tempDetail.ItemGuid = tempProduct.GuidFixed
			tempDetail.UnitCode = tempProduct.ItemUnitCode
			tempDetail.UnitNames = tempProduct.ItemUnitNames
			tempDetail.ManufacturerGUID = tempProduct.ManufacturerGUID

			tempDetail.ItemCode = tempProduct.ItemCode
			tempDetail.ItemNames = tempProduct.Names
			tempDetail.ItemType = tempProduct.ItemType
			tempDetail.TaxType = tempProduct.TaxType
			tempDetail.VatType = tempProduct.VatType
			tempDetail.Discount = tempProduct.Discount

			tempDetail.DivideValue = tempProduct.DivideValue
			tempDetail.StandValue = tempProduct.StandValue
			tempDetail.VatCal = tempProduct.VatCal
		}

		(*details)[i] = tempDetail
	}

	return nil
}

func (svc SaleInvoiceHttpService) UpdateSaleInvoice(shopID string, guid string, authUsername string, doc models.SaleInvoice) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	docData := findDoc
	docData.SaleInvoice = doc

	err = svc.prepareSaleInvoiceDetail(ctx, shopID, docData.SaleInvoice.Details)

	if err != nil {
		return err
	}

	docData.DocNo = findDoc.DocNo
	docData.UpdatedBy = authUsername
	docData.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, docData)

	if err != nil {
		return err
	}

	go func() {
		svc.repoMq.Update(docData)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc SaleInvoiceHttpService) DeleteSaleInvoice(shopID string, guid string, authUsername string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(ctx, shopID, guid, authUsername)
	if err != nil {
		return err
	}

	go func() {
		svc.repoMq.Delete(findDoc)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc SaleInvoiceHttpService) DeleteSaleInvoiceByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	func() {
		docs, _ := svc.repo.FindByGuids(ctx, shopID, GUIDs)
		svc.repoMq.DeleteInBatch(docs)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc SaleInvoiceHttpService) GetLastPOSDocNo(shopID, posID, maxDocNo string) (string, error) {
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	lastDocNo, err := svc.repo.FindLastPOSDocNo(ctx, shopID, posID, maxDocNo)

	if err != nil {
		return "", err
	}

	return lastDocNo, nil
}

func (svc SaleInvoiceHttpService) InfoSaleInvoice(shopID string, guid string) (models.SaleInvoiceInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.SaleInvoiceInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SaleInvoiceInfo{}, errors.New("document not found")
	}

	return findDoc.SaleInvoiceInfo, nil
}

func (svc SaleInvoiceHttpService) InfoSaleInvoiceByCode(shopID string, code string) (models.SaleInvoiceInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", code)

	if err != nil {
		return models.SaleInvoiceInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SaleInvoiceInfo{}, errors.New("document not found")
	}

	return findDoc.SaleInvoiceInfo, nil
}

func (svc SaleInvoiceHttpService) SearchSaleInvoice(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.SaleInvoiceInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc SaleInvoiceHttpService) SearchSaleInvoiceStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.SaleInvoiceInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc SaleInvoiceHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.SaleInvoice) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.SaleInvoice](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.DocNo)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(ctx, shopID, "docno", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.DocNo)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.SaleInvoice, models.SaleInvoiceDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.SaleInvoice) models.SaleInvoiceDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.SaleInvoiceDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.SaleInvoice = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.SaleInvoice, models.SaleInvoiceDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.SaleInvoiceDoc, error) {
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", guid)
		},
		func(doc models.SaleInvoiceDoc) bool {
			return doc.DocNo != ""
		},
		func(shopID string, authUsername string, data models.SaleInvoice, doc models.SaleInvoiceDoc) error {

			doc.SaleInvoice = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			err = svc.repo.Update(ctx, shopID, doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(ctx, createDataList)

		if err != nil {
			return common.BulkImport{}, err
		}

	}

	createDataKey := []string{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.DocNo)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.DocNo)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.DocNo)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, svc.getDocIDKey(doc))
	}

	svc.saveMasterSync(shopID)

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc SaleInvoiceHttpService) getDocIDKey(doc models.SaleInvoice) string {
	return doc.DocNo
}

func (svc SaleInvoiceHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc SaleInvoiceHttpService) GetModuleName() string {
	return "saleInvoice"
}

func (svc SaleInvoiceHttpService) Export(shopID string, languageCode string, languageHeader map[string]string) ([][]string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docs, err := svc.repo.Find(ctx, shopID, []string{}, "")

	if err != nil {
		return [][]string{}, err
	}

	keyCols := []string{
		"docdate",        //"วันที่",
		"docno",          //เลขที่เอกสาร",
		"barcode",        //บาร์โค้ด",
		"productname",    //"ชื่อสินค้า",
		"unitcode",       //"หน่วยนับ",
		"unitname",       //"ชื่อหน่วยนับ",
		"qty",            //"จำนวน",
		"price",          //ราคา",
		"discountamount", // "มูลค่าส่วนลด",
		"sumamount",      //"มูลค่าสินค้า",
	}

	headerRow := []string{}
	for _, keyCol := range keyCols {
		tempVal := keyCol
		if val, ok := languageHeader[keyCol]; ok && val != "" {
			tempVal = val
		}
		headerRow = append(headerRow, tempVal)
	}

	results := [][]string{}

	results = append(results, headerRow)

	for _, doc := range docs {
		tempResults := prepareDataToCSV(languageCode, doc)
		results = append(results, tempResults...)
	}

	return results, nil
}

func prepareDataToCSV(languageCode string, data models.SaleInvoiceInfo) [][]string {

	results := [][]string{}

	if data.Details == nil {
		return results
	}

	for _, value := range *data.Details {
		langCode := languageCode

		productName := getName(value.ItemNames, langCode)
		unitName := getName(value.UnitNames, langCode)

		qty := fmt.Sprintf("%.2f", value.Qty)
		price := fmt.Sprintf("%.2f", value.Price)
		discountAmount := fmt.Sprintf("%.2f", value.DiscountAmount)
		sumAmount := fmt.Sprintf("%.2f", value.SumAmount)

		dateLayout := "2006-01-02"
		docDate := data.DocDatetime.Format(dateLayout)

		results = append(results, []string{docDate, data.DocNo, value.Barcode, productName, value.UnitCode, unitName, qty, price, discountAmount, sumAmount})
	}

	return results
}

func getName(names *[]common.NameX, langCode string) string {
	if names == nil {
		return ""
	}

	for _, name := range *names {
		if *name.Code == langCode {
			return *name.Name
		}
	}

	return ""
}
