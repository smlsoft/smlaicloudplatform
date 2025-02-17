package repositories

import (
	"context"
	"smlaicloudplatform/internal/dimension/models"
	"smlaicloudplatform/pkg/microservice"
	"time"

	"gorm.io/gorm"
)

type IDimensionPGRepository interface {
	Get(ctx context.Context, shopID string, guidFixed string) (*models.DimensionPg, error)
	DimensionList(ctx context.Context, shopID string, name string, page int, pageSize int) ([]models.DimensionPg, int64, error)
	Create(ctx context.Context, doc *models.DimensionPg, items []models.DimensionItemPg) error
	Update(ctx context.Context, shopID string, guidFixed string, doc *models.DimensionPg, items []models.DimensionItemPg) error
	Delete(ctx context.Context, shopID string, guidFixed string) error
}

type DimensionPGRepository struct {
	pst microservice.IPersister
}

// ✅ **สร้าง New Repository**
func NewDimensionPGRepository(pst microservice.IPersister) *DimensionPGRepository {
	return &DimensionPGRepository{
		pst: pst,
	}
}

// ✅ **Get Dimension + Items**
func (repo *DimensionPGRepository) Get(ctx context.Context, shopID string, id string) (*models.DimensionPg, error) {
	var dimension models.DimensionPg
	err := repo.pst.DBClient().
		Where("guidfixed = ?", id).Where("shopid = ?", shopID).
		First(&dimension).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// ✅ ดึง Items แยกต่างหาก
	var items []models.DimensionItemPg
	err = repo.pst.DBClient().
		Where("dimension_guid = ?", id).Where("shopid = ?", shopID).
		Find(&items).
		Error
	if err != nil {
		return nil, err
	}

	dimension.Items = items
	return &dimension, nil
}

// ✅ **DimensionList (ดึงรายการ + Items)**
func (repo *DimensionPGRepository) DimensionList(ctx context.Context, shopID string, name string, page int, pageSize int) ([]models.DimensionPg, int64, error) {
	var results []models.DimensionPg
	var totalRecords int64

	query := repo.pst.DBClient().Model(&models.DimensionPg{}).Where("shopid = ?", shopID)

	if name != "" {
		searchPattern := "%" + name + "%"
		query = query.Where("names::TEXT ILIKE ?", searchPattern)
	}

	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Order("createdat DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&results).
		Error
	if err != nil {
		return nil, 0, err
	}

	// ✅ โหลด Items แยกต่างหาก
	for i := range results {
		var items []models.DimensionItemPg
		err = repo.pst.DBClient().
			Where("dimension_guid = ?", results[i].GuidFixed).Where("shopid = ?", shopID).
			Find(&items).
			Error
		if err != nil {
			return nil, 0, err
		}
		results[i].Items = items
	}

	return results, totalRecords, nil
}

// ✅ **Create (เพิ่ม Dimension พร้อม Items)**
func (repo *DimensionPGRepository) Create(ctx context.Context, doc *models.DimensionPg, items []models.DimensionItemPg) error {
	return repo.pst.Transaction(func(tx *microservice.Persister) error {
		err := tx.Create(doc)
		if err != nil {
			return err
		}

		if len(items) > 0 {
			for i := range items {
				items[i].DimensionGuid = doc.GuidFixed // 🔥 เชื่อมโยงกับ Dimension
			}
			err = tx.CreateInBatch(items, len(items))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ✅ **Update (อัปเดต Dimension และ Items)**
func (repo *DimensionPGRepository) Update(ctx context.Context, shopID string, id string, doc *models.DimensionPg, items []models.DimensionItemPg) error {
	return repo.pst.Transaction(func(tx *microservice.Persister) error {
		doc.UpdatedAt = time.Now()

		// ✅ อัปเดตข้อมูลหลัก
		err := tx.Update(doc, map[string]interface{}{
			"guidfixed": id,
			"shopid":    shopID,
		})
		if err != nil {
			return err
		}

		// ✅ ลบ Items เก่าของ Dimension
		err = tx.Delete(&models.DimensionItemPg{}, map[string]interface{}{
			"dimension_guid": id,
			"shopid":         shopID,
		})
		if err != nil {
			return err
		}

		// ✅ เพิ่ม Items ใหม่เข้าไป
		if len(items) > 0 {
			for i := range items {
				items[i].DimensionGuid = id
			}
			err = tx.CreateInBatch(items, len(items))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// ✅ **Delete (ลบ Dimension พร้อม Items)**
func (repo *DimensionPGRepository) Delete(ctx context.Context, shopID string, id string) error {
	return repo.pst.Transaction(func(tx *microservice.Persister) error {
		// ✅ ลบ Items ก่อน
		err := tx.Delete(&models.DimensionItemPg{}, map[string]interface{}{
			"dimension_guid": id,
			"shopid":         shopID,
		})
		if err != nil {
			return err
		}

		// ✅ ลบ Dimension หลัก
		err = tx.Delete(&models.DimensionPg{}, map[string]interface{}{
			"guidfixed": id,
			"shopid":    shopID,
		})
		if err != nil {
			return err
		}

		return nil
	})
}
