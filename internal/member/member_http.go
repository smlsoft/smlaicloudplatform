package member

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/config"
	mastersync "smlcloudplatform/internal/mastersync/repositories"
	"smlcloudplatform/internal/member/models"
	common "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/shop"
	"smlcloudplatform/internal/utils"
	"smlcloudplatform/pkg/microservice"
	"time"
)

type IMemberHttp interface {
	RegisterHttp()
	CreateMember(ctx microservice.IContext) error
	UpdateMember(ctx microservice.IContext) error
	DeleteMember(ctx microservice.IContext) error
	InfoMember(ctx microservice.IContext) error
	SearchMember(ctx microservice.IContext) error
}

type MemberHttp struct {
	ms      *microservice.Microservice
	cfg     config.IConfig
	service IMemberService
}

func NewMemberHttp(ms *microservice.Microservice, cfg config.IConfig) MemberHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstPg := ms.Persister(cfg.PersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	memberRepo := NewMemberRepository(pst)
	memberPgRepo := NewMemberPGRepository(pstPg)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)

	shopRepo := shop.NewShopRepository(pst)
	shopUserRepo := shop.NewShopUserRepository(pst)
	shopService := shop.NewShopService(shopRepo, shopUserRepo, utils.NewGUID, ms.TimeNow)

	authService := microservice.NewAuthServicePrefix("linemember:", "linememberrefresh:", ms.Cacher(cfg.CacherConfig()), 24*3*time.Hour, 24*30*time.Hour)

	service := NewMemberService(memberRepo, memberPgRepo, shopService, authService, masterSyncCacheRepo)

	return MemberHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h MemberHttp) RegisterLineHttp() {
	h.ms.POST("/member/line", h.MemberAuthLine)
	h.ms.GET("/member/line/profile", h.LineProfileInfo)
	h.ms.PUT("/member/line/profile", h.UpdateMemberProfileWithLine)
}

func (h MemberHttp) RegisterHttp() {

	h.ms.GET("/member/:id", h.InfoMember)
	h.ms.GET("/member", h.SearchMember)

	h.ms.POST("/member", h.CreateMember)
	h.ms.PUT("/member/:id", h.UpdateMember)
	h.ms.DELETE("/member/:id", h.DeleteMember)
}

// Auth Line Member godoc
// @Description Auth Line Member
// @Tags		Member
// @Param		LineAuthRequest  body      models.LineAuthRequest  true  "Line Auth Request"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /member/line [post]
func (h MemberHttp) MemberAuthLine(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	payload := models.LineAuthRequest{}
	err := json.Unmarshal([]byte(input), &payload)

	if err != nil {
		ctx.ResponseError(400, "payload invalid")
		return err
	}

	idx, err := h.service.AuthWithLine(payload)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

// Update Line Member godoc
// @Description Updat Line Member
// @Tags		Member
// @Param		Member  body      models.Member  true  "Member"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /member/profile [put]
func (h MemberHttp) UpdateMemberProfileWithLine(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	lineUID := userInfo.Username

	input := ctx.ReadInput()

	docReq := models.Member{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateProfileWithLine(shopID, lineUID, docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Get Member Line Profile Infomation godoc
// @Description Get Member Line Profile
// @Tags		Member
// @Accept 		json
// @Success		200	{object}	models.MemberInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /member/line/profile [get]
func (h MemberHttp) LineProfileInfo(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	lineUID := userInfo.Username

	doc, err := h.service.LineProfileInfo(shopID, lineUID)

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

	doc := models.Member{}
	err := json.Unmarshal([]byte(input), &doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.Create(shopID, authUsername, doc)

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

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.service.Search(shopID, pageable)

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
