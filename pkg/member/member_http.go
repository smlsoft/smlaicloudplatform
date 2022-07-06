package member

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/member/models"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"
)

type IMemberHttp interface {
	RouteSetup()
	CreateMember(ctx microservice.IContext) error
	UpdateMember(ctx microservice.IContext) error
	DeleteMember(ctx microservice.IContext) error
	InfoMember(ctx microservice.IContext) error
	SearchMember(ctx microservice.IContext) error
}

type MemberHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service IMemberService
}

func NewMemberHttp(ms *microservice.Microservice, cfg microservice.IConfig) MemberHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstPg := ms.Persister(cfg.PersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	memberRepo := NewMemberRepository(pst)
	memberPgRepo := NewMemberPGRepository(pstPg)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache, "member")
	service := NewMemberService(memberRepo, memberPgRepo, masterSyncCacheRepo)

	return MemberHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h MemberHttp) RouteSetup() {

	h.ms.GET("/member/:id", h.InfoMember)
	h.ms.GET("/member", h.SearchMember)

	h.ms.POST("/member", h.CreateMember)
	h.ms.PUT("/member/:id", h.UpdateMember)
	h.ms.DELETE("/member/:id", h.DeleteMember)
	h.ms.GET("/member/fetchupdate", h.LastActivityMember)
}

// Create Member godoc
// @Description Create Member
// @Tags		Member
// @Param		Member  body      models.Member  true  "Member"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /member [post]
func (h MemberHttp) CreateMember(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	doc := &models.Member{}
	err := json.Unmarshal([]byte(input), &doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.Create(shopID, authUsername, *doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

// Update Member godoc
// @Description Update Member
// @Tags		Member
// @Param		id  path      string  true  "Member ID"
// @Param		Member  body      models.Member  true  "Member"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /member/{id} [put]
func (h MemberHttp) UpdateMember(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Member{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.Update(shopID, id, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete Member godoc
// @Description Delete Member
// @Tags		Member
// @Param		id  path      string  true  "Member ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /member/{id} [delete]
func (h MemberHttp) DeleteMember(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.Delete(shopID, id, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Get Member Infomation godoc
// @Description Get Member
// @Tags		Member
// @Param		id  path      string  true  "Member Id"
// @Accept 		json
// @Success		200	{object}	models.MemberInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /member/{id} [get]
func (h MemberHttp) InfoMember(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.service.Info(shopID, id)

	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if len(doc.GuidFixed) == 0 {
		ctx.Response(http.StatusNotFound, common.ApiResponse{
			Success: false,
			Message: "document not found",
		})
		return nil
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// List Member godoc
// @Description List Member Category
// @Tags		Member
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.MemberPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /member [get]
func (h MemberHttp) SearchMember(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.service.Search(shopID, q, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})

	return nil
}

// Fetch Update Member By Date godoc
// @Description Fetch Update Member By Date
// @Tags		Member
// @Param		lastUpdate query string true "DateTime"
// @Accept		json
// @Success		200 {object} models.MemberFetchUpdateResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/member/fetchupdate [get]
func (h MemberHttp) LastActivityMember(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02T15:04" //
	lastUpdateStr := ctx.QueryParam("lastUpdate")

	lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(lastUpdateStr) < 1 {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return nil
	}

	lastUpdate, err := time.Parse(layout, lastUpdateStr)

	if err != nil {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return err
	}

	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.service.LastActivity(shopID, lastUpdate, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})

	return nil
}
