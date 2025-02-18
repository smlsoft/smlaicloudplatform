package repositories

import (
	"context"
	"smlaicloudplatform/internal/product/unit/models"
	"smlaicloudplatform/pkg/microservice"

	"gorm.io/gorm"
)

type IUnitPGRepository interface {
	Get(shopID string, unitcode string) (*models.UnitPg, error)
	FindByUnitCode(ctx context.Context, shopID string, unitcode string) (*models.UnitPg, error)
	FindByUnitCodes(ctx context.Context, shopID string, unitcodes []string) ([]models.UnitPg, error)
	UnitList(ctx context.Context, shopID string, name string, page int, pageSize int) ([]models.UnitPg, int64, error)
	Create(ctx context.Context, doc *models.UnitPg) error
	Update(ctx context.Context, shopID string, unitcode string, doc *models.UnitPg) error
	Delete(ctx context.Context, shopID string, unitcode string) error
}

type UnitPGRepository struct {
	pst microservice.IPersister
}

func NewUnitPGRepository(pst microservice.IPersister) *UnitPGRepository {
	return &UnitPGRepository{
		pst: pst,
	}
}

// ✅ **UnitList (SELECT + Pagination)**
func (repo *UnitPGRepository) UnitList(ctx context.Context, shopID string, name string, page int, pageSize int) ([]models.UnitPg, int64, error) {
	var units []models.UnitPg
	var totalRecords int64

	// ✅ เริ่มต้น Query
	query := repo.pst.DBClient().
		Model(&models.UnitPg{}).
		Where("shopid = ?", shopID)

	// ✅ ถ้า `name` ไม่เป็นค่าว่าง ให้เพิ่มเงื่อนไข ILIKE
	if name != "" {
		searchPattern := "%" + name + "%"
		query = query.Where("names::TEXT ILIKE ?", searchPattern)
	}

	// ✅ นับจำนวนทั้งหมด
	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	// ✅ ดึงข้อมูลตาม pagination
	err = query.Order("createdat DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&units).
		Error
	if err != nil {
		return nil, 0, err
	}

	return units, totalRecords, nil
}

func (repo *UnitPGRepository) Get(shopID string, unitcode string) (*models.UnitPg, error) {
	var result models.UnitPg
	_, err := repo.pst.First(&result, "shopid=? AND guidfixed=?", shopID, unitcode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &result, nil
}

// ✅ **FindByUnitCode (SELECT 1 รายการ)**
func (repo *UnitPGRepository) FindByUnitCode(ctx context.Context, shopID string, unitcode string) (*models.UnitPg, error) {
	var result models.UnitPg
	_, err := repo.pst.First(&result, "shopid=? AND unitcode = ?", shopID, unitcode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ✅ **FindByUnitCodes (SELECT หลายรายการ)**
func (repo *UnitPGRepository) FindByUnitCodes(ctx context.Context, shopID string, unitcodes []string) ([]models.UnitPg, error) {
	var results []models.UnitPg
	_, err := repo.pst.Where(&results, "shopid=? AND unitcode IN ?", shopID, unitcodes)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return results, nil
}

// ✅ **Create (INSERT)**
func (repo *UnitPGRepository) Create(ctx context.Context, doc *models.UnitPg) error {
	err := repo.pst.Create(doc)
	if err != nil {
		return err
	}
	return nil
}

// ✅ **Update (UPDATE)**
func (repo *UnitPGRepository) Update(ctx context.Context, shopID string, unitcode string, doc *models.UnitPg) error {
	err := repo.pst.Update(&doc, map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": unitcode,
	})
	if err != nil {
		return err
	}
	return nil
}

// ✅ **Delete (DELETE)**
func (repo *UnitPGRepository) Delete(ctx context.Context, shopID string, unitcode string) error {
	err := repo.pst.Delete(&models.UnitPg{}, map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": unitcode,
	})
	if err != nil {
		return err
	}
	return nil
}
