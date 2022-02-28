package transaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

func StartTransactionAPI(ms *microservice.Microservice, cfg microservice.IConfig) {

	repo := NewTransactionRepository(ms.MongoPersister(cfg.MongoPersisterConfig()))
	service := NewTransactionService(repo)
	ms.AsyncPOST("/trans", cfg.CacherConfig(), cfg.MQServer(), func(ctx microservice.IContext) error {
		userInfo := ctx.UserInfo()
		authUsername := userInfo.Username
		merchantId := userInfo.MerchantId
		input := ctx.ReadInput()

		prod := ctx.Producer(cfg.MQServer())

		fmt.Println(ctx.UserInfo())

		trans := &models.Transaction{}
		err := json.Unmarshal([]byte(input), &trans)

		if err != nil {
			fmt.Println(err.Error())
			ctx.ResponseError(400, err.Error())
			return err
		}

		idx, err := service.CreateTransaction(merchantId, authUsername, *trans)

		if err != nil {
			ctx.ResponseError(400, err.Error())
		}

		err = prod.SendMessage("when-transaction-created", "", trans)
		if err != nil {
			ctx.Log(err.Error())
			return err
		}

		ctx.Response(http.StatusOK, idx)
		return nil
	})
}
