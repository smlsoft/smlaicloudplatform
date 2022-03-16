package utils

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
)

func HasPermissionShop(pst microservice.IPersisterMongo, ctx microservice.IContext) (bool, error) {

	shopId := ctx.Param("shop_id")

	return HasPermissionShopById(pst, ctx, shopId)
}

func HasPermissionShopById(pst microservice.IPersisterMongo, ctx microservice.IContext, shopId string) (bool, error) {

	authUsername := ctx.UserInfo().Username

	if len(shopId) < 1 {
		return false, fmt.Errorf("shop not found")
	}

	shop := &models.Shop{}
	pst.FindOne(&models.Shop{}, bson.M{"guidFixed": shopId, "deleted": false}, shop)

	if len(shop.GuidFixed) < 1 {
		return false, fmt.Errorf("shop invalid")
	}

	if shop.CreatedBy != authUsername {
		return false, fmt.Errorf("username invalid")
	}

	return true, nil
}
