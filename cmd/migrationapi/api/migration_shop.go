package api

import (
	"smlcloudplatform/pkg/shop"
	"smlcloudplatform/pkg/shop/models"
	shopModel "smlcloudplatform/pkg/shop/models"
	"smlcloudplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func (m *MigrationService) ImportShop(shops []shopModel.ShopDoc) error {
	shopRepo := shop.NewShopRepository(m.mongoPersister)
	for _, shop := range shops {

		findShopDoc, err := shopRepo.FindByGuid(shop.GuidFixed)

		if err != nil {

		}

		if findShopDoc.GuidFixed == "" {
			_, err = shopRepo.Create(shop)
			if err != nil {
				return err
			}
		}

		err = m.ImportShopUser(shop)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MigrationService) ImportShopUser(shop shopModel.ShopDoc) error {

	username := shop.GuidFixed
	findUser := &shopModel.UserDoc{}

	err := m.mongoPersister.FindOne(&models.UserDoc{}, bson.M{"username": username}, findUser)

	if err != nil {
		return err
	}

	if findUser.Username == "" {
		userPassword, err := utils.HashPassword(username)
		if err != nil {
			return err
		}
		newUser := shopModel.UserDoc{
			UsernameCode: shopModel.UsernameCode{
				Username: username,
			},
			UserPassword: shopModel.UserPassword{
				Password: userPassword,
			},
		}
		_, err = m.mongoPersister.Create(&shopModel.UserDoc{}, newUser)
		if err != nil {
			return err
		}
	}

	findUserShop, err := m.userShopRepo.FindByShopIDAndUsername(shop.GuidFixed, username)

	if findUserShop.Username == "" {

		err = m.userShopRepo.Save(shop.GuidFixed, username, models.ROLE_OWNER)

		if err != nil {
			return err
		}
	}

	return nil
}
