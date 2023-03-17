package sysinfo

import (
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/product/models"
	"smlcloudplatform/pkg/repositories"
)

type SysInfoHttp struct {
	ms *microservice.Microservice
}

func NewSysInfoHttp(ms *microservice.Microservice, cfg microservice.IConfig) SysInfoHttp {

	return SysInfoHttp{
		ms: ms,
	}
}

func (h SysInfoHttp) RouteSetup() {

	h.ms.GET("/sys-info/version", h.Version)
}

func (h SysInfoHttp) Version(ctx microservice.IContext) error {

	repo := repositories.NewSearchRepository[models.ProductInfo](nil)

	txt := repo.CreateTextFilter([]string{"names.name"}, "A3331")

	ctx.Response(http.StatusOK, map[string]interface{}{
		"version": "1.0.1",
		"q":       fmt.Sprintf("%v", txt),
	})
	return nil
}
