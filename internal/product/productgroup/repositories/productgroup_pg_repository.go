package repositories

import (
	"context"

	"smlaicloudplatform/internal/product/productgroup/models"
	"smlaicloudplatform/pkg/microservice"

	"gorm.io/gorm"
)

type IProductGroupPGRepository interface {
	Get(shopID string, code string) (*models.ProductGroupPg, error)
	FindByProductGroupCode(ctx context.Context, shopID string, code string) (*models.ProductGroupPg, error)
	FindByProductGroupCodes(ctx context.Context, shopID string, codes []string) ([]models.ProductGroupPg, error)
	ProductGroupList(ctx context.Context, shopID string, name string, page int, pageSize int) ([]models.ProductGroupPg, int64, error)
	Create(ctx context.Context, doc *models.ProductGroupPg) error
	Update(ctx context.Context, shopID string, code string, doc *models.ProductGroupPg) error
	Delete(ctx context.Context, shopID string, code string) error
}

type ProductGroupPGRepository struct {
	pst microservice.IPersister
}

func NewProductGroupPGRepository(pst microservice.IPersister) *ProductGroupPGRepository {
	return &ProductGroupPGRepository{
		pst: pst,
	}
}

// ✅ **ProductGroupList (SELECT + Pagination)**
func (repo *ProductGroupPGRepository) ProductGroupList(ctx context.Context, shopID string, name string, page int, pageSize int) ([]models.ProductGroupPg, int64, error) {
	var units []models.ProductGroupPg
	var totalRecords int64

	// ✅ เริ่มต้น Query
	query := repo.pst.DBClient().
		Model(&models.ProductGroupPg{}).
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

func (repo *ProductGroupPGRepository) Get(shopID string, code string) (*models.ProductGroupPg, error) {
	var result models.ProductGroupPg
	_, err := repo.pst.First(&result, "shopid=? AND code=?", shopID, code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &result, nil
}

// ✅ **FindByProductGroupCode (SELECT 1 รายการ)**
func (repo *ProductGroupPGRepository) FindByProductGroupCode(ctx context.Context, shopID string, code string) (*models.ProductGroupPg, error) {
	var result models.ProductGroupPg
	_, err := repo.pst.First(&result, "shopid=? AND code = ?", shopID, code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ✅ **FindByProductGroupCodes (SELECT หลายรายการ)**
func (repo *ProductGroupPGRepository) FindByProductGroupCodes(ctx context.Context, shopID string, codes []string) ([]models.ProductGroupPg, error) {
	var results []models.ProductGroupPg
	_, err := repo.pst.Where(&results, "shopid=? AND code IN ?", shopID, codes)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return results, nil
}

// ✅ **Create (INSERT)**
func (repo *ProductGroupPGRepository) Create(ctx context.Context, doc *models.ProductGroupPg) error {
	err := repo.pst.Create(doc)
	if err != nil {
		return err
	}
	return nil
}

// ✅ **Update (UPDATE)**
func (repo *ProductGroupPGRepository) Update(ctx context.Context, shopID string, code string, doc *models.ProductGroupPg) error {
	err := repo.pst.Update(&doc, map[string]interface{}{
		"shopid": shopID,
		"code":   code,
	})
	if err != nil {
		return err
	}
	return nil
}

// ✅ **Delete (DELETE)**
func (repo *ProductGroupPGRepository) Delete(ctx context.Context, shopID string, code string) error {
	err := repo.pst.Delete(&models.ProductGroupPg{}, map[string]interface{}{
		"shopid": shopID,
		"code":   code,
	})
	if err != nil {
		return err
	}
	return nil
}
