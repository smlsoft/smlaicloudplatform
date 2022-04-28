package saleinvoice

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISaleinvoiceRepository interface {
	Create(trans models.SaleinvoiceDoc) (primitive.ObjectID, error)
	Update(guid string, trans models.SaleinvoiceDoc) error
	Delete(guid string, shopID string, username string) error
	FindByGuid(guid string, shopID string) (models.SaleinvoiceDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.SaleinvoiceInfo, paginate.PaginationData, error)
	FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.SaleinvoiceInfo, paginate.PaginationData, error)
}

type SaleinvoiceRepository struct {
	pst microservice.IPersisterMongo
}

func NewSaleinvoiceRepository(pst microservice.IPersisterMongo) SaleinvoiceRepository {
	return SaleinvoiceRepository{
		pst: pst,
	}
}

func (repo SaleinvoiceRepository) Create(trans models.SaleinvoiceDoc) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.SaleinvoiceDoc{}, trans)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo SaleinvoiceRepository) Update(guid string, trans models.SaleinvoiceDoc) error {
	err := repo.pst.UpdateOne(&models.SaleinvoiceDoc{}, "guidfixed", guid, trans)
	if err != nil {
		return err
	}
	return nil
}

func (repo SaleinvoiceRepository) Delete(guid string, shopID string, username string) error {
	err := repo.pst.SoftDelete(&models.SaleinvoiceDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo SaleinvoiceRepository) FindByGuid(guid string, shopID string) (models.SaleinvoiceDoc, error) {
	trans := &models.SaleinvoiceDoc{}
	err := repo.pst.FindOne(
		&models.SaleinvoiceDoc{},
		bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}},
		trans,
	)
	if err != nil {
		return *trans, err
	}
	return *trans, nil
}

func (repo SaleinvoiceRepository) FindPage(shopID string, q string, page int, limit int) ([]models.SaleinvoiceInfo, paginate.PaginationData, error) {

	transList := []models.SaleinvoiceInfo{}
	pagination, err := repo.pst.FindPage(&models.SaleinvoiceInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &transList)

	if err != nil {
		return []models.SaleinvoiceInfo{}, paginate.PaginationData{}, err
	}

	return transList, pagination, nil
}

func (repo SaleinvoiceRepository) FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.SaleinvoiceInfo, paginate.PaginationData, error) {

	transList := []models.SaleinvoiceInfo{}
	pagination, err := repo.pst.FindPage(&models.Saleinvoice{}, limit, page, bson.M{
		"shopid":    shopID,
		"guidfixed": guid,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &transList)

	if err != nil {
		return []models.SaleinvoiceInfo{}, paginate.PaginationData{}, err
	}

	return transList, pagination, nil
}
