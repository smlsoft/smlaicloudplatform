package company

import (
	"encoding/json"
	"net/http"
	"smlaicloudplatform/internal/config"
	mastersync "smlaicloudplatform/internal/mastersync/repositories"
	common "smlaicloudplatform/internal/models"
	businessTypeRepositories "smlaicloudplatform/internal/organization/businesstype/repositories"
	"smlaicloudplatform/internal/organization/company/models"
	"smlaicloudplatform/internal/organization/company/repositories"
	"smlaicloudplatform/internal/organization/company/services"
	deparmentRepositories "smlaicloudplatform/internal/organization/department/repositories"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/internal/utils/requestfilter"
	"smlaicloudplatform/pkg/microservice"
)

type ICompanyHttp interface{}

type CompanyHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.ICompanyHttpService
}

func NewCompanyHttp(ms *microservice.Microservice, cfg config.IConfig) CompanyHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewCompanyRepository(pst)

	repoDepartment := deparmentRepositories.NewDepartmentRepository(pst)
	repoBusinessType := businessTypeRepositories.NewBusinessTypeRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewCompanyHttpService(repo, repoDepartment, repoBusinessType, masterSyncCacheRepo)

	return CompanyHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h CompanyHttp) RegisterHttp() {

	h.ms.POST("/organization/company/bulk", h.SaveBulk)

	h.ms.GET("/organization/company", h.SearchCompanyPage)
	h.ms.GET("/organization/company/list", h.SearchCompanyStep)
	h.ms.POST("/organization/company", h.CreateCompany)
	h.ms.GET("/organization/company/:id", h.InfoCompany)
	h.ms.GET("/organization/company/code/:code", h.InfoCompanyByCode)
	h.ms.PUT("/organization/company/:id", h.UpdateCompany)
	h.ms.DELETE("/organization/company/:id", h.DeleteCompany)
	h.ms.DELETE("/organization/company", h.DeleteCompanyByGUIDs)
}

// Create Company godoc
// @Description Create Company
// @Tags		Company
// @Param		Company  body      models.Company  true  "Company"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/company [post]
func (h CompanyHttp) CreateCompany(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.Company{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateCompany(shopID, authUsername, *docReq)

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

// Update Company godoc
// @Description Update Company
// @Tags		Company
// @Param		id  path      string  true  "Company ID"
// @Param		Company  body      models.Company  true  "Company"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/company/{id} [put]
func (h CompanyHttp) UpdateCompany(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Company{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateCompany(shopID, id, authUsername, *docReq)

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

// Delete Company godoc
// @Description Delete Company
// @Tags		Company
// @Param		id  path      string  true  "Company ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/company/{id} [delete]
func (h CompanyHttp) DeleteCompany(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteCompany(shopID, id, authUsername)

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

// Delete Company godoc
// @Description Delete Company
// @Tags		Company
// @Param		Company  body      []string  true  "Company GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/company [delete]
func (h CompanyHttp) DeleteCompanyByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteCompanyByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get Company godoc
// @Description get Company info by guidfixed
// @Tags		Company
// @Param		id  path      string  true  "Company guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/company/{id} [get]
func (h CompanyHttp) InfoCompany(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Company %v", id)
	doc, err := h.svc.InfoCompany(shopID, id)

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

// Get Company By Code godoc
// @Description get Company info by Code
// @Tags		Company
// @Param		code  path      string  true  "Company Code"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/company/code/{code} [get]
func (h CompanyHttp) InfoCompanyByCode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	code := ctx.Param("code")

	doc, err := h.svc.InfoCompanyByCode(shopID, code)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// List Company step godoc
// @Description get list step
// @Tags		Company
// @Param		q		query	string		false  "Search Value"
// @Param		businesstypecode		query	string		false  "business type code"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/company [get]
func (h CompanyHttp) SearchCompanyPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "businesstypecode",
			Field: "businesstype.code",
			Type:  requestfilter.FieldTypeString,
		},
	})

	docList, pagination, err := h.svc.SearchCompany(shopID, filters, pageable)

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

// List Company godoc
// @Description search limit offset
// @Tags		Company
// @Param		q		query	string		false  "Search Value"
// @Param		businesstypecode		query	string		false  "business type code"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/company/list [get]
func (h CompanyHttp) SearchCompanyStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "businesstypecode",
			Field: "businesstype.code",
			Type:  requestfilter.FieldTypeString,
		},
	})

	docList, total, err := h.svc.SearchCompanyStep(shopID, lang, filters, pageableStep)

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

// Create Company Bulk godoc
// @Description Create Company
// @Tags		Company
// @Param		Company  body      []models.Company  true  "Company"
// @Accept 		json
// @Success		201	{object}	common.BulkResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/company/bulk [post]
func (h CompanyHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Company{}
	err := json.Unmarshal([]byte(input), &dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	bulkResponse, err := h.svc.SaveInBatch(shopID, authUsername, dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		common.BulkResponse{
			Success:    true,
			BulkImport: bulkResponse,
		},
	)

	return nil
}
