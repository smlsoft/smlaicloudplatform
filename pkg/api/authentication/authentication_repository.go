package authentication

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationRepository interface {
	FindUser(id string) (*models.User, error)
	CreateUser(models.User) (primitive.ObjectID, error)
	UpdateUser(username string, user models.User) error
}

type AuthenticationRepository struct {
	pst microservice.IPersisterMongo
}

func NewAuthenticationRepository(pst microservice.IPersisterMongo) AuthenticationRepository {
	return AuthenticationRepository{
		pst: pst,
	}
}

func (r AuthenticationRepository) FindUser(username string) (*models.User, error) {

	findUser := &models.User{}
	err := r.pst.FindOne(&models.User{}, bson.M{"username": username}, findUser)

	if err != nil {
		return nil, err
	}

	return findUser, nil
}

func (r AuthenticationRepository) CreateUser(user models.User) (primitive.ObjectID, error) {

	idx, err := r.pst.Create(&models.User{}, user)

	if err != nil {
		return primitive.NilObjectID, err
	}
	return idx, nil
}

func (r AuthenticationRepository) UpdateUser(username string, user models.User) error {

	err := r.pst.UpdateOne(&models.User{}, "username", username, user)

	if err != nil {
		return err
	}
	return nil
}
