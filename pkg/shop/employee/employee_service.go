package employee

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"

	common "smlcloudplatform/pkg/models"

	mastersync "smlcloudplatform/pkg/mastersync/repositories"

	"github.com/userplant/mongopagination"
)

type IEmployeeService interface {
	Login(shopID string, loginReq models.EmployeeRequestLogin) (*models.EmployeeInfo, error)
	Register(shopID string, authUsername string, emp models.EmployeeRequestRegister) (string, error)
	Get(shopID string, code string) (models.EmployeeInfo, error)
	SearchEmployee(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.EmployeeInfo, mongopagination.PaginationData, error)
	SearchEmployeeStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.EmployeeInfo, int, error)
	Update(shopID string, authUsername string, emp models.EmployeeRequestUpdate) error
	UpdatePassword(shopID string, authUsername string, emp models.EmployeeRequestPassword) error

	LastActivity(shopID string, action string, lastUpdatedDate time.Time, pageable micromodels.Pageable) (common.LastActivity, mongopagination.PaginationData, error)
	GetModuleName() string
}

type EmployeeService struct {
	repo          IEmployeeRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository

	services.ActivityService[models.EmployeeActivity, models.EmployeeDeleteActivity]
}

func NewEmployeeService(repo IEmployeeRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *EmployeeService {
	insSvc := &EmployeeService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.EmployeeActivity, models.EmployeeDeleteActivity](repo)
	return insSvc
}

func (svc EmployeeService) Login(shopID string, loginReq models.EmployeeRequestLogin) (*models.EmployeeInfo, error) {

	loginReq.Code = strings.TrimSpace(loginReq.Code)

	findUser, err := svc.repo.FindEmployeeByCode(shopID, loginReq.Code)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, errors.New("auth: database connect error")
	}

	if len(findUser.Code) < 1 {
		return nil, errors.New("user code is exists")
	}

	passwordInvalid := !utils.CheckHashPassword(loginReq.Password, findUser.Password)

	if passwordInvalid {
		return nil, errors.New("password is not invalid")
	}

	return &findUser.EmployeeInfo, nil
}

func (svc EmployeeService) Register(shopID string, authUsername string, emp models.EmployeeRequestRegister) (string, error) {

	findUserCode, err := svc.repo.FindEmployeeByCode(shopID, emp.Code)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return "", errors.New("auth: database connect error")
	}

	if len(findUserCode.Code) > 0 {
		return "", errors.New("user code is exists")
	}

	hashPassword, err := utils.HashPassword(emp.Password)

	if err != nil {
		return "", err
	}

	newGuid := utils.NewGUID()

	empDoc := models.EmployeeDoc{}

	empDoc.ShopID = shopID
	empDoc.GuidFixed = newGuid
	empDoc.Employee = emp.Employee
	empDoc.Password = hashPassword

	empDoc.CreatedBy = authUsername
	empDoc.CreatedAt = time.Now()

	_, err = svc.repo.Create(empDoc)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuid, nil
}

func (svc EmployeeService) Get(shopID string, code string) (models.EmployeeInfo, error) {
	doc, err := svc.repo.FindEmployeeByCode(shopID, code)

	if err != nil {
		return models.EmployeeInfo{}, err
	}

	return doc.EmployeeInfo, nil
}

func (svc EmployeeService) Update(shopID string, authUsername string, emp models.EmployeeRequestUpdate) error {

	userFind, err := svc.repo.FindEmployeeByCode(shopID, emp.Code)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.Code) < 1 {
		return errors.New("user code is exists")
	}

	userFind.Name = emp.Name
	userFind.Roles = emp.Roles
	userFind.ProfilePicture = emp.ProfilePicture

	userFind.UpdatedBy = authUsername
	userFind.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, userFind.GuidFixed, userFind)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc EmployeeService) UpdatePassword(shopID string, authUsername string, emp models.EmployeeRequestPassword) error {

	userFind, err := svc.repo.FindEmployeeByCode(shopID, emp.Code)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.Code) < 1 {
		return errors.New("user code is exists")
	}

	hashPassword, err := utils.HashPassword(emp.Password)

	if err != nil {
		return err
	}

	userFind.Password = hashPassword

	userFind.UpdatedBy = authUsername
	userFind.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, userFind.GuidFixed, userFind)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc EmployeeService) SearchEmployee(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.EmployeeInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"name",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.EmployeeInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc EmployeeService) SearchEmployeeStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.EmployeeInfo, int, error) {
	searchInFields := []string{
		"code",
		"name",
	}

	selectFields := map[string]interface{}{
		"guidfixed":      1,
		"code":           1,
		"email":          1,
		"profilepicture": 1,
		"roles":          1,
		"isenabled":      1,
	}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.EmployeeInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc EmployeeService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc EmployeeService) GetModuleName() string {
	return "employee"
}
