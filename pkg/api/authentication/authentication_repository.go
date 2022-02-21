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
}

type AuthenticationRepository struct {
	pst microservice.IPersisterMongo
}

func NewAuthenticationRepository(pst microservice.IPersisterMongo) *AuthenticationRepository {

	authenticationRepository := &AuthenticationRepository{
		pst: pst,
	}
	return authenticationRepository
}

func (r *AuthenticationRepository) FindUser(username string) (*models.User, error) {

	findUser := &models.User{}
	err := r.pst.FindOne(&models.User{}, bson.M{"username": username}, findUser)

	if err != nil {
		return nil, err
	}

	return findUser, nil
}

func (r *AuthenticationRepository) CreateUser(user models.User) (primitive.ObjectID, error) {

	idx, err := r.pst.Create(&models.User{}, user)

	if err != nil {
		return primitive.NilObjectID, err
	}
	return idx, nil
}
