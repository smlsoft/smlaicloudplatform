package employee

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
)

type IEmployeeHttp interface{}

type EmployeeHttp struct {
	ms         *microservice.Microservice
	cfg        microservice.IConfig
	empService IEmployeeService
}

func NewEmployeeHttp(ms *microservice.Microservice, cfg microservice.IConfig) EmployeeHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	empRepo := NewEmployeeRepository(pst)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	empService := NewEmployeeService(empRepo, masterSyncCacheRepo)

	return EmployeeHttp{
		ms:         ms,
		cfg:        cfg,
		empService: empService,
	}
}

func (h EmployeeHttp) RouteSetup() {
	// h.ms.POST("/employee/login", h.Login)
	h.ms.POST("/employee", h.Register)
	h.ms.GET("/employee/:username", h.InfoEmployee)
	h.ms.GET("/employee", h.SearchEmployee)
	h.ms.PUT("/employee/:username", h.Update)
	h.ms.PUT("/employee/password", h.UpdatePassword)
}

// Validate Employee godoc
// @Description Validate Employee
// @Tags		Employee
// @Param		EmployeeUserPassword  body      models.EmployeeRequestLogin  true  "EmployeeUserPassword"
// @Accept 		json
// @Success		201	{object}	models.EmployeeInfo
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /employee/login [post]
func (h EmployeeHttp) Login(ctx microservice.IContext) error {
	input := ctx.ReadInput()

	userReq := &models.EmployeeRequestLogin{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	employee, err := h.empService.Login(userReq.ShopID, *userReq)

	if err != nil {
		ctx.ResponseError(400, "login failed.")
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    employee,
	})

	return nil
}

// Create Employee godoc
// @Summary		Create Employee
// @Description	For Create Employee
// @Tags		Employee
// @Param		User  body      models.Employee  true  "Register Employee"
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		400 {object}	models.AuthResponseFailed
// @Accept 		json
// @Security     AccessToken
// @Router		/employee [post]
func (h EmployeeHttp) Register(ctx microservice.IContext) error {
	userAuthInfo := ctx.UserInfo()
	authUsername := userAuthInfo.Username
	shopID := userAuthInfo.ShopID
	input := ctx.ReadInput()

	userReq := models.EmployeeRequestRegister{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	idx, err := h.empService.Register(shopID, authUsername, userReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

// Update Employee godoc
// @Summary		Update Employee
// @Description	Update Employee
// @Tags		Employee
// @Param		username  path      string  true  "Employee username"
// @Param		Employee  body      models.EmployeeRequestUpdate  true  "Employee"
// @Success		200	{object}	models.ResponseSuccess
// @Failure		400 {object}	models.AuthResponseFailed
// @Accept 		json
// @Security     AccessToken
// @Router		/employee/{username} [put]
func (h EmployeeHttp) Update(ctx microservice.IContext) error {
	userAuthInfo := ctx.UserInfo()
	authUsername := userAuthInfo.Username
	shopID := userAuthInfo.ShopID
	input := ctx.ReadInput()

	username := ctx.Param("username")

	userReq := models.EmployeeRequestUpdate{}
	err := json.Unmarshal([]byte(input), &userReq)

	userReq.Username = username

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	err = h.empService.Update(shopID, authUsername, userReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
	})

	return nil
}

// Register Employee godoc
// @Summary		Register An Account
// @Description	For User Register Application
// @Tags		Employee
// @Param		id  path      string  true  "Employee ID"
// @Param		Employee  body      models.EmployeeRequestPassword  true  "Register Employee"
// @Success		200	{object}	models.ResponseSuccess
// @Failure		400 {object}	models.AuthResponseFailed
// @Accept 		json
// @Security     AccessToken
// @Router		/employee/password/{id} [put]
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

	err = h.empService.UpdatePassword(shopID, authUsername, userPwdReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
	})

	return nil
}

// Info Employee godoc
// @Description List Employee
// @Tags		Employee
// @Accept 		json
// @Success		200	{array}	models.EmployeePageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /employee/{username} [get]
func (h EmployeeHttp) InfoEmployee(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	username := ctx.Param("username")

	docList, err := h.empService.Get(shopID, username)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success: true,
			Data:    docList,
		})

	return nil
}

// List Employee godoc
// @Description List Employee
// @Tags		Employee
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{array}	models.EmployeePageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /employee [get]
func (h EmployeeHttp) SearchEmployee(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.empService.List(shopID, pageable)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})

	return nil
}
