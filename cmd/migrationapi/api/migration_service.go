package api

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shop"
	"smlcloudplatform/pkg/shop/models"
	shopModel "smlcloudplatform/pkg/shop/models"
	"smlcloudplatform/pkg/utils"
	journalModels "smlcloudplatform/pkg/vfgl/journal/models"
	journalRepo "smlcloudplatform/pkg/vfgl/journal/repositories"

	"go.mongodb.org/mongo-driver/bson"
)

type IMigrationService interface {
	ImportShop(shops []shopModel.ShopDoc) error
	ImportJournal(journals []journalModels.JournalDoc) error
}

type MigrationService struct {
	mongoPersister microservice.IPersisterMongo
	mqPersister    microservice.IProducer
	userShopRepo   shop.IShopUserRepository
}

func NewMigrationService(mongoPersister microservice.IPersisterMongo, mqPersister microservice.IProducer) *MigrationService {

	userShopRepo := shop.NewShopUserRepository(mongoPersister)

	return &MigrationService{
		mongoPersister: mongoPersister,
		mqPersister:    mqPersister,
		userShopRepo:   userShopRepo,
	}
}

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

func (m *MigrationService) ImportJournal(journals []journalModels.JournalDoc) error {
	journalpgRepo := journalRepo.NewJournalRepository(m.mongoPersister)
	journalMQRepo := journalRepo.NewJournalMqRepository(m.mqPersister)

	for _, journal := range journals {
		journal.GuidFixed = utils.NewGUID()
		_, err := journalpgRepo.Create(journal)
		if err != nil {
			return err
		}

		err = journalMQRepo.Create(journal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MigrationService) ImportShopUser(shop shopModel.ShopDoc) error {

	username := shop.GuidFixed
	findUser := &shopModel.UserDoc{}
	userPassword, err := utils.HashPassword(username)

	if err != nil {
		return err
	}
	err = m.mongoPersister.FindOne(&models.UserDoc{}, bson.M{"username": username}, findUser)

	if err != nil {
		return err
	}

	if findUser.Username == "" {
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
