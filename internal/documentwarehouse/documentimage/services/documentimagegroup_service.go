package services

import (
	"context"
	"errors"
	"fmt"
	"smlcloudplatform/internal/documentwarehouse/documentimage/models"
	micromodels "smlcloudplatform/pkg/microservice/models"
	"sort"
	"time"

	"github.com/samber/lo"
	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group

func (svc DocumentImageService) getDocumentImageNotReferencedInGroup(shopID string, currentGroupGUID string, docImageRefs []models.ImageReferenceBody) ([]models.ImageReference, []string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docImageGUIDs := []string{}

	for _, imageRef := range docImageRefs {
		docImageGUIDs = append(docImageGUIDs, imageRef.DocumentImageGUID)
	}

	passDocImagesRef := []models.ImageReference{}

	findGroups, err := svc.repoImageGroup.FindWithoutGUIDByDocumentImageGUIDs(ctx, shopID, currentGroupGUID, docImageGUIDs)

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

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

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

	newXOrder, _ := svc.newXOrderDocumentImageGroup(ctx, shopID, docImageGroup.TaskGUID)
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

	_, err = svc.repoImageGroup.Create(ctx, docImageGroupData)

	if err != nil {
		return "", err
	}

	return docImageGroupGUIDFixed, nil
}

func (svc DocumentImageService) UpdateStatusDocumentImageGroup(shopID string, authUsername string, groupGUID string, status int8) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repoImageGroup.FindByGuid(ctx, shopID, groupGUID)

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

	if status < models.IMAGE_PENDING || status > models.IMAGE_GL_COMPLETED {
		return errors.New("status out of range")
	}

	updateDoc := findDoc

	lastStatusHistory := models.StatusHistory{
		Status:    findDoc.Status,
		ChangedBy: authUsername,
		ChangedAt: svc.timeNowFnc(),
	}

	updateDoc.StatusHistories = append(updateDoc.StatusHistories, lastStatusHistory)
	updateDoc.Status = status
	svc.repoImageGroup.Update(ctx, shopID, groupGUID, updateDoc)

	_, err = svc.messageQueueReCountDocumentImageGroup(ctx, shopID, findDoc.TaskGUID)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (svc DocumentImageService) UpdateStatusDocumentImageGroupByTask(shopID string, authUsername string, taskGUID string, status int8) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	if status != models.IMAGE_PENDING && status != models.IMAGE_CHECKED {
		return errors.New("status out of range")
	}

	err := svc.repoImageGroup.UpdateStatusByTask(ctx, shopID, taskGUID, status)

	if err != nil {
		return err
	}

	_, err = svc.messageQueueReCountDocumentImageGroup(ctx, shopID, taskGUID)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (svc DocumentImageService) ReCountStatusDocumentImageGroupByTask(shopID string, authUsername string, taskGUID string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	_, err := svc.messageQueueReCountDocumentImageGroup(ctx, shopID, taskGUID)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (svc DocumentImageService) UpdateDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImageGroup models.DocumentImageGroup) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	if docImageGroup.ImageReferences == nil || len(*docImageGroup.ImageReferences) > svc.maxImageReferences {
		return fmt.Errorf("document image is over size %d", svc.maxImageReferences)
	}

	findDoc, err := svc.repoImageGroup.FindByGuid(ctx, shopID, groupGUID)

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

	findDocImages, err := svc.repoImage.FindInGUIDs(ctx, shopID, tempDocImageGUIDs)
	if err != nil {
		return err
	}

	docImgRefs := []models.ImageReference{}
	for _, docRefImage := range findDocImages {

		tempDocImgRef := models.ImageReference{}

		tempDocImgRef.DocumentImageGUID = docRefImage.GuidFixed
		tempDocImgRef.ImageURI = docRefImage.ImageURI
		tempDocImgRef.CloneImageFrom = docRefImage.GuidFixed
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

	updateDoc := findDoc
	timeAt := svc.timeNowFnc()

	updateDoc.DocumentImageGroup = docImageGroup
	updateDoc.ImageReferences = &docImgRefs
	updateDoc.References = findDoc.References

	updateDoc.UpdatedAt = timeAt
	updateDoc.UpdatedBy = authUsername

	updateDoc.Status = findDoc.Status
	updateDoc.StatusChangedBy = findDoc.StatusChangedBy
	updateDoc.StatusChangedAt = findDoc.StatusChangedAt
	updateDoc.StatusHistories = findDoc.StatusHistories

	if err = svc.repoImageGroup.Update(ctx, shopID, groupGUID, updateDoc); err != nil {
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

		newXOrderDocImgGroup, _ := svc.newXOrderDocumentImageGroup(ctx, shopID, findDoc.TaskGUID)

		docImageGroup.XOrder = newXOrderDocImgGroup

		_, err = svc.repoImageGroup.Create(ctx, docImageGroup)

		if err != nil {
			return err
		}
	}

	return nil
}

func (svc DocumentImageService) UpdateImageReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docImages []models.ImageReferenceBody) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repoImageGroup.FindByGuid(ctx, shopID, groupGUID)

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

	findDocImages, err := svc.repoImage.FindInGUIDs(ctx, shopID, tempDocImageGUIDs)

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

	if err = svc.repoImageGroup.Update(ctx, shopID, groupGUID, findDoc); err != nil {
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
		newXOrderDocImgGroup, _ := svc.newXOrderDocumentImageGroup(ctx, shopID, findDoc.TaskGUID)

		docImageGroup.XOrder = newXOrderDocImgGroup

		_, err = svc.repoImageGroup.Create(ctx, docImageGroup)

		if err != nil {
			return err
		}
	}

	return nil
}

func (svc DocumentImageService) UpdateReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docRef models.Reference) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repoImageGroup.FindByGuid(ctx, shopID, groupGUID)

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
	docImageClearRefs, err := svc.repoImage.FindByReference(ctx, shopID, docRef)
	if err != nil {
		return err
	}

	docImageGroupClearRefs, err := svc.repoImageGroup.FindByReference(ctx, shopID, docRef)
	if err != nil {
		return err
	}

	for _, docImage := range docImageClearRefs {
		tempDocRefs := lo.Filter[models.Reference](docImage.References, func(tempDocRef models.Reference, idx int) bool {
			return docRef.DocNo != tempDocRef.DocNo
		})

		docImage.References = tempDocRefs

		svc.repoImage.Update(ctx, shopID, docImage.GuidFixed, docImage)
	}

	for _, docImageGroup := range docImageGroupClearRefs {
		tempDocRefs := lo.Filter[models.Reference](docImageGroup.References, func(tempDocRef models.Reference, idx int) bool {
			return docRef.DocNo != tempDocRef.DocNo
		})

		docImageGroup.References = tempDocRefs

		svc.repoImageGroup.Update(ctx, shopID, docImageGroup.GuidFixed, docImageGroup)
	}

	tempDocImageGUIDs := []string{}

	for _, imageRef := range *findDoc.ImageReferences {
		tempDocImageGUIDs = append(tempDocImageGUIDs, imageRef.DocumentImageGUID)
	}

	findDocImages, err := svc.repoImage.FindInGUIDs(ctx, shopID, tempDocImageGUIDs)

	if err != nil {
		return err
	}

	findDoc.References = append(findDoc.References, docRef)

	timeAt := svc.timeNowFnc()

	findDoc.UpdatedAt = timeAt
	findDoc.UpdatedBy = authUsername

	if err = svc.repoImageGroup.Update(ctx, shopID, groupGUID, findDoc); err != nil {
		return err
	}

	for _, docImage := range findDocImages {

		docImage.UpdatedAt = timeAt
		docImage.UpdatedBy = authUsername

		docImage.References = append(docImage.References, docRef)

		if err = svc.repoImage.Update(ctx, shopID, docImage.GuidFixed, docImage); err != nil {
			return err
		}
	}

	return nil
}

func (svc DocumentImageService) UpdateTagsInDocumentImageGroup(shopID string, authUsername string, groupGUID string, tags []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repoImageGroup.FindByGuid(ctx, shopID, groupGUID)

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

	if err = svc.repoImageGroup.Update(ctx, shopID, groupGUID, findDoc); err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) DeleteReferenceByDocumentImageGroup(shopID string, authUsername string, groupGUID string, docRef models.Reference) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repoImageGroup.FindByGuid(ctx, shopID, groupGUID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document image group not found")
	}

	// Clear references
	docImageClearRefs, err := svc.repoImage.FindByReference(ctx, shopID, docRef)
	if err != nil {
		return err
	}

	docImageGroupClearRefs, err := svc.repoImageGroup.FindByReference(ctx, shopID, docRef)
	if err != nil {
		return err
	}

	for _, docImage := range docImageClearRefs {
		tempDocRefs := lo.Filter[models.Reference](docImage.References, func(tempDocRef models.Reference, idx int) bool {
			return docRef.Module != tempDocRef.Module && docRef.DocNo != tempDocRef.DocNo
		})

		docImage.References = tempDocRefs

		svc.repoImage.Update(ctx, shopID, docImage.GuidFixed, docImage)
	}

	for _, docImageGroup := range docImageGroupClearRefs {
		tempDocRefs := lo.Filter[models.Reference](docImageGroup.References, func(tempDocRef models.Reference, idx int) bool {
			return docRef.Module != tempDocRef.Module && docRef.DocNo != tempDocRef.DocNo
		})

		docImageGroup.References = tempDocRefs

		svc.repoImageGroup.Update(ctx, shopID, docImageGroup.GuidFixed, docImageGroup)
	}

	_, err = svc.messageQueueReCountDocumentImageGroup(ctx, shopID, findDoc.TaskGUID)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (svc DocumentImageService) DeleteDocumentImageGroupByGuid(shopID string, authUsername string, documentImageGroupGuidFixed string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDocGroup, err := svc.repoImageGroup.FindByGuid(ctx, shopID, documentImageGroupGuidFixed)

	if err != nil {
		return err
	}

	if len(findDocGroup.GuidFixed) < 1 {
		return nil
	}

	if svc.isDocumentImageGroupHasReferenced(findDocGroup) {
		return errors.New("document has referenced")
	}

	err = svc.repoImageGroup.Transaction(ctx, func(ctx context.Context) error {

		for _, docImage := range *findDocGroup.ImageReferences {
			err = svc.repoImage.DeleteByGuidfixed(ctx, shopID, docImage.DocumentImageGUID, authUsername)

			if err != nil {
				return err
			}
		}

		if err = svc.repoImageGroup.DeleteByGuidfixed(ctx, shopID, findDocGroup.GuidFixed); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	_, err = svc.messageQueueReCountDocumentImageGroup(ctx, shopID, findDocGroup.TaskGUID)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (svc DocumentImageService) DeleteDocumentImageGroupByGuids(shopID string, authUsername string, documentImageGroupGuidFixeds []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	err := svc.repoImageGroup.Transaction(ctx, func(ctx context.Context) error {
		for _, DocumentImageGroupGuidFixed := range documentImageGroupGuidFixeds {
			findDocGroup, err := svc.repoImageGroup.FindByGuid(ctx, shopID, DocumentImageGroupGuidFixed)

			if err != nil {
				return err
			}

			if svc.isDocumentImageGroupHasReferenced(findDocGroup) {
				return errors.New("document has referenced")
			}

			for _, docImage := range *findDocGroup.ImageReferences {
				err = svc.repoImage.DeleteByGuidfixed(ctx, shopID, docImage.DocumentImageGUID, authUsername)

				if err != nil {
					return err
				}
			}

			if err = svc.repoImageGroup.DeleteByGuidfixed(ctx, shopID, findDocGroup.GuidFixed); err != nil {
				return err
			}

			_, err = svc.messageQueueReCountDocumentImageGroup(ctx, shopID, findDocGroup.TaskGUID)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) UnGroupDocumentImageGroup(shopID string, authUsername string, groupGUID string) ([]string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDocGroup, err := svc.repoImageGroup.FindByGuid(ctx, shopID, groupGUID)

	if err != nil {
		return []string{}, err
	}

	if len(findDocGroup.GuidFixed) < 1 {
		return []string{}, nil
	}

	if svc.isDocumentImageGroupHasReferenced(findDocGroup) {
		return []string{}, errors.New("document has referenced")
	}

	updatedAt := svc.timeNowFnc()
	newImageGroupGUIDs := []string{}
	for _, imageRef := range *findDocGroup.ImageReferences {
		imageGroupGUID := svc.newDocumentImageGroupGUIDFnc()

		tags := []string{}
		if findDocGroup.Tags != nil {
			tags = *findDocGroup.Tags
		}

		docImageGroup := svc.createImageGroupByDocumentImage(shopID, authUsername, imageGroupGUID, imageRef, imageRef.ImageURI, tags, findDocGroup.TaskGUID, findDocGroup.PathTask, updatedAt)

		newXOrderDocImgGroup, _ := svc.newXOrderDocumentImageGroup(ctx, shopID, docImageGroup.TaskGUID)

		docImageGroup.XOrder = newXOrderDocImgGroup
		_, err = svc.repoImageGroup.Create(ctx, docImageGroup)

		if err != nil {
			return []string{}, err
		}

		newImageGroupGUIDs = append(newImageGroupGUIDs, imageGroupGUID)
	}

	svc.repoImageGroup.DeleteByGuidfixed(ctx, shopID, groupGUID)

	return newImageGroupGUIDs, nil
}

func (svc DocumentImageService) ListDocumentImageGroup(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentImageGroupInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{"title"}
	docList, pagination, err := svc.repoImageGroup.FindPageImageGroup(ctx, shopID, filters, searchInFields, pageable)

	return docList, pagination, err
}

func (svc DocumentImageService) GetDocumentImageDocRefGroup(shopID string, docImageGroupGUID string) (models.DocumentImageGroupInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	doc, err := svc.repoImageGroup.FindByGuid(ctx, shopID, docImageGroupGUID)

	if err != nil {
		return models.DocumentImageGroupInfo{}, err
	}

	if doc.ID.IsZero() {
		return models.DocumentImageGroupInfo{}, errors.New("document not found")
	}

	return doc.DocumentImageGroupInfo, nil
}

func (svc DocumentImageService) GetDocumentImageGroupByDocRef(shopID string, docRef string) (models.DocumentImageGroupInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repoImageGroup.FindOne(ctx, shopID, bson.M{"references.docno": docRef})

	if err != nil {
		return models.DocumentImageGroupInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DocumentImageGroupInfo{}, errors.New("document not found")
	}

	return findDoc.DocumentImageGroupInfo, nil

}

func (svc DocumentImageService) XSortsUpdate(ctx context.Context, shopID string, authUsername string, taskGUID string, xsorts []models.XSortDocumentImageGroupRequest) error {
	for _, xsort := range xsorts {
		if len(xsort.GUIDFixed) < 1 {
			continue
		}

		err := svc.repoImageGroup.UpdateXOrder(ctx, shopID, taskGUID, xsort.GUIDFixed, xsort.XOrder)

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
		ImageURI:       documentImage.ImageURI,
		CloneImageFrom: documentImage.CloneImageFrom,
		Name:           documentImage.Name,
		UploadedBy:     authUsername,
		UploadedAt:     createdAt,
		MetaFileAt:     documentImage.MetaFileAt,
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

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	err := svc.repoImageGroup.RemoveDocumentImageByDocumentImageGUIDs(ctx, shopID, docImageGUIDs)
	if err != nil {
		return err
	}

	err = svc.repoImageGroup.DeleteByGUIDsIsDocumentImageEmpty(ctx, shopID, docImageGroupGUIDs)
	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) clearUpdateDocumentImageGroupByDocumentGUIDs(shopID string, docGroupGUID string, clearDocImageGUIDs []string, docImageGUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	err := svc.repoImageGroup.RemoveDocumentImageByDocumentImageGUIDsWithoutDocumentImageGroupGUID(ctx, shopID, docGroupGUID, docImageGUIDs)
	if err != nil {
		return err
	}

	err = svc.repoImageGroup.DeleteByGUIDsIsDocumentImageEmptyWithoutDocumentImageGroupGUID(ctx, shopID, docGroupGUID, clearDocImageGUIDs)
	if err != nil {
		return err
	}

	err = svc.repoImageGroup.DeleteByGUIDIsDocumentImageEmpty(ctx, shopID, docGroupGUID)
	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentImageService) messageQueueReCountDocumentImageGroup(ctx context.Context, shopID string, taskGUID string) (int, error) {

	docList, err := svc.repoImageGroup.FindStatusByDocumentImageGroupTask(ctx, shopID, taskGUID)

	if err != nil {
		return 0, err
	}

	countDoc := len(docList)
	if countDoc < 1 {
		return 0, nil
	}

	totalStatus := map[int8]int{}

	for i := 0; i <= models.IMAGE_FROM_REJECT; i++ {
		totalStatus[int8(i)] = 0
	}

	for _, doc := range docList {
		if _, ok := totalStatus[doc.Status]; !ok {
			totalStatus[doc.Status] = 0
		}
		totalStatus[doc.Status] = totalStatus[doc.Status] + 1
	}

	countStatus := []models.CountStatus{}
	for status, count := range totalStatus {
		countStatus = append(countStatus, models.CountStatus{
			Status: status,
			Count:  count,
		})
	}

	taskMsg := models.DocumentImageTaskChangeMessage{
		ShopID:      shopID,
		TaskGUID:    taskGUID,
		Count:       countDoc,
		CountStatus: countStatus,
	}

	err = svc.repoMessagequeue.TaskChange(taskMsg)
	if err != nil {
		return 0, err
	}

	return countDoc, nil
}

func (svc DocumentImageService) newXOrderDocumentImageGroup(ctx context.Context, shopID string, taskGUID string) (int, error) {

	findDoc, err := svc.repoImageGroup.FindLastOneByTask(ctx, shopID, taskGUID)

	if err != nil {
		return 0, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return 0, nil
	}

	return findDoc.XOrder + 1, nil
}
