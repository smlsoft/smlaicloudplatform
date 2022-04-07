package inventory

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

func StartInventoryAsync(ms *microservice.Microservice, cfg microservice.IConfig) {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstPg := ms.Persister(cfg.PersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	invRepo := NewInventoryRepository(pst)
	invPgRepo := NewInventoryIndexPGRepository(pstPg)
	invMqRepo := NewInventoryMQRepository(prod)
	invService := NewInventoryService(invRepo, invPgRepo, invMqRepo)

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
