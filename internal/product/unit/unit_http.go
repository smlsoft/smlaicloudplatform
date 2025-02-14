package unit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"smlaicloudplatform/internal/config"
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/product/unit/models"
	"smlaicloudplatform/internal/product/unit/repositories"
	"smlaicloudplatform/internal/product/unit/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/pkg/microservice"
	"strings"
	"time"
)

type IUnitHttp interface{}

type UnitHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IUnitHttpService
}

func NewUnitHttp(ms *microservice.Microservice, cfg config.IConfig) UnitHttp {
	pst := ms.Persister(cfg.PersisterConfig())
	repo := repositories.NewUnitPGRepository(pst)
	svc := services.NewUnitHttpService(repo)

	return UnitHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h UnitHttp) RegisterHttp() {

	h.ms.GET("/unit", h.SearchUnit)
	h.ms.POST("/unit", h.CreateUnit)
	h.ms.GET("/unit/:id", h.InfoUnit)
	h.ms.PUT("/unit/:id", h.UpdateUnit)
	h.ms.DELETE("/unit/:id", h.DeleteUnit)
}

// SearchUnit godoc
// @Summary		Search units
// @Description Search units with pagination
// @Tags		Unit
// @Accept 		json
// @Produce 	json
// @Param		q query string false "keyword"
// @Param		page query int false "Page number"
// @Param		limit query int false "Items per page"
// @Success		200 {object} common.ApiResponse{data=[]models.UnitPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/unit [get]
func (h UnitHttp) SearchUnit(ctx microservice.IContext) error { // ✅ แก้ไขให้ return error
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	pageable := utils.GetPageable(ctx.QueryParam)

	units, pagination, err := h.svc.UnitList(shopID, pageable.Query, pageable.Page, pageable.Limit)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err // ✅ Return error แทนการ return เปล่า
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Pagination: pagination,
		Data:       units,
	})
	return nil
}

// CreateUnit godoc
// @Summary		Create a new unit
// @Description Create a new unit with details
// @Tags		Unit
// @Accept 		json
// @Produce 	json
// @Param		Unit  body      models.UnitPg  true  "Unit data"
// @Success		201 {object} common.ApiResponse{data=models.UnitPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/unit [post]
func (h UnitHttp) CreateUnit(ctx microservice.IContext) error {
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
	newUnit := &models.UnitPg{}
	err := json.Unmarshal([]byte(input), newUnit)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return err
	}

	// ✅ กำหนดค่า `ShopID` จาก `userInfo`
	newUnit.ShopID = shopID
	newUnit.GuidFixed = utils.NewGUID()

	// ✅ ตรวจสอบค่า `UnitCode`
	newUnit.UnitCode = strings.TrimSpace(newUnit.UnitCode)
	if newUnit.UnitCode == "" {
		ctx.ResponseError(http.StatusBadRequest, "UnitCode is required")
		return errors.New("UnitCode is required")
	}

	// ✅ ตรวจสอบ Validation
	if err = ctx.Validate(newUnit); err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Validation failed: "+err.Error())
		return err
	}

	// ✅ กำหนดค่า `CreatedBy` และ `CreatedAt`
	newUnit.CreatedBy = userInfo.Username
	newUnit.CreatedAt = time.Now()

	// ✅ Debug: ตรวจสอบค่าก่อนสร้าง Unit
	fmt.Println("Creating Unit:", newUnit)

	// ✅ เรียก Service เพื่อสร้าง Unit
	err = h.svc.Create(newUnit)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	// ✅ ตอบกลับ Client
	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Message: "Unit created successfully",
		Data:    newUnit,
	})
	return nil
}

// InfoUnit godoc
// @Summary		Get unit details
// @Description Get unit details by ID
// @Tags		Unit
// @Accept 		json
// @Produce 	json
// @Param		id path string true "Unit ID"
// @Success		200 {object} common.ApiResponse{data=models.UnitPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		404 {object} common.ApiResponse
// @Security	AccessToken
// @Router		/unit/{id} [get]
func (h UnitHttp) InfoUnit(ctx microservice.IContext) error {
	// ✅ ดึงค่า `id` จาก URL Parameter
	unitID := strings.TrimSpace(ctx.Param("id"))
	userInfo := ctx.UserInfo()
	shopID := strings.TrimSpace(userInfo.ShopID)

	// ✅ ตรวจสอบว่า `id` ถูกส่งมาหรือไม่
	if unitID == "" {
		ctx.ResponseError(http.StatusBadRequest, "Unit ID is required")
		return errors.New("Unit ID is required")
	}

	// ✅ Debug: ดูค่าที่รับเข้ามา
	fmt.Println("Fetching Unit - ShopID:", shopID, "UnitID:", unitID)

	// ✅ เรียก Service เพื่อดึงข้อมูล Unit
	unit, err := h.svc.GetUnit(shopID, unitID)
	if err != nil {
		ctx.ResponseError(http.StatusNotFound, "Unit not found")
		return err
	}

	// ✅ ตอบกลับข้อมูล
	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    unit,
	})
	return nil
}

// UpdateUnit godoc
// @Summary		Update an existing unit
// @Description Update an existing unit by ID
// @Tags		Unit
// @Accept 		json
// @Produce 	json
// @Param		id path string true "Unit ID"
// @Param		Unit body models.UnitPg true "Updated unit data"
// @Success		200 {object} common.ApiResponse{data=models.UnitPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		404 {object} common.ApiResponse
// @Security	AccessToken
// @Router		/unit/{id} [put]
func (h UnitHttp) UpdateUnit(ctx microservice.IContext) error {
	// ✅ ดึงค่า `id` จาก URL Parameter
	unitID := strings.TrimSpace(ctx.Param("id"))
	userInfo := ctx.UserInfo()
	shopID := strings.TrimSpace(userInfo.ShopID)

	// ✅ ตรวจสอบว่า `id` ถูกส่งมาหรือไม่
	if unitID == "" {
		ctx.ResponseError(http.StatusBadRequest, "Unit ID is required")
		return errors.New("Unit ID is required")
	}

	// ✅ อ่าน JSON Request Body
	input := strings.TrimSpace(ctx.ReadInput())
	if input == "" {
		ctx.ResponseError(http.StatusBadRequest, "Invalid input: Empty request body")
		return errors.New("Invalid input: Empty request body")
	}

	// ✅ Debug: ตรวจสอบ input ก่อนทำ Unmarshal
	fmt.Println("Received Input:", input)

	// ✅ แปลง JSON เป็น struct
	updateData := &models.UnitPg{}
	err := json.Unmarshal([]byte(input), updateData)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return err
	}

	// ✅ ตรวจสอบค่า `UnitCode`
	updateData.UnitCode = strings.TrimSpace(updateData.UnitCode)
	if updateData.UnitCode == "" {
		ctx.ResponseError(http.StatusBadRequest, "UnitCode is required")
		return errors.New("UnitCode is required")
	}

	// ✅ ตรวจสอบ Validation
	if err = ctx.Validate(updateData); err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Validation failed: "+err.Error())
		return err
	}

	// ✅ กำหนดค่า `UpdatedBy` และ `UpdatedAt`
	updateData.UpdatedBy = userInfo.Username
	updateData.UpdatedAt = time.Now()

	// ✅ Debug: ตรวจสอบค่าก่อนอัปเดต Unit
	fmt.Println("Updating Unit:", updateData)

	// ✅ เรียก Service เพื่ออัปเดต Unit
	err = h.svc.Update(shopID, unitID, updateData)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	// ✅ ตอบกลับ Client
	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Message: "Unit updated successfully",
		Data:    updateData,
	})
	return nil
}

// DeleteUnit godoc
// @Summary		Delete a unit
// @Description Delete a unit by ID
// @Tags		Unit
// @Accept 		json
// @Produce 	json
// @Param		id path string true "Unit ID"
// @Success		200 {object} common.ApiResponse
// @Failure		400 {object} common.ApiResponse
// @Failure		404 {object} common.ApiResponse
// @Security	AccessToken
// @Router		/unit/{id} [delete]
func (h UnitHttp) DeleteUnit(ctx microservice.IContext) error {
	// ✅ ดึงค่า `id` จาก URL Parameter
	unitID := strings.TrimSpace(ctx.Param("id"))
	userInfo := ctx.UserInfo()
	shopID := strings.TrimSpace(userInfo.ShopID)

	// ✅ ตรวจสอบว่า `id` ถูกส่งมาหรือไม่
	if unitID == "" {
		ctx.ResponseError(http.StatusBadRequest, "Unit ID is required")
		return errors.New("Unit ID is required")
	}

	// ✅ Debug: ดูค่าที่รับเข้ามา
	fmt.Println("Deleting Unit - ShopID:", shopID, "UnitID:", unitID)

	// ✅ เรียก Service เพื่อลบ Unit
	err := h.svc.Delete(shopID, unitID)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	// ✅ ตอบกลับ Client
	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Message: "Unit deleted successfully",
	})
	return nil
}
