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

type InventoryOperation interface {
	AddInventory() error
	DeleteInventory() error
	EditInventory() error
	GetInventoryInfo() error
	GetInventorySync() error
	SearchInventory() error
}

type InventoryServiceOld struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
}

func NewInventoryServiceOld(ms *microservice.Microservice, cfg microservice.IConfig) *InventoryServiceOld {

	inventoryapi := &InventoryServiceOld{
		ms:  ms,
		cfg: cfg,
	}
	return inventoryapi
}

func (svc *InventoryServiceOld) RouteSetup() {

	svc.ms.GET("/", microservice.MicroServiceHealthyHandler)
	svc.ms.GET("/inventory", svc.SearchInventory)
	svc.ms.POST("/inventory", svc.CreateInventory)
	svc.ms.GET("/inventory/:id", svc.InfoInventory)
	svc.ms.PUT("/inventory/:id", svc.EditInventory)
	svc.ms.GET("/inventorysync", svc.GetInventorySync)
	svc.ms.DELETE("/inventory/:id", svc.DeleteInventory)

	// 	svc.ms.GET("/category", svc.SearchCategory)
	// 	svc.ms.POST("/category", svc.CreateCategory)
	// 	svc.ms.GET("/category/:id", svc.InfoCategory)
	// 	svc.ms.PUT("/category/:id", svc.EditCategory)
	// 	svc.ms.DELETE("/category/:id", svc.DeleteCategory)

	// 	svc.ms.GET("/optgroup", svc.SearchOptionGroup)
	// 	svc.ms.POST("/optgroup", svc.CreateOptionGroup)
	// 	svc.ms.GET("/optgroup/:id", svc.InfoOptionGroup)
	// 	svc.ms.PUT("/optgroup/:id", svc.EditOptionGroup)
	// 	svc.ms.DELETE("/optgroup/:id", svc.DeleteOptionGroup)

	// 	svc.ms.POST("/option", svc.CreateInventoryOption)
	// 	svc.ms.GET("/option/:id", svc.InfoInventoryOption)
	// 	svc.ms.PUT("/option/:id", svc.EditInventoryOption)
	// 	svc.ms.DELETE("/option/:id", svc.DeleteInventoryOption)
}

func (svc *InventoryServiceOld) CreateInventory(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	input := ctx.ReadInput()

	inventory := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventory)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	inventory.MerchantId = merchantId
	inventory.GuidFixed = utils.NewGUID()
	inventory.WaitType = 0
	inventory.Deleted = false
	inventory.CreatedBy = authUsername
	inventory.CreatedAt = time.Now()

	idx, err := pst.Create(&models.Inventory{}, inventory)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK,
		models.ApiResponse{
			Success: true,
			Id:      idx,
		})
	return nil
}

func (svc *InventoryServiceOld) DeleteInventory(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	findDoc := &models.Inventory{}
	err := pst.FindOne(&models.Inventory{}, bson.M{"merchantId": merchantId, "guidFixed": id, "deleted": false}, findDoc)

	if err != nil && err.Error() != "mongo: no documents in result" {
		svc.ms.Log("merchant service", err.Error())
		ctx.ResponseError(400, "database error")
		return err
	}

	if findDoc.CreatedBy != authUsername {
		ctx.ResponseError(400, "username invalid")
		return err
	}

	inventory := &models.Inventory{}
	err = pst.SoftDeleteByID(inventory, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK,
		models.ApiResponse{
			Success: true,
		})

	return nil
}

func (svc *InventoryServiceOld) EditInventory(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")
	input := ctx.ReadInput()

	inventoryReq := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	findInv := &models.Inventory{}
	err = pst.FindOne(&models.Inventory{}, bson.M{"merchantId": merchantId, "guidFixed": id, "createdBy": authUsername, "deleted": false}, findInv)

	if err != nil {
		ctx.ResponseError(400, "guid invalid")
		return err
	}

	findInv.ItemSku = inventoryReq.ItemSku
	findInv.CategoryGuid = inventoryReq.CategoryGuid
	findInv.LineNumber = inventoryReq.LineNumber
	findInv.Price = inventoryReq.Price
	findInv.Recommended = inventoryReq.Recommended
	findInv.HaveImage = inventoryReq.HaveImage
	findInv.Activated = inventoryReq.Activated

	findInv.Name1 = inventoryReq.Name1
	findInv.Name2 = inventoryReq.Name2
	findInv.Name3 = inventoryReq.Name3
	findInv.Name4 = inventoryReq.Name4
	findInv.Name5 = inventoryReq.Name5

	findInv.Description1 = inventoryReq.Description1
	findInv.Description2 = inventoryReq.Description2
	findInv.Description3 = inventoryReq.Description3
	findInv.Description4 = inventoryReq.Description4
	findInv.Description5 = inventoryReq.Description5

	findInv.UpdatedBy = authUsername
	findInv.UpdatedAt = time.Now()

	err = pst.UpdateOne(&models.Inventory{}, "guidFixed", id, findInv)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK,
		models.ApiResponse{
			Success: true,
		})
	return nil
}

func (svc *InventoryServiceOld) InfoInventory(ctx microservice.IServiceContext) error {

	userInfo := ctx.UserInfo()
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	inventory := &models.Inventory{}
	err := pst.FindOne(&models.Inventory{}, bson.M{"merchantId": merchantId, "guidFixed": id, "deleted": false}, inventory)

	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if len(inventory.GuidFixed) < 1 {
		ctx.ResponseError(404, "Not found.")
		return err
	}

	ctx.Response(http.StatusOK,
		models.ApiResponse{
			Success: true,
			Data:    inventory,
		},
	)
	return nil
}

func (svc *InventoryServiceOld) GetInventorySync(ctx microservice.IServiceContext) error {

	ctx.Response(http.StatusOK,
		models.ApiResponse{
			Success: true,
		})
	return nil
}

func (svc *InventoryServiceOld) SearchInventory(ctx microservice.IServiceContext) error {

	userInfo := ctx.UserInfo()
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

	inventories := []models.Inventory{}

	pagination, err := pst.FindPage(&models.Inventory{}, limit, page, bson.M{
		"merchantId": merchantId,
		"deleted":    false,
		"$or": []interface{}{
			bson.M{"guidFixed": q},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &inventories)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success:    true,
		Pagination: pagination,
		Data:       inventories,
	})
	return nil
}
