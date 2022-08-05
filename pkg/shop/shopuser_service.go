package shop

import (
	"errors"
	"smlcloudplatform/pkg/shop/models"

	paginate "github.com/gobeam/mongo-go-pagination"
)

type IShopUserService interface {
	ListShopByUser(authUsername string, page int, limit int) ([]models.ShopUserInfo, paginate.PaginationData, error)
	SaveUserPermissionShop(shopID string, authUsername string, username string, role string) error
	DeleteUserPermissionShop(shopID string, authUsername string, username string) error
	ListUserInShop(shopID string, page int, limit int) ([]models.ShopUser, paginate.PaginationData, error)
}

type ShopUserService struct {
	repo IShopUserRepository
}

func NewShopUserService(shopUserRepo IShopUserRepository) ShopUserService {
	return ShopUserService{
		repo: shopUserRepo,
	}
}

func (svc ShopUserService) ListShopByUser(authUsername string, page int, limit int) ([]models.ShopUserInfo, paginate.PaginationData, error) {

	docList, pagination, err := svc.repo.FindByUsernamePage(authUsername, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, err
}

func (svc ShopUserService) ListUserInShop(shopID string, page int, limit int) ([]models.ShopUser, paginate.PaginationData, error) {

	docList, pagination, err := svc.repo.FindByUserInShopPage(shopID, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, err
}

func (svc ShopUserService) SaveUserPermissionShop(shopID string, authUsername string, username string, role string) error {

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

	err = svc.repo.Delete(shopID, username)

	if err != nil {
		return err
	}
	return nil
}
