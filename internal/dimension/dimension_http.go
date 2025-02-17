package dimension

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/dimension/models"
	"smlaicloudplatform/internal/dimension/repositories"
	"smlaicloudplatform/internal/dimension/services"
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/pkg/microservice"
	"strings"
	"time"
)

type IDimensionHttp interface{}

type DimensionHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IDimensionHttpService
}

func NewDimensionHttp(ms *microservice.Microservice, cfg config.IConfig) DimensionHttp {
	pst := ms.Persister(cfg.PersisterConfig())
	repo := repositories.NewDimensionPGRepository(pst)
	svc := services.NewDimensionHttpService(repo)

	return DimensionHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h DimensionHttp) RegisterHttp() {
	h.ms.GET("/dimension", h.SearchDimension)
	h.ms.POST("/dimension", h.CreateDimension)
	h.ms.GET("/dimension/:guidfixed", h.InfoDimension)
	h.ms.PUT("/dimension/:guidfixed", h.UpdateDimension)
	h.ms.DELETE("/dimension/:guidfixed", h.DeleteDimension)
}

// SearchDimension godoc
// @Summary		Search dimensions
// @Description Search dimensions with pagination
// @Tags		Dimension
// @Accept 		json
// @Produce 	json
// @Param		q query string false "Keyword for search"
// @Param		page query int false "Page number"
// @Param		limit query int false "Items per page"
// @Success		200 {object} common.ApiResponse{data=[]models.DimensionPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/dimension [get]
func (h DimensionHttp) SearchDimension(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	pageable := utils.GetPageable(ctx.QueryParam)

	dimensions, pagination, err := h.svc.DimensionList(shopID, pageable.Query, pageable.Page, pageable.Limit)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Pagination: pagination,
		Data:       dimensions,
	})
	return nil
}

// CreateDimension godoc
// @Summary		Create a new dimension
// @Description Create a new dimension with details and associated items
// @Tags		Dimension
// @Accept 		json
// @Produce 	json
// @Param		Dimension  body      models.DimensionPg  true  "Dimension data with items"
// @Success		201 {object} common.ApiResponse{data=models.DimensionPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/dimension [post]
func (h DimensionHttp) CreateDimension(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := strings.TrimSpace(userInfo.ShopID)

	// ✅ อ่าน JSON Request Body
	input := strings.TrimSpace(ctx.ReadInput())
	if input == "" {
		ctx.ResponseError(http.StatusBadRequest, "Invalid input: Empty request body")
		return errors.New("Invalid input: Empty request body")
	}

	// ✅ Debug JSON ที่รับมา
	fmt.Println("Received Input:", input)

	// ✅ แปลง JSON เป็น struct
	newDimension := &models.DimensionPg{}
	err := json.Unmarshal([]byte(input), newDimension)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return err
	}

	// ✅ ตรวจสอบว่ามี `items` หรือไม่
	if len(newDimension.Items) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "Items are required")
		return errors.New("Items are required")
	}

	// ✅ กำหนดค่า `ShopID` และ `GuidFixed`
	newDimension.ShopID = shopID
	newDimension.GuidFixed = utils.NewGUID()

	// ✅ กำหนด GUID และ ShopID ให้ Items
	for i := range newDimension.Items {
		newDimension.Items[i].ShopID = shopID
		newDimension.Items[i].DimensionGuid = newDimension.GuidFixed
		if newDimension.Items[i].GuidFixed == "" {
			newDimension.Items[i].GuidFixed = utils.NewGUID()
		}
	}

	// ✅ ตรวจสอบ Validation
	if err = ctx.Validate(newDimension); err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Validation failed: "+err.Error())
		return err
	}

	// ✅ กำหนดค่า `CreatedBy` และ `CreatedAt`
	newDimension.CreatedBy = userInfo.Username
	newDimension.CreatedAt = time.Now()

	// ✅ Debug
	fmt.Println("Creating Dimension:", newDimension)

	// ✅ เรียก Service เพื่อสร้าง Dimension และ Items
	err = h.svc.Create(newDimension, newDimension.Items)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	// ✅ ตอบกลับ Client
	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Message: "Dimension created successfully",
		Data:    newDimension,
	})
	return nil
}

// InfoDimension godoc
// @Summary		Get dimension details
// @Description Get dimension details by GUID
// @Tags		Dimension
// @Accept 		json
// @Produce 	json
// @Param		guidfixed path string true "Dimension GUID"
// @Success		200 {object} common.ApiResponse{data=models.DimensionPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		404 {object} common.ApiResponse
// @Security	AccessToken
// @Router		/dimension/{guidfixed} [get]
func (h DimensionHttp) InfoDimension(ctx microservice.IContext) error {
	guidFixed := strings.TrimSpace(ctx.Param("guidfixed"))
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	if guidFixed == "" {
		ctx.ResponseError(http.StatusBadRequest, "Dimension GUID is required")
		return errors.New("Dimension GUID is required")
	}

	dimension, err := h.svc.GetDimension(shopID, guidFixed)
	if err != nil {
		ctx.ResponseError(http.StatusNotFound, "Dimension not found")
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    dimension,
	})
	return nil
}

// UpdateDimension godoc
// @Summary		Update an existing dimension
// @Description Update an existing dimension by GUID, including associated items
// @Tags		Dimension
// @Accept 		json
// @Produce 	json
// @Param		guidfixed path string true "Dimension GUID"
// @Param		Dimension body models.DimensionPg true "Updated dimension data with items"
// @Success		200 {object} common.ApiResponse{data=models.DimensionPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		404 {object} common.ApiResponse
// @Security	AccessToken
// @Router		/dimension/{guidfixed} [put]
func (h DimensionHttp) UpdateDimension(ctx microservice.IContext) error {
	guidFixed := strings.TrimSpace(ctx.Param("guidfixed"))
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	if guidFixed == "" {
		ctx.ResponseError(http.StatusBadRequest, "Dimension GUID is required")
		return errors.New("Dimension GUID is required")
	}

	input := strings.TrimSpace(ctx.ReadInput())
	if input == "" {
		ctx.ResponseError(http.StatusBadRequest, "Invalid input: Empty request body")
		return errors.New("Invalid input: Empty request body")
	}

	// ✅ แปลง JSON เป็น struct
	updateData := &models.DimensionPg{}
	err := json.Unmarshal([]byte(input), updateData)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return err
	}

	updateData.UpdatedBy = userInfo.Username
	updateData.UpdatedAt = time.Now()

	// ✅ Debug
	fmt.Println("Updating Dimension:", updateData)

	// ✅ จัดการ DimensionItems
	for i := range updateData.Items {
		updateData.Items[i].ShopID = shopID
		updateData.Items[i].DimensionGuid = guidFixed
		if updateData.Items[i].GuidFixed == "" {
			updateData.Items[i].GuidFixed = utils.NewGUID()
		}
	}

	// ✅ อัปเดต Dimension พร้อม Items
	err = h.svc.Update(shopID, guidFixed, updateData, updateData.Items)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Message: "Dimension updated successfully",
		Data:    updateData,
	})
	return nil
}

// DeleteDimension godoc
// @Summary		Delete a dimension
// @Description Delete a dimension by GUID along with associated items
// @Tags		Dimension
// @Accept 		json
// @Produce 	json
// @Param		guidfixed path string true "Dimension GUID"
// @Success		200 {object} common.ApiResponse
// @Failure		400 {object} common.ApiResponse
// @Failure		404 {object} common.ApiResponse
// @Security	AccessToken
// @Router		/dimension/{guidfixed} [delete]
func (h DimensionHttp) DeleteDimension(ctx microservice.IContext) error {
	guidFixed := strings.TrimSpace(ctx.Param("guidfixed"))
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	if guidFixed == "" {
		ctx.ResponseError(http.StatusBadRequest, "Dimension GUID is required")
		return errors.New("Dimension GUID is required")
	}

	err := h.svc.Delete(shopID, guidFixed)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Message: "Dimension deleted successfully",
	})
	return nil
}
