package saleinvoice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

func StartSaleinvoiceAsync(ms *microservice.Microservice, cfg microservice.IConfig) {

	repo := NewSaleinvoiceRepository(ms.MongoPersister(cfg.MongoPersisterConfig()))
	prod := ms.Producer(cfg.MQConfig())

	mqRepo := NewSaleinvoiceMQRepository(prod)
	service := NewSaleinvoiceService(repo, mqRepo)
	err := ms.AsyncPOST("/trans", cfg.CacherConfig(), cfg.MQConfig(), func(ctx microservice.IContext) error {
		userInfo := ctx.UserInfo()
		authUsername := userInfo.Username
		shopID := userInfo.ShopID
		input := ctx.ReadInput()

		trans := models.Saleinvoice{}
		err := json.Unmarshal([]byte(input), &trans)

		if err != nil {
			fmt.Println(err.Error())
			ctx.ResponseError(400, err.Error())
			return err
		}

		idx, err := service.CreateSaleinvoice(shopID, authUsername, trans)

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
