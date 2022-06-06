package services_test

import (
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJournalRepsitory struct {
	mock.Mock
}

func (m *MockJournalRepsitory) CreateInBatch(docList []models.JournalPg) error {
	ret := m.Called(docList)
	return ret.Error(0)
}

func (m *MockJournalRepsitory) Create(doc models.JournalPg) error {
	ret := m.Called(doc)
	return ret.Error(0)
}

func (m *MockJournalRepsitory) Update(shopID string, docNo string, doc models.JournalPg) error {
	ret := m.Called(shopID, docNo, doc)
	return ret.Error(0)
}

func (m *MockJournalRepsitory) Delete(shopID string, docNo string) error {
	ret := m.Called(shopID, docNo)
	return ret.Error(0)
}

func (m *MockJournalRepsitory) Get(shopID string, docNo string) (*models.JournalPg, error) {
	ret := m.Called(shopID, docNo)
	return ret.Get(0).(*models.JournalPg), ret.Error(1)
}

func TestJournalConsumeServiceCreated(t *testing.T) {

	get := models.JournalPg{
		JournalBody: models.JournalBody{
			DocNo: "0001",
		},
	}

	give := models.JournalDoc{
		JournalData: models.JournalData{
			JournalInfo: models.JournalInfo{
				Journal: models.Journal{
					JournalBody: models.JournalBody{
						DocNo: "0001",
					},
				},
			},
		},
	}

	mockRepo := new(MockJournalRepsitory)
	mockRepo.On("Create", get).Return(nil)

	journalService := services.NewJournalConsumeService(mockRepo)
	err := journalService.Create(give)
	assert.Nil(t, err, "Error should be nil")
}

func TestJournalConsumeServiceUpdate(t *testing.T) {

	get := models.JournalPg{
		JournalBody: models.JournalBody{
			DocNo: "0001",
		},
	}

	give := models.JournalDoc{
		JournalData: models.JournalData{
			JournalInfo: models.JournalInfo{
				Journal: models.Journal{
					JournalBody: models.JournalBody{
						DocNo: "0001",
					},
				},
			},
		},
	}

	mockRepo := new(MockJournalRepsitory)
	mockRepo.On("Update", "SHOPID", "0001", get).Return(nil)

	journalService := services.NewJournalConsumeService(mockRepo)
	err := journalService.Update("SHOPID", "0001", give)
	assert.Nil(t, err, "Error should be nil")
}
