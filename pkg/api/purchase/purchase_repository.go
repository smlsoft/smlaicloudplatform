package purchase

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPurchaseRepository interface {
	Create(doc models.Purchase) (primitive.ObjectID, error)
	Update(guid string, doc models.Purchase) error
	Delete(guid string, merchantId string) error
	FindByGuid(guid string, merchantId string) (models.Purchase, error)
	FindPage(merchantId string, q string, page int, limit int) ([]models.Purchase, paginate.PaginationData, error)
	FindItemsByGuidPage(guid string, merchantId string, q string, page int, limit int) ([]models.Purchase, paginate.PaginationData, error)
}

type PurchaseRepository struct {
	pst microservice.IPersisterMongo
}

func NewPurchaseRepository(pst microservice.IPersisterMongo) IPurchaseRepository {
	return &PurchaseRepository{
		pst: pst,
	}
}

func (repo *PurchaseRepository) Create(doc models.Purchase) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.Purchase{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo *PurchaseRepository) Update(guid string, doc models.Purchase) error {
	err := repo.pst.UpdateOne(&models.Purchase{}, "guidFixed", guid, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PurchaseRepository) Delete(guid string, merchantId string) error {
	err := repo.pst.SoftDelete(&models.Purchase{}, bson.M{"guidFixed": guid, "merchantId": merchantId})
	if err != nil {
		return err
	}
	return nil
}

func (repo *PurchaseRepository) FindByGuid(guid string, merchantId string) (models.Purchase, error) {
	doc := &models.Purchase{}
	err := repo.pst.FindOne(&models.Purchase{}, bson.M{"merchantId": merchantId, "guidFixed": guid, "deleted": false}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo *PurchaseRepository) FindPage(merchantId string, q string, page int, limit int) ([]models.Purchase, paginate.PaginationData, error) {

	docList := []models.Purchase{}
	pagination, err := repo.pst.FindPage(&models.Purchase{}, limit, page, bson.M{
		"merchantId": merchantId,
		"deleted":    false,
		"$or": []interface{}{
			bson.M{"guidFixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.Purchase{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo *PurchaseRepository) FindItemsByGuidPage(guid string, merchantId string, q string, page int, limit int) ([]models.Purchase, paginate.PaginationData, error) {

	docList := []models.Purchase{}
	pagination, err := repo.pst.FindPage(&models.Purchase{}, limit, page, bson.M{
		"merchantId": merchantId,
		"guidFixed":  guid,
		"deleted":    false,
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.Purchase{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
