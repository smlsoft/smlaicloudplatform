package services_test

import (
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/product/category/models"
	"smlcloudplatform/pkg/product/category/repositories"
	"smlcloudplatform/pkg/product/category/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

const shopIDMock = "TEST_SHOP"
const authUserMock = "TEST_USER"

var mongoPersister *microservice.PersisterMongo
var serviceMock services.CategoryService
var repoMock repositories.CategoryRepository

func init() {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister = microservice.NewPersisterMongo(mongoPersisterConfig)
	repoMock = repositories.NewCategoryRepository(mongoPersister)
	serviceMock = services.NewCategoryService(repoMock, nil)
}

func TestCreateCategory(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	cateReq1 := models.Category{}
	cateReq1.CategoryGuid = "cate001"
	cateReq1.Name1 = "cate name 1"

	cateReq2 := models.Category{}
	cateReq2.CategoryGuid = "cate002"
	cateReq2.Name1 = "cate name 2"

	cases := []struct {
		name     string
		wantErr  bool
		req      models.Category
		wantData models.Category
	}{
		{
			name:    "Fail on create duplicate",
			wantErr: true,
			req:     cateReq1,
		},
		{
			name:     "Success",
			req:      cateReq2,
			wantData: cateReq2,
		},
	}
	//clear mock data
	mongoPersister.Delete(models.CategoryDoc{}, bson.M{"shopid": shopIDMock})

	//mock data
	mockCategory := []models.CategoryDoc{}

	cateDoc1 := models.CategoryDoc{}
	cateDoc1.ShopID = shopIDMock
	cateDoc1.GuidFixed = "TESTGUID001"
	cateDoc1.CategoryGuid = "cate001"
	cateDoc1.Name1 = "cate name 1"

	mockCategory = append(mockCategory, cateDoc1)

	for _, categoryDoc := range mockCategory {
		_, err := repoMock.Create(categoryDoc)

		if err != nil {
			t.Error(err)
		}
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			respGuid, err := serviceMock.CreateCategory(shopIDMock, authUserMock, tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr && respGuid == "" {
				t.Error("guid is empty")
				return
			}
		})
	}
}

func TestInfoCategory(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}

	cateWant1 := models.CategoryInfo{}

	cateWant1.CategoryGuid = "cate001"
	cateWant1.Name1 = "cate name 1"

	cases := []struct {
		name     string
		wantErr  bool
		id       string
		wantData models.CategoryInfo
	}{
		{
			name:    "category does not exist",
			wantErr: true,
			id:      "TESTGUID000",
		},
		{
			name:     "category does not exist",
			wantErr:  false,
			id:       "TESTGUID001",
			wantData: cateWant1,
		},
	}

	//mock data
	mockCategory := []models.CategoryDoc{}

	cateDoc1 := models.CategoryDoc{}
	cateDoc1.ShopID = shopIDMock
	cateDoc1.GuidFixed = "TESTGUID001"
	cateDoc1.CategoryGuid = "cate001"
	cateDoc1.Name1 = "cate name 1"

	mockCategory = append(mockCategory, cateDoc1)

	for _, categoryDoc := range mockCategory {
		_, err := repoMock.Create(categoryDoc)

		if err != nil {
			t.Error(err)
		}
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := serviceMock.InfoCategory(shopIDMock, tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData.CategoryGuid != "" {
				if resp.CategoryGuid == "" {
					t.Error("category is not found")
				} else {
					assert.Equal(t, tt.wantData.Name1, resp.Name1)
				}

			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	//clear mock data
	mongoPersister.Delete(models.CategoryDoc{}, bson.M{"shopid": shopIDMock})

	mockCategory := []models.CategoryDoc{}

	cateDoc1 := models.CategoryDoc{}
	cateDoc1.ShopID = shopIDMock
	cateDoc1.GuidFixed = "TESTGUID001"
	cateDoc1.CategoryGuid = "cate001"
	cateDoc1.Name1 = "cate name 1"

	mockCategory = append(mockCategory, cateDoc1)

	//mock data
	for _, categoryDoc := range mockCategory {
		_, err := repoMock.Create(categoryDoc)

		if err != nil {
			t.Error(err)
		}
	}

	err := serviceMock.UpdateCategory(shopIDMock, cateDoc1.GuidFixed, authUserMock, cateDoc1.Category)

	require.NoError(t, err, "update category error")

}

func TestDeleteCategory(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	//clear mock data
	err := mongoPersister.Delete(models.CategoryDoc{}, bson.M{"shopid": shopIDMock})

	cateReq1 := models.Category{}
	cateReq1.CategoryGuid = "cate001"
	cateReq1.Name1 = "cate name 1"

	guidfixed, err := serviceMock.CreateCategory(shopIDMock, authUserMock, cateReq1)

	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.NotEqual(t, guidfixed, "", "guidfixed not found.")

	err = serviceMock.DeleteCategory(shopIDMock, guidfixed, authUserMock)

	require.NoError(t, err, "delete failed.")
}
