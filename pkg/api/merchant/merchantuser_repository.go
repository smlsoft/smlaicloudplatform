package merchant

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMerchantUserRepository interface {
	Save(merchantId string, username string, role string) error
	Delete(merchantId string, username string) error
	FindByMerchantIdAndUsername(merchantId string, username string) (models.MerchantUser, error)
	FindRole(merchantId string, username string) (string, error)
	FindByMerchantId(merchantId string) (*[]models.MerchantUser, error)
	FindByUsername(username string) (*[]models.MerchantUser, error)
}

type MerchantUserRepository struct {
	pst microservice.IPersisterMongo
}

func NewMerchantUserRepository(pst microservice.IPersisterMongo) IMerchantUserRepository {
	return &MerchantUserRepository{
		pst: pst,
	}
}

func (svc *MerchantUserRepository) Save(merchantId string, username string, role string) error {

	optUpdate := options.Update().SetUpsert(true)
	err := svc.pst.Update(&models.MerchantUser{}, bson.M{"merchantId": merchantId, "username": username}, bson.M{"$set": bson.M{"role": role}}, optUpdate)

	if err != nil {
		return err
	}

	return nil
}

func (svc *MerchantUserRepository) Delete(merchantId string, username string) error {

	err := svc.pst.Delete(&models.MerchantUser{}, bson.M{"merchantId": merchantId, "username": username})

	if err != nil {
		return err
	}

	return nil
}

func (svc *MerchantUserRepository) FindByMerchantIdAndUsername(merchantId string, username string) (models.MerchantUser, error) {

	merchantUser := &models.MerchantUser{}

	err := svc.pst.FindOne(&models.MerchantUser{}, bson.M{"merchantId": merchantId, "username": username}, merchantUser)
	if err != nil {
		fmt.Println("err -> ", err.Error())
		return models.MerchantUser{}, err
	}

	return *merchantUser, nil
}

func (svc *MerchantUserRepository) FindRole(merchantId string, username string) (string, error) {

	merchantUser := &models.MerchantUser{}

	err := svc.pst.FindOne(&models.MerchantUser{}, bson.M{"merchantId": merchantId, "username": username}, merchantUser)

	if err != nil {
		return "", err
	}

	return merchantUser.Role, nil
}

func (svc *MerchantUserRepository) FindByMerchantId(merchantId string) (*[]models.MerchantUser, error) {
	merchantUsers := &[]models.MerchantUser{}

	err := svc.pst.Find(&models.MerchantUser{}, bson.M{"merchantId": merchantId}, merchantUsers)

	if err != nil {
		return nil, err
	}

	return merchantUsers, nil
}

func (svc *MerchantUserRepository) FindByUsername(username string) (*[]models.MerchantUser, error) {
	merchantUsers := &[]models.MerchantUser{}

	err := svc.pst.Find(&models.MerchantUser{}, bson.M{"username": username}, merchantUsers)

	if err != nil {
		return nil, err
	}

	return merchantUsers, nil
}
