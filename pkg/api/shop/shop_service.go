package shop

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopService interface {
	CreateShop(username string, shop models.Shop) (string, error)
	UpdateShop(guid string, username string, shop models.Shop) error
	DeleteShop(guid string, username string) error
	InfoShop(guid string) (models.ShopInfo, error)
	SearchShop(q string, page int, limit int) ([]models.ShopInfo, paginate.PaginationData, error)
}

type ShopService struct {
	shopRepo     IShopRepository
	shopUserRepo IShopUserRepository
}

func NewShopService(shopRepo IShopRepository, shopUserRepo IShopUserRepository) IShopService {
	return &ShopService{
		shopRepo:     shopRepo,
		shopUserRepo: shopUserRepo,
	}
}

func (svc *ShopService) CreateShop(username string, shop models.Shop) (string, error) {

	shopId := utils.NewGUID()
	shop.GuidFixed = shopId
	shop.CreatedBy = username
	shop.CreatedAt = time.Now()

	_, err := svc.shopRepo.Create(shop)

	if err != nil {
		return "", err
	}

	svc.shopUserRepo.Save(shopId, username, models.ROLE_OWNER)

	return shopId, nil
}

func (svc *ShopService) UpdateShop(guid string, username string, shop models.Shop) error {

	findShop, err := svc.shopRepo.FindByGuid(guid)

	if err != nil {
		return err
	}

	if findShop.Id == primitive.NilObjectID {
		return errors.New("shop not found")
	}

	findShop.Name1 = shop.Name1
	findShop.UpdatedBy = username
	findShop.UpdatedAt = time.Now()

	err = svc.shopRepo.Update(guid, findShop)

	if err != nil {
		return err
	}

	return nil
}

func (svc *ShopService) DeleteShop(guid string, username string) error {

	err := svc.shopRepo.Delete(guid)

	if err != nil {
		return err
	}
	return nil
}

func (svc *ShopService) InfoShop(guid string) (models.ShopInfo, error) {
	findShop, err := svc.shopRepo.FindByGuid(guid)

	if err != nil {
		return models.ShopInfo{}, err
	}

	return models.ShopInfo{
		Id:        findShop.Id,
		GuidFixed: findShop.GuidFixed,
		Name1:     findShop.Name1,
	}, nil
}

func (svc *ShopService) SearchShop(q string, page int, limit int) ([]models.ShopInfo, paginate.PaginationData, error) {
	shopList, pagination, err := svc.shopRepo.FindPage(q, page, limit)

	if err != nil {
		return shopList, pagination, err
	}

	return shopList, pagination, nil
}
