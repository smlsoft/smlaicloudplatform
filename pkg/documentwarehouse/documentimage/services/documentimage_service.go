package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/repositories"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDocumentImageService interface {
	CreateDocumentImage(shopID string, authUsername string, doc models.DocumentImageRequest) (string, string, error)
	BulkCreateDocumentImage(shopID string, authUsername string, docs []models.DocumentImageRequest) error
	InfoDocumentImage(shopID string, guid string) (models.DocumentImageInfo, error)
	SearchDocumentImage(shopID string, matchFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)
	UploadDocumentImage(shopID string, authUsername string, fh *multipart.FileHeader) (*models.DocumentImageInfo, error)
	CreateImageEdit(shopID string, authUsername string, docImageGUID string, docRequest models.ImageEditRequest) error

	CreateDocumentImageGroup(shopID string, authUsername string, docImageGroup models.DocumentImageGroup) (string, error)
	GetDocumentImageDocRefGroup(shopID string, docImageGroupGUID string) (models.DocumentImageGroupInfo, error)
	GetDocumentImageGroupByDocRef(shopID string, docRef string) (models.DocumentImageGroupInfo, error)
	UpdateDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImageGroup models.DocumentImageGroup) error
	UpdateImageReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImages []models.ImageReferenceBody) error
	UpdateReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docRef models.Reference) error
	UpdateTagsInDocumentImageGroup(shopID string, authUsername string, groupGUID string, tags []string) error
	UpdateStatusDocumentImageGroup(shopID string, authUsername string, groupGUID string, status int8) error
	UnGroupDocumentImageGroup(shopID string, authUsername string, groupGUID string) ([]string, error)
	ListDocumentImageGroup(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentImageGroupInfo, mongopagination.PaginationData, error)
	DeleteReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docRef models.Reference) error
	DeleteDocumentImageGroupByGuid(shopID string, authUsername string, DocumentImageGroupGuidFixed string) error
	DeleteDocumentImageGroupByGuids(shopID string, authUsername string, documentImageGroupGuidFixeds []string) error
	XSortsUpdate(shopID string, authUsername string, taskGUID string, xsorts []models.XSortDocumentImageGroupRequest) error

	UpdateDocumentImageReferenceGroup() error
}

type DocumentImageService struct {
	repoImageGroup               repositories.DocumentImageGroupRepository
	repoImage                    repositories.IDocumentImageRepository
	repoMessagequeue             repositories.DocumentImageMessageQueueRepository
	FilePersister                microservice.IPersisterFile
	maxImageReferences           int
	timeNowFnc                   func() time.Time
	newDocumentImageGUIDFnc      func() string
	newDocumentImageGroupGUIDFnc func() string
}

func NewDocumentImageService(repo repositories.IDocumentImageRepository, repoImageGroup repositories.DocumentImageGroupRepository, repoMessagequeue repositories.DocumentImageMessageQueueRepository, filePersister microservice.IPersisterFile) DocumentImageService {
	return DocumentImageService{
		maxImageReferences: 100,
		repoImageGroup:     repoImageGroup,
		repoImage:          repo,
		repoMessagequeue:   repoMessagequeue,
		FilePersister:      filePersister,
		timeNowFnc: func() time.Time {
			return time.Now()
		},
		newDocumentImageGUIDFnc:      utils.NewGUID,
		newDocumentImageGroupGUIDFnc: utils.NewGUID,
	}
}

func (svc DocumentImageService) CreateDocumentImage(shopID string, authUsername string, docRequest models.DocumentImageRequest) (string, string, error) {

	findDocImgGroup := models.DocumentImageGroupDoc{}
	if len(docRequest.DocumentImageGroupGUID) > 0 {
		_, err := svc.repoImageGroup.FindByGuid(shopID, docRequest.DocumentImageGroupGUID)

		if err != nil {
			return "", "", err
		}

		// if len(findDocImgGroup.GuidFixed) == 0 {
		// 	return "", "", errors.New("document image group not found")
		// }
	}

	// do upload first

	createdAt := svc.timeNowFnc()

	documentImageGUID := svc.newDocumentImageGUIDFnc()

	docData := models.DocumentImageDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = documentImageGUID
	docData.DocumentImage = docRequest.DocumentImage

	docData.Edits = []models.ImageEdit{}
	docData.References = []models.Reference{}
	docData.MetaFileAt = docRequest.MetaFileAt

	docData.CreatedBy = authUsername
	docData.CreatedAt = createdAt

	docData.UploadedBy = authUsername
	docData.UploadedAt = createdAt

	tags := []string{}
	if docRequest.Tags != nil {
		tags = *docRequest.Tags
	}

	// image group
	imageGroupGUID := svc.newDocumentImageGroupGUIDFnc()

	if len(findDocImgGroup.GuidFixed) > 0 {
		imageGroupGUID = findDocImgGroup.GuidFixed
	}

	docImageRef := svc.documentImageToImageReference(documentImageGUID, docRequest.DocumentImage, authUsername, createdAt)
	docDataImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, docImageRef, docRequest.ImageURI, tags, docRequest.TaskGUID, docRequest.PathTask, createdAt)

	newXOrderDocImgGroup, _ := svc.newXOrderDocumentImageGroup(shopID, docRequest.TaskGUID)
	docDataImageGroup.XOrder = newXOrderDocImgGroup

	err := svc.repoImageGroup.Transaction(func() error {

		_, err := svc.repoImage.Create(docData)

		if err != nil {
			return err
		}

		_, err = svc.repoImageGroup.Create(docDataImageGroup)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", "", err
	}

	_, err = svc.messageQueueReCountDocumentImageGroup(shopID, docRequest.TaskGUID)
	if err != nil {
		return "", "", err
	}

	return documentImageGUID, imageGroupGUID, nil
}

func (svc DocumentImageService) CreateDocumentImageWithTask(shopID string, authUsername string, docRequest models.DocumentImageRequest) (string, string, error) {

	// do upload first

	createdAt := svc.timeNowFnc()

	documentImageGUID := svc.newDocumentImageGUIDFnc()

	docData := models.DocumentImageDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = documentImageGUID
	docData.DocumentImage = docRequest.DocumentImage

	docData.Edits = []models.ImageEdit{}
	docData.References = []models.Reference{}
	docData.MetaFileAt = docRequest.MetaFileAt

	docData.CreatedBy = authUsername
	docData.CreatedAt = createdAt

	docData.UploadedBy = authUsername
	docData.UploadedAt = createdAt

	tags := []string{}
	if docRequest.Tags != nil {
		tags = *docRequest.Tags
	}

	// image group
	imageGroupGUID := svc.newDocumentImageGroupGUIDFnc()
	docImageRef := svc.documentImageToImageReference(documentImageGUID, docRequest.DocumentImage, authUsername, createdAt)
	docDataImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, docImageRef, docRequest.ImageURI, tags, docRequest.TaskGUID, docRequest.PathTask, createdAt)

	newXOrderDocImgGroup, _ := svc.newXOrderDocumentImageGroup(shopID, docRequest.TaskGUID)

	docDataImageGroup.XOrder = newXOrderDocImgGroup

	err := svc.repoImageGroup.Transaction(func() error {

		_, err := svc.repoImage.Create(docData)

		if err != nil {
			return err
		}

		_, err = svc.repoImageGroup.Create(docDataImageGroup)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", "", err
	}

	_, err = svc.messageQueueReCountDocumentImageGroup(shopID, docRequest.TaskGUID)
	if err != nil {
		return "", "", err
	}

	return documentImageGUID, imageGroupGUID, nil
}

func (svc DocumentImageService) CreateImageEdit(shopID string, authUsername string, docImageGUID string, docRequest models.ImageEditRequest) error {

	findDoc, err := svc.repoImage.FindByGuid(shopID, docImageGUID)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) == 0 {
		return errors.New("document image not found")
	}

	tempDoc := findDoc

	if tempDoc.Edits == nil {
		tempDoc.Edits = []models.ImageEdit{}
	}

	imageEdit := models.ImageEdit{}

	imageEdit.ImageURI = docRequest.ImageURI
	imageEdit.EditedBy = authUsername
	imageEdit.EditedAt = svc.timeNowFnc()

	tempDoc.Edits = append(tempDoc.Edits, imageEdit)

	err = svc.repoImage.Update(shopID, docImageGUID, tempDoc)
	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) BulkCreateDocumentImage(shopID string, authUsername string, docs []models.DocumentImageRequest) error {

	// do upload first

	createdAt := svc.timeNowFnc()
	docDataList := []models.DocumentImageDoc{}
	docDataImageGroupList := []models.DocumentImageGroupDoc{}

	taskLastXOrder := map[string]int{}

	for _, doc := range docs {
		if len(doc.TaskGUID) < 1 {
			return fmt.Errorf("job is empty")
		}

		_, ok := taskLastXOrder[doc.TaskGUID]
		if !ok {
			newXOrderDocImgGroup, _ := svc.newXOrderDocumentImageGroup(shopID, doc.TaskGUID)
			taskLastXOrder[doc.TaskGUID] = newXOrderDocImgGroup
		} else {
			taskLastXOrder[doc.TaskGUID]++
		}

		documentImageGUID := svc.newDocumentImageGUIDFnc()

		docData := models.DocumentImageDoc{}
		docData.ShopID = shopID
		docData.GuidFixed = documentImageGUID
		docData.DocumentImage = doc.DocumentImage

		docData.Edits = []models.ImageEdit{}
		docData.References = []models.Reference{}
		docData.MetaFileAt = doc.MetaFileAt

		docData.CreatedBy = authUsername
		docData.CreatedAt = createdAt

		// docData.UpdatedBy = authUsername
		// docData.UpdatedAt = createdAt

		docData.UploadedBy = authUsername
		// docData.UploadedAt = createdAt

		// image group
		imageGroupGUID := svc.newDocumentImageGroupGUIDFnc()
		docImageRef := svc.documentImageToImageReference(documentImageGUID, doc.DocumentImage, authUsername, createdAt)

		tags := []string{}
		if doc.Tags != nil {
			tags = *doc.Tags
		}

		docDataImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, docImageRef, doc.ImageURI, tags, doc.TaskGUID, doc.PathTask, createdAt)

		docDataImageGroup.XOrder = taskLastXOrder[doc.TaskGUID]

		docDataList = append(docDataList, docData)
		docDataImageGroupList = append(docDataImageGroupList, docDataImageGroup)
	}

	err := svc.repoImageGroup.Transaction(func() error {

		err := svc.repoImage.CreateInBatch(docDataList)

		if err != nil {
			return err
		}

		err = svc.repoImageGroup.CreateInBatch(docDataImageGroupList)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	taskGUIDsChanged := map[string]struct{}{}
	for _, tempDocGroup := range docDataImageGroupList {
		taskGUIDsChanged[tempDocGroup.TaskGUID] = struct{}{}
	}

	for taskGUID := range taskGUIDsChanged {
		_, err = svc.messageQueueReCountDocumentImageGroup(shopID, taskGUID)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (svc DocumentImageService) InfoDocumentImage(shopID string, guid string) (models.DocumentImageInfo, error) {

	findDoc, err := svc.repoImage.FindByGuid(shopID, guid)

	if err != nil {
		return models.DocumentImageInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DocumentImageInfo{}, errors.New("document not found")
	}

	return findDoc.DocumentImageInfo, nil
}

func (svc DocumentImageService) SearchDocumentImage(shopID string, matchFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentImageInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{"guidfixed", "documentref", "module"}
	docList, pagination, err := svc.repoImage.FindPageFilter(shopID, matchFilters, searchInFields, pageable)

	if err != nil {
		return []models.DocumentImageInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc DocumentImageService) UploadDocumentImage(shopID string, authUsername string, fh *multipart.FileHeader) (*models.DocumentImageInfo, error) {

	if fh.Filename == "" {
		return nil, errors.New("image file name not found")
	}
	// try upload
	fileUploadMetadataSlice := strings.Split(fh.Filename, ".")
	fileName := svc.newDocumentImageGUIDFnc() //fileUploadMetadataSlice[0]
	fileExtension := fileUploadMetadataSlice[1]

	fileNameWithShop := fmt.Sprintf("%s/%s", shopID, fileName)

	imageUri, err := svc.FilePersister.Save(fh, fileNameWithShop, fileExtension)
	if err != nil {
		return nil, err
	}

	// create document image
	doc := new(models.DocumentImageDoc)
	doc.GuidFixed = svc.newDocumentImageGUIDFnc()
	doc.ImageURI = imageUri
	doc.ShopID = shopID
	doc.UploadedBy = authUsername
	doc.UploadedAt = svc.timeNowFnc()
	doc.CreatedBy = authUsername
	doc.CreatedAt = doc.UploadedAt

	docRequest := models.DocumentImageRequest{
		DocumentImage: doc.DocumentImage,
		Tags:          &[]string{},
		TaskGUID:      "",
		PathTask:      "",
	}
	_, _, err = svc.CreateDocumentImage(shopID, authUsername, docRequest)

	if err != nil {
		return nil, err
	}

	return &doc.DocumentImageInfo, err
}

func (svc DocumentImageService) UpdateDocumentImageReferenceGroup() error {
	findDocList, err := svc.repoImage.FindAll()

	if err != nil {
		return err
	}

	for _, findDoc := range findDocList {
		findGroupDoc, err := svc.repoImageGroup.FindOneByDocumentImageGUIDAll(findDoc.GuidFixed)

		if err != nil {
			return err
		}

		refGroups := []models.ReferenceGroup{}
		if findGroupDoc.ID == primitive.NilObjectID {
			refGroup := models.ReferenceGroup{}
			refGroup.GroupType = ""
			refGroup.ParentGUID = ""
			refGroup.XOrder = 1
			refGroup.XType = 0

			refGroups = append(refGroups, refGroup)

		}

		findDoc.ReferenceGroups = refGroups

		err = svc.repoImage.UpdateAll(findDoc)

		if err != nil {
			return err
		}
	}

	return nil
}
