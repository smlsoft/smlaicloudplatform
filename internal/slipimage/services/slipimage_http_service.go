package services

import (
	"context"
	"errors"
	"fmt"
	mastersync "smlcloudplatform/internal/mastersync/repositories"
	"smlcloudplatform/internal/services"
	"smlcloudplatform/internal/slipimage/models"
	"smlcloudplatform/internal/slipimage/repositories"
	"smlcloudplatform/internal/utils"
	micromodels "smlcloudplatform/pkg/microservice/models"
	"strings"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISlipImageHttpService interface {
	CreateSlipImage(shopID string, authUsername string, doc models.SlipImageRequest) (models.SlipImageInfo, error)
	DeleteSlipImage(shopID string, guid string, authUsername string) error
	DeleteSlipImageByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoSlipImage(shopID string, guid string) (models.SlipImageInfo, error)
	SearchSlipImage(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SlipImageInfo, mongopagination.PaginationData, error)
	SearchSlipImageStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SlipImageInfo, int, error)
}

type SlipImageHttpService struct {
	repo             repositories.ISlipImageMongoRepository
	repoStorageImage repositories.ISlipImageStorageImageRepository
	syncCacheRepo    mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.SlipImageActivity, models.SlipImageDeleteActivity]
	contextTimeout time.Duration
}

func NewSlipImageHttpService(
	repo repositories.ISlipImageMongoRepository,
	repoStorageImage repositories.ISlipImageStorageImageRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,

	contextTimeout time.Duration,
) *SlipImageHttpService {

	insSvc := &SlipImageHttpService{
		repo:             repo,
		repoStorageImage: repoStorageImage,
		syncCacheRepo:    syncCacheRepo,

		contextTimeout: contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.SlipImageActivity, models.SlipImageDeleteActivity](repo)

	return insSvc
}

func (svc SlipImageHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc SlipImageHttpService) CreateSlipImage(shopID string, authUsername string, payload models.SlipImageRequest) (models.SlipImageInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	fileSize := payload.File.Size

	newGuidFixed := utils.NewGUID()

	layout := "2006-01-02"
	docDateStr := payload.DocDate.Format(layout)

	tempSplitFileName := strings.Split(payload.File.Filename, ".")
	fileExt := tempSplitFileName[len(tempSplitFileName)-1]

	fileName := fmt.Sprintf("%s/slip/%s/%s/%s", shopID, payload.PosID, docDateStr, payload.DocNo)

	uploadUri, err := svc.repoStorageImage.Upload(payload.File, fileName, fileExt)

	if err != nil {
		return models.SlipImageInfo{}, fmt.Errorf("upload file failed")
	}

	docData := models.SlipImageDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SlipImage = models.SlipImage{
		URI:     uploadUri,
		Size:    fileSize,
		DocNo:   payload.DocNo,
		DocDate: payload.DocDate,
		PosID:   payload.PosID,
	}

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return models.SlipImageInfo{}, err
	}

	return docData.SlipImageInfo, nil
}

func (svc SlipImageHttpService) DeleteSlipImage(shopID string, guid string, authUsername string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(ctx, shopID, guid, authUsername)
	if err != nil {
		return err
	}

	return nil
}

func (svc SlipImageHttpService) DeleteSlipImageByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc SlipImageHttpService) InfoSlipImage(shopID string, guid string) (models.SlipImageInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.SlipImageInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SlipImageInfo{}, errors.New("document not found")
	}

	return findDoc.SlipImageInfo, nil
}

func (svc SlipImageHttpService) SearchSlipImage(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SlipImageInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"uid",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.SlipImageInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc SlipImageHttpService) SearchSlipImageStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SlipImageInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"uid",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.SlipImageInfo{}, 0, err
	}

	return docList, total, nil
}
