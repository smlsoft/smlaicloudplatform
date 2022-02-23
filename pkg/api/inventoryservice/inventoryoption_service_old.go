package inventoryservice

/*
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
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	input := ctx.ReadInput()

	modelReq := &models.InventoryOption{}
	err := json.Unmarshal([]byte(input), &modelReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

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

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Id:      idx,
	})
	return nil
}

func (svc *InventoryService) EditInventoryOption(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")
	input := ctx.ReadInput()

	modelReq := &models.InventoryOption{}
	err := json.Unmarshal([]byte(input), &modelReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

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

	err = pst.UpdateOne(&models.InventoryOption{}, "guidFixed", id, findDoc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (svc *InventoryService) InfoInventoryOption(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	doc := &models.InventoryOption{}

	err := pst.FindOne(&models.InventoryOption{}, bson.M{"guidFixed": id, "merchantId": merchantId, "createdBy": authUsername, "deleted": false}, doc)

	if err != nil {
		ctx.ResponseError(400, "not found")
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

func (svc *InventoryService) DeleteInventoryOption(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

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

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})

	return nil
}

func (svc *InventoryService) SearchInventoryOption(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

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

	docList := []models.InventoryOption{}
	pagination, err := pst.FindPage(&models.InventoryOption{}, limit, page, bson.M{"merchantId": merchantId, "createdBy": authUsername, "deleted": false, "optionName1": bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + q + ".*",
		Options: "",
	}}}, &docList)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success:    true,
		Pagination: pagination,
		Data:       docList,
	})
	return nil
}
*/
