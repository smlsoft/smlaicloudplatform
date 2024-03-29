package services

import (
	"context"
	"errors"
	"fmt"
	mastersync "smlcloudplatform/internal/mastersync/repositories"
	common "smlcloudplatform/internal/models"
	productbarcode_models "smlcloudplatform/internal/product/productbarcode/models"
	productbarcode_repositories "smlcloudplatform/internal/product/productbarcode/repositories"
	"smlcloudplatform/internal/services"
	trans_models "smlcloudplatform/internal/transaction/models"
	trans_cache "smlcloudplatform/internal/transaction/repositories"
	"smlcloudplatform/internal/transaction/saleinvoice/models"
	"smlcloudplatform/internal/transaction/saleinvoice/repositories"
	"smlcloudplatform/internal/utils"
	"smlcloudplatform/internal/utils/importdata"
	micromodels "smlcloudplatform/pkg/microservice/models"
	"strconv"
	"strings"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISaleInvoiceService interface {
	CreateSaleInvoice(shopID string, authUsername string, doc models.SaleInvoice) (string, string, error)
	UpdateSaleInvoice(shopID string, guid string, authUsername string, doc models.SaleInvoice) error
	UpdateSlip(shopID string, authUsername string, docNo string, mode uint8, imageUrl string) error

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
	TRANS_FLAG  = 44
)

type ISaleInvocieParser interface {
	ParseProductBarcode(detail trans_models.Detail, productBarcodeInfo productbarcode_models.ProductBarcodeInfo) trans_models.Detail
}

type ISaleInvoiceExport interface {
	ParseCSV(languageCode string, data models.SaleInvoiceInfo) [][]string
}

type SaleInvoiceService struct {
	repoMq             repositories.ISaleInvoiceMessageQueueRepository
	repo               repositories.ISaleInvoiceRepository
	repoCache          trans_cache.ICacheRepository
	productbarcodeRepo productbarcode_repositories.IProductBarcodeRepository
	cacheExpireDocNo   time.Duration
	syncCacheRepo      mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity]
	parser         ISaleInvocieParser
	exporter       ISaleInvoiceExport
	contextTimeout time.Duration
}

func NewSaleInvoiceService(
	repo repositories.ISaleInvoiceRepository,
	repoCache trans_cache.ICacheRepository,
	productbarcodeRepo productbarcode_repositories.IProductBarcodeRepository,
	repoMq repositories.ISaleInvoiceMessageQueueRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,
	parser ISaleInvocieParser,
	exporter ISaleInvoiceExport,
) *SaleInvoiceService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &SaleInvoiceService{
		repo:               repo,
		repoMq:             repoMq,
		repoCache:          repoCache,
		productbarcodeRepo: productbarcodeRepo,
		syncCacheRepo:      syncCacheRepo,
		parser:             parser,
		exporter:           exporter,
		cacheExpireDocNo:   time.Hour * 24,
		contextTimeout:     contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity](repo)

	return insSvc
}

func (svc SaleInvoiceService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc SaleInvoiceService) getDocNoPrefix(docDate time.Time) string {
	docDateStr := docDate.Format("20060102")
	return fmt.Sprintf("%s%s", MODULE_NAME, docDateStr)
}

func (svc SaleInvoiceService) generateNewDocNo(ctx context.Context, shopID, prefixDocNo string, docNumber int) (string, int, error) {
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

func (svc SaleInvoiceService) CreateSaleInvoice(shopID string, authUsername string, doc models.SaleInvoice) (string, string, error) {

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

	if isGenerateDocNo {
		findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", docNo)

		if err != nil {
			return "", "", err
		}

		if len(findDoc.GuidFixed) > 0 {
			return "", "", errors.New("DocNo is exists")
		}
	}

	newGuidFixed := utils.NewGUID()

	dataDoc := models.SaleInvoiceDoc{}
	dataDoc.ShopID = shopID
	dataDoc.GuidFixed = newGuidFixed
	dataDoc.SaleInvoice = doc

	dataDoc.TransFlag = TRANS_FLAG
	dataDoc.DocNo = docNo
	if isGenerateDocNo && doc.TaxDocNo == "" {
		dataDoc.TaxDocNo = docNo
	}

	productBarcodes, err := svc.GetDetailProductBarcodes(ctx, shopID, *doc.Details)
	if err != nil {
		return "", "", err
	}

	details := svc.PrepareDetail(*doc.Details, productBarcodes)
	dataDoc.Details = &details

	dataDoc.CreatedBy = authUsername
	dataDoc.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, dataDoc)

	if err != nil {
		return "", "", err
	}

	go func() {
		err := svc.repoMq.Create(dataDoc)
		if err != nil {
			fmt.Printf("create mq error :: %s", err.Error())
		}
		svc.saveMasterSync(shopID)

		if isGenerateDocNo {
			svc.repoCache.Save(shopID, prefixDocNo, newDocNumber, svc.cacheExpireDocNo)
		}
	}()

	return newGuidFixed, docNo, nil
}

func (svc SaleInvoiceService) GetDetailProductBarcodes(ctx context.Context, shopID string, details []trans_models.Detail) ([]productbarcode_models.ProductBarcodeInfo, error) {
	var tempBarcodes []string
	for _, doc := range details {
		tempBarcodes = append(tempBarcodes, doc.Barcode)
	}
	return svc.productbarcodeRepo.FindByBarcodes(ctx, shopID, tempBarcodes)
}

func (svc SaleInvoiceService) PrepareDetail(details []trans_models.Detail, productBarcodes []productbarcode_models.ProductBarcodeInfo) []trans_models.Detail {

	productBarcodeDict := map[string]productbarcode_models.ProductBarcodeInfo{}
	for _, doc := range productBarcodes {
		productBarcodeDict[doc.Barcode] = doc
	}

	for i := 0; i < len(details); i++ {
		tempDetail := (details)[i]
		tempProduct := productBarcodeDict[tempDetail.Barcode]
		if _, ok := productBarcodeDict[tempDetail.Barcode]; ok {
			tempDetail = svc.parser.ParseProductBarcode(tempDetail, tempProduct)
		}

		(details)[i] = tempDetail
	}

	return details
}

func (svc SaleInvoiceService) UpdateSaleInvoice(shopID string, guid string, authUsername string, doc models.SaleInvoice) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	dataDoc := findDoc
	dataDoc.SaleInvoice = doc

	productBarcodes, err := svc.GetDetailProductBarcodes(ctx, shopID, *doc.Details)
	if err != nil {
		return err
	}

	details := svc.PrepareDetail(*doc.Details, productBarcodes)
	dataDoc.Details = &details

	dataDoc.SlipQrUrl = findDoc.SlipQrUrl
	dataDoc.SlipQrUrlHistories = findDoc.SlipQrUrlHistories
	dataDoc.SlipUrl = findDoc.SlipUrl
	dataDoc.SlipUrlHistories = findDoc.SlipUrlHistories

	dataDoc.DocNo = findDoc.DocNo
	dataDoc.TransFlag = TRANS_FLAG
	dataDoc.UpdatedBy = authUsername
	dataDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, dataDoc)

	if err != nil {
		return err
	}

	go func() {
		err := svc.repoMq.Update(dataDoc)
		if err != nil {
			fmt.Printf("create mq error :: %s", err.Error())
		}
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc SaleInvoiceService) UpdateSlip(shopID string, authUsername string, docNo string, mode uint8, imageUrl string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", docNo)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	dataDoc := findDoc

	if mode == 1 {
		if findDoc.SlipQrUrl != "" {
			dataDoc.SlipQrUrlHistories = append(dataDoc.SlipQrUrlHistories, findDoc.SlipQrUrl)
		}

		dataDoc.SlipQrUrl = imageUrl
	} else {
		if findDoc.SlipUrl != "" {
			dataDoc.SlipUrlHistories = append(dataDoc.SlipUrlHistories, findDoc.SlipUrl)
		}

		dataDoc.SlipUrl = imageUrl
	}

	dataDoc.UpdatedBy = authUsername
	dataDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, findDoc.GuidFixed, dataDoc)

	if err != nil {
		return err
	}

	go func() {
		err := svc.repoMq.Update(dataDoc)
		if err != nil {
			fmt.Printf("create mq error :: %s", err.Error())
		}
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc SaleInvoiceService) DeleteSaleInvoice(shopID string, guid string, authUsername string) error {

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

func (svc SaleInvoiceService) DeleteSaleInvoiceByGUIDs(shopID string, authUsername string, GUIDs []string) error {

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

func (svc SaleInvoiceService) GetLastPOSDocNo(shopID, posID, maxDocNo string) (string, error) {
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	lastDocNo, err := svc.repo.FindLastPOSDocNo(ctx, shopID, posID, maxDocNo)

	if err != nil {
		return "", err
	}

	return lastDocNo, nil
}

func (svc SaleInvoiceService) InfoSaleInvoice(shopID string, guid string) (models.SaleInvoiceInfo, error) {

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

func (svc SaleInvoiceService) InfoSaleInvoiceByCode(shopID string, code string) (models.SaleInvoiceInfo, error) {

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

func (svc SaleInvoiceService) SearchSaleInvoice(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceInfo, mongopagination.PaginationData, error) {

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

func (svc SaleInvoiceService) SearchSaleInvoiceStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceInfo, int, error) {

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

func (svc SaleInvoiceService) SaveInBatch(shopID string, authUsername string, dataList []models.SaleInvoice) (common.BulkImport, error) {

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

			dataDoc.TransFlag = TRANS_FLAG
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
			doc.TransFlag = TRANS_FLAG
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

func (svc SaleInvoiceService) getDocIDKey(doc models.SaleInvoice) string {
	return doc.DocNo
}

func (svc SaleInvoiceService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc SaleInvoiceService) GetModuleName() string {
	return "saleInvoice"
}

func (svc SaleInvoiceService) Export(shopID string, languageCode string, languageHeader map[string]string) ([][]string, error) {

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
		tempResults := svc.exporter.ParseCSV(languageCode, doc)
		results = append(results, tempResults...)
	}

	return results, nil
}
