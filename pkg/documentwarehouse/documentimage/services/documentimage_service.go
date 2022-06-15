package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/repositories"
	"smlcloudplatform/pkg/utils"
	"strings"
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
	UploadDocumentImage(shopID string, authUsername string, moduleName string, fh *multipart.FileHeader) (*models.DocumentImageInfo, error)
}

type DocumentImageService struct {
	Repo          repositories.IDocumentImageRepository
	FilePersister microservice.IPersisterFile
	NowFn         func() time.Time
	NewGUIDFn     func() string
}

func NewDocumentImageService(repo repositories.IDocumentImageRepository, filePersister microservice.IPersisterFile) DocumentImageService {
	return DocumentImageService{
		Repo:          repo,
		FilePersister: filePersister,
		NowFn: func() time.Time {
			return time.Now()
		},
		NewGUIDFn: func() string { return utils.NewGUID() },
	}
}

func (svc DocumentImageService) CreateDocumentImage(shopID string, authUsername string, doc models.DocumentImage) (string, error) {

	// do upload first

	newGuidFixed := utils.NewGUID()

	docData := models.DocumentImageDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.DocumentImage = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.Repo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc DocumentImageService) UpdateDocumentImage(shopID string, guid string, authUsername string, doc models.DocumentImage) error {

	findDoc, err := svc.Repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.DocumentImage = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.Repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc DocumentImageService) DeleteDocumentImage(shopID string, guid string, authUsername string) error {
	err := svc.Repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc DocumentImageService) InfoDocumentImage(shopID string, guid string) (models.DocumentImageInfo, error) {

	findDoc, err := svc.Repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.DocumentImageInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DocumentImageInfo{}, errors.New("document not found")
	}

	return findDoc.DocumentImageInfo, nil

}

func (svc DocumentImageService) SearchDocumentImage(shopID string, q string, page int, limit int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.Repo.FindPage(shopID, []string{"guidfixed", "documentref", "module"}, q, page, limit)

	if err != nil {
		return []models.DocumentImageInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc DocumentImageService) UploadDocumentImage(shopID string, authUsername string, moduleName string, fh *multipart.FileHeader) (*models.DocumentImageInfo, error) {

	if fh.Filename == "" {
		return nil, errors.New("Image Filename Not Found")
	}
	// try upload
	fileUploadMetadataSlice := strings.Split(fh.Filename, ".")
	fileName := svc.NewGUIDFn() //fileUploadMetadataSlice[0]
	fileExtension := fileUploadMetadataSlice[1]

	fileNameWithShop := fmt.Sprintf("%s/%s", shopID, fileName)

	imageUri, err := svc.FilePersister.Save(fh, fileNameWithShop, fileExtension)
	if err != nil {
		return nil, err
	}

	// create document image
	doc := new(models.DocumentImageDoc)
	doc.GuidFixed = svc.NewGUIDFn()
	doc.ImageUri = imageUri
	doc.ShopID = shopID
	doc.DocumentRef = svc.NewGUIDFn()
	doc.Module = moduleName
	doc.UploadedBy = authUsername
	doc.UploadedAt = svc.NowFn()
	doc.CreatedBy = authUsername
	doc.CreatedAt = doc.UploadedAt

	_, err = svc.Repo.Create(*doc)
	if err != nil {
		return nil, err
	}

	return &doc.DocumentImageInfo, err
}
