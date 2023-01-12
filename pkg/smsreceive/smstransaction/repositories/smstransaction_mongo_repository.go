package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/smsreceive/smstransaction/models"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISmsTransactionRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.SmsTransactionDoc) (string, error)
	CreateInBatch(docList []models.SmsTransactionDoc) error
	Update(shopID string, guid string, doc models.SmsTransactionDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.SmsTransactionInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SmsTransactionDoc, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.SmsTransactionDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.SmsTransactionInfo, mongopagination.PaginationData, error)

	FindFilterSms(shopID string, storefrontGUID string, address string, startTime time.Time, endTime time.Time) ([]models.SmsTransactionInfo, error)
}

type SmsTransactionRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SmsTransactionDoc]
	repositories.SearchRepository[models.SmsTransactionInfo]
	repositories.GuidRepository[models.SmsTransactionItemGuid]
}

func NewSmsTransactionRepository(pst microservice.IPersisterMongo) SmsTransactionRepository {

	insRepo := SmsTransactionRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SmsTransactionDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SmsTransactionInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SmsTransactionItemGuid](pst)

	return insRepo
}

func (repo SmsTransactionRepository) FindFilterSms(shopID string, storefrontGUID string, address string, startTime time.Time, endTime time.Time) ([]models.SmsTransactionInfo, error) {

	filters := bson.M{
		"shopid":         shopID,
		"storefrontguid": storefrontGUID,
		"deletedat":      bson.M{"$exists": false},
		"address":        address,
		"status":         0,
		"createdat": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	docList := []models.SmsTransactionInfo{}
	err := repo.pst.Find(models.SmsTransactionInfo{}, filters, &docList)

	if err != nil {
		return []models.SmsTransactionInfo{}, nil
	}

	return docList, nil
}
