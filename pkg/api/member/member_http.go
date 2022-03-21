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

func NewMemberHttp(ms *microservice.Microservice, cfg microservice.IConfig) IMemberHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	memberRepo := NewMemberRepository(pst)

	service := NewMemberService(memberRepo)

	return &MemberHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h *MemberHttp) RouteSetup() {

	h.ms.GET("/member/:id", h.InfoMember)
	h.ms.GET("/member", h.SearchMember)

	h.ms.POST("/member", h.CreateMember)
	h.ms.PUT("/member/:id", h.UpdateMember)
	h.ms.DELETE("/member/:id", h.DeleteMember)
}

func (h *MemberHttp) CreateMember(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopId := userInfo.ShopId

	input := ctx.ReadInput()

	doc := &models.Member{}
	err := json.Unmarshal([]byte(input), &doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreateMember(shopId, authUsername, *doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		Id:      idx,
	})

	return nil
}

func (h *MemberHttp) UpdateMember(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopId := userInfo.ShopId

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Member{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateMember(id, shopId, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h *MemberHttp) DeleteMember(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopId := userInfo.ShopId

	id := ctx.Param("id")

	err := h.service.DeleteMember(id, shopId, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h *MemberHttp) InfoMember(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopId := userInfo.ShopId

	id := ctx.Param("id")

	doc, err := h.service.InfoMember(id, shopId)

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

func (h *MemberHttp) SearchMember(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopId := userInfo.ShopId

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	docList, pagination, err := h.service.SearchMember(shopId, q, page, limit)

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
