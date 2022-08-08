package microservice

import (
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice/models"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func NewAuthService(cacher ICacher, expireHour int) *AuthService {

	return &AuthService{
		cacher:              cacher,
		expire:              time.Duration(expireHour) * time.Hour,
		prefixCacheKey:      "auth-",
		prefixAuthorization: "Bearer",
	}
}

type AuthService struct {
	cacher              ICacher
	expire              time.Duration
	prefixCacheKey      string
	prefixAuthorization string
}

func (authService *AuthService) MWFuncWithRedisMixShop(cacher ICacher, shopPath []string, publicPath ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			currentPath := c.Path()

			for _, publicPath := range publicPath {
				if strings.HasPrefix(currentPath, publicPath) {
					return next(c)
				} else if currentPath == publicPath {
					return next(c)
				}

			}

			tokenStr, err := authService.GetTokenFromContext(c)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "Token Invalid."})
			}

			cacheKey := authService.prefixCacheKey + tokenStr
			tempUserInfo, err := authService.cacher.HMGet(cacheKey, []string{"username", "name", "shopid", "role"})

			if err != nil || tempUserInfo[0] == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "Token Invalid."})
			}

			tempShopID := ""

			if tempUserInfo[2] != nil {
				tempShopID = fmt.Sprintf("%v", tempUserInfo[2])
			}

			// check accept shop path
			thisPathExceptShopSelected := false
			for _, publicPath := range shopPath {
				if currentPath == publicPath {
					thisPathExceptShopSelected = true
				}
			}

			if !thisPathExceptShopSelected && len(string(tempShopID)) < 1 {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "Shop not selected."})
			}

			userRole, err := strconv.Atoi(fmt.Sprintf("%v", tempUserInfo[3]))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": fmt.Sprintf("Mix Shop: User role invalid. %v", tempUserInfo[3])})
			}

			userInfo := models.UserInfo{
				Username: fmt.Sprintf("%v", tempUserInfo[0]),
				Name:     fmt.Sprintf("%v", tempUserInfo[1]),
				ShopID:   fmt.Sprintf("%v", tempUserInfo[2]),
				Role:     uint8(userRole),
			}

			cacher.Expire(cacheKey, authService.expire)
			c.Set("UserInfo", userInfo)

			return next(c)
		}
	}
}

func (authService *AuthService) MWFuncWithRedis(cacher ICacher, publicPath ...string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			currentPath := c.Path()

			for _, publicPath := range publicPath {
				if strings.HasPrefix(currentPath, publicPath) {
					return next(c)
				} else if currentPath == publicPath {
					return next(c)
				}

			}

			tokenStr, err := authService.GetTokenFromContext(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "Token Invalid."})
			}

			cacheKey := authService.prefixCacheKey + tokenStr
			tempUserInfo, err := authService.cacher.HMGet(cacheKey, []string{"username", "name", "shopid", "role"})

			if err != nil || tempUserInfo[0] == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "Token Invalid."})
			}

			tempShopID := ""

			if tempUserInfo[2] != nil {
				tempShopID = fmt.Sprintf("%v", tempUserInfo[2])
			}

			if len(string(tempShopID)) < 1 {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "Shop not selected."})
			}

			userRole, err := strconv.ParseUint(fmt.Sprintf("%v", tempUserInfo[3]), 10, 8)

			if err != nil {
				fmt.Println(err)
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": fmt.Sprintf("User role invalid. %v", tempUserInfo[3])})
			}

			userInfo := models.UserInfo{
				Username: fmt.Sprintf("%v", tempUserInfo[0]),
				Name:     fmt.Sprintf("%v", tempUserInfo[1]),
				ShopID:   fmt.Sprintf("%v", tempUserInfo[2]),
				Role:     uint8(userRole),
			}

			cacher.Expire(cacheKey, authService.expire)
			c.Set("UserInfo", userInfo)

			return next(c)
		}
	}
}

func (authService *AuthService) MWFuncWithShop(cacher ICacher, publicPath ...string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			currentPath := c.Path()

			for _, publicPath := range publicPath {
				if currentPath == publicPath {
					return next(c)
				}
			}

			tokenStr, err := authService.GetTokenFromContext(c)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "Token Invalid."})
			}

			cacheKey := authService.prefixCacheKey + tokenStr

			tempUserInfo, err := authService.cacher.HMGet(cacheKey, []string{"username", "name"})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "Token Invalid."})
			}

			if tempUserInfo[0] == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "Token Invalid."})
			}

			userInfo := models.UserInfo{
				Username: fmt.Sprintf("%v", tempUserInfo[0]),
				Name:     fmt.Sprintf("%v", tempUserInfo[1]),
			}

			c.Set("UserInfo", userInfo)

			return next(c)
		}
	}
}

func (authService *AuthService) GetPrefixCacheKey() string {
	return authService.prefixCacheKey
}

func (authService *AuthService) GetTokenFromContext(c echo.Context) (string, error) {

	// socket
	if c.IsWebSocket() {
		return c.QueryParam("apikey"), nil
	} else {

		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return "", fmt.Errorf("missing authorization header")
		}

		parts := strings.SplitN(tokenString, " ", 2)
		if !(len(parts) == 2 && parts[0] == authService.prefixAuthorization) {
			return "", fmt.Errorf("missing authorization bearer")
		}

		return parts[1], nil
	}
}

func (authService *AuthService) GetTokenFromAuthorizationHeader(tokenAuthorization string) (string, error) {

	parts := strings.SplitN(tokenAuthorization, " ", 2)
	if !(len(parts) == 2 && parts[0] == authService.prefixAuthorization) {
		return "", fmt.Errorf("missing authorization bearer")
	}

	return parts[1], nil
}

func (authService *AuthService) GenerateTokenWithRedis(userInfo models.UserInfo) (string, error) {

	tokenStr := NewUUID()

	cacheKey := authService.prefixCacheKey + tokenStr
	authService.cacher.HMSet(cacheKey, map[string]interface{}{
		"username": userInfo.Username,
		"name":     userInfo.Name,
	})
	authService.cacher.Expire(cacheKey, authService.expire)

	return tokenStr, nil
}

func (authService *AuthService) SelectShop(tokenStr string, shopID string, role string) error {
	cacheKey := authService.prefixCacheKey + tokenStr
	err := authService.cacher.HMSet(cacheKey, map[string]interface{}{
		"shopid": shopID,
		"role":   role,
	})

	if err != nil {
		return err
	}

	return nil

}

func (authService *AuthService) ExpireToken(tokenAuthorizationHeader string) error {
	tokenStr, err := authService.GetTokenFromAuthorizationHeader(tokenAuthorizationHeader)
	if err != nil {
		return err
	}
	cacheKey := authService.prefixCacheKey + tokenStr
	authService.cacher.Expire(cacheKey, -1)
	return nil
}
