package authentication

import (
	"encoding/json"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shop/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationMongoCacheRepository interface {
	FindUser(id string) (*models.UserDoc, error)
	CreateUser(models.UserDoc) (primitive.ObjectID, error)
	UpdateUser(username string, user models.UserDoc) error
}

type AuthenticationMongoCacheRepository struct {
	pst   microservice.IPersisterMongo
	cache microservice.ICacher
}

func NewAuthenticationMongoCacheRepository(pst microservice.IPersisterMongo, cache microservice.ICacher) AuthenticationMongoCacheRepository {
	return AuthenticationMongoCacheRepository{
		pst:   pst,
		cache: cache,
	}
}

func (r AuthenticationMongoCacheRepository) FindUser(username string) (*models.UserDoc, error) {

	cacheKey := fmt.Sprintf("user:%s", username)
	userCache, err := r.cache.Get(cacheKey)

	if err != nil {
		fmt.Println(err.Error())
	}

	if len(userCache) > 0 {
		fmt.Println("user cache")
		user := &models.UserDoc{}
		err := json.Unmarshal([]byte(userCache), user)
		if err == nil {
			return user, nil
		}
	}

	findUser := &models.UserDoc{}
	err = r.pst.FindOne(&models.UserDoc{}, bson.M{"username": username}, findUser)

	if err != nil {
		return nil, err
	}

	tempUser, err := json.Marshal(findUser)

	if err != nil {
		fmt.Println(err.Error())
	}

	if err == nil {
		err = r.cache.SetS(cacheKey, string(tempUser), time.Second*60)

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return findUser, nil
}

func (r AuthenticationMongoCacheRepository) CreateUser(user models.UserDoc) (primitive.ObjectID, error) {

	idx, err := r.pst.Create(&models.UserDoc{}, user)

	if err != nil {
		return primitive.NilObjectID, err
	}
	return idx, nil
}

func (r AuthenticationMongoCacheRepository) UpdateUser(username string, user models.UserDoc) error {

	filterDoc := map[string]interface{}{
		"username": username,
	}

	err := r.pst.UpdateOne(&models.UserDoc{}, filterDoc, user)

	if err != nil {
		return err
	}
	return nil
}
