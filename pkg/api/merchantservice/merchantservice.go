package merchantservice

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MerchantService struct {
	ms         *microservice.Microservice
	cfg        microservice.IConfig
	jwtService *microservice.JwtService
	signKey    *rsa.PrivateKey
	verifyKey  *rsa.PublicKey
}

func NewMerchantService(ms *microservice.Microservice, cfg microservice.IConfig) *MerchantService {
	signBytes, err := ioutil.ReadFile("./../../private.key")

	if err != nil {
		ms.Log("auth", err.Error())
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err != nil {
		ms.Log("auth", err.Error())
	}

	verifyBytes, err := ioutil.ReadFile("./../../public.key")

	if err != nil {
		ms.Log("auth", err.Error())
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	if err != nil {
		ms.Log("auth", err.Error())
	}

	jwtService := microservice.NewJwtService(signKey, verifyKey, 60*24*10)

	return &MerchantService{
		ms:         ms,
		cfg:        cfg,
		jwtService: jwtService,
		signKey:    signKey,
		verifyKey:  verifyKey,
	}
}

func (svc *MerchantService) RouteSetup() {

	svc.ms.GET("/merchant", svc.SearchMerchant, svc.jwtService.MWFunc())
	svc.ms.POST("/merchant", svc.CreateMerchant, svc.jwtService.MWFunc())
	svc.ms.GET("/merchant/:id", svc.InfoMerchant, svc.jwtService.MWFunc())
	svc.ms.PUT("/merchant/:id", svc.EditMerchant, svc.jwtService.MWFunc())
	svc.ms.DELETE("/merchant/:id", svc.DeleteMerchant, svc.jwtService.MWFunc())

	svc.ms.GET("/merchant/:merchant_id/member", svc.SearchMember, svc.jwtService.MWFunc())
	svc.ms.POST("/merchant/:merchant_id/member", svc.CreateMember, svc.jwtService.MWFunc())
	svc.ms.GET("/merchant/:merchant_id/member/:id", svc.GetMemberInfo, svc.jwtService.MWFunc())
	svc.ms.PUT("/merchant/:merchant_id/member/:id", svc.EditMember, svc.jwtService.MWFunc())
	svc.ms.PUT("/merchant/:merchant_id/member/:id/password", svc.ChangePasswordMember, svc.jwtService.MWFunc())
	svc.ms.DELETE("/merchant/:merchant_id/member/:id", svc.DeleteMember, svc.jwtService.MWFunc())

}

func (svc *MerchantService) SearchMerchant(ctx microservice.IServiceContext) error {

	username := ctx.UserInfo().Username

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	merchantList := []models.MerchantInfo{}
	pagination, err := pst.FindPage(&models.Merchant{}, limit, page, bson.M{"createdBy": username, "deleted": false, "name1": bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + q + ".*",
		Options: "",
	}}}, &merchantList)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	for mid := range merchantList {
		count, err := pst.Count(&models.Member{}, bson.M{"merchantId": merchantList[mid].Id})
		if err != nil {
			merchantList[mid].TotalMember = 0
		}
		merchantList[mid].TotalMember = count
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": merchantList})
	return nil
}

func (svc *MerchantService) CreateMerchant(ctx microservice.IServiceContext) error {

	if len(ctx.UserInfo().Username) < 1 {
		ctx.ResponseError(400, "user authentication invalid")
	}

	input := ctx.ReadInput()

	merchantReq := &models.Merchant{}
	err := json.Unmarshal([]byte(input), &merchantReq)

	if err != nil {
		ctx.ResponseError(400, "merchant payload invalid")
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	merchantReq.GuidFixed = utils.NewGUID()
	merchantReq.CreatedBy = ctx.UserInfo().Username
	merchantReq.CreatedAt = time.Now()

	_, err = pst.Create(&models.Merchant{}, merchantReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "guidFixed": merchantReq.GuidFixed})
	return nil
}

func (svc *MerchantService) DeleteMerchant(ctx microservice.IServiceContext) error {
	username := ctx.UserInfo().Username
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	findMerchant := &models.Merchant{}
	err := pst.FindOne(&models.Merchant{}, bson.M{"guidFixed": id, "deleted": false}, findMerchant)

	if err != nil && err.Error() != "mongo: no documents in result" {
		svc.ms.Log("merchant service", err.Error())
		ctx.ResponseError(400, "database error")
		return err
	}

	if findMerchant.CreatedBy != username {
		ctx.ResponseError(400, "username invalid")
		return err
	}

	err = pst.SoftDeleteByID(&models.Merchant{}, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *MerchantService) EditMerchant(ctx microservice.IServiceContext) error {

	username := ctx.UserInfo().Username
	id := ctx.Param("id")
	input := ctx.ReadInput()

	merchantRequest := &models.Merchant{}
	err := json.Unmarshal([]byte(input), &merchantRequest)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	findMerchant := &models.Merchant{}

	err = pst.FindOne(&models.Merchant{}, bson.M{"guidFixed": id, "createdBy": username, "deleted": false}, findMerchant)

	if len(findMerchant.CreatedBy) < 1 {
		ctx.ResponseError(400, "username invalid")
		return err
	}

	findMerchant.Name1 = merchantRequest.Name1
	findMerchant.UpdatedBy = username
	findMerchant.UpdatedAt = time.Now()

	err = pst.Update(&models.Merchant{}, findMerchant, "guidFixed", id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *MerchantService) InfoMerchant(ctx microservice.IServiceContext) error {
	username := ctx.UserInfo().Username
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	merchant := &models.Merchant{}

	err := pst.FindOne(&models.Merchant{}, bson.M{"guidFixed": id, "createdBy": username, "deleted": false}, merchant)

	if err != nil {
		ctx.ResponseError(400, "not found")
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "data": merchant})
	return nil
}
