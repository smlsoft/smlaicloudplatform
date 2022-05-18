package employee

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IEmployeeHttp interface{}

type EmployeeHttp struct {
	ms         *microservice.Microservice
	cfg        microservice.IConfig
	empService IEmployeeService
}

func NewEmployeeHttp(ms *microservice.Microservice, cfg microservice.IConfig) EmployeeHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	empRepo := NewEmployeeRepository(pst)
	empService := NewEmployeeService(empRepo)

	return EmployeeHttp{
		ms:         ms,
		cfg:        cfg,
		empService: empService,
	}
}

func (h EmployeeHttp) RouteSetup() {
	h.ms.POST("/employee/login", h.Login)
	h.ms.POST("/employee", h.Register)
	h.ms.GET("/employee", h.SearchEmployee)
	h.ms.PUT("/employee", h.Update)
	h.ms.PUT("/employee/password", h.UpdatePassword)
}

func (h EmployeeHttp) Login(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	userReq := &models.EmployeeRequestLogin{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	authToken, err := h.empService.Login(*userReq)

	if err != nil {
		ctx.ResponseError(400, "login failed.")
		return err
	}

	ctx.Response(http.StatusOK, models.AuthResponse{
		Success: true,
		Token:   authToken,
	})

	return nil
}

func (h EmployeeHttp) Register(ctx microservice.IContext) error {
	userAuthInfo := ctx.UserInfo()
	authUsername := userAuthInfo.Username
	shopID := userAuthInfo.ShopID
	input := ctx.ReadInput()

	userReq := models.Employee{}
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

func (h EmployeeHttp) Update(ctx microservice.IContext) error {
	userAuthInfo := ctx.UserInfo()
	authUsername := userAuthInfo.Username
	shopID := userAuthInfo.ShopID
	input := ctx.ReadInput()

	userReq := models.EmployeeRequestUpdate{}
	err := json.Unmarshal([]byte(input), &userReq)

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

func (h EmployeeHttp) SearchEmployee(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	docList, pagination, err := h.empService.List(shopID, q, page, limit)

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
