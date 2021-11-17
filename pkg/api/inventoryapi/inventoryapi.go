package inventoryapi

import (
	"net/http"
	"smlcloudplatform/internal/microservice"
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


type InventoryApi struct  {
	ms	*microservice.Microservice
}

func NewInventoryAPI(ms *microservice.Microservice, cfg microservice.IConfig) *InventoryApi {

	inventoryapi := &InventoryApi{
		ms: ms,
	}
	return inventoryapi
}

func (api *InventoryApi)RouteSetup() {
	
	api.ms.GET("/", microservice.MicroServiceHealthyHandler)
	api.ms.POST("/inventory", api.AddInventory)
	api.ms.GET("/inventory", api.SearchInventory)
	api.ms.GET("/inventory/:id", api.GetInventoryInfo)
	api.ms.PUT("/inventory/:id", api.EditInventory)
	api.ms.GET("/inventorysync", api.GetInventorySync)
	api.ms.DELETE("/inventory/:id", api.DeleteInventory)
}

func (api *InventoryApi)AddInventory(ctx microservice.IServiceContext) error {
	ctx.Log("AddInventory")
	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (api *InventoryApi)DeleteInventory(ctx microservice.IServiceContext) error {
	ctx.Log("DeleteInventory")
	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (api *InventoryApi)EditInventory(ctx microservice.IServiceContext) error {
	ctx.Log("EditInventory")
	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (api *InventoryApi)GetInventoryInfo(ctx microservice.IServiceContext) error {
	ctx.Log("GetInventoryInfo")
	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (api *InventoryApi)GetInventorySync(ctx microservice.IServiceContext) error {
	ctx.Log("GetInventorySync")
	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (api *InventoryApi)SearchInventory(ctx microservice.IServiceContext) error {
	ctx.Log("SearchInventory")
	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}