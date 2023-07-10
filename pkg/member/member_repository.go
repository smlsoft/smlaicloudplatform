package member

import (
	"context"
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
	Create(ctx context.Context, doc models.MemberDoc) (primitive.ObjectID, error)
	Update(ctx context.Context, shopID string, guid string, doc models.MemberDoc) error
	Delete(ctx context.Context, shopID string, guid string, username string) error
	FindByGuid(ctx context.Context, shopID string, guid string) (models.MemberDoc, error)
	FindPage(ctx context.Context, shopID string, pageable micromodels.Pageable) ([]models.MemberInfo, mongopagination.PaginationData, error)

	FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.MemberDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.MemberActivity, mongopagination.PaginationData, error)
	FindDeletedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.MemberDeleteActivity, error)
	FindCreatedOrUpdatedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.MemberActivity, error)
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

func (repo MemberRepository) Create(ctx context.Context, doc models.MemberDoc) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(ctx, &models.MemberDoc{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return idx, nil
}

func (repo MemberRepository) Update(ctx context.Context, shopID string, guid string, doc models.MemberDoc) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}
	err := repo.pst.UpdateOne(ctx, &models.MemberDoc{}, filterDoc, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo MemberRepository) Delete(ctx context.Context, shopID string, guid string, uername string) error {
	err := repo.pst.SoftDeleteLastUpdate(ctx, &models.MemberDoc{}, uername, bson.M{"guidfixed": guid, "shopid": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo MemberRepository) FindByGuid(ctx context.Context, shopID string, guid string) (models.MemberDoc, error) {
	doc := &models.MemberDoc{}
	err := repo.pst.FindOne(ctx, &models.MemberDoc{}, bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo MemberRepository) FindPage(ctx context.Context, shopID string, pageable micromodels.Pageable) ([]models.MemberInfo, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"name": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "i",
			}}},
		},
	}

	docList := []models.MemberInfo{}
	pagination, err := repo.pst.FindPage(ctx, &models.MemberInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.MemberInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
