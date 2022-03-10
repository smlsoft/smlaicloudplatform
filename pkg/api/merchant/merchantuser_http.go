package merchant

import (
	"smlcloudplatform/internal/microservice"
)

type IMerchantMemberHttp interface{}

type MerchantMemberHttp struct {
	ms  *microservice.Microservice
	svc IMerchantUserService
}

func NewMerchantMemberHttp(ms *microservice.Microservice, cfg microservice.IConfig) IMerchantMemberHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	repo := NewMerchantUserRepository(pst)
	svc := NewMerchantUserService(repo)
	return &MerchantMemberHttp{
		svc: svc,
		ms:  ms,
	}
}

func (h *MerchantMemberHttp) RouteSetup() {

}
