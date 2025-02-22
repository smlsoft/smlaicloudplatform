package services

import (
	"context"
	"errors"
	"fmt"
	"smlaicloudplatform/internal/debtaccount/debtor/models"
	"smlaicloudplatform/internal/debtaccount/debtor/repositories"
	groupModels "smlaicloudplatform/internal/debtaccount/debtorgroup/models"
	groupRepositories "smlaicloudplatform/internal/debtaccount/debtorgroup/repositories"
	"smlaicloudplatform/internal/logger"
	mastersync "smlaicloudplatform/internal/mastersync/repositories"
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/internal/utils/importdata"
	micromodels "smlaicloudplatform/pkg/microservice/models"
	"time"

	"github.com/samber/lo"
	"github.com/smlsoft/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDebtorHttpService interface {
	CreateDebtor(shopID string, authUsername string, doc models.DebtorRequest) (string, error)
	UpdateDebtor(shopID string, guid string, authUsername string, doc models.DebtorRequest) error
	DeleteDebtor(shopID string, guid string, authUsername string) error
	DeleteDebtorByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoDebtor(shopID string, guid string) (models.DebtorInfo, error)
	InfoDebtorByCode(shopID string, code string) (models.DebtorInfo, error)
	SearchDebtor(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DebtorInfo, mongopagination.PaginationData, error)
	SearchDebtorStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DebtorInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.DebtorRequest) (common.BulkImport, error)
	InfoAuthDebtor(shopID string, username string, password string) (models.DebtorInfo, error)

	GetModuleName() string
}

type DebtorHttpService struct {
	repo          repositories.IDebtorRepository
	repoMq        repositories.IDebtorMessageQueueRepository
	repoGroup     groupRepositories.IDebtorGroupRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.DebtorActivity, models.DebtorDeleteActivity]
	hashPassword      func(password string) (string, error)
	checkHashPassword func(password, hash string) bool
	contextTimeout    time.Duration
}

func NewDebtorHttpService(
	repo repositories.IDebtorRepository,
	repoMq repositories.IDebtorMessageQueueRepository,
	repoGroup groupRepositories.IDebtorGroupRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,
	hashPassword func(password string) (string, error),
	checkHashPassword func(password, hash string) bool,
) *DebtorHttpService {
	contextTimeout := time.Duration(15) * time.Second

	insSvc := &DebtorHttpService{
		repo:              repo,
		repoMq:            repoMq,
		repoGroup:         repoGroup,
		syncCacheRepo:     syncCacheRepo,
		hashPassword:      hashPassword,
		checkHashPassword: checkHashPassword,
		contextTimeout:    contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.DebtorActivity, models.DebtorDeleteActivity](repo)

	return insSvc
}

func (svc DebtorHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc DebtorHttpService) InfoAuthDebtor(shopID string, username string, password string) (models.DebtorInfo, error) {
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	if username == "" || password == "" {
		return models.DebtorInfo{}, errors.New("username or password incorrect")
	}

	findDoc, err := svc.repo.FindAuthByUsername(ctx, shopID, username)

	if err != nil {
		return models.DebtorInfo{}, err
	}

	if findDoc.Auth.Username == "" {
		return models.DebtorInfo{}, errors.New("username or password incorrect")
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DebtorInfo{}, errors.New("username or password incorrect")
	}

	if findDoc.Auth.Password == "" {
		return models.DebtorInfo{}, errors.New("username or password incorrect")
	}

	invalidPassword := !svc.checkHashPassword(password, findDoc.Auth.Password)

	if invalidPassword {
		return models.DebtorInfo{}, errors.New("username or password incorrect")
	}

	findDoc.Auth.Password = ""

	return findDoc.DebtorInfo, nil

}

func (svc DebtorHttpService) CreateDebtor(shopID string, authUsername string, doc models.DebtorRequest) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("code is exists")
	}

	if doc.Auth.Username != "" {
		findDocAuth, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "auth.username", doc.Code)

		if err != nil {
			return "", err
		}

		if findDocAuth.Auth.Username != "" {
			return "", errors.New("auth username is exists")
		}
	}

	newGuidFixed := utils.NewGUID()

	docData := models.DebtorDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Debtor = doc.Debtor
	docData.GroupGUIDs = &doc.Groups

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	if doc.Auth.Password != "" {
		hashedPassword, err := svc.hashPassword(doc.Auth.Password)
		if err != nil {
			return "", err
		}
		docData.Auth.Password = hashedPassword
	}

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", err
	}

	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.Create(docData)
		if err != nil {
			logger.GetLogger().Errorf("Create creditor message queue error :: %s", err.Error())
		}
	}()

	return newGuidFixed, nil
}

func (svc DebtorHttpService) UpdateDebtor(shopID string, guid string, authUsername string, doc models.DebtorRequest) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	if doc.Auth.Username != "" {
		findDocAuth, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "auth.username", doc.Auth.Username)

		if err != nil {
			return err
		}

		if findDoc.Auth.Username != findDocAuth.Auth.Username && findDocAuth.Auth.Username != "" {
			return errors.New("auth username is exists")
		}

	}

	dataDoc := findDoc

	dataDoc.Debtor = doc.Debtor
	dataDoc.GroupGUIDs = &doc.Groups

	dataDoc.UpdatedBy = authUsername
	dataDoc.UpdatedAt = time.Now()

	if doc.Auth.Password != "" {
		hashedPassword, err := svc.hashPassword(doc.Auth.Password)
		if err != nil {
			return err
		}
		dataDoc.Auth.Password = hashedPassword
	}

	err = svc.repo.Update(ctx, shopID, guid, dataDoc)

	if err != nil {
		return err
	}

	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.Update(dataDoc)
		if err != nil {
			logger.GetLogger().Errorf("Update creditor message queue error :: %s", err.Error())
		}
	}()

	return nil
}

func (svc DebtorHttpService) DeleteDebtor(shopID string, guid string, authUsername string) error {
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(ctx, shopID, guid, authUsername)
	if err != nil {
		return err
	}

	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.Delete(findDoc)
		if err != nil {
			logger.GetLogger().Errorf("Delete creditor message queue error :: %s", err.Error())
		}
	}()

	return nil
}

func (svc DebtorHttpService) DeleteDebtorByGUIDs(shopID string, authUsername string, GUIDs []string) error {
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDocs, err := svc.repo.FindByGuids(ctx, shopID, GUIDs)

	if err != nil {
		return err
	}

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err = svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.DeleteInBatch(findDocs)
		if err != nil {
			logger.GetLogger().Errorf("Delete creditor message queue error :: %s", err.Error())
		}
	}()

	return nil
}

func (svc DebtorHttpService) InfoDebtor(shopID string, guid string) (models.DebtorInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.DebtorInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DebtorInfo{}, errors.New("document not found")
	}

	findGroups, err := svc.repoGroup.FindByGuids(ctx, shopID, *findDoc.GroupGUIDs)

	if err != nil {
		return models.DebtorInfo{}, err
	}

	custGroupInfo := lo.Map[groupModels.DebtorGroupDoc, groupModels.DebtorGroupInfo](
		findGroups,
		func(docGroup groupModels.DebtorGroupDoc, idx int) groupModels.DebtorGroupInfo {
			return docGroup.DebtorGroupInfo
		})

	findDoc.DebtorInfo.Groups = &custGroupInfo

	docInfo := findDoc.DebtorInfo
	docInfo.Auth.Password = ""

	return docInfo, nil

}

func (svc DebtorHttpService) InfoDebtorByCode(shopID string, code string) (models.DebtorInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", code)

	if err != nil {
		return models.DebtorInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DebtorInfo{}, errors.New("document not found")
	}

	findGroups, err := svc.repoGroup.FindByGuids(ctx, shopID, *findDoc.GroupGUIDs)

	if err != nil {
		return models.DebtorInfo{}, err
	}

	custGroupInfo := lo.Map[groupModels.DebtorGroupDoc, groupModels.DebtorGroupInfo](
		findGroups,
		func(docGroup groupModels.DebtorGroupDoc, idx int) groupModels.DebtorGroupInfo {
			return docGroup.DebtorGroupInfo
		})

	findDoc.DebtorInfo.Groups = &custGroupInfo

	docInfo := findDoc.DebtorInfo
	docInfo.Auth.Password = ""

	return docInfo, nil

}

func (svc DebtorHttpService) SearchDebtor(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DebtorInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
		"names.name",
		"groups",
		"fundcode",
		"addressforbilling.address.0",
		"addressforbilling.phoneprimary",
		"addressforbilling.phonesecondary",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.DebtorInfo{}, pagination, err
	}

	for idx, doc := range docList {
		if doc.GroupGUIDs != nil {
			findCustGroups, err := svc.repoGroup.FindByGuids(ctx, shopID, *doc.GroupGUIDs)
			if err != nil {
				return []models.DebtorInfo{}, pagination, err
			}

			custGroupInfo := lo.Map[groupModels.DebtorGroupDoc, groupModels.DebtorGroupInfo](
				findCustGroups,
				func(docGroup groupModels.DebtorGroupDoc, idx int) groupModels.DebtorGroupInfo {
					return docGroup.DebtorGroupInfo
				})

			docList[idx].Groups = &custGroupInfo
		}
	}

	for i := 0; i < len(docList); i++ {
		docList[i].Auth.Password = ""
	}

	return docList, pagination, nil
}

func (svc DebtorHttpService) SearchDebtorStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DebtorInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
		"names.name",
		"groups",
		"fundcode",
		"addressforbilling.address.0",
		"addressforbilling.phoneprimary",
		"addressforbilling.phonesecondary",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.DebtorInfo{}, 0, err
	}

	for idx, doc := range docList {
		if doc.GroupGUIDs != nil {
			findCustGroups, err := svc.repoGroup.FindByGuids(ctx, shopID, *doc.GroupGUIDs)
			if err != nil {
				return []models.DebtorInfo{}, 0, err
			}

			custGroupInfo := lo.Map[groupModels.DebtorGroupDoc, groupModels.DebtorGroupInfo](
				findCustGroups,
				func(docGroup groupModels.DebtorGroupDoc, idx int) groupModels.DebtorGroupInfo {
					return docGroup.DebtorGroupInfo
				})

			docList[idx].Groups = &custGroupInfo
		}
	}

	for i := 0; i < len(docList); i++ {
		docList[i].Auth.Password = ""
	}

	return docList, total, nil
}

func (svc DebtorHttpService) SaveInBatch(shopID string, authUsername string, dataListReq []models.DebtorRequest) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	dataList := []models.Debtor{}
	for _, doc := range dataListReq {
		doc.GroupGUIDs = &doc.Groups
		dataList = append(dataList, doc.Debtor)
	}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Debtor](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Code)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(ctx, shopID, "code", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Debtor, models.DebtorDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Debtor) models.DebtorDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.DebtorDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Debtor = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Debtor, models.DebtorDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.DebtorDoc, error) {
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", guid)
		},
		func(doc models.DebtorDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Debtor, doc models.DebtorDoc) error {

			doc.Debtor = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			err = svc.repo.Update(ctx, shopID, doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(ctx, createDataList)

		if err != nil {
			return common.BulkImport{}, err
		}

	}

	createDataKey := []string{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.Code)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Code)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, svc.getDocIDKey(doc))
	}

	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.CreateInBatch(createDataList)
		if err != nil {
			logger.GetLogger().Errorf("Create creditor message queue error :: %s", err.Error())
		}
		svc.repoMq.UpdateInBatch(updateSuccessDataList)

		if err != nil {
			logger.GetLogger().Errorf("Update creditor message queue error :: %s", err.Error())
		}
	}()

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc DebtorHttpService) getDocIDKey(doc models.Debtor) string {
	return doc.Code
}

func (svc DebtorHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc DebtorHttpService) GetModuleName() string {
	return "debtor"
}
