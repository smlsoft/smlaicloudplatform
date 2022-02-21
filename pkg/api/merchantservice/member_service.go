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
	FindByMerchantId(merchantId string) (*[]models.MerchantMember, error)
	FindByUsername(username string) (*[]models.MerchantMember, error)
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
	err := svc.pst.Update(&models.MerchantMember{}, bson.M{"merchantId": merchantId, "username": username}, bson.M{"$set": bson.M{"role": role}}, optUpdate)

	if err != nil {
		return err
	}

	return nil
}

func (svc *MemberService) Delete(merchantId string, username string) error {

	err := svc.pst.Delete(&models.MerchantMember{}, bson.M{"merchantId": merchantId, "username": username})

	if err != nil {
		return err
	}

	return nil
}

func (svc *MemberService) FindRole(merchantId string, username string) (models.UserRole, error) {

	merchantMember := &models.MerchantMember{}

	err := svc.pst.FindOne(&models.MerchantMember{}, bson.M{"merchantId": merchantId, "username": username}, merchantMember)

	if err != nil {
		return "", err
	}

	return merchantMember.Role, nil
}

func (svc *MemberService) FindByMerchantId(merchantId string) (*[]models.MerchantMember, error) {
	merchantMembers := &[]models.MerchantMember{}

	err := svc.pst.Find(&models.MerchantMember{}, bson.M{"merchantId": merchantId}, merchantMembers)

	if err != nil {
		return nil, err
	}

	return merchantMembers, nil

}

func (svc *MemberService) FindByUsername(username string) (*[]models.MerchantMember, error) {
	merchantMembers := &[]models.MerchantMember{}

	err := svc.pst.Find(&models.MerchantMember{}, bson.M{"username": username}, merchantMembers)

	if err != nil {
		return nil, err
	}

	return merchantMembers, nil

}
