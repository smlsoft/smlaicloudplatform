package employee

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IEmployeeRepository interface {
	FindPageEmployeeByShopID(shopID string, pageable micromodels.Pageable) ([]models.EmployeeInfo, mongopagination.PaginationData, error)
	FindEmployeeByCode(shopID string, code string) (models.EmployeeDoc, error)
	Create(models.EmployeeDoc) (primitive.ObjectID, error)
	Update(shopID string, guidFixed string, employee models.EmployeeDoc) error

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.EmployeeDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.EmployeeActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.EmployeeDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.EmployeeActivity, error)
}

type EmployeeRepository struct {
	pst microservice.IPersisterMongo
	repositories.ActivityRepository[models.EmployeeActivity, models.EmployeeDeleteActivity]
}

func NewEmployeeRepository(pst microservice.IPersisterMongo) *EmployeeRepository {
	insRepo := &EmployeeRepository{
		pst: pst,
	}

	insRepo.ActivityRepository = repositories.NewActivityRepository[models.EmployeeActivity, models.EmployeeDeleteActivity](pst)

	return insRepo
}

func (r EmployeeRepository) FindPageEmployeeByShopID(shopID string, pageable micromodels.Pageable) ([]models.EmployeeInfo, mongopagination.PaginationData, error) {
	filterQueries := bson.M{
		"shopid": shopID,
		// "$or": []interface{}{
		// 	bson.M{"username": bson.M{"$regex": primitive.Regex{
		// 		Pattern: ".*" + q + ".*",
		// 		Options: "",
		// 	}}},
		// 	bson.M{"name": bson.M{"$regex": primitive.Regex{
		// 		Pattern: ".*" + q + ".*",
		// 		Options: "",
		// 	}}},
		// },
	}
	docList := []models.EmployeeInfo{}
	pagination, err := r.pst.FindPage(&models.EmployeeInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.EmployeeInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (r EmployeeRepository) FindEmployeeByCode(shopID string, code string) (models.EmployeeDoc, error) {

	findDoc := &models.EmployeeDoc{}
	err := r.pst.FindOne(&models.EmployeeDoc{}, bson.M{
		"shopid": shopID,
		"code":   code,
	}, findDoc)

	if err != nil {
		return models.EmployeeDoc{}, err
	}

	return *findDoc, nil
}

func (r EmployeeRepository) Create(employee models.EmployeeDoc) (primitive.ObjectID, error) {

	idx, err := r.pst.Create(&models.EmployeeDoc{}, employee)

	if err != nil {
		return primitive.NilObjectID, err
	}
	return idx, nil
}

func (r EmployeeRepository) Update(shopID string, guidFixed string, employee models.EmployeeDoc) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guidFixed,
	}

	err := r.pst.UpdateOne(&models.EmployeeDoc{}, filterDoc, employee)

	if err != nil {
		return err
	}
	return nil
}
