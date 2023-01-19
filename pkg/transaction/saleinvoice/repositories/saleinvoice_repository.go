package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISaleinvoiceRepository interface {
	Create(trans models.SaleinvoiceDoc) (primitive.ObjectID, error)
	Update(shopID string, guid string, trans models.SaleinvoiceDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.SaleinvoiceDoc, error)
	FindPage(shopID string, pageable micromodels.Pageable) ([]models.SaleinvoiceInfo, mongopagination.PaginationData, error)
	FindItemsByGuidPage(guid string, shopID string, pageable micromodels.Pageable) ([]models.SaleinvoiceInfo, mongopagination.PaginationData, error)
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

func (repo SaleinvoiceRepository) Update(shopID string, guid string, trans models.SaleinvoiceDoc) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}
	err := repo.pst.UpdateOne(&models.SaleinvoiceDoc{}, filterDoc, trans)
	if err != nil {
		return err
	}
	return nil
}

func (repo SaleinvoiceRepository) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&models.SaleinvoiceDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo SaleinvoiceRepository) FindByGuid(shopID string, guid string) (models.SaleinvoiceDoc, error) {
	trans := &models.SaleinvoiceDoc{}
	filterDoc := bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}}

	err := repo.pst.FindOne(
		&models.SaleinvoiceDoc{},
		filterDoc,
		trans,
	)
	if err != nil {
		return *trans, err
	}
	return *trans, nil
}

func (repo SaleinvoiceRepository) FindPage(shopID string, pageable micromodels.Pageable) ([]models.SaleinvoiceInfo, mongopagination.PaginationData, error) {
	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "",
			}}},
		},
	}

	transList := []models.SaleinvoiceInfo{}
	pagination, err := repo.pst.FindPage(&models.SaleinvoiceInfo{}, filterQueries, pageable, &transList)

	if err != nil {
		return []models.SaleinvoiceInfo{}, mongopagination.PaginationData{}, err
	}

	return transList, pagination, nil
}

func (repo SaleinvoiceRepository) FindItemsByGuidPage(guid string, shopID string, pageable micromodels.Pageable) ([]models.SaleinvoiceInfo, mongopagination.PaginationData, error) {
	filterQueries := bson.M{
		"shopid":    shopID,
		"guidfixed": guid,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "",
			}}},
		},
	}

	transList := []models.SaleinvoiceInfo{}
	pagination, err := repo.pst.FindPage(&models.Saleinvoice{}, filterQueries, pageable, &transList)

	if err != nil {
		return []models.SaleinvoiceInfo{}, mongopagination.PaginationData{}, err
	}

	return transList, pagination, nil
}
