package authentication

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationRepository interface {
	FindUser(id string) (*models.UserDoc, error)
	CreateUser(models.UserDoc) (primitive.ObjectID, error)
	UpdateUser(username string, user models.UserDoc) error
}

type AuthenticationRepository struct {
	pst microservice.IPersisterMongo
}

func NewAuthenticationRepository(pst microservice.IPersisterMongo) AuthenticationRepository {
	return AuthenticationRepository{
		pst: pst,
	}
}

func (r AuthenticationRepository) FindUser(username string) (*models.UserDoc, error) {

	findUser := &models.UserDoc{}
	err := r.pst.FindOne(&models.UserDoc{}, bson.M{"username": username}, findUser)

	if err != nil {
		return nil, err
	}

	return findUser, nil
}

func (r AuthenticationRepository) CreateUser(user models.UserDoc) (primitive.ObjectID, error) {

	idx, err := r.pst.Create(&models.UserDoc{}, user)

	if err != nil {
		return primitive.NilObjectID, err
	}
	return idx, nil
}

func (r AuthenticationRepository) UpdateUser(username string, user models.UserDoc) error {

	filterDoc := map[string]interface{}{
		"username": username,
	}

	err := r.pst.UpdateOne(&models.UserDoc{}, filterDoc, user)

	if err != nil {
		return err
	}
	return nil
}
