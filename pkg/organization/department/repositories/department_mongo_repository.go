package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/organization/department/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IDepartmentRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.DepartmentDoc) (string, error)
	CreateInBatch(docList []models.DepartmentDoc) error
	Update(shopID string, guid string, doc models.DepartmentDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.DepartmentInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.DepartmentDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.DepartmentItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.DepartmentDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.DepartmentInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.DepartmentInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DepartmentDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DepartmentActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DepartmentDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DepartmentActivity, error)

	FindOneByCode(shopID, branchCode, departmentCode string) (models.DepartmentDoc, error)
}

type DepartmentRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DepartmentDoc]
	repositories.SearchRepository[models.DepartmentInfo]
	repositories.GuidRepository[models.DepartmentItemGuid]
	repositories.ActivityRepository[models.DepartmentActivity, models.DepartmentDeleteActivity]
}

func NewDepartmentRepository(pst microservice.IPersisterMongo) *DepartmentRepository {

	insRepo := &DepartmentRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DepartmentDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DepartmentInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.DepartmentItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.DepartmentActivity, models.DepartmentDeleteActivity](pst)

	return insRepo
}

func (repo DepartmentRepository) FindOneByCode(shopID, branchCode, departmentCode string) (models.DepartmentDoc, error) {
	doc := models.DepartmentDoc{}
	err := repo.pst.FindOne(models.DepartmentDoc{}, bson.M{
		"shop_id":        shopID,
		"deletedat":      bson.M{"$exists": false},
		"branchcode":     branchCode,
		"departmentcode": departmentCode,
	}, &doc)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
