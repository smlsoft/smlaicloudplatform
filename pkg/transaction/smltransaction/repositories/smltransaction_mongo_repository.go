package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/transaction/smltransaction/models"

	"go.mongodb.org/mongo-driver/bson"
)

type ISMLTransactionRepository interface {
	FindByDocIndentityKey(collectionName string, shopID string, indentityField string, indentityValue interface{}) (map[string]interface{}, error)
	Create(collectionName string, doc map[string]interface{}) (string, error)
	Update(collectionName string, shopID string, guid string, doc map[string]interface{}) error
	CreateInBatch(collectionName string, docList []map[string]interface{}) error
	DeleteByGuidfixed(collectionName string, shopID string, guid string, username string) error
	Delete(collectionName string, shopID string, username string, filters map[string]interface{}) error
	Transaction(fnc func() error) error
}

type SMLTransactionRepository struct {
	pst microservice.IPersisterMongo
}

func NewSMLTransactionRepository(pst microservice.IPersisterMongo) *SMLTransactionRepository {

	insRepo := &SMLTransactionRepository{
		pst: pst,
	}

	return insRepo
}

func (repo SMLTransactionRepository) Create(collectionName string, doc map[string]interface{}) (string, error) {

	idx, err := repo.pst.Create(&models.DynamicCollection{Collection: collectionName}, doc)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo SMLTransactionRepository) Update(collectionName string, shopID string, guid string, doc map[string]interface{}) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}

	err := repo.pst.UpdateOne(&models.DynamicCollection{Collection: collectionName}, filterDoc, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo SMLTransactionRepository) CreateInBatch(collectionName string, docList []map[string]interface{}) error {
	var tempList []interface{}

	for _, inv := range docList {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(&models.DynamicCollection{Collection: collectionName}, tempList)

	if err != nil {
		return err
	}
	return nil
}

type Doc struct {
	DocNo string `bson:"docno"`
	Name  string `bson:"name"`
}

func (repo SMLTransactionRepository) FindByDocIndentityKey(collectionName string, shopID string, indentityField string, indentityValue interface{}) (map[string]interface{}, error) {

	doc := map[string]interface{}{}

	err := repo.pst.FindOne(&models.DynamicCollection{Collection: collectionName}, bson.M{"shopid": shopID, "deletedat": bson.M{"$exists": false}, indentityField: indentityValue}, &doc)

	if err != nil {
		return map[string]interface{}{}, err
	}

	return doc, nil
}

func (repo SMLTransactionRepository) DeleteByGuidfixed(collectionName string, shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&models.DynamicCollection{Collection: collectionName}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	return nil
}

func (repo SMLTransactionRepository) Delete(collectionName string, shopID string, username string, filters map[string]interface{}) error {

	filterQuery := bson.M{}

	for col, val := range filters {
		filterQuery[col] = val
	}

	filterQuery["shopid"] = shopID

	err := repo.pst.SoftDelete(&models.DynamicCollection{Collection: collectionName}, username, filterQuery)

	if err != nil {
		return err
	}

	return nil
}

func (repo SMLTransactionRepository) Transaction(fnc func() error) error {
	return repo.pst.Transaction(fnc)
}
