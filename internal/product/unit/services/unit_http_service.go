package services

import (
	"context"
	"errors"
	"math"
	"smlaicloudplatform/internal/product/unit/models"
	"smlaicloudplatform/internal/product/unit/repositories"
	"time"

	"github.com/smlsoft/mongopagination"
	"gorm.io/gorm"
)

type IUnitHttpService interface {
	GetModuleName() string
	GetUnit(shopID string, unitCode string) (*models.UnitPg, error)
	UnitList(shopID string, name string, page int, pageSize int) ([]models.UnitPg, mongopagination.PaginationData, error)
	Create(doc *models.UnitPg) error
	Update(shopID string, unitCode string, doc *models.UnitPg) error
	Delete(shopID string, unitCode string) error
}

type UnitHttpService struct {
	repo           repositories.IUnitPGRepository
	contextTimeout time.Duration
}

func NewUnitHttpService(repo repositories.IUnitPGRepository) *UnitHttpService {
	return &UnitHttpService{
		repo:           repo,
		contextTimeout: 15 * time.Second,
	}
}

func (svc UnitHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc UnitHttpService) GetModuleName() string {
	return "productunit"
}

// ✅ **GetUnit (ดึง Unit ตาม shopID และ unitCode)**
func (svc UnitHttpService) GetUnit(shopID string, unitCode string) (*models.UnitPg, error) {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()
	unit, err := svc.repo.FindByUnitCode(ctx, shopID, unitCode)
	if err != nil {
		return nil, err
	}

	return unit, nil
}

// ✅ **UnitList (ค้นหา Unit ตามชื่อ + Pagination)**
func (svc UnitHttpService) UnitList(shopID string, name string, page int, pageSize int) ([]models.UnitPg, mongopagination.PaginationData, error) {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	// ✅ ดึงข้อมูลจาก `repo.UnitList()`
	units, totalRecords, err := svc.repo.UnitList(ctx, shopID, name, page, pageSize)
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

// ✅ **Create (สร้าง Unit ใหม่)**
func (svc UnitHttpService) Create(doc *models.UnitPg) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if doc.ShopID == "" || doc.UnitCode == "" {
		return errors.New("ShopID and UnitCode are required")
	}

	// ✅ ตรวจสอบว่ามีอยู่แล้วหรือไม่
	existingUnit, err := svc.repo.Get(doc.ShopID, doc.UnitCode)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err // ❌ ถ้า error ไม่ใช่ "ไม่พบข้อมูล" ให้คืนค่า error ทันที
	}

	if existingUnit != nil {
		return errors.New("Unit already exists") // ❌ ถ้ามี unit อยู่แล้ว ไม่ต้องสร้างซ้ำ
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

// ✅ **Update (อัปเดต Unit)**
func (svc UnitHttpService) Update(shopID string, unitCode string, doc *models.UnitPg) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if shopID == "" || unitCode == "" {
		return errors.New("ShopID and UnitCode are required")
	}

	doc.UpdatedAt = time.Now()

	err := svc.repo.Update(ctx, shopID, unitCode, doc)
	if err != nil {
		return err
	}

	return nil
}

// ✅ **Delete (ลบ Unit)**
func (svc UnitHttpService) Delete(shopID string, unitCode string) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if shopID == "" || unitCode == "" {
		return errors.New("ShopID and UnitCode are required")
	}

	err := svc.repo.Delete(ctx, shopID, unitCode)
	if err != nil {
		return err
	}

	return nil
}
