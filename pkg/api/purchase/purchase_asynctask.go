package purchase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

func StartPurchaseAsync(ms *microservice.Microservice, cfg microservice.IConfig) {

	prod := ms.Producer(cfg.MQConfig())
	repo := NewPurchaseRepository(ms.MongoPersister(cfg.MongoPersisterConfig()))

	mqRepo := NewPurchaseMQRepository(prod)

	service := NewPurchaseService(repo, mqRepo)
	err := ms.AsyncPOST("/purchase/async", cfg.CacherConfig(), cfg.MQConfig(), func(ctx microservice.IContext) error {
		userInfo := ctx.UserInfo()
		authUsername := userInfo.Username
		shopID := userInfo.ShopID
		input := ctx.ReadInput()

		purchase := models.Purchase{}
		err := json.Unmarshal([]byte(input), &purchase)

		if err != nil {
			fmt.Println(err.Error())
			ctx.ResponseError(400, err.Error())
			return err
		}

		idx, err := service.CreatePurchase(shopID, authUsername, purchase)

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
