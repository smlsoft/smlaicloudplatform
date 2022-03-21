package shop

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopRepository interface {
	Create(shop models.ShopDoc) (string, error)
	Update(guid string, shop models.ShopDoc) error
	FindByGuid(guid string) (models.ShopDoc, error)
	FindPage(q string, page int, limit int) ([]models.ShopInfo, paginate.PaginationData, error)
	Delete(guid string) error
}

type ShopRepository struct {
	pst microservice.IPersisterMongo
}

func NewShopRepository(pst microservice.IPersisterMongo) ShopRepository {
	return ShopRepository{
		pst: pst,
	}
}

func (repo ShopRepository) Create(shop models.ShopDoc) (string, error) {
	idx, err := repo.pst.Create(&models.ShopDoc{}, shop)
	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo ShopRepository) Update(guid string, shop models.ShopDoc) error {
	err := repo.pst.UpdateOne(&models.ShopDoc{}, "guidFixed", guid, shop)

	if err != nil {
		return err
	}

	return nil
}

func (repo ShopRepository) FindByGuid(guid string) (models.ShopDoc, error) {
	findShop := &models.ShopDoc{}
	err := repo.pst.FindOne(&models.ShopDoc{}, bson.M{"guidFixed": guid, "deleted": false}, findShop)

	if err != nil {
		return models.ShopDoc{}, err
	}
	return *findShop, err
}

func (repo ShopRepository) FindPage(q string, page int, limit int) ([]models.ShopInfo, paginate.PaginationData, error) {

	shopList := []models.ShopInfo{}

	pagination, err := repo.pst.FindPage(&models.ShopInfo{}, limit, page, bson.M{
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

func (repo ShopRepository) Delete(guid string) error {
	err := repo.pst.SoftDeleteByID(&models.ShopInfo{}, guid)
	if err != nil {
		return err
	}
	return nil
}
