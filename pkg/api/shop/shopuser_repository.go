package shop

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IShopUserRepository interface {
	Save(shopId string, username string, role string) error
	Delete(shopId string, username string) error
	FindByShopIdAndUsername(shopId string, username string) (models.ShopUser, error)
	FindRole(shopId string, username string) (string, error)
	FindByShopId(shopId string) (*[]models.ShopUser, error)
	FindByUsername(username string) (*[]models.ShopUser, error)
	FindByUsernamePage(username string, page int, limit int) ([]models.ShopUserInfo, paginate.PaginationData, error)
}

type ShopUserRepository struct {
	pst microservice.IPersisterMongo
}

func NewShopUserRepository(pst microservice.IPersisterMongo) IShopUserRepository {
	return &ShopUserRepository{
		pst: pst,
	}
}

func (svc *ShopUserRepository) Save(shopId string, username string, role string) error {

	optUpdate := options.Update().SetUpsert(true)
	err := svc.pst.Update(&models.ShopUser{}, bson.M{"shopId": shopId, "username": username}, bson.M{"$set": bson.M{"role": role}}, optUpdate)

	if err != nil {
		return err
	}

	return nil
}

func (svc *ShopUserRepository) Delete(shopId string, username string) error {

	err := svc.pst.Delete(&models.ShopUser{}, bson.M{"shopId": shopId, "username": username})

	if err != nil {
		return err
	}

	return nil
}

func (svc *ShopUserRepository) FindByShopIdAndUsername(shopId string, username string) (models.ShopUser, error) {

	shopUser := &models.ShopUser{}

	err := svc.pst.FindOne(&models.ShopUser{}, bson.M{"shopId": shopId, "username": username}, shopUser)
	if err != nil {
		fmt.Println("err -> ", err.Error())
		return models.ShopUser{}, err
	}

	return *shopUser, nil
}

func (svc *ShopUserRepository) FindRole(shopId string, username string) (string, error) {

	shopUser := &models.ShopUser{}

	err := svc.pst.FindOne(&models.ShopUser{}, bson.M{"shopId": shopId, "username": username}, shopUser)

	if err != nil {
		return "", err
	}

	return shopUser.Role, nil
}

func (svc *ShopUserRepository) FindByShopId(shopId string) (*[]models.ShopUser, error) {
	shopUsers := &[]models.ShopUser{}

	err := svc.pst.Find(&models.ShopUser{}, bson.M{"shopId": shopId}, shopUsers)

	if err != nil {
		return nil, err
	}

	return shopUsers, nil
}

func (svc *ShopUserRepository) FindByUsername(username string) (*[]models.ShopUser, error) {
	shopUsers := &[]models.ShopUser{}

	err := svc.pst.Find(&models.ShopUser{}, bson.M{"username": username}, shopUsers)

	if err != nil {
		return nil, err
	}

	return shopUsers, nil
}

func (repo *ShopUserRepository) FindByUsernamePage(username string, page int, limit int) ([]models.ShopUserInfo, paginate.PaginationData, error) {

	docList := []models.ShopUserInfo{}

	aggPaginatedData, err := repo.pst.AggregatePage(&models.ShopUser{}, limit, page,
		bson.M{"$match": bson.M{"username": username}},
		bson.M{"$lookup": bson.M{
			"from":         "shop",
			"localField":   "shopId",
			"foreignField": "guidFixed",
			"as":           "shopInfo",
		}},
		bson.M{
			"$match": bson.M{"shopInfo.deleted": false},
		},
		bson.M{
			"$project": bson.M{
				"_id":    1,
				"role":   1,
				"shopId": 1,
				"name":   bson.M{"$first": "$shopInfo.name1"},
			},
		},
	)

	if err != nil {
		return []models.ShopUserInfo{}, paginate.PaginationData{}, err
	}

	for _, raw := range aggPaginatedData.Data {
		var doc *models.ShopUserInfo

		if marshallErr := bson.Unmarshal(raw, &doc); marshallErr == nil {
			docList = append(docList, *doc)
		}

	}

	return docList, aggPaginatedData.Pagination, nil
}
