package stockinout

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

func StartStockInOutAsync(ms *microservice.Microservice, cfg microservice.IConfig) {

	prod := ms.Producer(cfg.MQConfig())
	repo := NewStockInOutRepository(ms.MongoPersister(cfg.MongoPersisterConfig()))

	mqRepo := NewStockInOutMQRepository(prod)

	service := NewStockInOutService(repo, mqRepo)
	err := ms.AsyncPOST("/stockinout/async", cfg.CacherConfig(), cfg.MQConfig(), func(ctx microservice.IContext) error {
		userInfo := ctx.UserInfo()
		authUsername := userInfo.Username
		shopId := userInfo.ShopId
		input := ctx.ReadInput()

		stockinout := &models.StockInOut{}
		err := json.Unmarshal([]byte(input), &stockinout)

		if err != nil {
			fmt.Println(err.Error())
			ctx.ResponseError(400, err.Error())
			return err
		}

		idx, err := service.CreateStockInOut(shopId, authUsername, stockinout)

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
