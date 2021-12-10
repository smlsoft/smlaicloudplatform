package inventoryapi

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryOperation interface {
	// defint e
	AddInventory() error
	DeleteInventory() error
	EditInventory() error
	GetInventoryInfo() error
	GetInventorySync() error
	SearchInventory() error
}

type InventoryApi struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
}

func NewInventoryAPI(ms *microservice.Microservice, cfg microservice.IConfig) *InventoryApi {

	inventoryapi := &InventoryApi{
		ms:  ms,
		cfg: cfg,
	}
	return inventoryapi
}

func (api *InventoryApi) RouteSetup() {

	api.ms.GET("/", microservice.MicroServiceHealthyHandler)
	api.ms.POST("/inventory", api.AddInventory)
	api.ms.GET("/inventory", api.SearchInventory)
	api.ms.GET("/inventory/:id", api.GetInventoryInfo)
	api.ms.PUT("/inventory/:id", api.EditInventory)
	api.ms.GET("/inventorysync", api.GetInventorySync)
	api.ms.DELETE("/inventory/:id", api.DeleteInventory)
}

func (api *InventoryApi) AddInventory(ctx microservice.IServiceContext) error {
	ctx.Log("AddInventory")

	input := ctx.ReadInput()

	inventory := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventory)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	inventory.LatestUpdate = time.Now()

	pst := api.ms.MongoPersister(api.cfg.MongoPersisterConfig())

	idx, err := pst.Create(&models.Inventory{}, inventory)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "id": idx})
	return nil
}

func (api *InventoryApi) DeleteInventory(ctx microservice.IServiceContext) error {
	ctx.Log("DeleteInventory")
	id := ctx.Param("id")

	pst := api.ms.MongoPersister(api.cfg.MongoPersisterConfig())

	inventory := &models.Inventory{}
	err := pst.DeleteByID(inventory, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (api *InventoryApi) EditInventory(ctx microservice.IServiceContext) error {
	ctx.Log("EditInventory")

	id := ctx.Param("id")
	input := ctx.ReadInput()

	inventory := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventory)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	inventory.LatestUpdate = time.Now()

	pst := api.ms.MongoPersister(api.cfg.MongoPersisterConfig())

	err = pst.Update(inventory, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (api *InventoryApi) GetInventoryInfo(ctx microservice.IServiceContext) error {
	ctx.Log("GetInventoryInfo")

	id := ctx.Param("id")

	pst := api.ms.MongoPersister(api.cfg.MongoPersisterConfig())

	inventory := &models.Inventory{}
	err := pst.FindByID(&models.Inventory{}, id, inventory)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "data": inventory})
	return nil
}

func (api *InventoryApi) GetInventorySync(ctx microservice.IServiceContext) error {
	ctx.Log("GetInventorySync")
	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (api *InventoryApi) SearchInventory(ctx microservice.IServiceContext) error {
	ctx.Log("SearchInventory")

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		ctx.ResponseError(400, "page number invalid.")
		return err
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		ctx.ResponseError(400, "limit number invalid.")
		return err
	}

	pst := api.ms.MongoPersister(api.cfg.MongoPersisterConfig())

	inventories := []models.Inventory{}
	pagination, err := pst.FindPage(&models.Inventory{}, limit, page, bson.M{"product_name": bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + q + ".*",
		Options: "",
	}}}, &inventories)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": inventories})
	return nil
}
