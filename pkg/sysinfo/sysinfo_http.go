package sysinfo

import (
	"net/http"
	"smlcloudplatform/internal/microservice"
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
	ctx.Response(http.StatusOK, map[string]interface{}{
		"version": "1.0.1",
	})
	return nil
}
