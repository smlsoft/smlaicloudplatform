package member

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMemberRepository interface {
	Create(doc models.MemberDoc) (primitive.ObjectID, error)
	Update(shopID string, guid string, doc models.MemberDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.MemberDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.MemberDeleteActivity, paginate.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.MemberActivity, paginate.PaginationData, error)
}

type MemberRepository struct {
	pst microservice.IPersisterMongo
}

func NewMemberRepository(pst microservice.IPersisterMongo) MemberRepository {
	return MemberRepository{
		pst: pst,
	}
}

func (repo MemberRepository) Create(doc models.MemberDoc) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.MemberDoc{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return idx, nil
}

func (repo MemberRepository) Update(shopID string, guid string, doc models.MemberDoc) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}
	err := repo.pst.UpdateOne(&models.MemberDoc{}, filterDoc, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo MemberRepository) Delete(shopID string, guid string, uername string) error {
	err := repo.pst.SoftDeleteLastUpdate(&models.MemberDoc{}, uername, bson.M{"guidfixed": guid, "shopid": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo MemberRepository) FindByGuid(shopID string, guid string) (models.MemberDoc, error) {
	doc := &models.MemberDoc{}
	err := repo.pst.FindOne(&models.MemberDoc{}, bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo MemberRepository) FindPage(shopID string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error) {

	docList := []models.MemberInfo{}
	pagination, err := repo.pst.FindPage(&models.MemberInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"name": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.MemberInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo MemberRepository) FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.MemberDeleteActivity, paginate.PaginationData, error) {

	docList := []models.MemberDeleteActivity{}
	pagination, err := repo.pst.FindPage(&models.MemberDeleteActivity{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}, &docList)

	if err != nil {
		return []models.MemberDeleteActivity{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo MemberRepository) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.MemberActivity, paginate.PaginationData, error) {

	docList := []models.MemberActivity{}
	pagination, err := repo.pst.FindPage(&models.MemberActivity{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}, &docList)

	if err != nil {
		return []models.MemberActivity{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
