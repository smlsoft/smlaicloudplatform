package services

import (
	"smlaicloudplatform/pkg/microservice"
	micromodel "smlaicloudplatform/pkg/microservice/models"
	"time"
)

type IApiKeyServiceHttpService interface {
	CreateApiKeyService(userInfo micromodel.UserInfo, expireTime time.Duration) (string, error)
	RemoveApiKey(authorizationHeader string) error
}

type ApiKeyServiceHttpService struct {
	// repoCache repositories.IApiKeyServiceCacheRepository
	authService microservice.IAuthService
}

func NewApiKeyServiceHttpService(authService microservice.IAuthService) *ApiKeyServiceHttpService {

	return &ApiKeyServiceHttpService{
		authService: authService,
	}
}

func (svc ApiKeyServiceHttpService) CreateApiKeyService(userInfo micromodel.UserInfo, expireTime time.Duration) (string, error) {
	return svc.authService.GenerateTokenWithRedisExpire(microservice.AUTHTYPE_XAPIKEY, userInfo, expireTime)
}

func (svc ApiKeyServiceHttpService) RemoveApiKey(authorizationHeader string) error {
	return svc.authService.ExpireToken(microservice.AUTHTYPE_XAPIKEY, authorizationHeader)
}
