package utils

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
)

func HasPermissionMerchant(pst microservice.IPersisterMongo, ctx microservice.IServiceContext) (bool, error) {

	merchantId := ctx.Param("merchant_id")

	return HasPermissionMerchantById(pst, ctx, merchantId)
}

func HasPermissionMerchantById(pst microservice.IPersisterMongo, ctx microservice.IServiceContext, merchantId string) (bool, error) {

	authUsername := ctx.UserInfo().Username

	if len(merchantId) < 1 {
		return false, fmt.Errorf("merchant not found")
	}

	merchant := &models.Merchant{}
	pst.FindOne(&models.Merchant{}, bson.M{"guidFixed": merchantId, "deleted": false}, merchant)

	if len(merchant.GuidFixed) < 1 {
		return false, fmt.Errorf("merchant invalid")
	}

	if merchant.CreatedBy != authUsername {
		return false, fmt.Errorf("username invalid")
	}

	return true, nil
}
