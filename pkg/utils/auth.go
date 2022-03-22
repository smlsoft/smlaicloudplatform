package utils

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
)

func HasPermissionShop(pst microservice.IPersisterMongo, ctx microservice.IContext) (bool, error) {

	shopID := ctx.Param("shop_id")

	return HasPermissionShopByID(pst, ctx, shopID)
}

func HasPermissionShopByID(pst microservice.IPersisterMongo, ctx microservice.IContext, shopID string) (bool, error) {

	authUsername := ctx.UserInfo().Username

	if len(shopID) < 1 {
		return false, fmt.Errorf("shop not found")
	}

	shop := &models.ShopDoc{}
	pst.FindOne(&models.Shop{}, bson.M{"guidFixed": shopID, "deletedAt": bson.M{"$exists": false}}, shop)

	if len(shop.GuidFixed) < 1 {
		return false, fmt.Errorf("shop invalid")
	}

	if shop.CreatedBy != authUsername {
		return false, fmt.Errorf("username invalid")
	}

	return true, nil
}
