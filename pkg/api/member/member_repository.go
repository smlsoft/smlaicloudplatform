package member

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMemberRepository interface {
	Create(doc models.MemberDoc) (primitive.ObjectID, error)
	Update(guid string, doc models.MemberDoc) error
	Delete(guid string, shopId string, username string) error
	FindByGuid(guid string, shopId string) (models.MemberDoc, error)
	FindPage(shopId string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error)
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

func (repo MemberRepository) Update(guid string, doc models.MemberDoc) error {
	err := repo.pst.UpdateOne(&models.MemberDoc{}, "guidFixed", guid, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo MemberRepository) Delete(guid string, shopId string, uername string) error {
	err := repo.pst.SoftDelete(&models.MemberDoc{}, uername, bson.M{"guidFixed": guid, "shopID": shopId})
	if err != nil {
		return err
	}
	return nil
}

func (repo MemberRepository) FindByGuid(guid string, shopId string) (models.MemberDoc, error) {
	doc := &models.MemberDoc{}
	err := repo.pst.FindOne(&models.MemberDoc{}, bson.M{"shopID": shopId, "guidFixed": guid, "deletedAt": bson.M{"$exists": false}}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo MemberRepository) FindPage(shopId string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error) {

	docList := []models.MemberInfo{}
	pagination, err := repo.pst.FindPage(&models.MemberInfo{}, limit, page, bson.M{
		"shopID":    shopId,
		"deletedAt": bson.M{"$exists": false},
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
