package products

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"smlaicloudplatform/internal/config"
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/product/product/models"
	"smlaicloudplatform/internal/product/product/repositories"
	"smlaicloudplatform/internal/product/product/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/pkg/microservice"
	"strings"
	"time"
)

type IProductHttp interface{}

type ProductHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IProductHttpService
}

// ✅ **สร้าง New ProductHttp**
func NewProductHttp(ms *microservice.Microservice, cfg config.IConfig) ProductHttp {
	pst := ms.Persister(cfg.PersisterConfig())
	repo := repositories.NewProductPGRepository(pst)
	svc := services.NewProductHttpService(repo)

	return ProductHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

// ✅ **Register API Routes**
func (h ProductHttp) RegisterHttp() {
	h.ms.GET("/product", h.SearchProduct)
	h.ms.POST("/product", h.CreateProduct)
	h.ms.GET("/product/:code", h.InfoProduct)
	h.ms.PUT("/product/:code", h.UpdateProduct)
	h.ms.DELETE("/product/:code", h.DeleteProduct)
}

// @Summary		Search products
// @Description Search products with pagination
// @Tags		Product
// @Accept 		json
// @Produce 	json
// @Param		q query string false "Keyword to search"
// @Param		page query int false "Page number"
// @Param		limit query int false "Items per page"
// @Success		200 {object} common.ApiResponse{data=[]models.ProductPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/product [get]
func (h ProductHttp) SearchProduct(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	pageable := utils.GetPageable(ctx.QueryParam)

	products, pagination, err := h.svc.ProductList(shopID, pageable.Query, pageable.Page, pageable.Limit)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Pagination: pagination,
		Data:       products,
	})
	return nil
}

// @Summary		Create a new product
// @Description Create a new product with details
// @Tags		Product
// @Accept 		json
// @Produce 	json
// @Param		Product body models.ProductPg true "Product data"
// @Success		201 {object} common.ApiResponse{data=models.ProductPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/product [post]
func (h ProductHttp) CreateProduct(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := strings.TrimSpace(userInfo.ShopID)

	input := strings.TrimSpace(ctx.ReadInput())
	if input == "" {
		ctx.ResponseError(http.StatusBadRequest, "Invalid input: Empty request body")
		return errors.New("Invalid input: Empty request body")
	}

	// ✅ แปลง JSON เป็น struct
	newProduct := &models.ProductPg{}
	err := json.Unmarshal([]byte(input), newProduct)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return err
	}

	// ✅ กำหนดค่า `ShopID` และ `GuidFixed`
	newProduct.ShopID = shopID
	newProduct.GuidFixed = utils.NewGUID()

	// ✅ ตรวจสอบ Validation
	if err = ctx.Validate(newProduct); err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Validation failed: "+err.Error())
		return err
	}

	// ✅ กำหนดค่า `CreatedBy` และ `CreatedAt`
	newProduct.CreatedBy = userInfo.Username
	newProduct.CreatedAt = time.Now()

	// ✅ Debug
	fmt.Println("Creating Product:", newProduct)

	// ✅ เรียก Service เพื่อสร้าง Product
	err = h.svc.Create(newProduct)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    newProduct,
	})
	return nil
}

// @Summary		Get product details
// @Description Get product details by code
// @Tags		Product
// @Accept 		json
// @Produce 	json
// @Param		code path string true "Product Code"
// @Success		200 {object} common.ApiResponse{data=models.ProductPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		404 {object} common.ApiResponse
// @Security	AccessToken
// @Router		/product/{code} [get]
func (h ProductHttp) InfoProduct(ctx microservice.IContext) error {
	code := strings.TrimSpace(ctx.Param("code"))
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	if code == "" {
		ctx.ResponseError(http.StatusBadRequest, "Product Code is required")
		return errors.New("Product Code is required")
	}

	product, err := h.svc.GetProduct(shopID, code)
	if err != nil {
		ctx.ResponseError(http.StatusNotFound, "Product not found")
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    product,
	})
	return nil
}

// @Summary		Update an existing product
// @Description Update an existing product by code
// @Tags		Product
// @Accept 		json
// @Produce 	json
// @Param		code path string true "Product Code"
// @Param		Product body models.ProductPg true "Updated product data"
// @Success		200 {object} common.ApiResponse{data=models.ProductPg}
// @Failure		400 {object} common.ApiResponse
// @Failure		404 {object} common.ApiResponse
// @Security	AccessToken
// @Router		/product/{code} [put]
func (h ProductHttp) UpdateProduct(ctx microservice.IContext) error {
	code := strings.TrimSpace(ctx.Param("code"))
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	if code == "" {
		ctx.ResponseError(http.StatusBadRequest, "Product Code is required")
		return errors.New("Product Code is required")
	}

	input := strings.TrimSpace(ctx.ReadInput())
	if input == "" {
		ctx.ResponseError(http.StatusBadRequest, "Invalid input: Empty request body")
		return errors.New("Invalid input: Empty request body")
	}

	// ✅ แปลง JSON เป็น struct
	updateData := &models.ProductPg{}
	err := json.Unmarshal([]byte(input), updateData)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return err
	}

	updateData.UpdatedBy = userInfo.Username
	updateData.UpdatedAt = time.Now()

	// ✅ Debug
	fmt.Println("Updating Product:", updateData)

	// ✅ อัปเดต Product
	err = h.svc.Update(shopID, code, updateData)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Message: "Product updated successfully",
		Data:    updateData,
	})
	return nil
}

// @Summary		Delete a product
// @Description Delete a product by code
// @Tags		Product
// @Accept 		json
// @Produce 	json
// @Param		code path string true "Product Code"
// @Success		200 {object} common.ApiResponse
// @Failure		400 {object} common.ApiResponse
// @Failure		404 {object} common.ApiResponse
// @Security	AccessToken
// @Router		/product/{code} [delete]
func (h ProductHttp) DeleteProduct(ctx microservice.IContext) error {
	code := strings.TrimSpace(ctx.Param("code"))
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	if code == "" {
		ctx.ResponseError(http.StatusBadRequest, "Product Code is required")
		return errors.New("Product Code is required")
	}

	err := h.svc.Delete(shopID, code)
	if err != nil {
		ctx.ResponseError(http.StatusInternalServerError, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
	return nil
}
