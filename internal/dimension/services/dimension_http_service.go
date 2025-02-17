package services

import (
	"context"
	"errors"
	"math"
	"smlaicloudplatform/internal/dimension/models"
	"smlaicloudplatform/internal/dimension/repositories"
	"smlaicloudplatform/internal/utils"
	"time"

	"github.com/smlsoft/mongopagination"
	"gorm.io/gorm"
)

type IDimensionHttpService interface {
	GetModuleName() string
	GetDimension(shopID, guidFixed string) (*models.DimensionPg, error)
	DimensionList(shopID, name string, page, pageSize int) ([]models.DimensionPg, mongopagination.PaginationData, error)
	Create(doc *models.DimensionPg, items []models.DimensionItemPg) error
	Update(shopID, guidFixed string, doc *models.DimensionPg, items []models.DimensionItemPg) error
	Delete(shopID, guidFixed string) error
}

type DimensionHttpService struct {
	repo           repositories.IDimensionPGRepository
	contextTimeout time.Duration
}

// ✅ **สร้าง Service**
func NewDimensionHttpService(repo repositories.IDimensionPGRepository) *DimensionHttpService {
	return &DimensionHttpService{
		repo:           repo,
		contextTimeout: 15 * time.Second,
	}
}

// ✅ **ตั้งค่า Timeout**
func (svc DimensionHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc DimensionHttpService) GetModuleName() string {
	return "dimension"
}

// ✅ **GetDimension (ดึง Dimension ตาม ShopID และ GuidFixed)**
func (svc DimensionHttpService) GetDimension(shopID, guidFixed string) (*models.DimensionPg, error) {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	dimension, err := svc.repo.Get(ctx, shopID, guidFixed)
	if err != nil {
		return nil, err
	}

	return dimension, nil
}

// ✅ **DimensionList (ค้นหา Dimension ตามชื่อ + Pagination)**
func (svc DimensionHttpService) DimensionList(shopID, name string, page, pageSize int) ([]models.DimensionPg, mongopagination.PaginationData, error) {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	// ✅ ดึงข้อมูลจาก `repo.DimensionList()`
	dimensions, totalRecords, err := svc.repo.DimensionList(ctx, shopID, name, page, pageSize)
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

	return dimensions, pagination, nil
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

// ✅ **Create (สร้าง Dimension ใหม่)**
func (svc DimensionHttpService) Create(doc *models.DimensionPg, items []models.DimensionItemPg) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if doc.ShopID == "" {
		return errors.New("ShopID is required")
	}

	// ✅ สร้าง `GuidFixed` ถ้ายังไม่มีค่า
	if doc.GuidFixed == "" {
		doc.GuidFixed = utils.NewGUID()
	}

	// ✅ ตรวจสอบว่ามีอยู่แล้วหรือไม่
	existingDimension, err := svc.repo.Get(ctx, doc.ShopID, doc.GuidFixed)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existingDimension != nil {
		return errors.New("Dimension already exists")
	}

	// ✅ ตั้งค่าเวลาก่อนสร้าง
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	// ✅ สร้าง `GuidFixed` สำหรับ Items
	for i := range items {
		items[i].ShopID = doc.ShopID
		items[i].DimensionGuid = doc.GuidFixed
		items[i].GuidFixed = utils.NewGUID()
	}

	// ✅ เรียก `Create()`
	err = svc.repo.Create(ctx, doc, items)
	if err != nil {
		return err
	}

	return nil
}

// ✅ **Update (อัปเดต Dimension และ Items)**
func (svc DimensionHttpService) Update(shopID, guidFixed string, doc *models.DimensionPg, items []models.DimensionItemPg) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if shopID == "" || guidFixed == "" {
		return errors.New("ShopID and GuidFixed are required")
	}

	doc.UpdatedAt = time.Now()

	// ✅ สร้าง `GuidFixed` สำหรับ Items ใหม่
	for i := range items {
		items[i].ShopID = shopID
		items[i].DimensionGuid = guidFixed
		items[i].GuidFixed = utils.NewGUID()
	}

	// ✅ อัปเดต Dimension และ Items
	err := svc.repo.Update(ctx, shopID, guidFixed, doc, items)
	if err != nil {
		return err
	}

	return nil
}

// ✅ **Delete (ลบ Dimension และ Items)**
func (svc DimensionHttpService) Delete(shopID, guidFixed string) error {
	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	if shopID == "" || guidFixed == "" {
		return errors.New("ShopID and GuidFixed are required")
	}

	err := svc.repo.Delete(ctx, shopID, guidFixed)
	if err != nil {
		return err
	}

	return nil
}
