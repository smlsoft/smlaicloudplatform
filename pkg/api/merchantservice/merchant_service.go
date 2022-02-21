package merchantservice

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strconv"
	"time"

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
	signKey, verifyKey, err := utils.LoadKey(cfg.SignKeyPath(), cfg.VerifyKeyPath())

	if err != nil {
		fmt.Println("jwt key error :: " + err.Error())
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

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success:    true,
		Pagination: pagination,
		Data:       merchantList,
	})
	return nil
}

func (svc *MerchantService) CreateMerchant(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	if len(authUsername) < 1 {
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

	merchantId := utils.NewGUID()
	merchantReq.GuidFixed = merchantId
	merchantReq.CreatedBy = ctx.UserInfo().Username
	merchantReq.CreatedAt = time.Now()

	_, err = pst.Create(&models.Merchant{}, merchantReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	findUser := &models.User{}
	err = pst.FindOne(&models.User{}, bson.M{"username": authUsername}, findUser)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if findUser.ID == primitive.NilObjectID {
		ctx.ResponseError(400, "user not found.")
		return err
	}

	userMerchant := models.UserMerchant{
		MerchantId: merchantId,
		Role:       models.ROLE_OWNER,
	}

	findUser.Merchants = append(findUser.Merchants, userMerchant)

	pst.UpdateOne(&models.User{}, "username", authUsername, findUser)

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Id:      merchantReq.GuidFixed,
	})
	return nil
}

func (svc *MerchantService) DeleteMerchant(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	findMerchant := &models.Merchant{}
	err := pst.FindOne(&models.Merchant{}, bson.M{"guidFixed": id, "deleted": false}, findMerchant)

	if err != nil && err.Error() != "mongo: no documents in result" {
		svc.ms.Log("merchant service", err.Error())
		ctx.ResponseError(400, "database error")
		return err
	}

	if findMerchant.CreatedBy != authUsername {
		ctx.ResponseError(400, "username invalid")
		return err
	}

	err = pst.SoftDeleteByID(&models.Merchant{}, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = pst.Update(&models.User{}, bson.M{"username": authUsername}, bson.M{
		"$pull": bson.M{"merchants": bson.M{"merchantId": id}},
	})

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
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

	err = pst.UpdateOne(&models.Merchant{}, "guidFixed", id, findMerchant)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
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

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    merchant,
	})
	return nil
}
