package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/smsreceive/smspatterns/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISmsPatternsRepository interface {
	Create(doc models.SmsPatternsDoc) (string, error)
	UpdateByGuid(guid string, doc models.SmsPatternsDoc) error
	DeleteByGuid(guid string) error
	FindByCode(code string) (models.SmsPatternsDoc, error)
	FindByGuid(guid string) (models.SmsPatternsDoc, error)
	FindPage(colNameSearch []string, q string, page int, limit int) ([]models.SmsPatternsInfo, mongopagination.PaginationData, error)
}

type SmsPatternsRepository struct {
	pst microservice.IPersisterMongo
}

func NewSmsPatternsRepository(pst microservice.IPersisterMongo) *SmsPatternsRepository {

	insRepo := &SmsPatternsRepository{
		pst: pst,
	}

	return insRepo
}

func (repo SmsPatternsRepository) Create(doc models.SmsPatternsDoc) (string, error) {

	idx, err := repo.pst.Create(models.SmsPatternsDoc{}, doc)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo SmsPatternsRepository) UpdateByGuid(guid string, doc models.SmsPatternsDoc) error {
	filterDoc := map[string]interface{}{
		"guidfixed": guid,
	}

	err := repo.pst.UpdateOne(models.SmsPatternsDoc{}, filterDoc, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo SmsPatternsRepository) UpdateByCode(code string, doc models.SmsPatternsDoc) error {
	filterDoc := map[string]interface{}{
		"code": code,
	}

	err := repo.pst.UpdateOne(models.SmsPatternsDoc{}, filterDoc, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo SmsPatternsRepository) DeleteByGuid(guid string) error {
	filterDoc := map[string]interface{}{
		"guidfixed": guid,
	}

	err := repo.pst.Delete(models.SmsPatternsDoc{}, filterDoc)

	if err != nil {
		return err
	}

	return nil
}

func (repo SmsPatternsRepository) FindByCode(code string) (models.SmsPatternsDoc, error) {

	doc := models.SmsPatternsDoc{}

	filters := bson.M{
		"code": code,
	}

	err := repo.pst.FindOne(models.SmsPatternsDoc{}, filters, &doc)

	if err != nil {
		return models.SmsPatternsDoc{}, err
	}

	return doc, nil
}

func (repo SmsPatternsRepository) FindByGuid(guidFixed string) (models.SmsPatternsDoc, error) {

	doc := models.SmsPatternsDoc{}

	filters := bson.M{
		"guidfixed": guidFixed,
	}

	err := repo.pst.FindOne(models.SmsPatternsDoc{}, filters, &doc)

	if err != nil {
		return models.SmsPatternsDoc{}, err
	}

	return doc, nil
}

func (repo SmsPatternsRepository) FindPage(colNameSearch []string, q string, page int, limit int) ([]models.SmsPatternsInfo, mongopagination.PaginationData, error) {

	searchFilterList := []interface{}{}

	for _, colName := range colNameSearch {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}}})
	}

	docList := []models.SmsPatternsInfo{}
	pagination, err := repo.pst.FindPage(models.SmsPatternsInfo{}, limit, page, bson.M{
		"$or": searchFilterList,
	}, &docList)

	if err != nil {
		return []models.SmsPatternsInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
