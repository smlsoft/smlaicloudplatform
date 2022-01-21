package tools

import "smlcloudplatform/internal/microservice"

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
