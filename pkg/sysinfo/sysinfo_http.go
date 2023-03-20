package sysinfo

import (
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/product/models"
	"smlcloudplatform/pkg/repositories"

	m "github.com/veer66/mapkha"
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
	h.ms.GET("/sys-info/wordcut", h.Wordcut)
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

func (h SysInfoHttp) Wordcut(ctx microservice.IContext) error {

	wordCut := &m.Wordcut{}

	dict, err := m.LoadDefaultDict()
	if err != nil {
		ctx.Response(http.StatusOK, map[string]interface{}{
			"error": err.Error(),
		})
	}

	wordCut = m.NewWordcut(dict)

	results := wordCut.Segment("ขาAC001ขา")
	ctx.Response(http.StatusOK, results)

	return nil
}
