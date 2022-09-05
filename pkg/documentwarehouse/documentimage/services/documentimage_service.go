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
	UpdateDocumentImageStatus(shopID string, guid string, status int8) error
	UpdateDocumentImageStatusByDocumentRef(shopID string, docRef string, status int8) error
	DeleteDocumentImage(shopID string, guid string, authUsername string) error
	InfoDocumentImage(shopID string, guid string) (models.DocumentImageInfo, error)
	SearchDocumentImage(shopID string, matchFilters map[string]interface{}, q string, page int, limit int, sorts map[string]int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)
	UploadDocumentImage(shopID string, authUsername string, moduleName string, fh *multipart.FileHeader) (*models.DocumentImageInfo, error)

	SaveDocumentImageDocRefGroup(shopID string, docRef string, docImages []string) error
	GetDocumentImageDocRefGroup(shopID string, docRef string) (models.DocumentImageGroup, error)
	ListDocumentImageDocRefGroup(shopID string, filters map[string]interface{}, q string, page int, limit int) ([]models.DocumentImageGroup, mongopagination.PaginationData, error)
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

func (svc DocumentImageService) UpdateDocumentImageStatus(shopID string, guid string, status int8) error {

	if len(guid) < 1 {
		return errors.New("guid is not empty")
	}

	err := svc.Repo.UpdateDocumentImageStatus(shopID, guid, status)

	if err != nil {
		return err
	}
	return nil
}

func (svc DocumentImageService) UpdateDocumentImageStatusByDocumentRef(shopID string, docRef string, status int8) error {

	err := svc.Repo.UpdateDocumentImageStatusByDocumentRef(shopID, docRef, status)

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

func (svc DocumentImageService) SearchDocumentImage(shopID string, matchFilters map[string]interface{}, q string, page int, limit int, sorts map[string]int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error) {
	// docList, pagination, err := svc.Repo.FindPage(shopID, []string{"guidfixed", "documentref", "module"}, q, page, limit)
	docList, pagination, err := svc.Repo.FindPageFilterSort(shopID, matchFilters, []string{"guidfixed", "documentref", "module"}, q, page, limit, sorts)

	if err != nil {
		return []models.DocumentImageInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc DocumentImageService) UploadDocumentImage(shopID string, authUsername string, moduleName string, fh *multipart.FileHeader) (*models.DocumentImageInfo, error) {

	if fh.Filename == "" {
		return nil, errors.New("image file name not found")
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

func (svc DocumentImageService) SaveDocumentImageDocRefGroup(shopID string, docRef string, docImages []string) error {
	return svc.Repo.SaveDocumentImageDocRefGroup(shopID, docRef, docImages)
}

func (svc DocumentImageService) ListDocumentImageDocRefGroup(shopID string, filters map[string]interface{}, q string, page int, limit int) ([]models.DocumentImageGroup, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.Repo.ListDocumentImageGroup(shopID, filters, q, page, limit)

	return docList, pagination, err
}

func (svc DocumentImageService) GetDocumentImageDocRefGroup(shopID string, docRef string) (models.DocumentImageGroup, error) {
	docList, err := svc.Repo.GetDocumentImageGroup(shopID, docRef)

	return docList, err
}
