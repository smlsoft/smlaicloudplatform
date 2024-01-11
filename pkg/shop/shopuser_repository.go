package shop

import (
	"context"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/authentication/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IShopUserRepository interface {
	Create(ctx context.Context, shopUser *models.ShopUser) error
	Update(ctx context.Context, id primitive.ObjectID, shopID string, username string, role models.UserRole) error
	Save(ctx context.Context, shopID string, username string, role models.UserRole) error
	UpdateLastAccess(ctx context.Context, shopID string, username string, lastAccessedAt time.Time) error
	SaveFavorite(ctx context.Context, shopID string, username string, isFavorite bool) error
	Delete(ctx context.Context, shopID string, username string) error
	FindByShopIDAndUsernameInfo(ctx context.Context, shopID string, username string) (models.ShopUserInfo, error)
	FindByShopIDAndUsername(ctx context.Context, shopID string, username string) (models.ShopUser, error)
	FindRole(ctx context.Context, shopID string, username string) (models.UserRole, error)
	FindByShopID(ctx context.Context, shopID string) (*[]models.ShopUser, error)
	FindByUsername(ctx context.Context, username string) (*[]models.ShopUser, error)
	FindByUsernamePage(ctx context.Context, username string, pageable micromodels.Pageable) ([]models.ShopUserInfo, mongopagination.PaginationData, error)
	FindByUserInShopPage(ctx context.Context, shopID string, pageable micromodels.Pageable) ([]models.ShopUser, mongopagination.PaginationData, error)
	FindUserProfileByUsernames(ctx context.Context, usernames []string) ([]models.UserProfile, error)
}

type ShopUserRepository struct {
	pst microservice.IPersisterMongo
}

func NewShopUserRepository(pst microservice.IPersisterMongo) ShopUserRepository {
	return ShopUserRepository{
		pst: pst,
	}
}

func (svc ShopUserRepository) Create(ctx context.Context, shopUser *models.ShopUser) error {

	_, err := svc.pst.Create(ctx, &models.ShopUser{}, shopUser)

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopUserRepository) Update(ctx context.Context, id primitive.ObjectID, shopID string, username string, role models.UserRole) error {

	err := svc.pst.Update(ctx, &models.ShopUser{}, bson.M{"_id": id, "shopid": shopID}, bson.M{"$set": bson.M{"username": username, "role": role}})

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopUserRepository) Save(ctx context.Context, shopID string, username string, role models.UserRole) error {

	optUpdate := options.Update().SetUpsert(true)
	err := svc.pst.Update(ctx, &models.ShopUser{}, bson.M{"shopid": shopID, "username": username}, bson.M{"$set": bson.M{"role": role}}, optUpdate)

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopUserRepository) UpdateLastAccess(ctx context.Context, shopID string, username string, lastAccessedAt time.Time) error {

	optUpdate := options.Update().SetUpsert(true)
	err := svc.pst.Update(ctx, &models.ShopUser{}, bson.M{"shopid": shopID, "username": username}, bson.M{"$set": bson.M{"lastaccessedat": lastAccessedAt}}, optUpdate)

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopUserRepository) SaveFavorite(ctx context.Context, shopID string, username string, isFavorite bool) error {

	optUpdate := options.Update().SetUpsert(true)
	err := svc.pst.Update(ctx, &models.ShopUser{}, bson.M{"shopid": shopID, "username": username}, bson.M{"$set": bson.M{"isfavorite": isFavorite}}, optUpdate)

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopUserRepository) Delete(ctx context.Context, shopID string, username string) error {

	err := svc.pst.Delete(ctx, &models.ShopUser{}, bson.M{"shopid": shopID, "username": username})

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopUserRepository) FindByShopIDAndUsernameInfo(ctx context.Context, shopID string, username string) (models.ShopUserInfo, error) {

	shopUser := &models.ShopUserInfo{}

	err := svc.pst.FindOne(ctx, &models.ShopUserInfo{}, bson.M{"shopid": shopID, "username": username}, shopUser)
	if err != nil {
		return models.ShopUserInfo{}, err
	}

	return *shopUser, nil
}

func (svc ShopUserRepository) FindByShopIDAndUsername(ctx context.Context, shopID string, username string) (models.ShopUser, error) {

	shopUser := &models.ShopUser{}

	err := svc.pst.FindOne(ctx, &models.ShopUser{}, bson.M{"shopid": shopID, "username": username}, shopUser)
	if err != nil {
		return models.ShopUser{}, err
	}

	return *shopUser, nil
}

func (svc ShopUserRepository) FindRole(ctx context.Context, shopID string, username string) (models.UserRole, error) {

	shopUser := &models.ShopUser{}

	err := svc.pst.FindOne(ctx, &models.ShopUser{}, bson.M{"shopid": shopID, "username": username}, shopUser)

	if err != nil {
		return models.ROLE_USER, err
	}

	return shopUser.Role, nil
}

func (svc ShopUserRepository) FindByShopID(ctx context.Context, shopID string) (*[]models.ShopUser, error) {
	shopUsers := &[]models.ShopUser{}

	err := svc.pst.Find(ctx, &models.ShopUser{}, bson.M{"shopid": shopID}, shopUsers)

	if err != nil {
		return nil, err
	}

	return shopUsers, nil
}

func (svc ShopUserRepository) FindByUsername(ctx context.Context, username string) (*[]models.ShopUser, error) {
	shopUsers := &[]models.ShopUser{}

	err := svc.pst.Find(ctx, &models.ShopUser{}, bson.M{"username": username}, shopUsers)

	if err != nil {
		return nil, err
	}

	return shopUsers, nil
}

func (repo ShopUserRepository) FindByUsernamePage(ctx context.Context, username string, pageable micromodels.Pageable) ([]models.ShopUserInfo, mongopagination.PaginationData, error) {

	docList := []models.ShopUserInfo{}

	searchFilterList := []interface{}{}

	searchInFields := []string{
		"shopid",
		"names",
	}

	for _, colName := range searchInFields {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + pageable.Query + ".*",
			Options: "",
		}}})
	}

	aggPaginatedData, err := repo.pst.AggregatePage(ctx, &models.ShopUser{}, pageable,
		bson.M{"$match": bson.M{
			"username": username,
			"deletedat": bson.M{
				"$exists": false,
			},
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
				"_id":            1,
				"role":           1,
				"shopid":         1,
				"isfavorite":     1,
				"lastaccessedat": 1,
				"names":          bson.M{"$first": "$shopInfo.names"},
				"branchcode":     bson.M{"$first": "$shopInfo.branchcode"},
				"createdby":      bson.M{"$first": "$shopInfo.createdby"},
			},
		},
		bson.M{
			"$match": bson.M{
				"$or": searchFilterList,
			},
		},
		bson.M{
			"$sort": bson.M{
				"lastaccessedat": -1,
			},
		},
	)

	if err != nil {
		return []models.ShopUserInfo{}, mongopagination.PaginationData{}, err
	}

	for _, raw := range aggPaginatedData.Data {
		var doc *models.ShopUserInfo

		if marshallErr := bson.Unmarshal(raw, &doc); marshallErr == nil {
			docList = append(docList, *doc)
		}

	}

	return docList, aggPaginatedData.Pagination, nil
}

func (repo ShopUserRepository) FindByUserInShopPage(ctx context.Context, shopID string, pageable micromodels.Pageable) ([]models.ShopUser, mongopagination.PaginationData, error) {

	docList := []models.ShopUser{}

	searchInFields := []string{
		"username",
	}

	searchFilterList := []interface{}{}

	for _, colName := range searchInFields {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + pageable.Query + ".*",
			Options: "",
		}}})
	}

	filtter := bson.M{
		"shopid": shopID,
		"$or":    searchFilterList,
	}

	paginattion, err := repo.pst.FindPage(ctx, &models.ShopUser{}, filtter, pageable, &docList)

	if err != nil {
		return []models.ShopUser{}, mongopagination.PaginationData{}, err
	}

	return docList, paginattion, nil
}

func (repo ShopUserRepository) FindUserProfileByUsernames(ctx context.Context, usernames []string) ([]models.UserProfile, error) {

	docList := []models.UserProfile{}

	filters := bson.M{
		"username": bson.M{
			"$in": usernames,
		},
	}

	err := repo.pst.Find(ctx, &models.UserProfile{}, filters, &docList)

	if err != nil {
		return []models.UserProfile{}, err
	}

	return docList, nil
}
