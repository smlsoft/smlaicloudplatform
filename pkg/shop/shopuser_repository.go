package shop

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shop/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IShopUserRepository interface {
	Save(shopID string, username string, role models.UserRole) error
	Delete(shopID string, username string) error
	FindByShopIDAndUsername(shopID string, username string) (models.ShopUser, error)
	FindRole(shopID string, username string) (models.UserRole, error)
	FindByShopID(shopID string) (*[]models.ShopUser, error)
	FindByUsername(username string) (*[]models.ShopUser, error)
	FindByUsernamePage(username string, q string, page int, limit int) ([]models.ShopUserInfo, paginate.PaginationData, error)
	FindByUserInShopPage(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ShopUser, paginate.PaginationData, error)
}

type ShopUserRepository struct {
	pst microservice.IPersisterMongo
}

func NewShopUserRepository(pst microservice.IPersisterMongo) ShopUserRepository {
	return ShopUserRepository{
		pst: pst,
	}
}

func (svc ShopUserRepository) Save(shopID string, username string, role models.UserRole) error {

	optUpdate := options.Update().SetUpsert(true)
	err := svc.pst.Update(&models.ShopUser{}, bson.M{"shopid": shopID, "username": username}, bson.M{"$set": bson.M{"role": role}}, optUpdate)

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopUserRepository) Delete(shopID string, username string) error {

	err := svc.pst.Delete(&models.ShopUser{}, bson.M{"shopid": shopID, "username": username})

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopUserRepository) FindByShopIDAndUsername(shopID string, username string) (models.ShopUser, error) {

	shopUser := &models.ShopUser{}

	err := svc.pst.FindOne(&models.ShopUser{}, bson.M{"shopid": shopID, "username": username}, shopUser)
	if err != nil {
		fmt.Println("err -> ", err.Error())
		return models.ShopUser{}, err
	}

	return *shopUser, nil
}

func (svc ShopUserRepository) FindRole(shopID string, username string) (models.UserRole, error) {

	shopUser := &models.ShopUser{}

	err := svc.pst.FindOne(&models.ShopUser{}, bson.M{"shopid": shopID, "username": username}, shopUser)

	if err != nil {
		return models.ROLE_USER, err
	}

	return shopUser.Role, nil
}

func (svc ShopUserRepository) FindByShopID(shopID string) (*[]models.ShopUser, error) {
	shopUsers := &[]models.ShopUser{}

	err := svc.pst.Find(&models.ShopUser{}, bson.M{"shopid": shopID}, shopUsers)

	if err != nil {
		return nil, err
	}

	return shopUsers, nil
}

func (svc ShopUserRepository) FindByUsername(username string) (*[]models.ShopUser, error) {
	shopUsers := &[]models.ShopUser{}

	err := svc.pst.Find(&models.ShopUser{}, bson.M{"username": username}, shopUsers)

	if err != nil {
		return nil, err
	}

	return shopUsers, nil
}

func (repo ShopUserRepository) FindByUsernamePage(username string, q string, page int, limit int) ([]models.ShopUserInfo, paginate.PaginationData, error) {

	docList := []models.ShopUserInfo{}

	searchFilterList := []interface{}{}

	searchCols := []string{
		"shopid",
		"name",
	}

	for _, colName := range searchCols {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}}})
	}

	aggPaginatedData, err := repo.pst.AggregatePage(&models.ShopUser{}, limit, page,
		bson.M{"$match": bson.M{
			"username": username,
		}},
		bson.M{"$lookup": bson.M{
			"from":         "shops",
			"localField":   "shopid",
			"foreignField": "guidfixed",
			"as":           "shopInfo",
		}},
		bson.M{
			"$match": bson.M{"shopInfo.deletedAt": bson.M{"$exists": false}},
		},
		bson.M{
			"$project": bson.M{
				"_id":    1,
				"role":   1,
				"shopid": 1,
				"name":   bson.M{"$first": "$shopInfo.name1"},
			},
		},
		bson.M{
			"$match": bson.M{
				"$or": searchFilterList,
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

func (repo ShopUserRepository) FindByUserInShopPage(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ShopUser, paginate.PaginationData, error) {

	docList := []models.ShopUser{}

	searchCols := []string{
		"username",
	}

	searchFilterList := []interface{}{}

	for _, colName := range searchCols {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}}})
	}

	filtter := bson.M{
		"shopid": shopID,
		"$or":    searchFilterList,
	}

	paginattion, err := repo.pst.FindPageSort(&models.ShopUser{}, limit, page, filtter, sort, &docList)

	if err != nil {
		return []models.ShopUser{}, paginate.PaginationData{}, err
	}

	return docList, paginattion, nil
}
