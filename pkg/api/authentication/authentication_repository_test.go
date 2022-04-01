package authentication

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getRepo() AuthenticationRepository {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repository := NewAuthenticationRepository(mongoPersister)
	return repository
}

func TestFindUser(t *testing.T) {

	// os.Setenv("MONGODB_URI", "mongodb://root:rootx@localhost:27017/")
	// defer os.Unsetenv("MONGODB_URI")

	repository := getRepo()

	password, _ := utils.HashPassword("test")

	createAt := time.Now()
	give := &models.User{
		Username:  "test",
		Name:      "test",
		Password:  password,
		CreatedAt: createAt,
	}
	want := &models.User{
		Username:  "test",
		Name:      "test",
		Password:  password,
		CreatedAt: createAt,
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
