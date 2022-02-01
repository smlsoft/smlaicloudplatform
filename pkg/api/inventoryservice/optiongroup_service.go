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

func (svc *InventoryService) CreateOptionGroup(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	input := ctx.ReadInput()

	inventoryOperationGroup := &models.InventoryOptionGroup{}
	err := json.Unmarshal([]byte(input), &inventoryOperationGroup)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	inventoryOperationGroup.MerchantId = merchantId
	inventoryOperationGroup.GuidFixed = utils.NewGUID()
	inventoryOperationGroup.CreatedBy = authUsername
	inventoryOperationGroup.CreatedAt = time.Now()
	inventoryOperationGroup.Deleted = false

	for i := range inventoryOperationGroup.Details {
		inventoryOperationGroup.Details[i].GuidFixed = utils.NewGUID()
	}

	idx, err := pst.Create(&models.InventoryOptionGroup{}, inventoryOperationGroup)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "id": idx})
	return nil
}

func (svc *InventoryService) EditOptionGroup(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")
	input := ctx.ReadInput()

	inventoryOperationGroupReq := &models.InventoryOptionGroup{}
	err := json.Unmarshal([]byte(input), &inventoryOperationGroupReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	findDoc := &models.InventoryOptionGroup{}
	err = pst.FindOne(&models.InventoryOptionGroup{}, bson.M{"guidFixed": id, "merchantId": merchantId, "deleted": false}, findDoc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	findDoc.OptionName1 = inventoryOperationGroupReq.OptionName1
	findDoc.ProductSelectOption1 = inventoryOperationGroupReq.ProductSelectOption1
	findDoc.ProductSelectOption2 = inventoryOperationGroupReq.ProductSelectOption2
	findDoc.ProductSelectOptionMin = inventoryOperationGroupReq.ProductSelectOptionMin
	findDoc.ProductSelectOptionMax = inventoryOperationGroupReq.ProductSelectOptionMax
	findDoc.Details = inventoryOperationGroupReq.Details
	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = pst.Update(&models.InventoryOptionGroup{}, findDoc, "guidFixed", id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *InventoryService) InfoOptionGroup(ctx microservice.IServiceContext) error {
	username := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	category := &models.InventoryOptionGroup{}

	err := pst.FindOne(&models.InventoryOptionGroup{}, bson.M{"guidFixed": id, "merchantId": merchantId, "createdBy": username, "deleted": false}, category)

	if err != nil {
		ctx.ResponseError(400, "not found")
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "data": category})
	return nil
}

func (svc *InventoryService) DeleteOptionGroup(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	findCategory := &models.InventoryOptionGroup{}
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

	category := &models.InventoryOptionGroup{}
	err = pst.SoftDeleteByID(category, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *InventoryService) SearchOptionGroup(ctx microservice.IServiceContext) error {
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

	inventoryOptionGroupList := []models.InventoryOptionGroup{}
	pagination, err := pst.FindPage(&models.InventoryOptionGroup{}, limit, page, bson.M{"merchantId": merchantId, "createdBy": username, "deleted": false, "optionName1": bson.M{"$regex": primitive.Regex{
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
