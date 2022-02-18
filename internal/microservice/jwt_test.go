package microservice

import (
	"smlcloudplatform/internal/microservice/models"
	"testing"
)

const (
	jwtKey = "946796991b2ece76900bfbc65612debc2e54554ef692cef6ec52181abe063d4d"
)

type TestCacherConfig struct{}

func (cfg *TestCacherConfig) Endpoint() string {
	return "127.0.0.1:6379"
}

func (cfg *TestCacherConfig) Password() string {
	return ""
}

func (cfg *TestCacherConfig) DB() int {
	return 0
}

func (cfg *TestCacherConfig) ConnectionSettings() ICacherConnectionSettings {
	return NewDefaultCacherConnectionSettings()
}
func TestGenerateToken(t *testing.T) {

	cacher := NewCacher(&TestCacherConfig{})

	jwtService := NewJwtService(cacher, jwtKey, 60*24*10)

	token, err := jwtService.GenerateToken(models.UserInfo{Username: "u001", Name: "My Name"})

	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(token) == 0 {
		t.Error("token is empty")
		return
	}

	t.Log(token)

}

func TestParseToken(t *testing.T) {

	cacher := NewCacher(&TestCacherConfig{})

	jwtService := NewJwtService(cacher, jwtKey, 60*24*10)

	tokenString, err := jwtService.GenerateToken(models.UserInfo{Username: "u001", Name: "My Name"})

	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(tokenString) == 0 {
		t.Error("token is empty")
		return
	}

	token, err := jwtService.ParseToken(tokenString)

	if err != nil {
		t.Error(err.Error())
		return
	}

	claims := token.Claims.(*CustomClaims)

	if claims.UserInfo.Username != "u001" {
		t.Error("username in token invalid")
		return
	}

	t.Log(claims.UserInfo.Name)
}

func TestParseTokenReal(t *testing.T) {

	cacher := NewCacher(&TestCacherConfig{})
	jwtService := NewJwtService(cacher, jwtKey, 60*24*10)

	tokenGen, err := jwtService.GenerateToken(models.UserInfo{Username: "u001", Name: "My Name"})

	if err != nil {
		t.Error(err.Error())
		return
	}

	token, err := jwtService.ParseToken(tokenGen)

	if err != nil {
		t.Error(err.Error())
		return
	}

	claims := token.Claims.(*CustomClaims)

	t.Log(claims.UserInfo.Name)
}

/*
func getKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	signBytes, err := ioutil.ReadFile("./../../private.key")

	if err != nil {
		return nil, nil, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err != nil {
		return nil, nil, err
	}

	verifyBytes, err := ioutil.ReadFile("./../../public.key")

	if err != nil {
		return nil, nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	if err != nil {
		return nil, nil, err
	}

	return signKey, verifyKey, nil
}

func TestGenerateToken(t *testing.T) {

	signKey, verifyKey, err := getKey()

	if err != nil {
		t.Error(err.Error())
		return
	}

	jwtService := NewJwtService(signKey, verifyKey, 60*24*10)

	token, err := jwtService.GenerateToken(models.UserInfo{Username: "u001", Name: "My Name"}, time.Duration(60*24*10)*time.Minute)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(token) == 0 {
		t.Error("token is empty")
		return
	}

	t.Log(token)

}

func TestParseToken(t *testing.T) {
	signKey, verifyKey, err := getKey()

	if err != nil {
		t.Error(err.Error())
		return
	}

	jwtService := NewJwtService(signKey, verifyKey, 60*24*10)

	tokenString, err := jwtService.GenerateToken(models.UserInfo{Username: "u001", Name: "My Name"}, time.Duration(60*24*10)*time.Minute)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(tokenString) == 0 {
		t.Error("token is empty")
		return
	}

	token, err := jwtService.ParseToken(tokenString)

	if err != nil {
		t.Error(err.Error())
		return
	}

	claims := token.Claims.(*CustomClaims)

	if claims.UserInfo.Username != "u001" {
		t.Error("username in token invalid")
		return
	}

	t.Log(claims.UserInfo.Name)
}

func TestParseTokenReal(t *testing.T) {
	signKey, verifyKey, err := getKey()

	if err != nil {
		t.Error(err.Error())
		return
	}

	jwtService := NewJwtService(signKey, verifyKey, 60*24*10)

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyODIyMjUsInVzZXJuYW1lIjoidTAwMSIsIm5hbWUiOiJNeSBOYW1lIn0.XGdnkpkr0sHo5PBGJUn-sA4pSwqsJ86B8i6lN-EiD_A"

	token, err := jwtService.ParseToken(tokenString)

	if err != nil {
		t.Error(err.Error())
		return
	}

	claims := token.Claims.(*CustomClaims)

	t.Log(claims.UserInfo.Name)
}
*/
