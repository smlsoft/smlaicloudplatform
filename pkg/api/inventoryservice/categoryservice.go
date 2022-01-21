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

func (svc *InventoryService) CreateCategory(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	input := ctx.ReadInput()

	categoryReq := &models.Category{}
	err := json.Unmarshal([]byte(input), &categoryReq)

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

	countCategory, err := pst.Count(&models.Category{}, bson.M{"merchantId": merchantId})

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	categoryReq.MerchantId = merchantId
	categoryReq.GuidFixed = utils.NewGUID()
	categoryReq.LineNumber = int(countCategory) + 1
	categoryReq.CreatedBy = authUsername
	categoryReq.CreatedAt = time.Now()
	categoryReq.Deleted = false

	idx, err := pst.Create(&models.Category{}, categoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "id": idx})
	return nil
}

func (svc *InventoryService) EditCategory(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")
	input := ctx.ReadInput()

	categoryReq := &models.Category{}
	err := json.Unmarshal([]byte(input), &categoryReq)

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

	findCategory := &models.Category{}
	err = pst.FindOne(&models.Category{}, bson.M{"guidFixed": id, "deleted": false}, findCategory)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	findCategory.Name1 = categoryReq.Name1
	findCategory.HaveImage = categoryReq.HaveImage
	findCategory.UpdatedBy = authUsername
	findCategory.UpdatedAt = time.Now()

	err = pst.Update(&models.Category{}, findCategory, "guidFixed", id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *InventoryService) InfoCategory(ctx microservice.IServiceContext) error {
	username := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	category := &models.Category{}

	err := pst.FindOne(&models.Category{}, bson.M{"guidFixed": id, "merchantId": merchantId, "createdBy": username, "deleted": false}, category)

	if err != nil {
		ctx.ResponseError(400, "not found")
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "data": category})
	return nil
}

func (svc *InventoryService) DeleteCategory(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	findCategory := &models.Category{}
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

	category := &models.Category{}
	err = pst.SoftDeleteByID(category, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *InventoryService) SearchCategory(ctx microservice.IServiceContext) error {
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

	categoryList := []models.Category{}
	pagination, err := pst.FindPage(&models.Category{}, limit, page, bson.M{"merchantId": merchantId, "createdBy": username, "deleted": false, "name1": bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + q + ".*",
		Options: "",
	}}}, &categoryList)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": categoryList})
	return nil
}
