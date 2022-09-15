package inventorysearchconsumer_test

import (
	common "smlcloudplatform/pkg/models"
	inventoryModel "smlcloudplatform/pkg/product/inventory/models"
	"smlcloudplatform/pkg/product/inventorysearchconsumer"
	"smlcloudplatform/pkg/product/inventorysearchconsumer/models"
	"smlcloudplatform/pkg/utils"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

type MockInventorySearchConsumerRepository struct {
	mock.Mock
}

func (m *MockInventorySearchConsumerRepository) UpSert(inventory *models.InventorySearch) error {
	ret := m.Called(inventory)
	return ret.Error(0)
}
func (m *MockInventorySearchConsumerRepository) Delete(guidfixed string) error {
	ret := m.Called(guidfixed)
	return ret.Error(0)
}

func TestConsumerServiceOnConsumeInsert(t *testing.T) {

	itemName := "itemName"
	giveString := ""

	wantCreateOpenSearchData := &models.InventorySearch{
		InventoryInfo: inventoryModel.InventoryInfo{
			DocIdentity: common.DocIdentity{
				GuidFixed: utils.NewGUID(),
			},

			Inventory: inventoryModel.Inventory{
				Name: common.Name{
					Name1: &itemName,
				},
			},
		},
		ShopIdentity: common.ShopIdentity{
			ShopID: "TOETEST",
		},
	}

	mockRepo := new(MockInventorySearchConsumerRepository)
	mockRepo.On("Create", wantCreateOpenSearchData).Return(nil)

	svc := inventorysearchconsumer.NewInventorySearchConsumerService(mockRepo)

	wantErrNil := svc.Create(giveString)
	assert.Nil(t, wantErrNil, "Error should be nil")
}

func TestConsumerServiceOnConsumeUpdate(t *testing.T) {
	giveString := ""

	mockRepo := new(MockInventorySearchConsumerRepository)
	svc := inventorysearchconsumer.NewInventorySearchConsumerService(mockRepo)

	wantErrNil := svc.Update(giveString)
	assert.Nil(t, wantErrNil, "Error should be nil")
}

func TestConsumerServiceOnConsumeDelete(t *testing.T) {
	giveString := ""

	mockRepo := new(MockInventorySearchConsumerRepository)
	svc := inventorysearchconsumer.NewInventorySearchConsumerService(mockRepo)

	wantErrNil := svc.Create(giveString)
	assert.Nil(t, wantErrNil, "Error should be nil")
}
