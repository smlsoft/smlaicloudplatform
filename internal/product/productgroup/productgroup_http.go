package productgroup

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"smlaicloudplatform/internal/config"
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/product/productgroup/models"
	"smlaicloudplatform/internal/product/productgroup/repositories"
	"smlaicloudplatform/internal/product/productgroup/services"

	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/pkg/microservice"
	"strings"
	"time"
)

type IProductGroupHttp interface{}

type ProductGroupHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IProductGroupHttpService
}

func NewProductGroupHttp(ms *microservice.Microservice, cfg config.IConfig) ProductGroupHttp {
	pst := ms.Persister(cfg.PersisterConfig())
	repo := repositories.NewProductGroupPGRepository(pst)
	svc := services.NewProductGroupHttpService(repo)

	return ProductGroupHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ProductGroupHttp) RegisterHttp() {

	h.ms.GET("/productgroup", h.Search)
	h.ms.POST("/productgroup", h.Create)
	h.ms.GET("/productgroup/:id", h.Info)
	h.ms.PUT("/productgroup/:id", h.Update)
	h.ms.DELETE("/productgroup/:id", h.Delete)
}

// Search Product Groups godoc
// @Summary     Search Product Groups
// @Description ค้นหารายการ Product Group
// @Tags        ProductGroup
// @Param       q    query    string  false "คำค้นหา (optional)"
// @Param       page     query    int     false "หมายเลขหน้า (default: 1)"
// @Param       limit    query    int     false "จำนวนรายการต่อหน้า (default: 10)"
// @Accept      json
// @Produce     json
// @Success     200 {object} common.ApiResponse{data=[]models.ProductGroupPg}
// @Failure     400 {object} common.ApiResponse
// @Failure     401 {object} common.ApiResponse
// @Security    AccessToken
// @Router      /productgroup [get]
func (h ProductGroupHttp) Search(ctx microservice.IContext) error { // ✅ แก้ไขให้ return error
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	pageable := utils.GetPageable(ctx.QueryParam)

	datas, pagination, err := h.svc.ProductGroupList(shopID, pageable.Query, pageable.Page, pageable.Limit)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err // ✅ Return error แทนการ return เปล่า
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Pagination: pagination,
		Data:       datas,
	})
	return nil
}

// Create Product Group godoc
// @Summary     Create Product Group
// @Description สร้าง Product Group ใหม่
// @Tags        ProductGroup
// @Param       request body models.ProductGroupPg true "รายละเอียด Product Group"
// @Accept      json
// @Produce     json
// @Success     201 {object} common.ApiResponse{data=models.ProductGroupPg}
// @Failure     400 {object} common.ApiResponse
// @Failure     401 {object} common.ApiResponse
// @Security    AccessToken
// @Router      /productgroup [post]
func (h ProductGroupHttp) Create(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := strings.TrimSpace(userInfo.ShopID)

	// ✅ อ่าน JSON Request Body
	input := strings.TrimSpace(ctx.ReadInput())
	if input == "" {
		ctx.ResponseError(http.StatusBadRequest, "Invalid input: Empty request body")
		return errors.New("Invalid input: Empty request body")
	}

	// ✅ Debug: ตรวจสอบ input ก่อนทำ Unmarshal
	fmt.Println("Received Input:", input)

	// ✅ แปลง JSON เป็น struct
	newProductGroup := &models.ProductGroupPg{}
	err := json.Unmarshal([]byte(input), newProductGroup)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return err
	}

	// ✅ กำหนดค่า `ShopID` จาก `userInfo`
	newProductGroup.ShopID = shopID
	newProductGroup.GuidFixed = utils.NewGUID() // สร้าง GUID สำหรับ Product Group

	// ✅ ตรวจสอบค่า `Code`
	newProductGroup.Code = strings.TrimSpace(newProductGroup.Code)
	if newProductGroup.Code == "" {
		ctx.ResponseError(http.StatusBadRequest, "Code is required")
		return errors.New("Code is required")
	}

	// ✅ ตรวจสอบ Validation
	if err = ctx.Validate(newProductGroup); err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Validation failed: "+err.Error())
		return err
	}

	// ✅ กำหนดค่า `CreatedBy` และ `CreatedAt`
	newProductGroup.CreatedBy = userInfo.Username
	newProductGroup.CreatedAt = time.Now()

	// ✅ Debug: ตรวจสอบค่าก่อนสร้าง Product Group
	fmt.Println("Creating Product Group:", newProductGroup)

	// ✅ เรียก Service เพื่อสร้าง Product Group
	err = h.svc.Create(newProductGroup)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	// ✅ ตอบกลับ Client
	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Message: "Product Group created successfully",
		Data:    newProductGroup,
	})
	return nil
}

// Get Product Group godoc
// @Summary     Get Product Group
// @Description ดึงข้อมูล Product Group ตามรหัส
// @Tags        ProductGroup
// @Param       id path string true "รหัส Product Group"
// @Accept      json
// @Produce     json
// @Success     200 {object} common.ApiResponse{data=models.ProductGroupPg}
// @Failure     400 {object} common.ApiResponse
// @Failure     404 {object} common.ApiResponse
// @Security    AccessToken
// @Router      /productgroup/{id} [get]
func (h ProductGroupHttp) Info(ctx microservice.IContext) error {
	// ✅ ดึง `Code` จาก URL
	code := strings.TrimSpace(ctx.Param("id"))
	userInfo := ctx.UserInfo()
	shopID := strings.TrimSpace(userInfo.ShopID)

	// ✅ ตรวจสอบว่ามีค่า `Code` หรือไม่
	if code == "" {
		ctx.ResponseError(http.StatusBadRequest, "Product Group Code is required")
		return errors.New("Product Group Code is required")
	}

	// ✅ ดึงข้อมูลจาก Service
	productGroup, err := h.svc.Get(shopID, code)
	if err != nil {
		ctx.ResponseError(http.StatusNotFound, "Product Group not found")
		return err
	}

	// ✅ ส่ง Response กลับ
	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    productGroup,
	})
	return nil
}

// Update Product Group godoc
// @Summary     Update Product Group
// @Description อัปเดตข้อมูล Product Group
// @Tags        ProductGroup
// @Param       id    path  string  true  "รหัส Product Group"
// @Param       request body  models.ProductGroupPg true "ข้อมูล Product Group ที่ต้องการอัปเดต"
// @Accept      json
// @Produce     json
// @Success     200 {object} common.ApiResponse{data=models.ProductGroupPg}
// @Failure     400 {object} common.ApiResponse
// @Failure     401 {object} common.ApiResponse
// @Failure     404 {object} common.ApiResponse
// @Security    AccessToken
// @Router      /productgroup/{id} [put]
func (h ProductGroupHttp) Update(ctx microservice.IContext) error {
	// ✅ ดึง `Code` จาก URL
	code := strings.TrimSpace(ctx.Param("id"))
	userInfo := ctx.UserInfo()
	shopID := strings.TrimSpace(userInfo.ShopID)

	// ✅ ตรวจสอบว่ามีค่า `Code` หรือไม่
	if code == "" {
		ctx.ResponseError(http.StatusBadRequest, "Product Group Code is required")
		return errors.New("Product Group Code is required")
	}

	// ✅ อ่าน JSON Request Body
	input := strings.TrimSpace(ctx.ReadInput())
	if input == "" {
		ctx.ResponseError(http.StatusBadRequest, "Invalid input: Empty request body")
		return errors.New("Invalid input: Empty request body")
	}

	// ✅ แปลง JSON เป็น struct
	updateData := &models.ProductGroupPg{}
	err := json.Unmarshal([]byte(input), updateData)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return err
	}

	// ✅ ตรวจสอบ Validation
	if err = ctx.Validate(updateData); err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Validation failed: "+err.Error())
		return err
	}

	// ✅ กำหนดค่า `UpdatedBy` และ `UpdatedAt`
	updateData.UpdatedBy = userInfo.Username
	updateData.UpdatedAt = time.Now()

	// ✅ อัปเดตข้อมูล
	err = h.svc.Update(shopID, code, updateData)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	// ✅ ส่ง Response กลับ
	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Message: "Product Group updated successfully",
		Data:    updateData,
	})
	return nil
}

// Delete Product Group godoc
// @Summary     Delete Product Group
// @Description ลบ Product Group ตามรหัส
// @Tags        ProductGroup
// @Param       id path string true "รหัส Product Group"
// @Accept      json
// @Produce     json
// @Success     200 {object} common.ApiResponse
// @Failure     400 {object} common.ApiResponse
// @Failure     401 {object} common.ApiResponse
// @Failure     404 {object} common.ApiResponse
// @Security    AccessToken
// @Router      /productgroup/{id} [delete]
func (h ProductGroupHttp) Delete(ctx microservice.IContext) error {
	// ✅ ดึง `Code` จาก URL
	code := strings.TrimSpace(ctx.Param("id"))
	userInfo := ctx.UserInfo()
	shopID := strings.TrimSpace(userInfo.ShopID)

	// ✅ ตรวจสอบว่ามีค่า `Code` หรือไม่
	if code == "" {
		ctx.ResponseError(http.StatusBadRequest, "Product Group Code is required")
		return errors.New("Product Group Code is required")
	}

	// ✅ ลบข้อมูล
	err := h.svc.Delete(shopID, code)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	// ✅ ส่ง Response กลับ
	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Message: "Product Group deleted successfully",
	})
	return nil
}
