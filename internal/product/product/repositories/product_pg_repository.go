package repositories

import (
	"context"
	"errors"
	"fmt"
	dimension "smlaicloudplatform/internal/dimension/models"
	"smlaicloudplatform/internal/product/product/models"
	group "smlaicloudplatform/internal/product/productgroup/models"
	"smlaicloudplatform/pkg/microservice"
	"time"

	"gorm.io/gorm"
)

type IProductPGRepository interface {
	Get(ctx context.Context, shopID string, code string) (*models.ProductPg, error)
	ProductList(ctx context.Context, shopID string, name string, page int, pageSize int) ([]models.ProductPg, int64, error)
	Create(ctx context.Context, doc *models.ProductPg) error
	Update(ctx context.Context, shopID string, code string, doc *models.ProductPg) error
	Delete(ctx context.Context, shopID string, code string) error
}

type ProductPGRepository struct {
	pst microservice.IPersister
}

// ✅ **สร้าง Repository**
func NewProductPGRepository(pst microservice.IPersister) *ProductPGRepository {
	return &ProductPGRepository{
		pst: pst,
	}
}

// ✅ **Get (ดึง Product ตาม ShopID และ Code)**
func (repo *ProductPGRepository) Get(ctx context.Context, shopID string, code string) (*models.ProductPg, error) {
	var product models.ProductPg

	// ✅ ดึง Product ตาม ShopID และ Code
	err := repo.pst.DBClient().
		Where("shopid = ? AND guidfixed = ?", shopID, code).
		First(&product).Error
	if err != nil {
		return nil, err
	}

	// ดึง group name ถ้ามี groupguid
	if *product.GroupGuid != "" {
		var group group.ProductGroupPg
		err = repo.pst.DBClient().
			Where("shopid = ? AND guidfixed = ?", shopID, product.GroupGuid).
			First(&group).Error
		if err != nil {
			return nil, err
		}
		product.GroupCode = &group.Code
		product.GroupName = group.Names
	}

	// ✅ ดึง Dimensions ที่สัมพันธ์กับ Product
	var productDimensions []models.ProductDimensionPg
	err = repo.pst.DBClient().
		Where("shopid = ? AND product_guid = ?", shopID, product.GuidFixed).
		Find(&productDimensions).Error
	if err != nil {
		return nil, err
	}

	// ✅ ถ้าไม่มี Dimensions ให้คืนค่า `product` ทันที
	if len(productDimensions) == 0 {
		return &product, nil
	}

	// ✅ ดึง GUID ของ Dimensions
	var dimensionGuids []string
	for _, pd := range productDimensions {
		dimensionGuids = append(dimensionGuids, pd.DimensionGuid)
	}

	// ✅ ดึงข้อมูล Dimensions หลัก
	var dimensions []dimension.DimensionPg
	err = repo.pst.DBClient().
		Where("guidfixed IN ?", dimensionGuids).
		Where("shopid = ?", shopID).
		Find(&dimensions).Error
	if err != nil {
		return nil, err
	}

	// ✅ ดึง Items ของแต่ละ Dimension แยกต่างหาก
	for i, dim := range dimensions {
		var items []dimension.DimensionItemPg
		err = repo.pst.DBClient().
			Where("shopid = ? AND dimension_guid = ?", shopID, dim.GuidFixed).
			Find(&items).Error
		if err != nil {
			return nil, err
		}

		// ✅ ใส่ Items เข้าไปใน Dimension
		dimensions[i].Items = items
	}

	// ✅ ใส่ข้อมูล Dimensions กลับเข้าไปใน Product
	product.Dimensions = dimensions

	return &product, nil
}

func (repo *ProductPGRepository) GetProductDimension(ctx context.Context, shopID string, code string) ([]dimension.DimensionPg, error) {

	// ✅ ดึง Dimensions ที่สัมพันธ์กับ Product
	var productDimensions []models.ProductDimensionPg
	err := repo.pst.DBClient().
		Where("shopid = ? AND product_guid = ?", shopID, code).
		Find(&productDimensions).Error
	if err != nil {
		return nil, err
	}

	// ✅ ถ้าไม่มี Dimension ให้คืนค่าเป็น `nil`
	if len(productDimensions) == 0 {
		return nil, nil
	}

	// ✅ ดึง GUID ของ Dimensions
	var dimensionGuids []string
	for _, pd := range productDimensions {
		dimensionGuids = append(dimensionGuids, pd.DimensionGuid)
	}

	// ✅ ดึงข้อมูล Dimensions หลัก
	var dimensions []dimension.DimensionPg
	err = repo.pst.DBClient().
		Where("guidfixed IN ?", dimensionGuids).
		Where("shopid = ?", shopID).
		Find(&dimensions).Error
	if err != nil {
		return nil, err
	}

	// ✅ ดึง Items ของแต่ละ Dimension แยกต่างหาก
	for i, dim := range dimensions {
		var items []dimension.DimensionItemPg
		err = repo.pst.DBClient().
			Where("shopid = ? AND dimension_guid = ?", shopID, dim.GuidFixed).
			Find(&items).Error
		if err != nil {
			return nil, err
		}

		// ✅ ใส่ Items เข้าไปใน Dimension
		dimensions[i].Items = items
	}

	return dimensions, nil
}

// ✅ **ProductList (ค้นหา Products + Pagination)**
func (repo *ProductPGRepository) ProductList(ctx context.Context, shopID string, name string, page int, pageSize int) ([]models.ProductPg, int64, error) {
	var products []models.ProductPg
	var totalRecords int64

	// ✅ เริ่มต้น Query
	query := repo.pst.DBClient().
		Model(&models.ProductPg{}).
		Where("shopid = ?", shopID)

	// ✅ ถ้ามีการค้นหาชื่อ ให้เพิ่มเงื่อนไข ILIKE
	if name != "" {
		searchPattern := "%" + name + "%"
		query = query.Where("names::TEXT ILIKE ?", searchPattern)
	}

	// ✅ นับจำนวนทั้งหมด
	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	// ✅ คำนวณ offset และดึงข้อมูลตาม pagination
	err = query.Order("createdat DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&products).
		Error
	if err != nil {
		return nil, 0, err
	}

	// ✅ ดึง Dimensions สำหรับแต่ละ Product
	productGuids := []string{}
	for _, product := range products {
		productGuids = append(productGuids, product.GuidFixed)
	}

	// ✅ ดึงข้อมูลจาก Pivot Table `product_dimensions`
	var productDimensions []models.ProductDimensionPg
	if len(productGuids) > 0 {
		err = repo.pst.DBClient().
			Where("shopid = ? AND product_guid IN ?", shopID, productGuids).
			Find(&productDimensions).Error
		if err != nil {
			return nil, 0, err
		}
	}

	// ✅ รวมข้อมูล Dimensions จริง ๆ
	dimensionGuids := []string{}
	for _, pd := range productDimensions {
		dimensionGuids = append(dimensionGuids, pd.DimensionGuid)
	}

	var dimensions []dimension.DimensionPg
	if len(dimensionGuids) > 0 {
		err = repo.pst.DBClient().
			Where("guidfixed IN ?", dimensionGuids).
			Find(&dimensions).Error
		if err != nil {
			return nil, 0, err
		}
	}

	// ✅ Map Dimensions กลับไปยัง Products
	dimensionMap := make(map[string]dimension.DimensionPg)
	for _, dim := range dimensions {
		dimensionMap[dim.GuidFixed] = dim
	}

	for i := range products {
		for _, pd := range productDimensions {
			if products[i].GuidFixed == pd.ProductGuid {
				if dim, ok := dimensionMap[pd.DimensionGuid]; ok {
					products[i].Dimensions = append(products[i].Dimensions, dim)
				}
			}
		}
	}

	return products, totalRecords, nil
}

func (repo *ProductPGRepository) Create(ctx context.Context, doc *models.ProductPg) error {
	// ✅ ตรวจสอบค่า `shopid` ก่อนบันทึก
	if doc.ShopID == "" {
		return errors.New("ShopID cannot be empty")
	}
	if doc.GuidFixed == "" {
		return errors.New("Product GuidFixed cannot be empty")
	}

	// ✅ บันทึก Product
	err := repo.pst.Create(doc)
	if err != nil {
		return err
	}

	// ✅ Debug ก่อน Insert
	fmt.Printf("✅ Creating product dimensions for ShopID: %s, ProductGuid: %s\n", doc.ShopID, doc.GuidFixed)

	// ✅ เพิ่ม Dimensions ลงใน Pivot Table
	for _, dimension := range doc.Dimensions {
		if dimension.GuidFixed == "" {
			fmt.Println("❌ DimensionGuid is empty, skipping dimension:", dimension)
			continue
		}

		// ✅ ตรวจสอบว่าค่า `shopid` ถูกต้อง
		if dimension.ShopID == "" {
			dimension.ShopID = doc.ShopID
		}

		productDimension := models.ProductDimensionPg{
			ShopID:        doc.ShopID, // ✅ ตรวจสอบให้แน่ใจว่า `shopid` มีค่า
			ProductGuid:   doc.GuidFixed,
			DimensionGuid: dimension.GuidFixed,
		}

		// ✅ Debug Log ก่อน Insert
		fmt.Printf("📝 Before Insert: %+v\n", productDimension)

		// ✅ บันทึกค่าเข้า `product_dimensions`
		err = repo.pst.Create(&productDimension)
		if err != nil {
			fmt.Printf("❌ Error inserting product_dimension: %v\n", err)
			return err
		}
	}

	return nil
}

// ✅ **Update (อัปเดต Product)**
func (repo *ProductPGRepository) Update(ctx context.Context, shopID string, code string, doc *models.ProductPg) error {
	// ✅ ดึงข้อมูล Product ก่อนอัปเดต
	var existingProduct models.ProductPg
	err := repo.pst.DBClient().
		Where("shopid = ? AND guidfixed = ?", shopID, code).
		First(&existingProduct).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Product not found")
		}
		return err
	}

	// ✅ ตั้งค่า UpdatedAt
	doc.UpdatedAt = time.Now()

	// ✅ อัปเดตข้อมูลโดยใช้ `map[string]interface{}` เพื่อให้แน่ใจว่าค่าที่เป็น "" และ 0 จะถูกอัปเดต
	updateData := map[string]interface{}{
		"names":            doc.Names,
		"groupguid":        doc.GroupGuid,
		"unitguid":         doc.UnitGuid,
		"itemtype":         doc.ItemType,
		"manufacturerguid": doc.ManufacturerGUID,
		"updatedat":        doc.UpdatedAt,
		"updatedby":        doc.UpdatedBy,
	}

	err = repo.pst.DBClient().
		Model(&models.ProductPg{}).
		Where("shopid = ? AND guidfixed = ?", shopID, code).
		Updates(updateData).Error
	if err != nil {
		return err
	}

	// ✅ ลบ Dimensions เก่าออกก่อน
	err = repo.pst.DBClient().
		Where("shopid = ? AND product_guid = ?", shopID, code).
		Delete(&models.ProductDimensionPg{}).Error
	if err != nil {
		return err
	}

	// ✅ เพิ่ม Dimensions ใหม่
	for _, dimension := range doc.Dimensions {
		if dimension.GuidFixed == "" {
			return errors.New("DimensionGuid cannot be empty")
		}

		productDimension := models.ProductDimensionPg{
			ShopID:        shopID,
			ProductGuid:   doc.GuidFixed,
			DimensionGuid: dimension.GuidFixed,
		}

		err = repo.pst.DBClient().Create(&productDimension).Error
		if err != nil {
			fmt.Printf("Error inserting product_dimension: %v\n", err)
			return err
		}
	}

	return nil
}

// ✅ **Delete (ลบ Product)**
func (repo *ProductPGRepository) Delete(ctx context.Context, shopID string, code string) error {
	// ✅ ลบจากตารางหลัก
	err := repo.pst.Delete(&models.ProductPg{}, map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": code,
	})
	if err != nil {
		return err
	}

	// ✅ ลบจาก Pivot Table
	err = repo.pst.Delete(&models.ProductDimensionPg{}, map[string]interface{}{
		"shopid":       shopID,
		"product_guid": code,
	})
	if err != nil {
		return err
	}

	return nil
}
