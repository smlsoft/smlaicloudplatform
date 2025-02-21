package services

import (
	"context"
	"errors"
	"math"
	creditorRepo "smlaicloudplatform/internal/debtaccount/creditor/repositories"
	"smlaicloudplatform/internal/product/product/models"
	"smlaicloudplatform/internal/product/product/repositories"
	barcodeModel "smlaicloudplatform/internal/product/productbarcode/models"
	productBarcodeRepo "smlaicloudplatform/internal/product/productbarcode/repositories"
	"smlaicloudplatform/internal/utils"
	"strings"
	"time"

	"github.com/smlsoft/mongopagination"
)

type IProductHttpService interface {
	GetModuleName() string
	GetProduct(shopID string, code string) (*models.ProductPg, error)
	ProductList(shopID string, name string, page int, pageSize int) ([]models.ProductPg, mongopagination.PaginationData, error)
	Create(doc *models.ProductPg) error
	Update(shopID string, code string, doc *models.ProductPg) error
	Delete(shopID string, code string) error
}

type ProductHttpService struct {
	repo                 repositories.IProductPGRepository
	repomgCreditror      creditorRepo.CreditorRepository
	repomgProductBarcode productBarcodeRepo.ProductBarcodeRepository
	contextTimeout       time.Duration
}

// ✅ **สร้าง Service**
func NewProductHttpService(repo repositories.IProductPGRepository, repomgCreditror creditorRepo.CreditorRepository, repomgProductBarcode productBarcodeRepo.ProductBarcodeRepository) *ProductHttpService {
	return &ProductHttpService{
		repo:                 repo,
		repomgCreditror:      repomgCreditror,
		repomgProductBarcode: repomgProductBarcode,
		contextTimeout:       15 * time.Second,
	}
}

// ✅ **ตั้งค่า Timeout**
func (svc ProductHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc ProductHttpService) GetModuleName() string {
	return "product"
}

// ✅ **GetProduct (ดึงข้อมูล Product)**
func (svc ProductHttpService) GetProduct(shopID string, code string) (*models.ProductPg, error) {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	// ✅ ดึงข้อมูล Product จาก PostgreSQL
	product, err := svc.repo.Get(ctx, shopID, code)
	if err != nil {
		return nil, err
	}

	// ✅ ดึงข้อมูล Manufacturer ถ้ามีค่า `ManufacturerGUID`
	if product.ManufacturerGUID != nil && strings.TrimSpace(*product.ManufacturerGUID) != "" {
		findDoc, err := svc.repomgCreditror.FindByGuid(ctx, shopID, *product.ManufacturerGUID)
		if err == nil { // ไม่คืนค่า error ถ้าไม่เจอข้อมูล
			product.ManufacturerCode = &findDoc.Code
			product.ManufacturerName = *findDoc.Names
		}
	}

	// ✅ ดึงข้อมูล Barcode จาก MongoDB
	barcodes, err := svc.repomgProductBarcode.FindByItemCode(ctx, shopID, product.Code)
	if err != nil || barcodes == nil {
		// ถ้าไม่เจอข้อมูล หรือเกิดข้อผิดพลาด ให้ตั้งค่า barcodes = []
		barcodes = []barcodeModel.ProductBarcodeDoc{}
	}

	tempBarcodes := []models.Barcodes{}
	for _, barcode := range barcodes {
		tempPrices := []models.ProductPrice{}
		if barcode.Prices != nil { // ตรวจสอบก่อน loop
			for _, price := range *barcode.Prices {
				tempPrices = append(tempPrices, models.ProductPrice{
					KeyNumber: price.KeyNumber,
					Price:     price.Price,
				})
			}
		}

		tempBarcodes = append(tempBarcodes, models.Barcodes{
			Barcode:       barcode.Barcode,
			ItemUnitCode:  barcode.ItemUnitCode,
			ItemUnitNames: barcode.ItemUnitNames,
			Prices:        &tempPrices,
			GuidFixed:     barcode.GuidFixed,
			Condition:     barcode.Condition,
			DivideValue:   barcode.DivideValue,
			StandValue:    barcode.StandValue,
			Qty:           barcode.Qty,
			IsMainBarcode: barcode.IsMainBarcode,
		})
	}

	product.Barcodes = tempBarcodes // ✅ กำหนดค่า Barcodes ที่เป็น `[]` ถ้าไม่มีข้อมูล

	return product, nil
}

// ✅ **ProductList (ค้นหา Product ตามชื่อ + Pagination)**
func (svc ProductHttpService) ProductList(shopID string, name string, page int, pageSize int) ([]models.ProductPg, mongopagination.PaginationData, error) {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	// ✅ ดึงข้อมูลจาก `repo.ProductList()`
	products, totalRecords, err := svc.repo.ProductList(ctx, shopID, name, page, pageSize)
	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	// ✅ คำนวณ pagination
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	// ✅ สร้าง `PaginationData`
	pagination := mongopagination.PaginationData{
		Total:     totalRecords,
		Page:      int64(page),
		PerPage:   int64(pageSize),
		Prev:      int64(max(1, page-1)),
		Next:      int64(min(page+1, totalPages)),
		TotalPage: int64(totalPages),
	}

	return products, pagination, nil
}

// ✅ ฟังก์ชันช่วยคำนวณค่า Min/Max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ✅ **Create (สร้าง Product ใหม่)**
func (svc ProductHttpService) Create(doc *models.ProductPg) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if doc.ShopID == "" || doc.Code == "" {
		return errors.New("ShopID and Code are required")
	}

	// ✅ สร้าง `GuidFixed` ถ้ายังไม่มีค่า
	if doc.GuidFixed == "" {
		doc.GuidFixed = utils.NewGUID() // 🔥 สร้าง GUID ใหม่
	}

	// ✅ ตรวจสอบค่าที่เป็น "" และตั้งค่าให้เป็น nil
	if doc.GroupGuid != nil && *doc.GroupGuid == "" {
		doc.GroupGuid = nil
	}
	if doc.UnitGuid != nil && *doc.UnitGuid == "" {
		doc.UnitGuid = nil
	}
	if doc.ManufacturerGUID != nil && *doc.ManufacturerGUID == "" {
		doc.ManufacturerGUID = nil
	}

	// ✅ กำหนดค่าเริ่มต้นให้ `itemtype` หากไม่ได้ส่งมา
	if doc.ItemType == 0 {
		doc.ItemType = 0
	}

	// ✅ ตั้งค่าเวลาก่อนสร้าง
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	// ✅ เรียก `Create()`
	err := svc.repo.Create(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

// ✅ **Update (อัปเดต Product)**
func (svc ProductHttpService) Update(shopID string, code string, doc *models.ProductPg) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if shopID == "" || code == "" {
		return errors.New("ShopID and Code are required")
	}

	// ✅ ตั้งค่า UpdatedAt
	doc.UpdatedAt = time.Now()

	// ✅ เรียก Repository เพื่ออัปเดตข้อมูล
	err := svc.repo.Update(ctx, shopID, code, doc)
	if err != nil {
		return err
	}

	return nil
}

// ✅ **Delete (ลบ Product)**
func (svc ProductHttpService) Delete(shopID string, code string) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if shopID == "" || code == "" {
		return errors.New("ShopID and Code are required")
	}

	err := svc.repo.Delete(ctx, shopID, code)
	if err != nil {
		return err
	}

	return nil
}
