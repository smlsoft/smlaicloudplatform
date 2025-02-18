package services

import (
	"context"
	"errors"
	"math"
	"smlaicloudplatform/internal/product/productgroup/models"
	"smlaicloudplatform/internal/product/productgroup/repositories"

	"time"

	"github.com/smlsoft/mongopagination"
	"gorm.io/gorm"
)

type IProductGroupHttpService interface {
	GetModuleName() string
	Get(shopID string, unitCode string) (*models.ProductGroupPg, error)
	ProductGroupList(shopID string, name string, page int, pageSize int) ([]models.ProductGroupPg, mongopagination.PaginationData, error)
	Create(doc *models.ProductGroupPg) error
	Update(shopID string, unitCode string, doc *models.ProductGroupPg) error
	Delete(shopID string, unitCode string) error
}

type ProductGroupHttpService struct {
	repo           repositories.IProductGroupPGRepository
	contextTimeout time.Duration
}

func NewProductGroupHttpService(repo repositories.IProductGroupPGRepository) *ProductGroupHttpService {
	return &ProductGroupHttpService{
		repo:           repo,
		contextTimeout: 15 * time.Second,
	}
}

func (svc ProductGroupHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc ProductGroupHttpService) GetModuleName() string {
	return "productgroup"
}

// ✅ **GetProductGroup (ดึง ProductGroup ตาม shopID และ unitCode)**
func (svc ProductGroupHttpService) Get(shopID string, unitCode string) (*models.ProductGroupPg, error) {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()
	unit, err := svc.repo.Get(ctx, shopID, unitCode)
	if err != nil {
		return nil, err
	}

	return unit, nil
}

// ✅ **ProductGroupList (ค้นหา ProductGroup ตามชื่อ + Pagination)**
func (svc ProductGroupHttpService) ProductGroupList(shopID string, name string, page int, pageSize int) ([]models.ProductGroupPg, mongopagination.PaginationData, error) {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	// ✅ ดึงข้อมูลจาก `repo.ProductGroupList()`
	units, totalRecords, err := svc.repo.ProductGroupList(ctx, shopID, name, page, pageSize)
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

	return units, pagination, nil
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

// ✅ **Create (สร้าง ProductGroup ใหม่)**
func (svc ProductGroupHttpService) Create(doc *models.ProductGroupPg) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if doc.ShopID == "" || doc.Code == "" {
		return errors.New("ShopID and ProductGroupCode are required")
	}

	// ✅ ตรวจสอบว่ามีอยู่แล้วหรือไม่
	existingProductGroup, err := svc.repo.FindByProductGroupCode(ctx, doc.ShopID, doc.Code)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err // ❌ ถ้า error ไม่ใช่ "ไม่พบข้อมูล" ให้คืนค่า error ทันที
	}

	if existingProductGroup != nil {
		return errors.New("ProductGroup already exists") // ❌ ถ้ามี unit อยู่แล้ว ไม่ต้องสร้างซ้ำ
	}

	// ✅ ตั้งค่าเวลาก่อนสร้าง
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	// ✅ เรียก `Create()`
	err = svc.repo.Create(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

// ✅ **Update (อัปเดต ProductGroup)**
func (svc ProductGroupHttpService) Update(shopID string, unitCode string, doc *models.ProductGroupPg) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if shopID == "" || unitCode == "" {
		return errors.New("ShopID and ProductGroupCode are required")
	}

	doc.UpdatedAt = time.Now()

	err := svc.repo.Update(ctx, shopID, unitCode, doc)
	if err != nil {
		return err
	}

	return nil
}

// ✅ **Delete (ลบ ProductGroup)**
func (svc ProductGroupHttpService) Delete(shopID string, unitCode string) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if shopID == "" || unitCode == "" {
		return errors.New("ShopID and ProductGroupCode are required")
	}

	err := svc.repo.Delete(ctx, shopID, unitCode)
	if err != nil {
		return err
	}

	return nil
}
