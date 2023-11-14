package slipimage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/slipimage/models"
	"smlcloudplatform/pkg/slipimage/repositories"
	"smlcloudplatform/pkg/slipimage/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/requestfilter"
	"time"
)

type ISlipImageHttp interface{}

type SlipImageHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.ISlipImageHttpService
}

func NewSlipImageHttp(ms *microservice.Microservice, cfg config.IConfig) SlipImageHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewSlipImageMongoRepository(pst)

	azureFileBlob := microservice.NewPersisterAzureBlob()
	imagePersister := microservice.NewPersisterImage(azureFileBlob)
	repoStorageImage := repositories.NewSlipImageStorageImageRepository(imagePersister)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewSlipImageHttpService(repo, repoStorageImage, masterSyncCacheRepo, 30*time.Second)

	return SlipImageHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h SlipImageHttp) RegisterHttp() {

	h.ms.GET("/slipimage", h.SearchSlipImagePage)
	h.ms.GET("/slipimage/list", h.SearchSlipImageStep)
	h.ms.POST("/slipimage", h.UploadSlipImage)
	h.ms.GET("/slipimage/:id", h.InfoSlipImage)
	h.ms.DELETE("/slipimage/:id", h.DeleteSlipImage)
	h.ms.DELETE("/slipimage", h.DeleteSlipImageByGUIDs)
}

// Create SlipImage godoc
// @Description Create SlipImage
// @Tags		SlipImage
// @Param		file  formData      file  true  "Image"
// @Param		docno  formData      string  true  "DocNo"
// @Param		posid  formData      string  true  "POS ID"
// @Param		docdate  formData      string  true  "Doc Date (yyyy-mm-dd)"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /slipimage [post]
func (h SlipImageHttp) UploadSlipImage(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	fileHeader, err := ctx.FormFile("file")

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "file is required")
		return err
	}

	docNo := ctx.FormValue("docno")
	posID := ctx.FormValue("posid")
	docDateRaw := ctx.FormValue("docdate")

	if docNo == "" {
		ctx.ResponseError(http.StatusBadRequest, "docno is required")
		return err
	}

	if posID == "" {
		ctx.ResponseError(http.StatusBadRequest, "posid is required")
		return err
	}

	if docDateRaw == "" {
		ctx.ResponseError(http.StatusBadRequest, "docdate is required")
		return err
	}

	layout := "2006-01-02"
	docDate, err := time.Parse(layout, docDateRaw)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, "invalid docdate. require format (yyyy-mm-dd)")
		return err
	}

	payload := models.SlipImageRequest{
		File:    fileHeader,
		DocNo:   docNo,
		PosID:   posID,
		DocDate: docDate,
	}

	data, err := h.svc.CreateSlipImage(shopID, authUsername, payload)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    data,
	})
	return nil
}

// Delete SlipImage godoc
// @Description Delete SlipImage
// @Tags		SlipImage
// @Param		id  path      string  true  "SlipImage ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /slipimage/{id} [delete]
func (h SlipImageHttp) DeleteSlipImage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteSlipImage(shopID, id, authUsername)

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

// Delete SlipImage godoc
// @Description Delete SlipImage
// @Tags		SlipImage
// @Param		SlipImage  body      []string  true  "SlipImage GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /slipimage [delete]
func (h SlipImageHttp) DeleteSlipImageByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteSlipImageByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get SlipImage godoc
// @Description get SlipImage info by guidfixed
// @Tags		SlipImage
// @Param		id  path      string  true  "SlipImage guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /slipimage/{id} [get]
func (h SlipImageHttp) InfoSlipImage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get SlipImage %v", id)
	doc, err := h.svc.InfoSlipImage(shopID, id)

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

// List SlipImage step godoc
// @Description get list step
// @Tags		SlipImage
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /slipimage [get]
func (h SlipImageHttp) SearchSlipImagePage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	docdate := ctx.QueryParam("docdate")

	fmt.Println("docdate", docdate)

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "posid",
			Field: "posid",
			Type:  requestfilter.FieldTypeString,
		},
		{
			Param: "docno",
			Field: "docno",
			Type:  requestfilter.FieldTypeString,
		},
		{
			Param: "docdate",
			Field: "docdate",
			Type:  requestfilter.FieldTypeDate,
		},
	})

	docList, pagination, err := h.svc.SearchSlipImage(shopID, filters, pageable)

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

// List SlipImage godoc
// @Description search limit offset
// @Tags		SlipImage
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /slipimage/list [get]
func (h SlipImageHttp) SearchSlipImageStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "posid",
			Field: "posid",
			Type:  requestfilter.FieldTypeString,
		},
		{
			Param: "docno",
			Field: "docno",
			Type:  requestfilter.FieldTypeString,
		},
		{
			Param: "docdate",
			Field: "docdate",
			Type:  requestfilter.FieldTypeDate,
		},
	})

	docList, total, err := h.svc.SearchSlipImageStep(shopID, lang, filters, pageableStep)

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
