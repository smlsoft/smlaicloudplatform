package member

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
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

	memberRepo := NewMemberRepository(pst)
	memberPgRepo := NewMemberPGRepository(pstPg)

	service := NewMemberService(memberRepo, memberPgRepo)

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
}

// Create Member godoc
// @Description Create Member
// @Tags		Member
// @Param		Member  body      models.Member  true  "Member"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
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

	idx, err := h.service.CreateMember(shopID, authUsername, *doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
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
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
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

	err = h.service.UpdateMember(id, shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete Member godoc
// @Description Delete Member
// @Tags		Member
// @Param		id  path      string  true  "Member ID"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
// @Security     AccessToken
// @Router /member/{id} [delete]
func (h MemberHttp) DeleteMember(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.DeleteMember(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

// Get Member Infomation godoc
// @Description Get Member Category
// @Tags		Member
// @Param		id  path      string  true  "Member Id"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
// @Security     AccessToken
// @Router /member/{id} [get]
func (h MemberHttp) InfoMember(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.service.InfoMember(id, shopID)

	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
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
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
// @Security     AccessToken
// @Router /member [get]
func (h MemberHttp) SearchMember(ctx microservice.IContext) error {

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

	docList, pagination, err := h.service.SearchMember(shopID, q, page, limit)

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
