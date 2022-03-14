package shop

import (
	"smlcloudplatform/internal/microservice"
)

type IShopMemberHttp interface{}

type ShopMemberHttp struct {
	ms  *microservice.Microservice
	svc IShopUserService
}

func NewShopMemberHttp(ms *microservice.Microservice, cfg microservice.IConfig) IShopMemberHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	repo := NewShopUserRepository(pst)
	svc := NewShopUserService(repo)
	return &ShopMemberHttp{
		svc: svc,
		ms:  ms,
	}
}

func (h *ShopMemberHttp) RouteSetup() {

}
