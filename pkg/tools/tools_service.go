package tools

import (
	"os"
	"smlcloudplatform/internal/microservice"
)

type ToolsService struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
}

func NewToolsService(ms *microservice.Microservice, cfg microservice.IConfig) *ToolsService {

	return &ToolsService{
		ms:  ms,
		cfg: cfg,
	}
}

func (svc *ToolsService) RouteSetup() {
	svc.ms.GET("tool/mongo", svc.CheckMongodbConnect)
	svc.ms.GET("tool/env", svc.AllEnv)
}

func (svc *ToolsService) CheckMongodbConnect(ctx microservice.IContext) error {
	mongoCfg := svc.cfg.MongoPersisterConfig()
	repStr := mongoCfg.MongodbURI()

	repStr = repStr + " \n " + mongoCfg.DB()

	ctx.ResponseS(200, repStr)

	return nil
}

func (svc *ToolsService) AllEnv(ctx microservice.IContext) error {

	ctx.Response(200, os.Environ())

	return nil
}
