package merchant

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMerchantRepository interface {
	Create(merchant models.Merchant) (string, error)
	Update(guid string, merchant models.Merchant) error
	FindByGuid(guid string) (models.Merchant, error)
	FindPage(q string, page int, limit int) ([]models.MerchantInfo, paginate.PaginationData, error)
	Delete(guid string) error
}

type MerchantRepository struct {
	pst microservice.IPersisterMongo
}

func NewMerchantRepository(pst microservice.IPersisterMongo) IMerchantRepository {
	return &MerchantRepository{
		pst: pst,
	}
}

func (repo *MerchantRepository) Create(merchant models.Merchant) (string, error) {
	idx, err := repo.pst.Create(&models.Merchant{}, merchant)
	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo *MerchantRepository) Update(guid string, merchant models.Merchant) error {
	err := repo.pst.UpdateOne(&models.Merchant{}, "guidFixed", guid, merchant)

	if err != nil {
		return err
	}

	return nil
}

func (repo *MerchantRepository) FindByGuid(guid string) (models.Merchant, error) {
	findMerchant := &models.Merchant{}
	err := repo.pst.FindOne(&models.Merchant{}, bson.M{"guidFixed": guid, "deleted": false}, findMerchant)

	if err != nil {
		return models.Merchant{}, err
	}
	return *findMerchant, err
}

func (repo *MerchantRepository) FindPage(q string, page int, limit int) ([]models.MerchantInfo, paginate.PaginationData, error) {

	merchantList := []models.MerchantInfo{}

	pagination, err := repo.pst.FindPage(&models.Merchant{}, limit, page, bson.M{
		"deleted": false,
		"name1": bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}}}, &merchantList)

	if err != nil {
		return []models.MerchantInfo{}, paginate.PaginationData{}, err
	}

	return merchantList, pagination, nil
}

func (repo *MerchantRepository) Delete(guid string) error {
	err := repo.pst.SoftDeleteByID(&models.Merchant{}, guid)
	if err != nil {
		return err
	}
	return nil
}
