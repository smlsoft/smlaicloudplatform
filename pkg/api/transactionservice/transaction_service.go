package transactionservice

import (
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

type TransactionService struct {
	ms         *microservice.Microservice
	cfg        microservice.IConfig
	jwtService *microservice.JwtService
}

func NewTransactionService(ms *microservice.Microservice, cfg microservice.IConfig) *TransactionService {

	signKey, verifyKey, err := utils.LoadKey(cfg.SignKeyPath(), cfg.VerifyKeyPath())

	if err != nil {
		fmt.Println("jwt key error :: " + err.Error())
	}

	jwtService := microservice.NewJwtService(signKey, verifyKey, 60*24*10)

	inventoryapi := &TransactionService{
		ms:         ms,
		cfg:        cfg,
		jwtService: jwtService,
	}
	return inventoryapi
}

func (svc *TransactionService) RouteSetup() {
	svc.ms.HttpMiddleware(svc.jwtService.MWFunc())
	svc.ms.GET("/merchant/:merchant_id/transaction", svc.SearchTransaction)
	svc.ms.POST("/merchant/:merchant_id/transaction", svc.CreateTransaction)
	svc.ms.GET("/merchant/:merchant_id/transaction/:id", svc.InfoTransaction)
	svc.ms.PUT("/merchant/:merchant_id/transaction/:id", svc.EditTransaction)
	svc.ms.DELETE("/merchant/:merchant_id/transaction/:id", svc.DeleteTransaction)
}

func (svc *TransactionService) CreateTransaction(ctx microservice.IServiceContext) error {
	input := ctx.ReadInput()
	merchantId := ctx.Param("merchant_id")
	authUsername := ctx.UserInfo().Username

	trans := &models.Transaction{}
	err := json.Unmarshal([]byte(input), &trans)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	merchant := &models.Merchant{}
	pst.FindOne(&models.Merchant{}, bson.M{"guidFixed": merchantId, "deleted": false}, merchant)

	if len(merchant.GuidFixed) < 1 {
		ctx.ResponseError(400, "merchant invalid")
		return nil
	}

	if merchant.CreatedBy != authUsername {
		ctx.ResponseError(400, "username invalid")
		return err
	}

	sumAmount := 0.0
	for i, transDetail := range trans.Items {
		trans.Items[i].LineNumber = i + 1
		sumAmount += transDetail.Price * transDetail.Qty
	}

	trans.SumAmount = sumAmount
	trans.Deleted = false
	trans.CreatedBy = authUsername
	trans.CreatedAt = time.Now()

	idx, err := pst.Create(&models.Transaction{}, trans)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "id": idx})
	return nil
}

func (svc *TransactionService) DeleteTransaction(ctx microservice.IServiceContext) error {
	id := ctx.Param("id")
	merchantId := ctx.Param("merchant_id")
	authUsername := ctx.UserInfo().Username

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	findDoc := &models.Transaction{}
	err := pst.FindOne(&models.Transaction{}, bson.M{"merchantId": merchantId, "guidFixed": id, "deleted": false}, findDoc)

	if err != nil && err.Error() != "mongo: no documents in result" {
		svc.ms.Log("merchant service", err.Error())
		ctx.ResponseError(400, "database error")
		return err
	}

	if findDoc.CreatedBy != authUsername {
		ctx.ResponseError(400, "username invalid")
		return err
	}

	trans := &models.Transaction{}
	err = pst.SoftDeleteByID(trans, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *TransactionService) EditTransaction(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")
	input := ctx.ReadInput()

	transReq := &models.Transaction{}
	err := json.Unmarshal([]byte(input), &transReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	findTrans := &models.Transaction{}
	err = pst.FindOne(&models.Transaction{}, bson.M{"merchantId": merchantId, "guidFixed": id, "createdBy": authUsername, "deleted": false}, findTrans)

	if err != nil {
		ctx.ResponseError(400, "guid invalid")
		return err
	}

	sumAmount := 0.0
	for i, transDetail := range findTrans.Items {
		findTrans.Items[i].LineNumber = i + 1
		sumAmount += transDetail.Price * transDetail.Qty
	}

	findTrans.Items = transReq.Items
	findTrans.SumAmount = sumAmount
	findTrans.UpdatedBy = authUsername
	findTrans.UpdatedAt = time.Now()

	err = pst.Update(&models.Transaction{}, findTrans, "guidFixed", id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *TransactionService) InfoTransaction(ctx microservice.IServiceContext) error {

	id := ctx.Param("id")
	merchantId := ctx.Param("merchant_id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	trans := &models.Transaction{}
	err := pst.FindOne(&models.Transaction{}, bson.M{"merchantId": merchantId, "guidFixed": id, "deleted": false}, trans)

	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "data": trans})
	return nil
}

func (svc *TransactionService) SearchTransaction(ctx microservice.IServiceContext) error {

	merchantId := ctx.Param("merchant_id")

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

	inventories := []models.Inventory{}

	pagination, err := pst.FindPage(&models.Inventory{}, limit, page, bson.M{
		"merchantId": merchantId,
		"deleted":    false,
		"$or": []interface{}{
			bson.M{"guidFixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &inventories)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": inventories})
	return nil
}

func (svc *TransactionService) SearchTransactionItems(ctx microservice.IServiceContext) error {

	merchantId := ctx.Param("merchant_id")
	transId := ctx.Param("trans_id")

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

	inventories := []models.Inventory{}

	pagination, err := pst.FindPage(&models.Inventory{}, limit, page, bson.M{
		"merchantId": merchantId,
		"guidFixed":  transId,
		"deleted":    false,
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &inventories)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": inventories})
	return nil
}
