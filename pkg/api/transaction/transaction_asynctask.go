package transaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

func StartTransactionAsync(ms *microservice.Microservice, cfg microservice.IConfig) {

	repo := NewTransactionRepository(ms.MongoPersister(cfg.MongoPersisterConfig()))
	prod := ms.Producer(cfg.MQConfig())

	mqRepo := NewTransactionMQRepository(prod)
	service := NewTransactionService(repo, mqRepo)
	err := ms.AsyncPOST("/trans", cfg.CacherConfig(), cfg.MQConfig(), func(ctx microservice.IContext) error {
		userInfo := ctx.UserInfo()
		authUsername := userInfo.Username
		merchantId := userInfo.MerchantId
		input := ctx.ReadInput()

		trans := &models.Transaction{}
		err := json.Unmarshal([]byte(input), &trans)

		if err != nil {
			fmt.Println(err.Error())
			ctx.ResponseError(400, err.Error())
			return err
		}

		idx, err := service.CreateTransaction(merchantId, authUsername, trans)

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
