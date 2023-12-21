package shop

import (
	"encoding/json"
	"errors"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	"smlcloudplatform/pkg/logger"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	branchModel "smlcloudplatform/pkg/organization/branch/models"
	branchRepositories "smlcloudplatform/pkg/organization/branch/repositories"
	branchServices "smlcloudplatform/pkg/organization/branch/services"
	businessTypeModels "smlcloudplatform/pkg/organization/businesstype/models"
	businessTypeRepositories "smlcloudplatform/pkg/organization/businesstype/repositories"
	businessTypeServices "smlcloudplatform/pkg/organization/businesstype/services"
	deparmentRepositories "smlcloudplatform/pkg/organization/department/repositories"
	"smlcloudplatform/pkg/shop/models"
	"smlcloudplatform/pkg/utils"
	"time"

	warehouseModels "smlcloudplatform/pkg/warehouse/models"
	warehouseRepositories "smlcloudplatform/pkg/warehouse/repositories"
	warehouseServices "smlcloudplatform/pkg/warehouse/services"
)

type IShopHttp interface {
	RegisterHttp()
	CreateShop(ctx microservice.IContext) error
	UpdateShop(ctx microservice.IContext) error
	DeleteShop(ctx microservice.IContext) error
	InfoShop(ctx microservice.IContext) error
	SearchShop(ctx microservice.IContext) error
}

type ShopHttp struct {
	ms                  *microservice.Microservice
	cfg                 config.IConfig
	service             IShopService
	serviceBranch       branchServices.IBranchHttpService
	serviceWarehouse    warehouseServices.IWarehouseHttpService
	servicebusinessType businessTypeServices.IBusinessTypeHttpService
	authService         *microservice.AuthService
}

func NewShopHttp(ms *microservice.Microservice, cfg config.IConfig) ShopHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	repo := NewShopRepository(pst)
	cache := ms.Cacher(cfg.CacherConfig())
	shopUserRepo := NewShopUserRepository(pst)
	service := NewShopService(repo, shopUserRepo, utils.NewGUID, ms.TimeNow)

	authService := microservice.NewAuthService(ms.Cacher(cfg.CacherConfig()), 24*3*time.Hour, 24*30*time.Hour)

	repoBrach := branchRepositories.NewBranchRepository(pst)

	repoDepartment := deparmentRepositories.NewDepartmentRepository(pst)
	repoBusinessType := businessTypeRepositories.NewBusinessTypeRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	serviceBranch := branchServices.NewBranchHttpService(repoBrach, repoDepartment, repoBusinessType, masterSyncCacheRepo)

	serviceBusinessType := businessTypeServices.NewBusinessTypeHttpService(repoBusinessType, masterSyncCacheRepo)

	repoWarehouse := warehouseRepositories.NewWarehouseRepository(pst)
	svcWarehouse := warehouseServices.NewWarehouseHttpService(repoWarehouse, masterSyncCacheRepo)

	return ShopHttp{
		ms:                  ms,
		cfg:                 cfg,
		service:             service,
		serviceBranch:       serviceBranch,
		serviceWarehouse:    svcWarehouse,
		servicebusinessType: serviceBusinessType,
		authService:         authService,
	}
}

func (h ShopHttp) RegisterHttp() {
	h.ms.GET("/shop/:id", h.InfoShop)
	// h.ms.GET("/shop", h.SearchShop)

	h.ms.POST("/shop", h.CreateShop, h.authService.MWFuncWithShop(h.ms.Cacher(h.cfg.CacherConfig())))
	h.ms.PUT("/shop/:id", h.UpdateShop)
	h.ms.DELETE("/shop/:id", h.DeleteShop)
}

// Create Shop On login  godoc
// @Description Create Shop on login
// @Tags		Authentication
// @Accept 		json
// @Param		Shop  body      models.Shop  true  "Add Shop"
// @Success		200	{object}		models.Shop
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /create-shop [post]
func Docs() {

}

// Create Shop godoc
// @Description Create Shop
// @Tags		Shop
// @Accept 		json
// @Param		Shop  body      models.Shop  true  "Add Shop"
// @Success		200	{object}		models.Shop
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop [post]
func (h ShopHttp) CreateShop(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	if len(authUsername) < 1 {
		ctx.ResponseError(400, "user authentication invalid")
		return nil
	}

	input := ctx.ReadInput()

	shopPayload := &models.ShopRequest{}
	err := json.Unmarshal([]byte(input), &shopPayload)

	if err != nil {
		ctx.ResponseError(400, "shop payload invalid")
		return err
	}

	shopTemp := shopPayload.Shop

	shopID, err := h.service.CreateShop(authUsername, shopTemp)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	err = h.initialShop(shopID, authUsername, *shopPayload)

	if err != nil {
		err = h.service.DeleteShop(shopID, authUsername)

		if err != nil {
			logger.GetLogger().Error("HTTP:: Error Rollback Shop " + err.Error())
		}

		ctx.Response(http.StatusBadRequest, &common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, &common.ApiResponse{
		Success: true,
		ID:      shopID,
	})

	return nil
}
func (h ShopHttp) initialShop(shopID string, authUsername string, shopReq models.ShopRequest) (err error) {

	branchDefault := branchModel.Branch{}

	if len(shopReq.Settings.LanguageConfigs) > 0 {
		primaryLanguageConfigs := shopReq.Settings.LanguageConfigs[0]

		for _, langConf := range shopReq.Settings.LanguageConfigs {
			if langConf.IsDefault {
				primaryLanguageConfigs = langConf
				break
			}
		}

		for _, tempName := range shopReq.Names {
			if *tempName.Code == primaryLanguageConfigs.Code {
				branchDefault.CompanyNames = &[]common.NameX{
					{
						Code: tempName.Code,
						Name: tempName.Name,
					},
				}

				break
			}
		}

	}

	branchDefault.Code = "00000"

	branchMainCodeTH := "th"
	branchMainNameTH := "สำนักงานใหญ่"

	branchMainCodeEN := "en"
	branchMainNameEN := "Head Office"

	branchDefault.Names = &[]common.NameX{
		{
			Code: &branchMainCodeTH,
			Name: &branchMainNameTH,
		},
		{
			Code: &branchMainCodeEN,
			Name: &branchMainNameEN,
		},
	}

	branchGUIDFixed, err := h.serviceBranch.CreateBranch(shopID, authUsername, branchDefault)

	if err != nil {
		return err
	}

	businessTypeDefault := businessTypeModels.BusinessType{}

	businessTypeDefault.Code = shopReq.BusinessType.Code
	businessTypeDefault.Names = shopReq.BusinessType.Names
	businessTypeDefault.IsDefault = true

	if len(businessTypeDefault.Code) < 1 {

		businessTypeDefault.Code = "00000"

		businessTypeMainCodeTH := "th"
		businessTypeMainNameTH := "ธุรกิจหลัก"

		businessTypeMainCodeEN := "en"
		businessTypeMainNameEN := "Main Business"

		businessTypeDefault.Names = &[]common.NameX{
			{
				Code: &businessTypeMainCodeTH,
				Name: &businessTypeMainNameTH,
			},
			{
				Code: &businessTypeMainCodeEN,
				Name: &businessTypeMainNameEN,
			},
		}
	}

	businessTypeGUIDFixed, err := h.servicebusinessType.CreateBusinessType(shopID, authUsername, businessTypeDefault)

	if err != nil {
		err = h.serviceBranch.DeleteBranch(shopID, branchGUIDFixed, authUsername)

		if err != nil {
			logger.GetLogger().Error("HTTP:: Error Rollback Branch " + err.Error())
		}

		return err
	}

	warehouseDefault := warehouseModels.Warehouse{}
	warehouseDefault.Code = "00000"

	warehouseMainCodeTH := "th"
	warehouseMainNameTH := "สำนักงานใหญ่"

	warehouseMainCodeEN := "en"
	warehouseMainNameEN := "Head Office"

	warehouseDefault.Names = &[]common.NameX{
		{
			Code: &warehouseMainCodeTH,
			Name: &warehouseMainNameTH,
		},
		{
			Code: &warehouseMainCodeEN,
			Name: &warehouseMainNameEN,
		},
	}

	_, err = h.serviceWarehouse.CreateWarehouse(shopID, authUsername, warehouseDefault)

	if err != nil {

		err = h.serviceBranch.DeleteBranch(shopID, branchGUIDFixed, authUsername)

		if err != nil {
			logger.GetLogger().Error("HTTP:: Error Rollback Branch " + err.Error())
		}

		err = h.servicebusinessType.DeleteBusinessType(shopID, businessTypeGUIDFixed, authUsername)

		if err != nil {
			logger.GetLogger().Error("HTTP:: Error Rollback BusinessType " + err.Error())
		}

		return err
	}

	return nil
}

// Update Shop godoc
// @Description Update Shop
// @Tags		Shop
// @Accept 		json
// @Param		id	path     string  true  "Shop ID"
// @Param		Shop  body      models.Shop  true  "Shop Body"
// @Success		200	{object}		models.Shop
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/{id} [put]
func (h ShopHttp) UpdateShop(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	id := ctx.Param("id")
	input := ctx.ReadInput()

	shopRequest := &models.Shop{}
	err := json.Unmarshal([]byte(input), &shopRequest)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &common.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	err = h.service.UpdateShop(id, authUsername, *shopRequest)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, &common.ApiResponse{
		Success: true,
		ID:      id,
	})
	return nil
}

// Delete Shop godoc
// @Description Delete Shop
// @Tags		Shop
// @Accept 		json
// @Param		id	path     string  true  "Shop ID"
// @Success		200	{object}		models.Shop
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/{id} [delete]
func (h ShopHttp) DeleteShop(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username

	id := ctx.Param("id")

	if userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &common.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	err := h.service.DeleteShop(id, authUsername)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}
	ctx.Response(http.StatusOK, &common.ApiResponse{
		Success: true,
		ID:      id,
	})
	return nil
}

// Info Shop godoc
// @Description Infomation Shop Profile
// @Tags		Shop
// @Accept 		json
// @Param		id	path     string  true  "Shop ID"
// @Success		200	{array}	models.ShopInfo
// @Failure		401 {object}	common.ApiResponse
// @Security     AccessToken
// @Router /shop/{id} [get]
func (h ShopHttp) InfoShop(ctx microservice.IContext) error {
	id := ctx.Param("id")

	shopInfo, err := h.service.InfoShop(id)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, &common.ApiResponse{
		Success: true,
		Data:    shopInfo,
	})
	return nil
}

// List Shop godoc
// @Description Access to Shop
// @Tags		Shop
// @Accept 		json
// @Success		200	{array}	models.ShopInfo
// @Failure		401 {object}	common.ApiResponse
// @Security     AccessToken
// @Router /shop [get]
func (h ShopHttp) SearchShop(ctx microservice.IContext) error {

	pageable := utils.GetPageable(ctx.QueryParam)

	shopList, pagination, err := h.service.SearchShop(pageable)

	if err != nil {
		ctx.ResponseError(400, "database error")
		h.ms.Logger.Error("HTTP:: SearchShop " + err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": shopList})
	return nil
}
