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

type InventoryService struct {
	ms         *microservice.Microservice
	cfg        microservice.IConfig
	jwtService *microservice.JwtService
}

func NewInventoryService(ms *microservice.Microservice, cfg microservice.IConfig) *InventoryService {

	// signKey, verifyKey, err := utils.LoadKey(cfg.SignKeyPath(), cfg.VerifyKeyPath())

	// if err != nil {
	// 	fmt.Println("jwt key error :: " + err.Error())
	// }

	// jwtService := microservice.NewJwtService(signKey, verifyKey, 60*24*10)

	jwtService := microservice.NewJwtService(cfg.JwtSecretKey(), 60*24*10)

	inventoryapi := &InventoryService{
		ms:         ms,
		cfg:        cfg,
		jwtService: jwtService,
	}
	return inventoryapi
}

func (svc *InventoryService) RouteSetup() {

	svc.ms.HttpMiddleware(svc.jwtService.MWFunc())
	svc.ms.GET("/", microservice.MicroServiceHealthyHandler)
	svc.ms.GET("/merchant/:merchant_id/inventory", svc.SearchInventory)
	svc.ms.POST("/merchant/:merchant_id/inventory", svc.CreateInventory)
	svc.ms.GET("/merchant/:merchant_id/inventory/:id", svc.InfoInventory)
	svc.ms.PUT("/merchant/:merchant_id/inventory/:id", svc.EditInventory)
	svc.ms.GET("/merchant/:merchant_id/inventorysync", svc.GetInventorySync)
	svc.ms.DELETE("/merchant/:merchant_id/inventory/:id", svc.DeleteInventory)

	svc.ms.GET("/merchant/:merchant_id/category", svc.SearchCategory)
	svc.ms.POST("/merchant/:merchant_id/category", svc.CreateCategory)
	svc.ms.GET("/merchant/:merchant_id/category/:id", svc.InfoCategory)
	svc.ms.PUT("/merchant/:merchant_id/category/:id", svc.EditCategory)
	svc.ms.DELETE("/merchant/:merchant_id/category/:id", svc.DeleteCategory)

	svc.ms.GET("/merchant/:merchant_id/optgroup", svc.SearchOptionGroup)
	svc.ms.POST("/merchant/:merchant_id/optgroup", svc.CreateOptionGroup)
	svc.ms.GET("/merchant/:merchant_id/optgroup/:id", svc.InfoOptionGroup)
	svc.ms.PUT("/merchant/:merchant_id/optgroup/:id", svc.EditOptionGroup)
	svc.ms.DELETE("/merchant/:merchant_id/optgroup/:id", svc.DeleteOptionGroup)

	svc.ms.POST("/merchant/:merchant_id/option", svc.CreateInventoryOption)
	svc.ms.GET("/merchant/:merchant_id/option/:id", svc.InfoInventoryOption)
	svc.ms.PUT("/merchant/:merchant_id/option/:id", svc.EditInventoryOption)
	svc.ms.DELETE("/merchant/:merchant_id/option/:id", svc.DeleteInventoryOption)
}

func (svc *InventoryService) CreateInventory(ctx microservice.IServiceContext) error {
	input := ctx.ReadInput()
	merchantId := ctx.Param("merchant_id")
	authUsername := ctx.UserInfo().Username

	inventory := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventory)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

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

func (svc *InventoryService) DeleteInventory(ctx microservice.IServiceContext) error {
	id := ctx.Param("id")
	merchantId := ctx.Param("merchant_id")
	authUsername := ctx.UserInfo().Username

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

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

func (svc *InventoryService) EditInventory(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.Param("merchant_id")
	id := ctx.Param("id")
	input := ctx.ReadInput()

	inventoryReq := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

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

	err = pst.Update(&models.Inventory{}, findInv, "guidFixed", id)

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

func (svc *InventoryService) InfoInventory(ctx microservice.IServiceContext) error {

	id := ctx.Param("id")
	merchantId := ctx.Param("merchant_id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

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

func (svc *InventoryService) GetInventorySync(ctx microservice.IServiceContext) error {

	ctx.Response(http.StatusOK,
		models.ApiResponse{
			Success: true,
		})
	return nil
}

func (svc *InventoryService) SearchInventory(ctx microservice.IServiceContext) error {

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

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

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
