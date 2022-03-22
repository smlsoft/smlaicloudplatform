package stockadjustment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

func StartStockAdjustmentAsync(ms *microservice.Microservice, cfg microservice.IConfig) {

	prod := ms.Producer(cfg.MQConfig())
	repo := NewStockAdjustmentRepository(ms.MongoPersister(cfg.MongoPersisterConfig()))

	mqRepo := NewStockAdjustmentMQRepository(prod)

	service := NewStockAdjustmentService(repo, mqRepo)
	err := ms.AsyncPOST("/stockadjustment/async", cfg.CacherConfig(), cfg.MQConfig(), func(ctx microservice.IContext) error {
		userInfo := ctx.UserInfo()
		authUsername := userInfo.Username
		shopID := userInfo.ShopID
		input := ctx.ReadInput()

		stockadjustment := models.StockAdjustment{}
		err := json.Unmarshal([]byte(input), &stockadjustment)

		if err != nil {
			fmt.Println(err.Error())
			ctx.ResponseError(400, err.Error())
			return err
		}

		idx, err := service.CreateStockAdjustment(shopID, authUsername, stockadjustment)

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
