package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IDocumentImageRepository interface {
	Create(doc models.DocumentImageDoc) (string, error)
	CreateInBatch(doc []models.DocumentImageDoc) error
	Update(shopID string, guid string, doc models.DocumentImageDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters interface{}) (models.DocumentImageDoc, error)
	FindByReferenceDocNo(shopID string, docNo string) ([]models.DocumentImageDoc, error)
	FindByReference(shopID string, reference models.Reference) ([]models.DocumentImageDoc, error)
	FindByGuid(shopID string, guid string) (models.DocumentImageDoc, error)

	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.DocumentImageItemGuid, error)
	FindInGUIDs(shopID string, docImageGUIDs []string) ([]models.DocumentImageDoc, error)

	FindAll() ([]models.DocumentImageDoc, error)
	UpdateAll(doc models.DocumentImageDoc) error

	// UpdateDocumentImageStatus(shopID string, guid string, docnoGUIDRef string, status int8) error
	// UpdateDocumentImageStatusByDocumentRef(shopID string, docRef string, docnoGUIDRef string, status int8) error
	// SaveDocumentImageDocRefGroup(shopID string, docRef string, docImages []models.DocumentImageGroup) error
	// ListDocumentImageGroup(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentImageGroup, mongopagination.PaginationData, error)
	// GetDocumentImageGroup(shopID string, docRef string) (models.DocumentImageGroup, error)
}

type DocumentImageRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DocumentImageDoc]
	repositories.SearchRepository[models.DocumentImageInfo]
	repositories.GuidRepository[models.DocumentImageItemGuid]
}

func NewDocumentImageRepository(pst microservice.IPersisterMongo) DocumentImageRepository {
	insRepo := DocumentImageRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DocumentImageDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DocumentImageInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.DocumentImageItemGuid](pst)

	return insRepo
}

func (repo DocumentImageRepository) FindInGUIDs(shopID string, docImageGUIDs []string) ([]models.DocumentImageDoc, error) {
	docList := []models.DocumentImageDoc{}
	err := repo.pst.Find(models.DocumentImageDoc{}, bson.M{
		"guidfixed": bson.M{
			"$in": docImageGUIDs,
		},
		"deletedat": bson.M{"$exists": false},
	}, &docList)

	if err != nil {
		return []models.DocumentImageDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageRepository) FindByReferenceDocNo(shopID string, docNo string) ([]models.DocumentImageDoc, error) {
	docList := []models.DocumentImageDoc{}
	err := repo.pst.Find(models.DocumentImageDoc{}, bson.M{
		"references.docno": docNo,
		"deletedat":        bson.M{"$exists": false},
	}, &docList)

	if err != nil {
		return []models.DocumentImageDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageRepository) FindByReference(shopID string, reference models.Reference) ([]models.DocumentImageDoc, error) {
	docList := []models.DocumentImageDoc{}
	err := repo.pst.Find(models.DocumentImageDoc{}, bson.M{
		"references.module": reference.Module,
		"references.docno":  reference.DocNo,
		"deletedat":         bson.M{"$exists": false},
	}, &docList)

	if err != nil {
		return []models.DocumentImageDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageRepository) UpdateReject(shopID string, authUsername string, updatedAt time.Time, docImageGUID string, isReject bool) error {
	fillter := bson.M{
		"shopid":    shopID,
		"guidfixed": docImageGUID,
	}

	data := bson.M{
		"$set": bson.M{"isreject": isReject, "updatedby": authUsername, "updatedat": updatedAt},
	}

	return repo.pst.Update(models.DocumentImageDoc{}, fillter, data)
}

func (repo DocumentImageRepository) FindAll() ([]models.DocumentImageDoc, error) {
	docList := []models.DocumentImageDoc{}
	err := repo.pst.Find(models.DocumentImageDoc{}, bson.M{}, &docList)

	if err != nil {
		return []models.DocumentImageDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageRepository) UpdateAll(doc models.DocumentImageDoc) error {
	filterDoc := map[string]interface{}{
		"guidfixed": doc.GuidFixed,
	}

	err := repo.pst.UpdateOne(models.DocumentImageDoc{}, filterDoc, doc)

	if err != nil {
		return err
	}

	return nil
}
