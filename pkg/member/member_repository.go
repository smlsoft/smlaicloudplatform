package member

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/member/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMemberRepository interface {
	Create(doc models.MemberDoc) (primitive.ObjectID, error)
	Update(shopID string, guid string, doc models.MemberDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.MemberDoc, error)
	FindPage(shopID string, pageable micromodels.Pageable) ([]models.MemberInfo, mongopagination.PaginationData, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.MemberDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.MemberActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.MemberDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.MemberActivity, error)
}

type MemberRepository struct {
	pst microservice.IPersisterMongo
	repositories.ActivityRepository[models.MemberActivity, models.MemberDeleteActivity]
}

func NewMemberRepository(pst microservice.IPersisterMongo) *MemberRepository {

	insRepo := &MemberRepository{
		pst: pst,
	}

	insRepo.ActivityRepository = repositories.NewActivityRepository[models.MemberActivity, models.MemberDeleteActivity](pst)

	return insRepo
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

func (repo MemberRepository) FindPage(shopID string, pageable micromodels.Pageable) ([]models.MemberInfo, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"name": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "",
			}}},
		},
	}

	docList := []models.MemberInfo{}
	pagination, err := repo.pst.FindPage(&models.MemberInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.MemberInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
