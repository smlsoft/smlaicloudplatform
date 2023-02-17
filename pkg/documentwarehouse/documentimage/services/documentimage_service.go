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
	"sort"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDocumentImageService interface {
	CreateDocumentImage(shopID string, authUsername string, doc models.DocumentImageRequest) (string, string, error)
	BulkCreateDocumentImage(shopID string, authUsername string, docs []models.DocumentImageRequest) error
	InfoDocumentImage(shopID string, guid string) (models.DocumentImageInfo, error)
	SearchDocumentImage(shopID string, matchFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)
	UploadDocumentImage(shopID string, authUsername string, fh *multipart.FileHeader) (*models.DocumentImageInfo, error)

	CreateDocumentImageGroup(shopID string, authUsername string, docImageGroup models.DocumentImageGroup) (string, error)
	GetDocumentImageDocRefGroup(shopID string, docImageGroupGUID string) (models.DocumentImageGroupInfo, error)
	GetDocumentImageGroupByDocRef(shopID string, docRef string) (models.DocumentImageGroupInfo, error)
	UpdateDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImageGroup models.DocumentImageGroup) error
	UpdateImageReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImages []models.ImageReferenceBody) error
	UpdateReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docRef models.Reference) error
	UpdateTagsInDocumentImageGroup(shopID string, authUsername string, groupGUID string, tags []string) error
	UpdateStatusDocumentImageGroup(shopID string, authUsername string, groupGUID string, status int8) error
	UnGroupDocumentImageGroup(shopID string, authUsername string, groupGUID string) error
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
		docDataImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, docImageRef, doc.ImageURI, *doc.Tags, doc.TaskGUID, doc.PathTask, createdAt)

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

		// if findGroupDoc.ID != primitive.NilObjectID {
		// 	fmt.Printf("image:: %v\n", findDoc.GuidFixed)
		// 	fmt.Printf("group:: %v\n", findGroupDoc.GuidFixed)
		// }

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

// Group

func (svc DocumentImageService) getDocumentImageNotReferencedInGroup(shopID string, currentGroupGUID string, docImageRefs []models.ImageReferenceBody) ([]models.ImageReference, []string, error) {

	docImageGUIDs := []string{}

	for _, imageRef := range docImageRefs {
		docImageGUIDs = append(docImageGUIDs, imageRef.DocumentImageGUID)
	}

	passDocImagesRef := []models.ImageReference{}

	findGroups, err := svc.repoImageGroup.FindWithoutGUIDByDocumentImageGUIDs(shopID, currentGroupGUID, docImageGUIDs)

	if err != nil {
		return []models.ImageReference{}, []string{}, err
	}

	for _, imageGroup := range findGroups {
		for _, imageRef := range *imageGroup.ImageReferences {
			foundImageRef, isFound := lo.Find[models.ImageReferenceBody](docImageRefs, func(tempImageRef models.ImageReferenceBody) bool {
				return imageRef.DocumentImageGUID == tempImageRef.DocumentImageGUID
			})

			if isFound && (imageGroup.References == nil || len(imageGroup.References) < 1) {
				imageRef.XOrder = foundImageRef.XOrder
				passDocImagesRef = append(passDocImagesRef, imageRef)
			} else if imageGroup.References != nil && len(imageGroup.References) > 0 {
				return []models.ImageReference{}, []string{}, fmt.Errorf("document image guid %s has referenced in %s", imageRef.DocumentImageGUID, imageGroup.GuidFixed)
			} else {
				return []models.ImageReference{}, []string{}, fmt.Errorf("document image guid \"%s\" has referenced in %s", imageRef.DocumentImageGUID, imageGroup.GuidFixed)
			}
		}
	}

	tempDocumentImageGroupGUIDs := []string{}
	for _, imageGroup := range findGroups {
		tempDocumentImageGroupGUIDs = append(tempDocumentImageGroupGUIDs, imageGroup.GuidFixed)
	}

	return passDocImagesRef, tempDocumentImageGroupGUIDs, nil
}

func (svc DocumentImageService) CreateDocumentImageGroup(shopID string, authUsername string, docImageGroup models.DocumentImageGroup) (string, error) {

	if docImageGroup.ImageReferences == nil || len(*docImageGroup.ImageReferences) < 1 {
		return "", errors.New("document image is size 0")
	}

	if docImageGroup.ImageReferences == nil || len(*docImageGroup.ImageReferences) > svc.maxImageReferences {
		return "", fmt.Errorf("document images is over size %d", svc.maxImageReferences)
	}

	docImageGroupData := models.DocumentImageGroupDoc{}

	tempImageRefs := lo.Map[models.ImageReference, models.ImageReferenceBody](
		*docImageGroup.ImageReferences,
		func(temp models.ImageReference, index int) models.ImageReferenceBody {
			return temp.ImageReferenceBody
		})

	passDocImagesRef, docImageGroupGUIDs, err := svc.getDocumentImageNotReferencedInGroup(shopID, "", tempImageRefs)
	if err != nil {
		return "", err
	}

	if len(passDocImagesRef) < 1 {
		return "", fmt.Errorf("document images invalid")
	}

	createdAt := svc.timeNowFnc()
	docImageGroupGUIDFixed := svc.newDocumentImageGroupGUIDFnc()

	docImageGroupData.ShopID = shopID
	docImageGroupData.DocumentImageGroup = docImageGroup
	docImageGroupData.GuidFixed = docImageGroupGUIDFixed
	docImageGroupData.Status = models.IMAGE_PENDING

	docImageGroupData.References = []models.Reference{}

	newXOrder, _ := svc.newXOrderDocumentImageGroup(shopID, docImageGroup.TaskGUID)
	docImageGroupData.XOrder = newXOrder

	docImageGroupData.CreatedBy = authUsername
	docImageGroupData.CreatedAt = createdAt

	sort.Slice(passDocImagesRef, func(i, j int) bool {
		return passDocImagesRef[i].XOrder < passDocImagesRef[j].XOrder
	})

	// if len(passDocImagesRef) > 0 {
	// 	docImageGroupData.UploadedBy = passDocImagesRef[0].UploadedBy
	// 	docImageGroupData.UploadedAt = passDocImagesRef[0].UploadedAt
	// }

	docImageGroupData.UploadedBy = authUsername
	docImageGroupData.UploadedAt = docImageGroup.UploadedAt

	docImageGroupData.ImageReferences = &passDocImagesRef

	tempGUIDDocumentImages := []string{}
	for _, imageRef := range passDocImagesRef {
		tempGUIDDocumentImages = append(tempGUIDDocumentImages, imageRef.DocumentImageGUID)

	}

	err = svc.clearCreateDocumentImageGroupByDocumentGUIDs(shopID, docImageGroupGUIDs, tempGUIDDocumentImages)
	if err != nil {
		return "", err
	}

	_, err = svc.repoImageGroup.Create(docImageGroupData)

	if err != nil {
		return "", err
	}

	return docImageGroupGUIDFixed, nil
}

func (svc DocumentImageService) UpdateStatusDocumentImageGroup(shopID string, authUsername string, groupGUID string, status int8) error {
	findDoc, err := svc.repoImageGroup.FindByGuid(shopID, groupGUID)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	if svc.isDocumentImageGroupHasReferenced(findDoc) {
		return errors.New("document has referenced")
	}

	if findDoc.Status == status {
		return nil
	}

	if status < models.IMAGE_PENDING || status > models.IMAGE_REJECT_KEYING {
		return errors.New("status out of range")
	}

	findDoc.Status = status
	svc.repoImageGroup.Update(shopID, groupGUID, findDoc)

	if status == models.IMAGE_REJECT_KEYING || status == models.IMAGE_REJECT {

		_, err = svc.messageQueueReCountRejectDocumentImageGroup(shopID, findDoc.TaskGUID)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if findDoc.Status == models.IMAGE_REJECT_KEYING || findDoc.Status == models.IMAGE_REJECT && status != models.IMAGE_REJECT_KEYING && status != models.IMAGE_REJECT {
		_, err = svc.messageQueueReCountRejectDocumentImageGroup(shopID, findDoc.TaskGUID)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func (svc DocumentImageService) UpdateDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImageGroup models.DocumentImageGroup) error {
	if docImageGroup.ImageReferences == nil || len(*docImageGroup.ImageReferences) > svc.maxImageReferences {
		return fmt.Errorf("document image is over size %d", svc.maxImageReferences)
	}

	findDoc, err := svc.repoImageGroup.FindByGuid(shopID, groupGUID)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	if svc.isDocumentImageGroupHasReferenced(findDoc) {
		return errors.New("document has referenced")
	}

	docImages := *docImageGroup.ImageReferences

	tempImageRefs := lo.Map[models.ImageReference, models.ImageReferenceBody](*docImageGroup.ImageReferences, func(temp models.ImageReference, index int) models.ImageReferenceBody {
		return temp.ImageReferenceBody
	})

	passDocImagesRef, docImageGroupGUIDs, err := svc.getDocumentImageNotReferencedInGroup(shopID, groupGUID, tempImageRefs)
	if err != nil {
		return err
	}

	updateDocImagesRef := map[string]models.ImageReference{}
	for _, imageRef := range passDocImagesRef {
		updateDocImagesRef[imageRef.DocumentImageGUID] = imageRef
	}

	tempRemoveDocImageFromGroup := []models.ImageReference{}
	// check exist from current document image group
	if findDoc.ImageReferences != nil {
		for _, imageRef := range *findDoc.ImageReferences {
			tempImageRef, isFound := lo.Find[models.ImageReference](docImages, func(tempImageRef models.ImageReference) bool {
				return imageRef.DocumentImageGUID == tempImageRef.DocumentImageGUID
			})

			if isFound {
				updateDocImagesRef[tempImageRef.DocumentImageGUID] = tempImageRef
			} else {
				tempRemoveDocImageFromGroup = append(tempRemoveDocImageFromGroup, imageRef)
			}
		}
	}

	tempDocImageGUIDs := []string{}
	tempDocImageRef := map[string]models.ImageReference{}
	for _, docImageRef := range updateDocImagesRef {
		tempDocImageGUIDs = append(tempDocImageGUIDs, docImageRef.DocumentImageGUID)
		tempDocImageRef[docImageRef.DocumentImageGUID] = docImageRef
	}

	findDocImages, err := svc.repoImage.FindInGUIDs(shopID, tempDocImageGUIDs)
	if err != nil {
		return err
	}

	docImgRefs := []models.ImageReference{}
	for _, docRefImage := range findDocImages {

		tempDocImgRef := models.ImageReference{}

		tempDocImgRef.DocumentImageGUID = docRefImage.GuidFixed
		tempDocImgRef.ImageURI = docRefImage.ImageURI
		tempDocImgRef.ImageEditURI = docRefImage.ImageEditURI
		tempDocImgRef.Name = docRefImage.Name
		tempDocImgRef.MetaFileAt = docRefImage.MetaFileAt
		tempDocImgRef.UploadedAt = docRefImage.UploadedAt
		tempDocImgRef.UploadedBy = docRefImage.UploadedBy

		if temp, ok := tempDocImageRef[docRefImage.GuidFixed]; ok {
			tempDocImgRef.XOrder = temp.XOrder
		}

		docImgRefs = append(docImgRefs, tempDocImgRef)
	}

	sort.Slice(docImgRefs, func(i, j int) bool {
		return docImgRefs[i].XOrder < docImgRefs[j].XOrder
	})

	if len(tempDocImageRef) > 0 {
		findDoc.UploadedBy = docImgRefs[0].UploadedBy
		findDoc.UploadedAt = docImgRefs[0].UploadedAt
	}

	tempDoc := findDoc
	timeAt := svc.timeNowFnc()

	tempStatus := findDoc.Status
	findDoc.DocumentImageGroup = docImageGroup
	findDoc.ImageReferences = &docImgRefs
	findDoc.References = tempDoc.References

	findDoc.UpdatedAt = timeAt
	findDoc.UpdatedBy = authUsername

	findDoc.Status = tempStatus

	if err = svc.repoImageGroup.Update(shopID, groupGUID, findDoc); err != nil {
		return err
	}

	if err = svc.clearUpdateDocumentImageGroupByDocumentGUIDs(shopID, groupGUID, docImageGroupGUIDs, tempDocImageGUIDs); err != nil {
		return err
	}

	for _, imageRef := range tempRemoveDocImageFromGroup {
		imageGroupGUID := svc.newDocumentImageGroupGUIDFnc()
		tempTags := []string{}
		if findDoc.Tags != nil {
			tempTags = *findDoc.Tags
		}

		docImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, imageRef, imageRef.ImageURI, tempTags, findDoc.TaskGUID, findDoc.PathTask, timeAt)

		newXOrderDocImgGroup, _ := svc.newXOrderDocumentImageGroup(shopID, findDoc.TaskGUID)

		docImageGroup.XOrder = newXOrderDocImgGroup

		_, err = svc.repoImageGroup.Create(docImageGroup)

		if err != nil {
			return err
		}
	}

	return nil
}

func (svc DocumentImageService) UpdateImageReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImages []models.ImageReferenceBody) error {
	findDoc, err := svc.repoImageGroup.FindByGuid(shopID, groupGUID)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	if svc.isDocumentImageGroupHasReferenced(findDoc) {
		return errors.New("document has referenced")
	}

	tempDocImageGUIDs := []string{}

	for _, imageRef := range docImages {
		tempDocImageGUIDs = append(tempDocImageGUIDs, imageRef.DocumentImageGUID)
	}

	findDocImages, err := svc.repoImage.FindInGUIDs(shopID, tempDocImageGUIDs)

	if err != nil {
		return err
	}

	maxSizeInvalid := 50
	if len(docImages) > (len(findDocImages) + maxSizeInvalid) {
		return errors.New("document image invalid")
	}

	passDocImagesRef, docImageGroupGUIDs, err := svc.getDocumentImageNotReferencedInGroup(shopID, groupGUID, docImages)
	if err != nil {
		return err
	}

	updateDocImagesRef := map[string]models.ImageReference{}
	for _, imageRef := range passDocImagesRef {
		updateDocImagesRef[imageRef.DocumentImageGUID] = imageRef
	}

	tempRemoveDocImageFromGroup := []models.ImageReference{}
	// check exist from current document image group
	if findDoc.ImageReferences != nil {
		for _, imageRef := range *findDoc.ImageReferences {
			tempImageRef, isFound := lo.Find[models.ImageReferenceBody](docImages, func(tempImageRef models.ImageReferenceBody) bool {
				return imageRef.DocumentImageGUID == tempImageRef.DocumentImageGUID
			})

			if isFound {
				imageRef.XOrder = tempImageRef.XOrder
				updateDocImagesRef[tempImageRef.DocumentImageGUID] = imageRef
			} else {
				tempRemoveDocImageFromGroup = append(tempRemoveDocImageFromGroup, imageRef)
			}
		}
	}

	tempGUIDDocumentImages := []string{}
	tempDocImageRef := []models.ImageReference{}
	for _, docImageRef := range updateDocImagesRef {
		tempGUIDDocumentImages = append(tempGUIDDocumentImages, docImageRef.DocumentImageGUID)
		tempDocImageRef = append(tempDocImageRef, docImageRef)

	}

	sort.Slice(tempDocImageRef, func(i, j int) bool {
		return tempDocImageRef[i].XOrder < tempDocImageRef[j].XOrder
	})

	if len(tempDocImageRef) > 0 {
		findDoc.UploadedBy = tempDocImageRef[0].UploadedBy
		findDoc.UploadedAt = tempDocImageRef[0].UploadedAt
	}

	timeAt := svc.timeNowFnc()

	findDoc.ImageReferences = &tempDocImageRef

	findDoc.UpdatedAt = timeAt
	findDoc.UpdatedBy = authUsername

	if err = svc.repoImageGroup.Update(shopID, groupGUID, findDoc); err != nil {
		return err
	}

	if err = svc.clearUpdateDocumentImageGroupByDocumentGUIDs(shopID, groupGUID, docImageGroupGUIDs, tempGUIDDocumentImages); err != nil {
		return err
	}

	for _, imageRef := range tempRemoveDocImageFromGroup {
		imageGroupGUID := svc.newDocumentImageGroupGUIDFnc()

		tempTags := []string{}
		if findDoc.Tags != nil {
			tempTags = *findDoc.Tags
		}

		docImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, imageRef, imageRef.ImageURI, tempTags, findDoc.TaskGUID, findDoc.PathTask, timeAt)
		newXOrderDocImgGroup, _ := svc.newXOrderDocumentImageGroup(shopID, findDoc.TaskGUID)

		docImageGroup.XOrder = newXOrderDocImgGroup

		_, err = svc.repoImageGroup.Create(docImageGroup)

		if err != nil {
			return err
		}
	}

	return nil
}

func (svc DocumentImageService) UpdateReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docRef models.Reference) error {
	findDoc, err := svc.repoImageGroup.FindByGuid(shopID, groupGUID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document image group not found")
	}

	if findDoc.References != nil {
		_, isExistsModule := lo.Find[models.Reference](findDoc.References, func(tempDoc models.Reference) bool {
			return tempDoc.Module == docRef.Module
		})

		if isExistsModule {
			return errors.New("document has referenced")
		}
	}

	// Clear references
	docImageClearRefs, err := svc.repoImage.FindByReference(shopID, docRef)
	if err != nil {
		return err
	}

	docImageGroupClearRefs, err := svc.repoImageGroup.FindByReference(shopID, docRef)
	if err != nil {
		return err
	}

	for _, docImage := range docImageClearRefs {
		tempDocRefs := lo.Filter[models.Reference](docImage.References, func(tempDocRef models.Reference, idx int) bool {
			return docRef.DocNo != tempDocRef.DocNo
		})

		docImage.References = tempDocRefs

		svc.repoImage.Update(shopID, docImage.GuidFixed, docImage)
	}

	for _, docImageGroup := range docImageGroupClearRefs {
		tempDocRefs := lo.Filter[models.Reference](docImageGroup.References, func(tempDocRef models.Reference, idx int) bool {
			return docRef.DocNo != tempDocRef.DocNo
		})

		docImageGroup.References = tempDocRefs

		svc.repoImageGroup.Update(shopID, docImageGroup.GuidFixed, docImageGroup)
	}

	tempDocImageGUIDs := []string{}

	for _, imageRef := range *findDoc.ImageReferences {
		tempDocImageGUIDs = append(tempDocImageGUIDs, imageRef.DocumentImageGUID)
	}

	findDocImages, err := svc.repoImage.FindInGUIDs(shopID, tempDocImageGUIDs)

	if err != nil {
		return err
	}

	findDoc.References = append(findDoc.References, docRef)

	timeAt := svc.timeNowFnc()

	findDoc.UpdatedAt = timeAt
	findDoc.UpdatedBy = authUsername

	if err = svc.repoImageGroup.Update(shopID, groupGUID, findDoc); err != nil {
		return err
	}

	for _, docImage := range findDocImages {

		docImage.UpdatedAt = timeAt
		docImage.UpdatedBy = authUsername

		docImage.References = append(docImage.References, docRef)

		if err = svc.repoImage.Update(shopID, docImage.GuidFixed, docImage); err != nil {
			return err
		}
	}

	return nil
}

func (svc DocumentImageService) UpdateTagsInDocumentImageGroup(shopID string, authUsername string, groupGUID string, tags []string) error {

	findDoc, err := svc.repoImageGroup.FindByGuid(shopID, groupGUID)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	if svc.isDocumentImageGroupHasReferenced(findDoc) {
		return errors.New("document has referenced")
	}

	findDoc.Tags = &tags

	findDoc.UpdatedAt = svc.timeNowFnc()
	findDoc.UpdatedBy = authUsername

	if err = svc.repoImageGroup.Update(shopID, groupGUID, findDoc); err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) DeleteReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docRef models.Reference) error {
	findDoc, err := svc.repoImageGroup.FindByGuid(shopID, groupGUID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document image group not found")
	}

	// Clear references
	docImageClearRefs, err := svc.repoImage.FindByReference(shopID, docRef)
	if err != nil {
		return err
	}

	docImageGroupClearRefs, err := svc.repoImageGroup.FindByReference(shopID, docRef)
	if err != nil {
		return err
	}

	for _, docImage := range docImageClearRefs {
		tempDocRefs := lo.Filter[models.Reference](docImage.References, func(tempDocRef models.Reference, idx int) bool {
			return docRef.Module != tempDocRef.Module && docRef.DocNo != tempDocRef.DocNo
		})

		docImage.References = tempDocRefs

		svc.repoImage.Update(shopID, docImage.GuidFixed, docImage)
	}

	for _, docImageGroup := range docImageGroupClearRefs {
		tempDocRefs := lo.Filter[models.Reference](docImageGroup.References, func(tempDocRef models.Reference, idx int) bool {
			return docRef.Module != tempDocRef.Module && docRef.DocNo != tempDocRef.DocNo
		})

		docImageGroup.References = tempDocRefs

		svc.repoImageGroup.Update(shopID, docImageGroup.GuidFixed, docImageGroup)
	}

	return nil
}

func (svc DocumentImageService) DeleteDocumentImageGroupByGuid(shopID string, authUsername string, documentImageGroupGuidFixed string) error {

	findDocGroup, err := svc.repoImageGroup.FindByGuid(shopID, documentImageGroupGuidFixed)

	if err != nil {
		return err
	}

	if len(findDocGroup.GuidFixed) < 1 {
		return nil
	}

	if svc.isDocumentImageGroupHasReferenced(findDocGroup) {
		return errors.New("document has referenced")
	}

	err = svc.repoImageGroup.Transaction(func() error {

		for _, docImage := range *findDocGroup.ImageReferences {
			err = svc.repoImage.DeleteByGuidfixed(shopID, docImage.DocumentImageGUID, authUsername)

			if err != nil {
				return err
			}
		}

		if err = svc.repoImageGroup.DeleteByGuidfixed(shopID, findDocGroup.GuidFixed); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) DeleteDocumentImageGroupByGuids(shopID string, authUsername string, documentImageGroupGuidFixeds []string) error {

	err := svc.repoImageGroup.Transaction(func() error {
		for _, DocumentImageGroupGuidFixed := range documentImageGroupGuidFixeds {
			findDocGroup, err := svc.repoImageGroup.FindByGuid(shopID, DocumentImageGroupGuidFixed)

			if err != nil {
				return err
			}

			if svc.isDocumentImageGroupHasReferenced(findDocGroup) {
				return errors.New("document has referenced")
			}

			for _, docImage := range *findDocGroup.ImageReferences {
				err = svc.repoImage.DeleteByGuidfixed(shopID, docImage.DocumentImageGUID, authUsername)

				if err != nil {
					return err
				}
			}

			if err = svc.repoImageGroup.DeleteByGuidfixed(shopID, findDocGroup.GuidFixed); err != nil {
				return err
			}

			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) UnGroupDocumentImageGroup(shopID string, authUsername string, groupGUID string) error {
	findDocGroup, err := svc.repoImageGroup.FindByGuid(shopID, groupGUID)

	if err != nil {
		return err
	}

	if svc.isDocumentImageGroupHasReferenced(findDocGroup) {
		return errors.New("document has referenced")
	}

	updatedAt := svc.timeNowFnc()
	for _, imageRef := range *findDocGroup.ImageReferences {
		imageGroupGUID := svc.newDocumentImageGroupGUIDFnc()

		tags := []string{}
		if findDocGroup.Tags != nil {
			tags = *findDocGroup.Tags
		}

		docImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, imageRef, imageRef.ImageURI, tags, findDocGroup.TaskGUID, findDocGroup.PathTask, updatedAt)

		newXOrderDocImgGroup, _ := svc.newXOrderDocumentImageGroup(shopID, docImageGroup.TaskGUID)

		docImageGroup.XOrder = newXOrderDocImgGroup
		_, err = svc.repoImageGroup.Create(docImageGroup)

		if err != nil {
			return err
		}
	}

	svc.repoImageGroup.DeleteByGuidfixed(shopID, groupGUID)

	return nil
}

func (svc DocumentImageService) ListDocumentImageGroup(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentImageGroupInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{"title"}
	docList, pagination, err := svc.repoImageGroup.FindPageImageGroup(shopID, filters, searchInFields, pageable)

	return docList, pagination, err
}

func (svc DocumentImageService) GetDocumentImageDocRefGroup(shopID string, docImageGroupGUID string) (models.DocumentImageGroupInfo, error) {
	doc, err := svc.repoImageGroup.FindByGuid(shopID, docImageGroupGUID)

	if err != nil {
		return models.DocumentImageGroupInfo{}, err
	}

	if doc.ID.IsZero() {
		return models.DocumentImageGroupInfo{}, errors.New("document not found")
	}

	return doc.DocumentImageGroupInfo, nil
}

func (svc DocumentImageService) GetDocumentImageGroupByDocRef(shopID string, docRef string) (models.DocumentImageGroupInfo, error) {

	findDoc, err := svc.repoImageGroup.FindOne(shopID, bson.M{"references.docno": docRef})

	if err != nil {
		return models.DocumentImageGroupInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DocumentImageGroupInfo{}, errors.New("document not found")
	}

	return findDoc.DocumentImageGroupInfo, nil

}

func (svc DocumentImageService) XSortsUpdate(shopID string, authUsername string, taskGUID string, xsorts []models.XSortDocumentImageGroupRequest) error {
	for _, xsort := range xsorts {
		if len(xsort.GUIDFixed) < 1 {
			continue
		}

		err := svc.repoImageGroup.UpdateXOrder(shopID, taskGUID, xsort.GUIDFixed, xsort.XOrder)

		if err != nil {
			return err
		}
	}

	return nil
}

func (svc DocumentImageService) isDocumentImageGroupHasReferenced(doc models.DocumentImageGroupDoc) bool {
	return doc.References != nil && len(doc.References) > 0
}

func (svc DocumentImageService) isDocumentImageHasReferenced(doc models.DocumentImageDoc) bool {
	return doc.References != nil && len(doc.References) > 0
}

func (svc DocumentImageService) documentImageToImageReference(documentImageGUID string, documentImage models.DocumentImage, authUsername string, createdAt time.Time) models.ImageReference {
	return models.ImageReference{
		ImageReferenceBody: models.ImageReferenceBody{
			XOrder:            1,
			DocumentImageGUID: documentImageGUID,
		},
		ImageURI:     documentImage.ImageURI,
		ImageEditURI: documentImage.ImageEditURI,
		Name:         documentImage.Name,
		UploadedBy:   authUsername,
		UploadedAt:   createdAt,
		MetaFileAt:   documentImage.MetaFileAt,
	}
}

func (svc DocumentImageService) createImageGroupByDocumentImage(shopID string, authUsername string, imageGroupGUID string, documentImageRef models.ImageReference, imageURI string, tags []string, fileFolderGUID string, pathTask string, createdAt time.Time) models.DocumentImageGroupDoc {
	docDataImageGroup := models.DocumentImageGroupDoc{}
	docDataImageGroup.ShopID = shopID
	docDataImageGroup.GuidFixed = imageGroupGUID
	docDataImageGroup.Title = documentImageRef.Name
	docDataImageGroup.References = []models.Reference{}
	docDataImageGroup.Tags = &tags
	docDataImageGroup.ImageReferences = &[]models.ImageReference{
		documentImageRef,
	}

	docDataImageGroup.TaskGUID = fileFolderGUID
	docDataImageGroup.PathTask = pathTask
	docDataImageGroup.Status = models.IMAGE_PENDING

	docDataImageGroup.CreatedBy = authUsername
	docDataImageGroup.CreatedAt = createdAt

	docDataImageGroup.UploadedBy = documentImageRef.UploadedBy
	docDataImageGroup.UploadedAt = documentImageRef.UploadedAt

	return docDataImageGroup
}

func (svc DocumentImageService) clearCreateDocumentImageGroupByDocumentGUIDs(shopID string, docImageGroupGUIDs []string, docImageGUIDs []string) error {

	err := svc.repoImageGroup.RemoveDocumentImageByDocumentImageGUIDs(shopID, docImageGUIDs)
	if err != nil {
		return err
	}

	err = svc.repoImageGroup.DeleteByGUIDsIsDocumentImageEmpty(shopID, docImageGroupGUIDs)
	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) clearUpdateDocumentImageGroupByDocumentGUIDs(shopID string, docGroupGUID string, clearDocImageGUIDs []string, docImageGUIDs []string) error {

	err := svc.repoImageGroup.RemoveDocumentImageByDocumentImageGUIDsWithoutDocumentImageGroupGUID(shopID, docGroupGUID, docImageGUIDs)
	if err != nil {
		return err
	}

	err = svc.repoImageGroup.DeleteByGUIDsIsDocumentImageEmptyWithoutDocumentImageGroupGUID(shopID, docGroupGUID, clearDocImageGUIDs)
	if err != nil {
		return err
	}

	err = svc.repoImageGroup.DeleteByGUIDIsDocumentImageEmpty(shopID, docGroupGUID)
	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) messageQueueReCountDocumentImageGroup(shopID string, taskGUID string) (int, error) {

	count, err := svc.repoImageGroup.CountByTask(shopID, taskGUID)

	if err != nil {
		return 0, err
	}

	taskMsg := models.DocumentImageTaskChangeMessage{
		ShopID:   shopID,
		TaskGUID: taskGUID,
		// Event:    models.TaskChangePlus,
		Count: count,
	}

	err = svc.repoMessagequeue.TaskChange(taskMsg)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (svc DocumentImageService) messageQueueReCountRejectDocumentImageGroup(shopID string, taskGUID string) (int, error) {

	count, err := svc.repoImageGroup.CountRejectByTask(shopID, taskGUID)

	if err != nil {
		return 0, err
	}

	taskMsg := models.DocumentImageTaskRejectMessage{
		ShopID:   shopID,
		TaskGUID: taskGUID,
		// Event:    models.TaskChangePlus,
		Count: count,
	}

	err = svc.repoMessagequeue.TaskReject(taskMsg)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (svc DocumentImageService) newXOrderDocumentImageGroup(shopID string, taskGUID string) (int, error) {

	findDoc, err := svc.repoImageGroup.FindLastOneByTask(shopID, taskGUID)

	if err != nil {
		return 0, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return 0, nil
	}

	return findDoc.XOrder + 1, nil
}
