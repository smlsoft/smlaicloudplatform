package member

import (
	"context"
	"smlcloudplatform/internal/member/models"
	"smlcloudplatform/pkg/microservice"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMemberRepository interface {
	Create(ctx context.Context, doc models.MemberDoc) (primitive.ObjectID, error)
	Update(ctx context.Context, guid string, doc models.MemberDoc) error
	FindByGuid(ctx context.Context, guid string) (models.MemberDoc, error)

	FindByLineUID(ctx context.Context, lineUID string) (models.MemberDoc, error)
}

type MemberRepository struct {
	pst microservice.IPersisterMongo
}

func NewMemberRepository(pst microservice.IPersisterMongo) *MemberRepository {

	insRepo := &MemberRepository{
		pst: pst,
	}

	return insRepo
}

func (repo MemberRepository) Create(ctx context.Context, doc models.MemberDoc) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(ctx, &models.MemberDoc{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return idx, nil
}

func (repo MemberRepository) Update(ctx context.Context, guid string, doc models.MemberDoc) error {
	filterDoc := map[string]interface{}{
		"guidfixed": guid,
	}
	err := repo.pst.UpdateOne(ctx, &models.MemberDoc{}, filterDoc, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo MemberRepository) FindByGuid(ctx context.Context, guid string) (models.MemberDoc, error) {
	doc := &models.MemberDoc{}
	err := repo.pst.FindOne(ctx, &models.MemberDoc{}, bson.M{"guidfixed": guid, "deletedat": bson.M{"$exists": false}}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo MemberRepository) FindByLineUID(ctx context.Context, lineUID string) (models.MemberDoc, error) {
	doc := &models.MemberDoc{}
	err := repo.pst.FindOne(ctx, &models.MemberDoc{}, bson.M{"lineuid": lineUID, "deletedat": bson.M{"$exists": false}}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}
