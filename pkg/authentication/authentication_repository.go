package authentication

import (
	"context"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shop/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationRepository interface {
	FindUser(ctx context.Context, id string) (*models.UserDoc, error)
	CreateUser(ctx context.Context, doc models.UserDoc) (primitive.ObjectID, error)
	UpdateUser(ctx context.Context, username string, user models.UserDoc) error
}

type AuthenticationRepository struct {
	pst microservice.IPersisterMongo
}

func NewAuthenticationRepository(pst microservice.IPersisterMongo) AuthenticationRepository {
	return AuthenticationRepository{
		pst: pst,
	}
}

func (r AuthenticationRepository) FindUser(ctx context.Context, username string) (*models.UserDoc, error) {

	findUser := &models.UserDoc{}
	err := r.pst.FindOne(ctx, &models.UserDoc{}, bson.M{"username": username}, findUser)

	if err != nil {
		return nil, err
	}

	return findUser, nil
}

func (r AuthenticationRepository) CreateUser(ctx context.Context, user models.UserDoc) (primitive.ObjectID, error) {

	idx, err := r.pst.Create(ctx, &models.UserDoc{}, user)

	if err != nil {
		return primitive.NilObjectID, err
	}
	return idx, nil
}

func (r AuthenticationRepository) UpdateUser(ctx context.Context, username string, user models.UserDoc) error {

	filterDoc := map[string]interface{}{
		"username": username,
	}

	err := r.pst.UpdateOne(ctx, &models.UserDoc{}, filterDoc, user)

	if err != nil {
		return err
	}
	return nil
}
