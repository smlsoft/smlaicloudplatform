package employee

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IEmployeeRepository interface {
	FindEmployeeByShopIDPage(shopID string, q string, page int, limit int) ([]models.EmployeeInfo, paginate.PaginationData, error)
	FindEmployeeByUsername(shopID string, username string) (models.EmployeeDoc, error)
	Create(models.EmployeeDoc) (primitive.ObjectID, error)
	Update(guidFixed string, employee models.EmployeeDoc) error
}

type EmployeeRepository struct {
	pst microservice.IPersisterMongo
}

func NewEmployeeRepository(pst microservice.IPersisterMongo) EmployeeRepository {
	return EmployeeRepository{
		pst: pst,
	}
}

func (r EmployeeRepository) FindEmployeeByShopIDPage(shopID string, q string, page int, limit int) ([]models.EmployeeInfo, paginate.PaginationData, error) {

	docList := []models.EmployeeInfo{}
	pagination, err := r.pst.FindPage(&models.InventoryInfo{}, limit, page, bson.M{
		"shopid": shopID,
		"$or": []interface{}{
			bson.M{"username": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
			bson.M{"name": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.EmployeeInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (r EmployeeRepository) FindEmployeeByUsername(shopID string, username string) (models.EmployeeDoc, error) {

	findDoc := &models.EmployeeDoc{}
	err := r.pst.FindOne(&models.EmployeeDoc{}, bson.M{
		"shopid":   shopID,
		"username": username,
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

func (r EmployeeRepository) Update(guidFixed string, employee models.EmployeeDoc) error {

	err := r.pst.UpdateOne(&models.EmployeeDoc{}, "guidFixed", guidFixed, employee)

	if err != nil {
		return err
	}
	return nil
}
