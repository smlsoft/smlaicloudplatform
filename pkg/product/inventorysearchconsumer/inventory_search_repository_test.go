package inventorysearchconsumer_test

import (
	"os"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	inventoryModel "smlcloudplatform/pkg/product/inventory/models"
	"smlcloudplatform/pkg/product/inventorysearchconsumer"
	"smlcloudplatform/pkg/product/inventorysearchconsumer/models"
	"smlcloudplatform/pkg/utils"

	"testing"

	"github.com/tj/assert"
)

type RealEngineOpenSearchConfig struct{}

func (r RealEngineOpenSearchConfig) Address() []string {
	return []string{
		"http://103.212.36.91:19200",
	}
}

func (r RealEngineOpenSearchConfig) Username() string {
	return "admin"
}

func (r RealEngineOpenSearchConfig) Password() string {
	return "admin"
}

func InitInventorySearchRepositoryRealEngine() *inventorysearchconsumer.InventorySearchRepository {

	config := &RealEngineOpenSearchConfig{}
	openSearchPersister := microservice.NewPersisterOpenSearch(config)
	repo := inventorysearchconsumer.NewInventorySearchRepository(openSearchPersister)
	return repo
}

func TestInventorySearchRepositoryUpsertRealEngineOpensearch(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	itemName := "itemTest"

	give := &models.InventorySearch{
		InventoryInfo: inventoryModel.InventoryInfo{
			DocIdentity: common.DocIdentity{
				GuidFixed: utils.NewGUID(),
			},

			Inventory: inventoryModel.Inventory{
				Name: common.Name{
					Name1: itemName,
				},
			},
		},
		ShopIdentity: common.ShopIdentity{
			ShopID: "TOETEST",
		},
	}

	repo := InitInventorySearchRepositoryRealEngine()
	get := repo.UpSert(give)
	assert.Nil(t, get)

}

func TestInventorySearchRepositoryDeleteRealEngineOpensearch(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	repo := InitInventorySearchRepositoryRealEngine()
	get := repo.Delete("2EnwbzO00CcMe92gEtRvVHcQ2S1")
	assert.Nil(t, get)

}
