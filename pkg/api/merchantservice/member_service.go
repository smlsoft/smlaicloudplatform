package merchantservice

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMemberService interface {
	Save(merchantId string, username string, role models.UserRole) error
	Delete(merchantId string, username string) error
	FindRole(merchantId string, username string) (models.UserRole, error)
	FindByMerchantId(merchantId string) (*[]models.MerchantUser, error)
	FindByUsername(username string) (*[]models.MerchantUser, error)
}

type MemberService struct {
	pst microservice.IPersisterMongo
}

func NewMemberService(pst microservice.IPersisterMongo) *MemberService {
	return &MemberService{
		pst: pst,
	}
}

func (svc *MemberService) Save(merchantId string, username string, role models.UserRole) error {

	optUpdate := options.Update().SetUpsert(true)
	err := svc.pst.Update(&models.MerchantUser{}, bson.M{"merchantId": merchantId, "username": username}, bson.M{"$set": bson.M{"role": role}}, optUpdate)

	if err != nil {
		return err
	}

	return nil
}

func (svc *MemberService) Delete(merchantId string, username string) error {

	err := svc.pst.Delete(&models.MerchantUser{}, bson.M{"merchantId": merchantId, "username": username})

	if err != nil {
		return err
	}

	return nil
}

func (svc *MemberService) FindRole(merchantId string, username string) (models.UserRole, error) {

	merchantUser := &models.MerchantUser{}

	err := svc.pst.FindOne(&models.MerchantUser{}, bson.M{"merchantId": merchantId, "username": username}, merchantUser)

	if err != nil {
		return "", err
	}

	return merchantUser.Role, nil
}

func (svc *MemberService) FindByMerchantId(merchantId string) (*[]models.MerchantUser, error) {
	merchantUsers := &[]models.MerchantUser{}

	err := svc.pst.Find(&models.MerchantUser{}, bson.M{"merchantId": merchantId}, merchantUsers)

	if err != nil {
		return nil, err
	}

	return merchantUsers, nil

}

func (svc *MemberService) FindByUsername(username string) (*[]models.MerchantUser, error) {
	merchantUsers := &[]models.MerchantUser{}

	err := svc.pst.Find(&models.MerchantUser{}, bson.M{"username": username}, merchantUsers)

	if err != nil {
		return nil, err
	}

	return merchantUsers, nil

}
