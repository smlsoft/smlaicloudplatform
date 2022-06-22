package inventory

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/models"
)

func StartInventoryAsync(ms *microservice.Microservice, cfg microservice.IConfig) {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	invRepo := NewInventoryRepository(pst)
	invMqRepo := NewInventoryMQRepository(prod)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache, "inventory")
	invService := NewInventoryService(invRepo, invMqRepo, masterSyncCacheRepo)

	err := ms.AsyncPOST("/inv", cfg.CacherConfig(), cfg.MQConfig(), func(ctx microservice.IContext) error {
		userInfo := ctx.UserInfo()
		authUsername := userInfo.Username
		shopID := userInfo.ShopID
		input := ctx.ReadInput()

		trans := models.Inventory{}
		err := json.Unmarshal([]byte(input), &trans)

		if err != nil {
			fmt.Println(err.Error())
			ctx.ResponseError(400, err.Error())
			return err
		}

		_, idx, err := invService.CreateInventory(shopID, authUsername, trans)

		if err != nil {
			ctx.ResponseError(400, err.Error())
		}

		ctx.Response(http.StatusOK, idx)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
