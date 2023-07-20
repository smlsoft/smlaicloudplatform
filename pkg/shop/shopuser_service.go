package shop

import (
	"context"
	"errors"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/shop/models"
	"smlcloudplatform/pkg/utils"

	"github.com/userplant/mongopagination"
)

type IShopUserService interface {
	SaveUserPermissionShop(shopID string, authUsername string, editusername string, username string, role models.UserRole) error
	DeleteUserPermissionShop(shopID string, authUsername string, username string) error

	InfoShopByUser(shopID string, username string) (models.ShopUserProfile, error)
	ListShopByUser(authUsername string, pageable micromodels.Pageable) ([]models.ShopUserInfo, mongopagination.PaginationData, error)
	ListUserInShop(shopID string, pageable micromodels.Pageable) ([]models.ShopUserProfile, mongopagination.PaginationData, error)
}

type ShopUserService struct {
	repo IShopUserRepository
}

func NewShopUserService(shopUserRepo IShopUserRepository) ShopUserService {
	return ShopUserService{
		repo: shopUserRepo,
	}
}

func (svc ShopUserService) InfoShopByUser(shopID string, username string) (models.ShopUserProfile, error) {

	shopUserProfile := models.ShopUserProfile{}

	shopUser, err := svc.repo.FindByShopIDAndUsernameInfo(context.Background(), shopID, username)

	if err != nil {
		return models.ShopUserProfile{}, err
	}

	userProfiles, err := svc.repo.FindUserProfileByUsernames(context.Background(), []string{username})

	if err != nil {
		return models.ShopUserProfile{}, err
	}

	shopUserProfile.ShopID = shopUser.ShopID
	shopUserProfile.Username = username
	shopUserProfile.Role = shopUser.Role

	if len(userProfiles) > 0 {
		shopUserProfile.UserProfileName = userProfiles[0].Name
	}

	return shopUserProfile, err
}

func (svc ShopUserService) ListShopByUser(authUsername string, pageable micromodels.Pageable) ([]models.ShopUserInfo, mongopagination.PaginationData, error) {

	docList, pagination, err := svc.repo.FindByUsernamePage(context.Background(), authUsername, pageable)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, err
}

func (svc ShopUserService) ListUserInShop(shopID string, pageable micromodels.Pageable) ([]models.ShopUserProfile, mongopagination.PaginationData, error) {
	shopUserProfiles := []models.ShopUserProfile{}

	shopUsers, pagination, err := svc.repo.FindByUserInShopPage(context.Background(), shopID, pageable)

	if err != nil {
		return shopUserProfiles, pagination, err
	}

	usernames := []string{}

	for _, doc := range shopUsers {
		usernames = append(usernames, doc.Username)
	}

	userProfiles, err := svc.repo.FindUserProfileByUsernames(context.Background(), usernames)

	if err != nil {
		return shopUserProfiles, pagination, err
	}

	for _, doc := range shopUsers {
		shopUserProfile := models.ShopUserProfile{}

		shopUserProfile.ShopID = doc.ShopID
		shopUserProfile.Username = doc.Username
		shopUserProfile.Role = doc.Role

		shopUserProfiles = append(shopUserProfiles, shopUserProfile)
	}

	dictUserProfiles := map[string]models.UserProfile{}
	for _, doc := range userProfiles {
		dictUserProfiles[doc.Username] = doc
	}

	for idx, doc := range userProfiles {
		tempUserProfile := dictUserProfiles[doc.Username]

		shopUserProfiles[idx].UserProfileName = tempUserProfile.Name
	}

	return shopUserProfiles, pagination, err
}

func (svc ShopUserService) SaveUserPermissionShop(shopID string, authUsername string, editusername string, username string, role models.UserRole) error {

	username = utils.NormalizeUsername(username)

	if authUsername == editusername {
		return errors.New("can not edit self permission")
	}

	authUser, err := svc.repo.FindByShopIDAndUsername(context.Background(), shopID, authUsername)

	if err != nil {
		return err
	}

	if authUser.Role != models.ROLE_OWNER {
		return errors.New("permission denied")
	}

	editusername = utils.NormalizeUsername(editusername)

	if editusername != "" {

		findEditUser, err := svc.repo.FindByShopIDAndUsername(context.Background(), shopID, editusername)

		if err != nil {
			return err
		}

		if findEditUser.Username != "" {
			tempID := findEditUser.ID

			err = svc.repo.Update(context.Background(), tempID, shopID, username, role)
			if err != nil {
				return err
			}
		} else {
			err = svc.repo.Save(context.Background(), shopID, username, role)

			if err != nil {
				return err
			}
		}

	} else {

		err = svc.repo.Save(context.Background(), shopID, username, role)

		if err != nil {
			return err
		}
	}
	return nil
}

func (svc ShopUserService) create(ctx context.Context, shopID string, username string, role models.UserRole) error {

	tempShopUser := models.ShopUser{}
	tempShopUser.ShopID = shopID
	tempShopUser.Username = username
	tempShopUser.Role = role

	err := svc.repo.Create(ctx, &tempShopUser)
	if err != nil {
		return err
	}
	return nil
}

func (svc ShopUserService) DeleteUserPermissionShop(shopID string, authUsername string, username string) error {

	authUser, err := svc.repo.FindByShopIDAndUsername(context.Background(), shopID, authUsername)

	if err != nil {
		return err
	}

	if authUser.Role != models.ROLE_OWNER && authUser.Role != models.ROLE_ADMIN {
		return errors.New("permission denied")
	}

	findUser, err := svc.repo.FindByShopIDAndUsername(context.Background(), shopID, username)

	if err != nil {
		return err
	}

	if authUser.Role == models.ROLE_ADMIN && findUser.Role == models.ROLE_OWNER {
		return errors.New("permission denied")
	}

	if findUser.Username == authUsername {
		return errors.New("can't delete your permission")
	}

	err = svc.repo.Delete(context.Background(), shopID, username)

	if err != nil {
		return err
	}
	return nil
}
