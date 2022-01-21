package microservice

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	*jwt.RegisteredClaims
	models.UserInfo
}

func NewJwtService(signKey *rsa.PrivateKey, verifyKey *rsa.PublicKey, expireMinute int) *JwtService {

	return &JwtService{
		signKey:   signKey,
		verifyKey: verifyKey,
		duration:  time.Duration(expireMinute) * time.Minute,
	}
}

// Service provides a Json-Web-Token authentication implementation
type JwtService struct {
	// Secret key used for signing.
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey

	// Duration for which the jwt token is valid.
	duration time.Duration
}

// MWFunc makes JWT implement the Middleware interface.
func (jwtService *JwtService) MWFunc() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := jwtService.ParseTokenFromContext(c)

			if err != nil || !token.Valid {

				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "jwt: " + err.Error()})
			}

			claims := token.Claims.(*CustomClaims)

			c.Set("UserInfo", claims.UserInfo)

			return next(c)
		}
	}
}

// ParseToken parses token from Authorization header
func (jwtService *JwtService) ParseTokenFromContext(c echo.Context) (*jwt.Token, error) {

	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	parts := strings.SplitN(tokenString, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, fmt.Errorf("missing authorization bearer")
	}

	return jwtService.ParseToken(parts[1])
}

func (jwtService *JwtService) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return jwtService.verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil

	// token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("35WPe96UooYISzbsn5Q5QsPax1mzc3F0kkkSp7ecVGNlHR0Lbf7jpVYvunv9xbasLBJRzCG5PuFMeHIKoJ3P7EwF4PQWfnxbdWWF"), nil
	// })

	// if err != nil {
	// 	return nil, err
	// }

	// return token, nil
}

func (jwtService *JwtService) GenerateToken(userInfo models.UserInfo, expire time.Duration) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// set claims
	t.Claims = &CustomClaims{
		&jwt.RegisteredClaims{
			// set the expire time
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
		userInfo,
	}

	// Creat token string
	return t.SignedString(jwtService.signKey)

	// claims := &CustomClaims{
	// 	&jwt.RegisteredClaims{
	// 		// set the expire time
	// 		// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
	// 	},
	// 	userInfo,
	// }

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// return token.SignedString([]byte("35WPe96UooYISzbsn5Q5QsPax1mzc3F0kkkSp7ecVGNlHR0Lbf7jpVYvunv9xbasLBJRzCG5PuFMeHIKoJ3P7EwF4PQWfnxbdWWF"))
}
