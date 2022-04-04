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
	pstPg := ms.Persister(cfg.PersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	invRepo := inventory.NewInventoryRepository(pst)
	invPgRepo := inventory.NewInventoryPGRepository(pstPg)
	invMqRepo := inventory.NewInventoryMQRepository(prod)
	invService := inventory.NewInventoryService(invRepo, invPgRepo, invMqRepo)

	memberRepo := member.NewMemberRepository(pst)
	memberPgRepo := member.NewMemberPGRepository(pstPg)
	memberService := member.NewMemberService(memberRepo, memberPgRepo)

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
		payload, err := h.getPayload(syncData.Data)

		if err != nil {
			return errors.New("data payload invalid")
		}

		inv := models.Inventory{}
		err = json.Unmarshal(payload, &inv)
		if err != nil {
			return err
		}

		err = h.syncData(
			models.Member{},
			userInfo,
			syncData,
			func() (bool, error) {
				return h.inventoryService.IsExistsGuid(userInfo.ShopID, syncData.MyGuid)
			},
			func() (string, error) {
				idx, err := h.inventoryService.CreateWithGuid(userInfo.ShopID, userInfo.Username, syncData.MyGuid, inv)
				if err != nil {
					return "", err
				}
				return idx, nil
			},
			func(idx string) error {

				docIdx := models.InventoryIndex{}
				docIdx.ID = idx
				docIdx.ShopID = userInfo.ShopID
				docIdx.GuidFixed = syncData.MyGuid

				err := h.inventoryService.CreateIndex(docIdx)

				if err != nil {
					return err
				}

				return nil
			},
			func() error {
				return h.inventoryService.UpdateInventory(userInfo.ShopID, syncData.MyGuid, userInfo.Username, inv)
			},
		)

		if err != nil {
			ctx.ResponseError(400, err.Error())
			return err
		}
	case "member":
		payload, err := h.getPayload(syncData.Data)

		if err != nil {
			return errors.New("data payload invalid")
		}

		member := models.Member{}
		err = json.Unmarshal(payload, &member)
		if err != nil {
			return err
		}

		err = h.syncData(
			models.Member{},
			userInfo,
			syncData,
			func() (bool, error) {
				return h.memberService.IsExistsGuid(userInfo.ShopID, syncData.MyGuid)
			},
			func() (string, error) {
				idx, err := h.memberService.CreateWithGuid(userInfo.ShopID, userInfo.Username, syncData.MyGuid, member)
				if err != nil {
					return "", err
				}
				return idx, nil
			},
			func(idx string) error {

				docIdx := models.MemberIndex{}
				docIdx.ID = idx
				docIdx.ShopID = userInfo.ShopID
				docIdx.GuidFixed = syncData.MyGuid

				err := h.memberService.CreateIndex(docIdx)
				if err != nil {
					return err
				}

				return nil
			},
			func() error {
				return h.memberService.Update(userInfo.ShopID, syncData.MyGuid, userInfo.Username, member)
			},
		)

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

func (h SyncDataHttp) syncData(
	model interface{},
	userInfo micro.UserInfo,
	syncData models.SyncData,
	fnIsExists func() (bool, error),
	fnCreate func() (string, error),
	fnCreateIndex func(string) error,
	fnUpdate func() error,
) error {
	payload, err := h.getPayload(syncData.Data)

	if err != nil {
		return errors.New("data payload invalid")
	}

	err = json.Unmarshal(payload, &model)
	if err != nil {
		return err
	}

	switch syncData.Mode {
	case 0:
		isExistsGuid, err := fnIsExists()

		if err != nil {
			return err
		}

		if isExistsGuid {
			return errors.New("guid '" + syncData.MyGuid + "' is exists")
		}

		idx, err := fnCreate()

		if err != nil {
			return err
		}

		fnCreateIndex(idx)

		return err
	case 1:
		isExistsGuid, err := fnIsExists()

		if err != nil {
			return err
		}

		if isExistsGuid {
			idx, err := fnCreate()

			if err != nil {
				return err
			}

			fnCreateIndex(idx)

			return err
		} else {
			return fnUpdate()
		}

	default:
		return errors.New("mode invalid")
	}
}
