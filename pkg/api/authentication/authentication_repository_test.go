package authentication

import (
	"errors"
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFindUser(t *testing.T) {

	os.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	defer os.Unsetenv("MONGODB_URI")

	mongoPersisterConfig := microservice.NewMongoPersisterConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repository := NewAuthenticationRepository(mongoPersister)
	give := &models.User{
		Username: "test",
	}
	want := &models.User{
		Username: "test",
	}
	get, err := repository.CreateUser(*give)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if get == primitive.NilObjectID {
		t.Error(errors.New("Create User Failed"))
	}

	t.Logf("Create User Success With ID %v", get)

	getUser, err := repository.FindUser(want.Username)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if getUser.Username != want.Username {
		t.Error(errors.New("Create User And Find Not Match"))
		return
	}

}
