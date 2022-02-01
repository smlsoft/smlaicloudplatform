package inventoryservice

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (svc *InventoryService) CreateInventoryOption(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	input := ctx.ReadInput()

	modelReq := &models.InventoryOption{}
	err := json.Unmarshal([]byte(input), &modelReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	modelReq.MerchantId = merchantId
	modelReq.GuidFixed = utils.NewGUID()
	modelReq.Deleted = false
	modelReq.CreatedBy = authUsername
	modelReq.CreatedAt = time.Now()

	idx, err := pst.Create(&models.InventoryOption{}, modelReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "id": idx})
	return nil
}

func (svc *InventoryService) EditInventoryOption(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")
	input := ctx.ReadInput()

	modelReq := &models.InventoryOption{}
	err := json.Unmarshal([]byte(input), &modelReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	findDoc := &models.InventoryOption{}
	err = pst.FindOne(&models.InventoryOption{}, bson.M{"guidFixed": id, "merchantId": merchantId, "deleted": false}, findDoc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	findDoc.InventoryId = modelReq.InventoryId
	findDoc.OptionGroupId = modelReq.OptionGroupId
	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = pst.Update(&models.InventoryOption{}, findDoc, "guidFixed", id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *InventoryService) InfoInventoryOption(ctx microservice.IServiceContext) error {
	username := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	category := &models.InventoryOption{}

	err := pst.FindOne(&models.InventoryOption{}, bson.M{"guidFixed": id, "merchantId": merchantId, "createdBy": username, "deleted": false}, category)

	if err != nil {
		ctx.ResponseError(400, "not found")
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "data": category})
	return nil
}

func (svc *InventoryService) DeleteInventoryOption(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	findCategory := &models.InventoryOption{}
	err := pst.FindOne(&models.Category{}, bson.M{"guidFixed": id, "merchantId": merchantId}, findCategory)

	if err != nil && err.Error() != "mongo: no documents in result" {
		svc.ms.Log("merchant service", err.Error())
		ctx.ResponseError(400, "database error")
		return err
	}

	if findCategory.CreatedBy != authUsername {
		ctx.ResponseError(400, "username invalid")
		return err
	}

	category := &models.InventoryOption{}
	err = pst.SoftDeleteByID(category, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *InventoryService) SearchInventoryOption(ctx microservice.IServiceContext) error {
	merchantId := ctx.Param("merchant_id")
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

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	inventoryOptionGroupList := []models.InventoryOption{}
	pagination, err := pst.FindPage(&models.InventoryOption{}, limit, page, bson.M{"merchantId": merchantId, "createdBy": username, "deleted": false, "optionName1": bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + q + ".*",
		Options: "",
	}}}, &inventoryOptionGroupList)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": inventoryOptionGroupList})
	return nil
}
