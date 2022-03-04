package merchant

import (
	"errors"
	"smlcloudplatform/pkg/models"
)

type IMerchantUserService interface {
	SaveUserPermissionMerchant(merchantId string, authUsername string, username string, role string) error
	DeleteUserPermissionMerchant(merchantId string, authUsername string, username string, guid string) error
}

type MerchantUserService struct {
	repo IMerchantUserRepository
}

func NewMerchantUserService(merchantUserRepo IMerchantUserRepository) IMerchantUserService {
	return &MerchantUserService{
		repo: merchantUserRepo,
	}
}

func (svc *MerchantUserService) SaveUserPermissionMerchant(merchantId string, authUsername string, username string, role string) error {

	authUser, err := svc.repo.FindByMerchantIdAndUsername(merchantId, authUsername)

	if err != nil {
		return err
	}

	if authUser.Role != models.ROLE_OWNER {
		return errors.New("permission denied")
	}

	err = svc.repo.Save(merchantId, username, role)

	if err != nil {
		return err
	}
	return nil
}

func (svc *MerchantUserService) DeleteUserPermissionMerchant(merchantId string, authUsername string, username string, guid string) error {

	authUser, err := svc.repo.FindByMerchantIdAndUsername(merchantId, authUsername)

	if err != nil {
		return err
	}

	if authUser.Role != models.ROLE_OWNER {
		return errors.New("permission denied")
	}

	err = svc.repo.Delete(merchantId, username)

	if err != nil {
		return err
	}
	return nil
}
