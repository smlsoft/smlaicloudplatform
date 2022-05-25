package authentication_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/authentication"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repoMock authentication.AuthenticationRepository

func init() {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repoMock = authentication.NewAuthenticationRepository(mongoPersister)
}

func TestFindUser(t *testing.T) {
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

	createUserID, err := repoMock.CreateUser(*give)
	if err != nil {
		t.Error(err.Error())
		return
	}

	require.NotEqual(t, createUserID, primitive.NilObjectID, "Create User Failed")

	getUser, err := repoMock.FindUser(want.Username)

	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, getUser.Username, want.Username, "Create User And Find Not Match")

}
