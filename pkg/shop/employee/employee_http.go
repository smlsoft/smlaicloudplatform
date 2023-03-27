package employee

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/shop/employee/models"
	"smlcloudplatform/pkg/shop/employee/repositories"
	"smlcloudplatform/pkg/shop/employee/services"
	"smlcloudplatform/pkg/utils"
)

type IEmployeeHttp interface{}

type EmployeeHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IEmployeeHttpService
}

func NewEmployeeHttp(ms *microservice.Microservice, cfg microservice.IConfig) EmployeeHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewEmployeeRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewEmployeeHttpService(repo, masterSyncCacheRepo, utils.HashPassword)

	return EmployeeHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h EmployeeHttp) RouteSetup() {

	h.ms.GET("/shop/employee", h.SearchEmployeePage)
	h.ms.GET("/shop/employee/list", h.SearchEmployeeStep)
	h.ms.POST("/shop/employee", h.CreateEmployee)
	h.ms.GET("/shop/employee/:id", h.InfoEmployee)
	h.ms.PUT("/shop/employee/:id", h.UpdateEmployee)
	h.ms.PUT("/shop/employee/password", h.UpdatePassword)
	h.ms.DELETE("/shop/employee/:id", h.DeleteEmployee)
	h.ms.DELETE("/shop/employee", h.DeleteEmployeeByGUIDs)
}

// Create Employee godoc
// @Description Create Employee
// @Tags		Employee
// @Param		Employee  body      models.EmployeeRequestRegister  true  "Employee"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/employee [post]
func (h EmployeeHttp) CreateEmployee(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.EmployeeRequestRegister{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateEmployee(shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})
	return nil
}

// Update Employee godoc
// @Description Update Employee
// @Tags		Employee
// @Param		id  path      string  true  "Employee ID"
// @Param		Employee  body      models.EmployeeRequestUpdate  true  "Employee"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/employee/{id} [put]
func (h EmployeeHttp) UpdateEmployee(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.EmployeeRequestUpdate{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateEmployee(shopID, id, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

// Delete Employee godoc
// @Description Delete Employee
// @Tags		Employee
// @Param		id  path      string  true  "Employee ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/employee/{id} [delete]
func (h EmployeeHttp) DeleteEmployee(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteEmployee(shopID, id, authUsername)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

// Delete Employee godoc
// @Description Delete Employee
// @Tags		Employee
// @Param		Employee  body      []string  true  "Employee GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/employee [delete]
func (h EmployeeHttp) DeleteEmployeeByGUIDs(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	input := ctx.ReadInput()

	docReq := []string{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.DeleteEmployeeByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get Employee godoc
// @Description get struct array by ID
// @Tags		Employee
// @Param		id  path      string  true  "Employee ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/employee/{id} [get]
func (h EmployeeHttp) InfoEmployee(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Employee %v", id)
	doc, err := h.svc.InfoEmployee(shopID, id)

	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v: %v", id, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// List Employee godoc
// @Description get struct array by ID
// @Tags		Employee
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/employee [get]
func (h EmployeeHttp) SearchEmployeePage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.svc.SearchEmployee(shopID, map[string]interface{}{}, pageable)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Data:       docList,
		Pagination: pagination,
	})
	return nil
}

// List Employee godoc
// @Description search limit offset
// @Tags		Employee
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/employee/list [get]
func (h EmployeeHttp) SearchEmployeeStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchEmployeeStep(shopID, lang, pageableStep)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    docList,
		Total:   total,
	})
	return nil
}

// Update Password Employee godoc
// @Summary		Update Password Employee
// @Description	Update Password Employee
// @Tags		Employee
// @Param		id  path      string  true  "Employee ID"
// @Param		Employee  body      models.EmployeeRequestPassword  true  "Register Employee"
// @Success		200	{object}	models.ResponseSuccess
// @Failure		400 {object}	models.AuthResponseFailed
// @Accept 		json
// @Security     AccessToken
// @Router		/employee/password [put]
func (h EmployeeHttp) UpdatePassword(ctx microservice.IContext) error {
	userAuthInfo := ctx.UserInfo()
	authUsername := userAuthInfo.Username
	shopID := userAuthInfo.ShopID

	input := ctx.ReadInput()

	userPwdReq := models.EmployeeRequestPassword{}
	err := json.Unmarshal([]byte(input), &userPwdReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	err = h.svc.UpdatePassword(shopID, authUsername, userPwdReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})

	return nil
}
