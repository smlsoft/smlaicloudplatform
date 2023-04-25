package services

import (
	"testing"
	"time"

	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/productsection/sectionbranch/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/userplant/mongopagination"
)

// Create a mock for the SectionBranchRepository
type MockSectionBranchRepository struct {
	mock.Mock
}

func (m *MockSectionBranchRepository) Count(shopID string) (int, error) {
	args := m.Called(shopID)
	return args.Int(0), args.Error(1)
}

func (m *MockSectionBranchRepository) Create(doc models.SectionBranchDoc) (string, error) {
	args := m.Called(doc)
	return args.String(0), args.Error(1)
}

func (m *MockSectionBranchRepository) CreateInBatch(docList []models.SectionBranchDoc) error {
	args := m.Called(docList)
	return args.Error(0)
}

func (m *MockSectionBranchRepository) Update(shopID string, guid string, doc models.SectionBranchDoc) error {
	args := m.Called(shopID, guid, doc)
	return args.Error(0)
}

func (m *MockSectionBranchRepository) DeleteByGuidfixed(shopID string, guid string, username string) error {
	args := m.Called(shopID, guid, username)
	return args.Error(0)
}

func (m *MockSectionBranchRepository) Delete(shopID string, username string, filters map[string]interface{}) error {
	args := m.Called(shopID, username, filters)
	return args.Error(0)
}

func (m *MockSectionBranchRepository) FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SectionBranchInfo, mongopagination.PaginationData, error) {
	args := m.Called(shopID, searchInFields, pageable)
	return args.Get(0).([]models.SectionBranchInfo), args.Get(1).(mongopagination.PaginationData), args.Error(2)
}

func (m *MockSectionBranchRepository) FindByGuid(shopID string, guid string) (models.SectionBranchDoc, error) {
	args := m.Called(shopID, guid)
	return args.Get(0).(models.SectionBranchDoc), args.Error(1)
}

func (m *MockSectionBranchRepository) FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.SectionBranchItemGuid, error) {
	args := m.Called(shopID, columnName, itemGuidList)
	return args.Get(0).([]models.SectionBranchItemGuid), args.Error(1)
}

func (m *MockSectionBranchRepository) FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.SectionBranchDoc, error) {
	args := m.Called(shopID, indentityField, indentityValue)
	return args.Get(0).(models.SectionBranchDoc), args.Error(1)
}

func (m *MockSectionBranchRepository) FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.SectionBranchInfo, mongopagination.PaginationData, error) {
	args := m.Called(shopID, filters, searchInFields, pageable)
	return args.Get(0).([]models.SectionBranchInfo), args.Get(1).(mongopagination.PaginationData), args.Error(2)
}

func (m *MockSectionBranchRepository) FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.SectionBranchInfo, int, error) {
	args := m.Called(shopID, filters, searchInFields, projects, pageableLimit)
	return args.Get(0).([]models.SectionBranchInfo), args.Int(1), args.Error(2)
}

func (m *MockSectionBranchRepository) FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBranchDeleteActivity, mongopagination.PaginationData, error) {
	args := m.Called(shopID, lastUpdatedDate, filters, pageable)
	return args.Get(0).([]models.SectionBranchDeleteActivity), args.Get(1).(mongopagination.PaginationData), args.Error(2)
}

func (m *MockSectionBranchRepository) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBranchActivity, mongopagination.PaginationData, error) {
	args := m.Called(shopID, lastUpdatedDate, filters, pageable)
	return args.Get(0).([]models.SectionBranchActivity), args.Get(1).(mongopagination.PaginationData), args.Error(2)
}

func (m *MockSectionBranchRepository) FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SectionBranchDeleteActivity, error) {
	args := m.Called(shopID, lastUpdatedDate, filters, pageableStep)
	return args.Get(0).([]models.SectionBranchDeleteActivity), args.Error(1)
}

func (m *MockSectionBranchRepository) FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SectionBranchActivity, error) {
	args := m.Called(shopID, lastUpdatedDate, filters, pageableStep)
	return args.Get(0).([]models.SectionBranchActivity), args.Error(1)
}

func mockNewGUID() string {
	return "testGuidFixed"
}

func TestSaveSectionBranch(t *testing.T) {
	mockRepo := new(MockSectionBranchRepository)
	// ... setup other mocks ...

	// Create an instance of SectionBranchHttpService
	svc := NewSectionBranchHttpService(mockRepo, mockNewGUID, nil)

	t.Run("Test SaveSectionBranch - Create", func(t *testing.T) {
		// Setup
		shopID := "testShopID"
		authUsername := "testUser"
		branchCode := "testBranchCode"

		doc := models.SectionBranch{
			BranchCode: branchCode,
		}

		emptyDoc := models.SectionBranchDoc{}
		mockRepo.On("FindByDocIndentityGuid", shopID, "branchcode", branchCode).Return(emptyDoc, nil)
		mockRepo.On("Create", mock.Anything).Return("testGuidFixed", nil)

		// Execute
		guidFixed, err := svc.SaveSectionBranch(shopID, authUsername, doc)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, guidFixed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Test SaveSectionBranch - Update", func(t *testing.T) {
		// Setup
		shopID := "testShopID"
		authUsername := "testUser"
		branchCode := "testBranchCode"

		doc := models.SectionBranch{
			BranchCode: branchCode,
		}

		existingDoc := models.SectionBranchDoc{}

		existingDoc.ShopID = shopID
		existingDoc.GuidFixed = "testGuidFixed"
		existingDoc.SectionBranch = doc

		mockRepo.On("FindByDocIndentityGuid", shopID, "branchcode", branchCode).Return(existingDoc, nil)
		mockRepo.On("Update", shopID, existingDoc.GuidFixed, mock.Anything).Return(nil)

		// Execute
		guidFixed, err := svc.SaveSectionBranch(shopID, authUsername, doc)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, existingDoc.GuidFixed, guidFixed)
		// mockRepo.AssertExpectations(t)
	})
}
