package shop

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopRepository interface {
	Create(shop models.Shop) (string, error)
	Update(guid string, shop models.Shop) error
	FindByGuid(guid string) (models.Shop, error)
	FindPage(q string, page int, limit int) ([]models.ShopInfo, paginate.PaginationData, error)
	Delete(guid string) error
}

type ShopRepository struct {
	pst microservice.IPersisterMongo
}

func NewShopRepository(pst microservice.IPersisterMongo) IShopRepository {
	return &ShopRepository{
		pst: pst,
	}
}

func (repo *ShopRepository) Create(shop models.Shop) (string, error) {
	idx, err := repo.pst.Create(&models.Shop{}, shop)
	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo *ShopRepository) Update(guid string, shop models.Shop) error {
	err := repo.pst.UpdateOne(&models.Shop{}, "guidFixed", guid, shop)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ShopRepository) FindByGuid(guid string) (models.Shop, error) {
	findShop := &models.Shop{}
	err := repo.pst.FindOne(&models.Shop{}, bson.M{"guidFixed": guid, "deleted": false}, findShop)

	if err != nil {
		return models.Shop{}, err
	}
	return *findShop, err
}

func (repo *ShopRepository) FindPage(q string, page int, limit int) ([]models.ShopInfo, paginate.PaginationData, error) {

	shopList := []models.ShopInfo{}

	pagination, err := repo.pst.FindPage(&models.Shop{}, limit, page, bson.M{
		"deleted": false,
		"name1": bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}}}, &shopList)

	if err != nil {
		return []models.ShopInfo{}, paginate.PaginationData{}, err
	}

	return shopList, pagination, nil
}

func (repo *ShopRepository) Delete(guid string) error {
	err := repo.pst.SoftDeleteByID(&models.Shop{}, guid)
	if err != nil {
		return err
	}
	return nil
}
