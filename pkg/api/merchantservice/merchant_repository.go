package merchantservice

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
)

type IMerchantRepository interface {
	Create(merchant models.Merchant) (string, error)
	Update(guid string, merchant models.Merchant) error
	FindByGuid(guid string) (models.Merchant, error)
	Delete(guid string) error
}

type MerchantRepository struct {
	pst microservice.IPersisterMongo
}

func NewMerchantRepository(pst microservice.IPersisterMongo) IMerchantRepository {
	return &MerchantRepository{
		pst: pst,
	}
}

func (repo *MerchantRepository) Create(merchant models.Merchant) (string, error) {
	idx, err := repo.pst.Create(&models.Merchant{}, merchant)
	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo *MerchantRepository) Update(guid string, merchant models.Merchant) error {
	err := repo.pst.UpdateOne(&models.Merchant{}, "guidFixed", guid, merchant)

	if err != nil {
		return err
	}

	return nil
}

func (repo *MerchantRepository) FindByGuid(guid string) (models.Merchant, error) {
	findMerchant := &models.Merchant{}
	err := repo.pst.FindOne(&models.Merchant{}, bson.M{"guidFixed": guid, "deleted": false}, findMerchant)

	if err != nil {
		return models.Merchant{}, err
	}
	return *findMerchant, err

}

func (repo *MerchantRepository) Delete(guid string) error {
	err := repo.pst.SoftDeleteByID(&models.Merchant{}, guid)
	if err != nil {
		return err
	}
	return nil
}
