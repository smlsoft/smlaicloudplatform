package sectiondepartment

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/productsection/sectiondepartment/models"
	"smlcloudplatform/pkg/productsection/sectiondepartment/repositories"
	"smlcloudplatform/pkg/productsection/sectiondepartment/services"
	"smlcloudplatform/pkg/utils"
)

type ISectionDepartmentHttp interface{}

type SectionDepartmentHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.ISectionDepartmentHttpService
}

func NewSectionDepartmentHttp(ms *microservice.Microservice, cfg microservice.IConfig) SectionDepartmentHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewSectionDepartmentRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewSectionDepartmentHttpService(repo, masterSyncCacheRepo)

	return SectionDepartmentHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h SectionDepartmentHttp) RouteSetup() {

	h.ms.POST("/product-section/deparment/bulk", h.SaveBulk)

	h.ms.GET("/product-section/deparment", h.SearchSectionDepartmentPage)
	h.ms.GET("/product-section/deparment/list", h.SearchSectionDepartmentStep)
	h.ms.GET("/product-section/deparment/:id", h.InfoSectionDepartment)
	h.ms.PUT("/product-section/deparment", h.SaveSectionDepartment)
	h.ms.DELETE("/product-section/deparment/:id", h.DeleteSectionDepartment)
	h.ms.DELETE("/product-section/deparment", h.DeleteSectionDepartmentByGUIDs)
}

// Save SectionDepartment godoc
// @Description Save SectionDepartment
// @Tags		SectionDepartment
// @Param		id  path      string  true  "SectionDepartment ID"
// @Param		SectionDepartment  body      models.SectionDepartment  true  "SectionDepartment"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product-section/deparment [put]
func (h SectionDepartmentHttp) SaveSectionDepartment(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	docReq := &models.SectionDepartment{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	id, err := h.svc.SaveSectionDepartment(shopID, authUsername, *docReq)

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

// Delete SectionDepartment godoc
// @Description Delete SectionDepartment
// @Tags		SectionDepartment
// @Param		id  path      string  true  "SectionDepartment ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product-section/deparment/{id} [delete]
func (h SectionDepartmentHttp) DeleteSectionDepartment(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteSectionDepartment(shopID, id, authUsername)

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

// Delete SectionDepartment godoc
// @Description Delete SectionDepartment
// @Tags		SectionDepartment
// @Param		SectionDepartment  body      []string  true  "SectionDepartment GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product-section/deparment [delete]
func (h SectionDepartmentHttp) DeleteSectionDepartmentByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteSectionDepartmentByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get SectionDepartment godoc
// @Description get SectionDepartment info by guidfixed
// @Tags		SectionDepartment
// @Param		id  path      string  true  "SectionDepartment guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product-section/deparment/{id} [get]
func (h SectionDepartmentHttp) InfoSectionDepartment(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get SectionDepartment %v", id)
	doc, err := h.svc.InfoSectionDepartment(shopID, id)

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

// List SectionDepartment step godoc
// @Description get list step
// @Tags		SectionDepartment
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product-section/deparment [get]
func (h SectionDepartmentHttp) SearchSectionDepartmentPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.svc.SearchSectionDepartment(shopID, map[string]interface{}{}, pageable)

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

// List SectionDepartment godoc
// @Description search limit offset
// @Tags		SectionDepartment
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product-section/deparment/list [get]
func (h SectionDepartmentHttp) SearchSectionDepartmentStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchSectionDepartmentStep(shopID, lang, pageableStep)

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

// Create SectionDepartment Bulk godoc
// @Description Create SectionDepartment
// @Tags		SectionDepartment
// @Param		SectionDepartment  body      []models.SectionDepartment  true  "SectionDepartment"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product-section/deparment/bulk [post]
func (h SectionDepartmentHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.SectionDepartment{}
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
		common.BulkReponse{
			Success:    true,
			BulkImport: bulkResponse,
		},
	)

	return nil
}
