package purchase

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPurchaseRepository interface {
	Create(doc models.PurchaseDoc) (primitive.ObjectID, error)
	Update(guid string, doc models.PurchaseDoc) error
	Delete(guid string, shopID string) error
	FindByGuid(guid string, shopID string) (models.PurchaseDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error)
	FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error)
}

type PurchaseRepository struct {
	pst microservice.IPersisterMongo
}

func NewPurchaseRepository(pst microservice.IPersisterMongo) PurchaseRepository {
	return PurchaseRepository{
		pst: pst,
	}
}

func (repo PurchaseRepository) Create(doc models.PurchaseDoc) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.PurchaseDoc{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo PurchaseRepository) Update(guid string, doc models.PurchaseDoc) error {
	err := repo.pst.UpdateOne(&models.PurchaseDoc{}, "guidFixed", guid, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo PurchaseRepository) Delete(guid string, shopID string) error {
	err := repo.pst.SoftDelete(&models.PurchaseDoc{}, bson.M{"guidFixed": guid, "shopID": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo PurchaseRepository) FindByGuid(guid string, shopID string) (models.PurchaseDoc, error) {
	doc := &models.PurchaseDoc{}
	err := repo.pst.FindOne(&models.PurchaseDoc{}, bson.M{"shopID": shopID, "guidFixed": guid, "deleted": false}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo PurchaseRepository) FindPage(shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error) {

	docList := []models.PurchaseInfo{}
	pagination, err := repo.pst.FindPage(&models.PurchaseInfo{}, limit, page, bson.M{
		"shopID":  shopID,
		"deleted": false,
		"$or": []interface{}{
			bson.M{"guidFixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.PurchaseInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo PurchaseRepository) FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error) {

	docList := []models.PurchaseInfo{}
	pagination, err := repo.pst.FindPage(&models.PurchaseInfo{}, limit, page, bson.M{
		"shopID":    shopID,
		"guidFixed": guid,
		"deleted":   false,
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.PurchaseInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
