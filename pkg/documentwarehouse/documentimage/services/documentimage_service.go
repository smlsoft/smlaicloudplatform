package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/repositories"
	common "smlcloudplatform/pkg/models"
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
	UpdateDocumentImage(shopID string, guid string, authUsername string, doc models.DocumentImage) error
	UpdateDocumentImageReject(shopID string, guid string, authUsername string, isReject bool) error
	DeleteDocumentImage(shopID string, authUsername string, imageGUID string) error

	InfoDocumentImage(shopID string, guid string) (models.DocumentImageInfo, error)
	SearchDocumentImage(shopID string, matchFilters map[string]interface{}, q string, page int, limit int, sorts map[string]int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)
	UploadDocumentImage(shopID string, authUsername string, fh *multipart.FileHeader) (*models.DocumentImageInfo, error)

	// SaveDocumentImageDocRefGroup(shopID string, authUsername string, docRef string, docImages []models.DocumentImageGroup) error
	// GetDocumentImageDocRefGroup(shopID string, docRef string) (models.DocumentImageGroup, error)
	// ListDocumentImageDocRefGroup(shopID string, filters map[string]interface{}, q string, page int, limit int) ([]models.DocumentImageGroup, mongopagination.PaginationData, error)

	CreateDocumentImageGroup(shopID string, authUsername string, docImageGroup models.DocumentImageGroup) (string, error)
	GetDocumentImageDocRefGroup(shopID string, docImageGroupGUID string) (models.DocumentImageGroupInfo, error)
	GetDocumentImageGroupByDocRef(shopID string, docRef string) (models.DocumentImageGroupInfo, error)
	UpdateDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImageGroup models.DocumentImageGroup) error
	UpdateImageReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImages []models.ImageReferenceBody) error
	UpdateReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docRef models.Reference) error
	UnGroupDocumentImageGroup(shopID string, authUsername string, groupGUID string) error
	ListDocumentImageGroup(shopID string, filters map[string]interface{}, pageable common.Pageable) ([]models.DocumentImageGroupInfo, mongopagination.PaginationData, error)
	DeleteReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docRef models.Reference) error

	UpdateDocumentImageRederenceGroup() error
}

type DocumentImageService struct {
	repoImageGroup               repositories.DocumentImageGroupRepository
	repoImage                    repositories.IDocumentImageRepository
	FilePersister                microservice.IPersisterFile
	maxImageReferences           int
	timeNowFnc                   func() time.Time
	newDocumentImageGUIDFnc      func() string
	newDocumentImageGroupGUIDFnc func() string
}

func NewDocumentImageService(repo repositories.IDocumentImageRepository, repoImageGroup repositories.DocumentImageGroupRepository, filePersister microservice.IPersisterFile) DocumentImageService {
	return DocumentImageService{
		maxImageReferences: 100,
		repoImageGroup:     repoImageGroup,
		repoImage:          repo,
		FilePersister:      filePersister,
		timeNowFnc: func() time.Time {
			return time.Now()
		},
		newDocumentImageGUIDFnc:      utils.NewGUID,
		newDocumentImageGroupGUIDFnc: utils.NewGUID,
	}
}

func (svc DocumentImageService) CreateDocumentImage(shopID string, authUsername string, docRequest models.DocumentImageRequest) (string, string, error) {

	// do upload first

	createdAt := svc.timeNowFnc()

	documentImageGUID := svc.newDocumentImageGUIDFnc()

	docData := models.DocumentImageDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = documentImageGUID
	docData.DocumentImage = docRequest.DocumentImage

	docData.IsReject = false
	docData.References = []models.Reference{}
	docData.MetaFileAt = docRequest.MetaFileAt

	docData.CreatedBy = authUsername
	docData.CreatedAt = createdAt

	// docData.UpdatedBy = authUsername
	// docData.UpdatedAt = createdAt

	docData.UploadedBy = authUsername
	// docData.UploadedAt = createdAt

	// image group
	imageGroupGUID := svc.newDocumentImageGroupGUIDFnc()
	docImageRef := svc.documentImageToImageReference(documentImageGUID, docRequest.DocumentImage, authUsername, createdAt)
	docDataImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, docImageRef, docRequest.ImageURI, *docRequest.Tags, docRequest.FileFolderGUID, docRequest.PathFileFolder, createdAt)

	svc.repoImageGroup.Transaction(func() error {

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

	return documentImageGUID, imageGroupGUID, nil
}

func (svc DocumentImageService) BulkCreateDocumentImage(shopID string, authUsername string, docs []models.DocumentImageRequest) error {

	// do upload first

	createdAt := svc.timeNowFnc()
	docDataList := []models.DocumentImageDoc{}
	docDataImageGroupList := []models.DocumentImageGroupDoc{}

	for _, doc := range docs {
		documentImageGUID := svc.newDocumentImageGUIDFnc()

		docData := models.DocumentImageDoc{}
		docData.ShopID = shopID
		docData.GuidFixed = documentImageGUID
		docData.DocumentImage = doc.DocumentImage

		docData.IsReject = false
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
		docDataImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, docImageRef, doc.ImageURI, *doc.Tags, doc.FileFolderGUID, doc.PathFileFolder, createdAt)

		docDataList = append(docDataList, docData)
		docDataImageGroupList = append(docDataImageGroupList, docDataImageGroup)
	}

	svc.repoImageGroup.Transaction(func() error {

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

	return nil
}

func (svc DocumentImageService) UpdateDocumentImage(shopID string, guid string, authUsername string, doc models.DocumentImage) error {

	findDoc, err := svc.repoImage.FindByGuid(shopID, guid)

	if svc.isDocumentImageHasReferenced(findDoc) {
		return errors.New("document has referenced")
	}

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	updatedAt := svc.timeNowFnc()

	findDoc.DocumentImage = doc

	findDoc.References = []models.Reference{}

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = updatedAt

	findDoc.UploadedBy = authUsername
	findDoc.UploadedAt = doc.UploadedAt

	err = svc.repoImage.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	findGroupDoc, err := svc.repoImageGroup.FindOneByDocumentImageGUID(shopID, guid)

	if err != nil {
		return err
	}

	groupIsReject := false
	tempImageRefs := *findGroupDoc.ImageReferences
	for idx, tempDoc := range tempImageRefs {

		if tempDoc.DocumentImageGUID == guid {
			tempDocImage := &tempImageRefs[idx]

			tempDocImage.ImageURI = doc.ImageURI
			tempDocImage.Name = doc.Name
			tempDocImage.UploadedBy = authUsername
			tempDocImage.UploadedAt = updatedAt
			tempDocImage.MetaFileAt = doc.MetaFileAt
			tempDocImage.IsReject = doc.IsReject

			if doc.IsReject {
				groupIsReject = doc.IsReject
			}

			if groupIsReject {
				break
			}
		} else if tempDoc.IsReject {
			groupIsReject = true
		}
	}
	findGroupDoc.IsReject = groupIsReject
	// findGroupDoc.ImageReferences = &tempImageRefs
	err = svc.repoImageGroup.Update(shopID, findGroupDoc.GuidFixed, findGroupDoc)
	if err != nil {
		return err
	}
	return nil
}

func (svc DocumentImageService) UpdateDocumentImageReject(shopID string, guid string, authUsername string, isReject bool) error {

	findDoc, err := svc.repoImage.FindByGuid(shopID, guid)

	if svc.isDocumentImageHasReferenced(findDoc) {
		return errors.New("document has referenced")
	}

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	updatedAt := svc.timeNowFnc()

	findDoc.IsReject = isReject
	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = updatedAt

	err = svc.repoImage.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	findGroupDoc, err := svc.repoImageGroup.FindOneByDocumentImageGUID(shopID, guid)

	if err != nil {
		return err
	}

	groupIsReject := false
	tempImageRefs := *findGroupDoc.ImageReferences
	for idx, tempDoc := range tempImageRefs {

		if tempDoc.DocumentImageGUID == guid {
			if isReject {
				groupIsReject = isReject
			}

			tempDocGroup := &tempImageRefs[idx]
			tempDocGroup.UploadedBy = authUsername
			tempDocGroup.UploadedAt = updatedAt

			tempDocGroup.IsReject = isReject

			if groupIsReject {
				break
			}
		} else if tempDoc.IsReject {
			groupIsReject = tempDoc.IsReject
		}
	}

	findGroupDoc.IsReject = groupIsReject

	err = svc.repoImageGroup.Update(shopID, findGroupDoc.GuidFixed, findGroupDoc)
	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) DeleteDocumentImage(shopID string, authUsername string, imageGUID string) error {
	findDoc, err := svc.repoImage.FindByGuid(shopID, imageGUID)

	if svc.isDocumentImageHasReferenced(findDoc) {
		return errors.New("document has referenced")
	}

	if err != nil {
		return err
	}

	findDocGroup, err := svc.repoImageGroup.FindByGuid(shopID, imageGUID)

	if err != nil {
		return err
	}

	if svc.isDocumentImageGroupHasReferenced(findDocGroup) {
		return errors.New("document has referenced")
	}

	groupIsReject := false
	tempImageRef := []models.ImageReference{}

	if findDocGroup.ImageReferences != nil {
		for _, imageRef := range *findDocGroup.ImageReferences {
			if imageRef.DocumentImageGUID != imageGUID {
				tempImageRef = append(tempImageRef, imageRef)
				if imageRef.IsReject {
					groupIsReject = true
				}
			}
		}
	}

	findDocGroup.IsReject = groupIsReject
	findDocGroup.ImageReferences = &tempImageRef

	err = svc.repoImage.DeleteByGuidfixed(shopID, imageGUID, authUsername)

	if err != nil {
		return err
	}

	if len(tempImageRef) > 0 {
		if err = svc.repoImageGroup.Update(shopID, findDocGroup.GuidFixed, findDocGroup); err != nil {
			return err
		}
	} else {
		if err = svc.repoImageGroup.DeleteByGuidfixed(shopID, findDocGroup.GuidFixed); err != nil {
			return err
		}
	}

	return nil
}

func (svc DocumentImageService) RejectDocumentImage(shopID string, guid string, authUsername string, rejectStatus bool) error {

	findDoc, err := svc.repoImage.FindByGuid(shopID, guid)

	if svc.isDocumentImageHasReferenced(findDoc) {
		return errors.New("document has referenced")
	}

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.IsReject = rejectStatus

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repoImage.Update(shopID, guid, findDoc)

	if err != nil {
		return err
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

func (svc DocumentImageService) SearchDocumentImage(shopID string, matchFilters map[string]interface{}, q string, page int, limit int, sorts map[string]int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repoImage.FindPageFilterSort(shopID, matchFilters, []string{"guidfixed", "documentref", "module"}, q, page, limit, sorts)

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
		DocumentImage:  doc.DocumentImage,
		Tags:           &[]string{},
		FileFolderGUID: "",
		PathFileFolder: "",
	}
	_, _, err = svc.CreateDocumentImage(shopID, authUsername, docRequest)

	if err != nil {
		return nil, err
	}

	return &doc.DocumentImageInfo, err
}

func (svc DocumentImageService) UpdateDocumentImageRederenceGroup() error {
	// findGroupDoc, err := svc.repoImageGroup.FindOneByDocumentImageGUIDAll("2GcTSCfk0JCzxclxidLyCFUx8wo")

	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("%v", findGroupDoc.GuidFixed)

	findDocList, err := svc.repoImage.FindAll()

	if err != nil {
		return err
	}

	// for _, findDoc := range findDocList {
	// 	if findDoc.GuidFixed == "2GcTSCfk0JCzxclxidLyCFUx8wo" {
	// 		fmt.Printf("%v", findDoc)
	// 	}

	// }

	for _, findDoc := range findDocList {
		findGroupDoc, err := svc.repoImageGroup.FindOneByDocumentImageGUIDAll(findDoc.GuidFixed)

		if err != nil {
			return err
		}

		if findGroupDoc.ID != primitive.NilObjectID {
			fmt.Printf("image:: %v\n", findDoc.GuidFixed)
			fmt.Printf("group:: %v\n", findGroupDoc.GuidFixed)
		}

		refGroups := []models.ReferenceGroup{}
		if findGroupDoc.ID == primitive.NilObjectID {
			refGroup := models.ReferenceGroup{}
			refGroup.GroupType = ""
			refGroup.ParentGUID = ""
			refGroup.XOder = 1
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

	tempImageRefs := lo.Map[models.ImageReference, models.ImageReferenceBody](*docImageGroup.ImageReferences, func(temp models.ImageReference, index int) models.ImageReferenceBody {
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

	docImageGroupData.References = []models.Reference{}

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

	groupIsReject := false
	tempGUIDDocumentImages := []string{}
	for _, imageRef := range passDocImagesRef {
		tempGUIDDocumentImages = append(tempGUIDDocumentImages, imageRef.DocumentImageGUID)
		if imageRef.IsReject {
			groupIsReject = true
		}
	}

	docImageGroup.IsReject = groupIsReject

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

func (svc DocumentImageService) UpdateDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImageGroup models.DocumentImageGroup) error {
	if docImageGroup.ImageReferences == nil || len(*docImageGroup.ImageReferences) > svc.maxImageReferences {
		return fmt.Errorf("document image is over size %d", svc.maxImageReferences)
	}

	findDoc, err := svc.repoImageGroup.FindByGuid(shopID, groupGUID)

	if err != nil {
		return err
	}

	if findDoc.ID.IsZero() {
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

	tempGUIDDocumentImages := []string{}
	groupIsReject := false
	tempDocImageRef := []models.ImageReference{}

	for _, docImageRef := range updateDocImagesRef {
		tempGUIDDocumentImages = append(tempGUIDDocumentImages, docImageRef.DocumentImageGUID)
		tempDocImageRef = append(tempDocImageRef, docImageRef)
		if docImageRef.IsReject {
			groupIsReject = true
		}
	}

	sort.Slice(tempDocImageRef, func(i, j int) bool {
		return tempDocImageRef[i].XOrder < tempDocImageRef[j].XOrder
	})

	if len(tempDocImageRef) > 0 {
		findDoc.UploadedBy = tempDocImageRef[0].UploadedBy
		findDoc.UploadedAt = tempDocImageRef[0].UploadedAt
	}

	timeAt := svc.timeNowFnc()

	findDoc.DocumentImageGroup = docImageGroup
	findDoc.ImageReferences = &tempDocImageRef
	findDoc.UpdatedAt = timeAt
	findDoc.UpdatedBy = authUsername

	findDoc.IsReject = groupIsReject

	if err = svc.repoImageGroup.Update(shopID, groupGUID, findDoc); err != nil {
		return err
	}

	if err = svc.clearUpdateDocumentImageGroupByDocumentGUIDs(shopID, groupGUID, docImageGroupGUIDs, tempGUIDDocumentImages); err != nil {
		return err
	}

	for _, imageRef := range tempRemoveDocImageFromGroup {
		imageGroupGUID := svc.newDocumentImageGroupGUIDFnc()
		docImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, imageRef, imageRef.ImageURI, *findDoc.Tags, findDoc.FileFolderGUID, findDoc.PathFileFolder, timeAt)
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

	groupIsReject := false
	tempGUIDDocumentImages := []string{}
	tempDocImageRef := []models.ImageReference{}
	for _, docImageRef := range updateDocImagesRef {
		tempGUIDDocumentImages = append(tempGUIDDocumentImages, docImageRef.DocumentImageGUID)
		tempDocImageRef = append(tempDocImageRef, docImageRef)
		if docImageRef.IsReject {
			groupIsReject = true
		}
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
	findDoc.IsReject = groupIsReject

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
		docImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, imageRef, imageRef.ImageURI, *findDoc.Tags, findDoc.FileFolderGUID, findDoc.PathFileFolder, timeAt)
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
		docImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, imageRef, imageRef.ImageURI, *findDocGroup.Tags, findDocGroup.FileFolderGUID, findDocGroup.PathFileFolder, updatedAt)
		_, err = svc.repoImageGroup.Create(docImageGroup)

		if err != nil {
			return err
		}
	}

	svc.repoImageGroup.DeleteByGuidfixed(shopID, groupGUID)

	return nil
}

func (svc DocumentImageService) ListDocumentImageGroup(shopID string, filters map[string]interface{}, pageable common.Pageable) ([]models.DocumentImageGroupInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repoImageGroup.FindPageImageGroup(shopID, "", filters, []string{"title"}, pageable.Q, pageable.Page, pageable.Limit, pageable.Sorts)

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
		ImageURI:   documentImage.ImageURI,
		Name:       documentImage.Name,
		UploadedBy: authUsername,
		UploadedAt: createdAt,
		MetaFileAt: documentImage.MetaFileAt,
	}
}

func (svc DocumentImageService) createImageGroupByDocumentImage(shopID string, authUsername string, imageGroupGUID string, documentImageRef models.ImageReference, imageURI string, tags []string, fileFolderGUID string, pathFileFolder string, createdAt time.Time) models.DocumentImageGroupDoc {
	docDataImageGroup := models.DocumentImageGroupDoc{}
	docDataImageGroup.ShopID = shopID
	docDataImageGroup.GuidFixed = imageGroupGUID
	docDataImageGroup.Title = documentImageRef.Name
	docDataImageGroup.References = []models.Reference{}
	docDataImageGroup.Tags = &tags
	docDataImageGroup.ImageReferences = &[]models.ImageReference{
		documentImageRef,
	}

	docDataImageGroup.FileFolderGUID = fileFolderGUID
	docDataImageGroup.PathFileFolder = pathFileFolder
	docDataImageGroup.IsReject = documentImageRef.IsReject

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

	return nil
}
