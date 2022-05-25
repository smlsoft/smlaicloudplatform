package services

import (
	"errors"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDocumentImageService interface {
	CreateDocumentImage(shopID string, authUsername string, doc models.DocumentImage) (string, error)
	UpdateDocumentImage(shopID string, guid string, authUsername string, doc models.DocumentImage) error
	DeleteDocumentImage(shopID string, guid string, authUsername string) error
	InfoDocumentImage(shopID string, guid string) (models.DocumentImageInfo, error)
	SearchDocumentImage(shopID string, q string, page int, limit int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)
}

type DocumentImageService struct {
	repo repositories.IDocumentImageRepository
}

func NewDocumentImageService(repo repositories.IDocumentImageRepository) DocumentImageService {
	return DocumentImageService{
		repo: repo,
	}
}

func (svc DocumentImageService) CreateDocumentImage(shopID string, authUsername string, doc models.DocumentImage) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.DocumentImageDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.DocumentImage = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc DocumentImageService) UpdateDocumentImage(shopID string, guid string, authUsername string, doc models.DocumentImage) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.DocumentImage = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc DocumentImageService) DeleteDocumentImage(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc DocumentImageService) InfoDocumentImage(shopID string, guid string) (models.DocumentImageInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.DocumentImageInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DocumentImageInfo{}, errors.New("document not found")
	}

	return findDoc.DocumentImageInfo, nil

}

func (svc DocumentImageService) SearchDocumentImage(shopID string, q string, page int, limit int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, []string{"guidfixed", "documentref", "module"}, q, page, limit)

	if err != nil {
		return []models.DocumentImageInfo{}, pagination, err
	}

	return docList, pagination, nil
}
