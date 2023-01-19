package shop

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/shop/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopRepository interface {
	Create(shop models.ShopDoc) (string, error)
	Update(guid string, shop models.ShopDoc) error
	FindByGuid(guid string) (models.ShopDoc, error)
	FindPage(pageable micromodels.Pageable) ([]models.ShopInfo, mongopagination.PaginationData, error)
	Delete(guid string, username string) error
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
	filterDoc := map[string]interface{}{
		"guidfixed": guid,
	}
	err := repo.pst.UpdateOne(&models.ShopDoc{}, filterDoc, shop)

	if err != nil {
		return err
	}

	return nil
}

func (repo ShopRepository) FindByGuid(guid string) (models.ShopDoc, error) {
	findShop := &models.ShopDoc{}
	err := repo.pst.FindOne(&models.ShopDoc{}, bson.M{"guidfixed": guid, "deletedat": bson.M{"$exists": false}}, findShop)

	if err != nil {
		return models.ShopDoc{}, err
	}
	return *findShop, err
}

func (repo ShopRepository) FindPage(pageable micromodels.Pageable) ([]models.ShopInfo, mongopagination.PaginationData, error) {
	filterQueries := bson.M{
		"deletedat": bson.M{"$exists": false},
		"name1": bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + pageable.Query + ".*",
			Options: "",
		}}}

	shopList := []models.ShopInfo{}

	pagination, err := repo.pst.FindPage(&models.ShopInfo{}, filterQueries, pageable, &shopList)

	if err != nil {
		return []models.ShopInfo{}, mongopagination.PaginationData{}, err
	}

	return shopList, pagination, nil
}

func (repo ShopRepository) Delete(guid string, username string) error {
	err := repo.pst.SoftDeleteByID(&models.ShopInfo{}, guid, username)
	if err != nil {
		return err
	}
	return nil
}
