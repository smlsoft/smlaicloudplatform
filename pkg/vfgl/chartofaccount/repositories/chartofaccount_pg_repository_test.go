package repositories_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/vfgl/chartofaccount/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

var repo repositories.ChartOfAccountPgRepository

func init() {
	persisterConfig := mock.NewPersisterPostgresqlConfig()
	pst := microservice.NewPersister(persisterConfig)
	repo = repositories.NewChartOfAccountPgRepository(pst)
}

func TestChartOfAccountRepositoryCreate(t *testing.T) {

	assert := assert.New(t)
	assert.NotNil(repo)

	give := &vfgl.ChartOfAccountPG{
		AccountCode: "10000",
		AccountName: "เงินสด",
	}

	err := repo.Create(*give)
	assert.Nil(err)

	get, err := repo.Get("10000")
	assert.NotNil(get)

}
