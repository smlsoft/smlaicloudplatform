package shop

import (
	"errors"
	"smlcloudplatform/pkg/shop/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopService interface {
	CreateShop(username string, shop models.Shop) (string, error)
	UpdateShop(guid string, username string, shop models.Shop) error
	DeleteShop(guid string, username string) error
	InfoShop(guid string) (models.ShopInfo, error)
	SearchShop(q string, page int, limit int) ([]models.ShopInfo, mongopagination.PaginationData, error)
}

type ShopService struct {
	shopRepo     IShopRepository
	shopUserRepo IShopUserRepository
	newGUID      func() string
	timeNow      func() time.Time
}

func NewShopService(shopRepo IShopRepository, shopUserRepo IShopUserRepository, newGUID func() string, timeNow func() time.Time) ShopService {
	return ShopService{
		shopRepo:     shopRepo,
		shopUserRepo: shopUserRepo,
		newGUID:      newGUID,
		timeNow:      timeNow,
	}
}

func (svc ShopService) CreateShop(username string, doc models.Shop) (string, error) {

	dataDoc := models.ShopDoc{}
	shopID := svc.newGUID()
	dataDoc.GuidFixed = shopID
	dataDoc.CreatedBy = username
	dataDoc.CreatedAt = svc.timeNow()
	dataDoc.Shop = doc

	_, err := svc.shopRepo.Create(dataDoc)

	if err != nil {
		return "", err
	}

	err = svc.shopUserRepo.Save(shopID, username, models.ROLE_OWNER)

	if err != nil {
		return "", err
	}

	return shopID, nil
}

func (svc ShopService) UpdateShop(guid string, username string, shop models.Shop) error {

	findShop, err := svc.shopRepo.FindByGuid(guid)

	if err != nil {
		return err
	}

	if findShop.ID == primitive.NilObjectID {
		return errors.New("shop not found")
	}

	guidx := findShop.GuidFixed

	findShop.UpdatedBy = username
	findShop.UpdatedAt = time.Now()

	findShop.Name1 = shop.Name1

	findShop.GuidFixed = guidx

	err = svc.shopRepo.Update(guid, findShop)

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopService) DeleteShop(guid string, username string) error {

	err := svc.shopRepo.Delete(guid, username)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopService) InfoShop(guid string) (models.ShopInfo, error) {
	findShop, err := svc.shopRepo.FindByGuid(guid)

	if err != nil {
		return models.ShopInfo{}, err
	}

	return findShop.ShopInfo, nil
}

func (svc ShopService) SearchShop(q string, page int, limit int) ([]models.ShopInfo, mongopagination.PaginationData, error) {
	shopList, pagination, err := svc.shopRepo.FindPage(q, page, limit)

	if err != nil {
		return shopList, pagination, err
	}

	return shopList, pagination, nil
}
