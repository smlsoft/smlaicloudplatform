package barcodemaster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/product/barcodemaster/models"
	"smlcloudplatform/pkg/product/barcodemaster/repositories"
	"smlcloudplatform/pkg/product/barcodemaster/services"
)

func StartBarcodeMasterAsync(ms *microservice.Microservice, cfg microservice.IConfig) {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	invRepo := repositories.NewBarcodeMasterRepository(pst)
	invMqRepo := repositories.NewBarcodeMasterMQRepository(prod)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache, "barcodemaster")
	invService := services.NewBarcodeMasterService(invRepo, invMqRepo, masterSyncCacheRepo)

	err := ms.AsyncPOST("/barcodemaster-async", cfg.CacherConfig(), cfg.MQConfig(), func(ctx microservice.IContext) error {
		userInfo := ctx.UserInfo()
		authUsername := userInfo.Username
		shopID := userInfo.ShopID
		input := ctx.ReadInput()

		trans := models.BarcodeMaster{}
		err := json.Unmarshal([]byte(input), &trans)

		if err != nil {
			fmt.Println(err.Error())
			ctx.ResponseError(400, err.Error())
			return err
		}

		_, idx, err := invService.CreateBarcodeMaster(shopID, authUsername, trans)

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
