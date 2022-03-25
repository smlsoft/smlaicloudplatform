package syncdata

import (
	"encoding/json"
	"errors"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/api/member"
	"smlcloudplatform/pkg/models"

	micro "smlcloudplatform/internal/microservice/models"
)

type SyncDataHttp struct {
	ms               *microservice.Microservice
	cfg              microservice.IConfig
	inventoryService inventory.IInventoryService
	memberService    member.IMemberService
}

func NewSyncDataHttp(ms *microservice.Microservice, cfg microservice.IConfig) SyncDataHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	invRepo := inventory.NewInventoryRepository(pst)
	invMqRepo := inventory.NewInventoryMQRepository(prod)
	invService := inventory.NewInventoryService(invRepo, invMqRepo)

	memberRepo := member.NewMemberRepository(pst)
	memberService := member.NewMemberService(memberRepo)

	return SyncDataHttp{
		ms:               ms,
		cfg:              cfg,
		inventoryService: invService,
		memberService:    memberService,
	}
}

func (h SyncDataHttp) RouteSetup() {

	h.ms.POST("/sync", h.Save)
}

func (h SyncDataHttp) Save(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()

	input := ctx.ReadInput()

	syncData := models.SyncData{}

	err := json.Unmarshal([]byte(input), &syncData)

	if err != nil {
		ctx.ResponseError(400, "payload invalid")
		return err
	}

	switch syncData.TableCode {
	case "inventory":
		err := h.inventorySync(userInfo, syncData)
		if err != nil {
			ctx.ResponseError(400, err.Error())
			return err
		}
	case "member":
		h.memberSync(userInfo, syncData)
		if err != nil {
			ctx.ResponseError(400, err.Error())
			return err
		}
	default:
		ctx.ResponseError(400, "table invalid")
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h SyncDataHttp) getPayload(data interface{}) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (h SyncDataHttp) inventorySync(userInfo micro.UserInfo, syncData models.SyncData) error {
	payload, err := h.getPayload(syncData.Data)

	if err != nil {
		return errors.New("data payload invalid")
	}

	inv := models.Inventory{}
	err = json.Unmarshal(payload, &inv)
	if err != nil {
		return err
	}

	switch syncData.Mode {
	case 0:
		_, err = h.inventoryService.CreateWithGuid(userInfo.ShopID, userInfo.Username, syncData.MyGuid, inv)
		return err
	case 1:
		return h.inventoryService.UpdateInventory(syncData.MyGuid, userInfo.ShopID, userInfo.Username, inv)
	case 3:
		invInfo, err := h.inventoryService.InfoInventory(syncData.MyGuid, userInfo.ShopID)
		if err != nil {
			return err
		}

		if len(invInfo.GuidFixed) < 1 {
			_, err = h.inventoryService.CreateWithGuid(userInfo.ShopID, userInfo.Username, syncData.MyGuid, inv)
			return err
		} else {
			return h.inventoryService.UpdateInventory(syncData.MyGuid, userInfo.ShopID, userInfo.Username, inv)
		}

	default:
		return errors.New("mode invalid")
	}
}

func (h SyncDataHttp) memberSync(userInfo micro.UserInfo, syncData models.SyncData) error {
	payload, err := h.getPayload(syncData.Data)

	if err != nil {
		return err
	}

	member := models.Member{}
	err = json.Unmarshal(payload, &member)
	if err != nil {
		return err
	}

	switch syncData.Mode {
	case 0:
		_, err = h.memberService.CreateWithGuid(userInfo.ShopID, userInfo.Username, syncData.MyGuid, member)
		return err
	case 1:
		return h.memberService.UpdateMember(syncData.MyGuid, userInfo.ShopID, userInfo.Username, member)
	case 3:
		invInfo, err := h.memberService.InfoMember(syncData.MyGuid, userInfo.ShopID)
		if err != nil {
			return err
		}

		if len(invInfo.GuidFixed) < 1 {
			_, err = h.memberService.CreateWithGuid(userInfo.ShopID, userInfo.Username, syncData.MyGuid, member)
			return err
		} else {
			return h.memberService.UpdateMember(syncData.MyGuid, userInfo.ShopID, userInfo.Username, member)
		}

	default:
		return errors.New("mode invalid")
	}
}
