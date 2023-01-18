package shop

import (
	"errors"
	"smlcloudplatform/pkg/shop/models"

	"github.com/userplant/mongopagination"
)

type IShopUserService interface {
	SaveUserPermissionShop(shopID string, authUsername string, username string, role models.UserRole) error
	DeleteUserPermissionShop(shopID string, authUsername string, username string) error

	InfoShopByUser(shopID string, infoUsername string) (models.ShopUserInfo, error)
	ListShopByUser(authUsername string, q string, page int, limit int) ([]models.ShopUserInfo, mongopagination.PaginationData, error)
	ListUserInShop(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ShopUser, mongopagination.PaginationData, error)
}

type ShopUserService struct {
	repo IShopUserRepository
}

func NewShopUserService(shopUserRepo IShopUserRepository) ShopUserService {
	return ShopUserService{
		repo: shopUserRepo,
	}
}

func (svc ShopUserService) InfoShopByUser(shopID string, infoUsername string) (models.ShopUserInfo, error) {

	doc, err := svc.repo.FindByShopIDAndUsernameInfo(shopID, infoUsername)

	if err != nil {
		return models.ShopUserInfo{}, err
	}

	return doc, err
}

func (svc ShopUserService) ListShopByUser(authUsername string, q string, page int, limit int) ([]models.ShopUserInfo, mongopagination.PaginationData, error) {

	docList, pagination, err := svc.repo.FindByUsernamePage(authUsername, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, err
}

func (svc ShopUserService) ListUserInShop(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ShopUser, mongopagination.PaginationData, error) {

	docList, pagination, err := svc.repo.FindByUserInShopPage(shopID, q, page, limit, sort)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, err
}

func (svc ShopUserService) SaveUserPermissionShop(shopID string, authUsername string, username string, role models.UserRole) error {

	authUser, err := svc.repo.FindByShopIDAndUsername(shopID, authUsername)

	if err != nil {
		return err
	}

	if authUser.Role != models.ROLE_OWNER {
		return errors.New("permission denied")
	}

	// check username is exists

	err = svc.repo.Save(shopID, username, role)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopUserService) DeleteUserPermissionShop(shopID string, authUsername string, username string) error {

	authUser, err := svc.repo.FindByShopIDAndUsername(shopID, authUsername)

	if err != nil {
		return err
	}

	if authUser.Role != models.ROLE_OWNER {
		return errors.New("permission denied")
	}

	if authUser.Username == authUsername {
		return errors.New("can't delete your permission")
	}

	err = svc.repo.Delete(shopID, username)

	if err != nil {
		return err
	}
	return nil
}
