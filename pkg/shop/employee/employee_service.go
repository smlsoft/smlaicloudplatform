package employee

import (
	"errors"
	"fmt"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"sync"
	"time"

	common "smlcloudplatform/pkg/models"

	mastersync "smlcloudplatform/pkg/mastersync/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IEmployeeService interface {
	Login(shopID string, loginReq models.EmployeeRequestLogin) (*models.EmployeeInfo, error)
	Register(shopID string, authUsername string, emp models.EmployeeRequestRegister) (string, error)
	Get(shopID string, username string) (models.EmployeeInfo, error)
	List(shopID string, q string, page int, limit int) ([]models.EmployeeInfo, mongopagination.PaginationData, error)
	Update(shopID string, authUsername string, emp models.EmployeeRequestUpdate) error
	UpdatePassword(shopID string, authUsername string, emp models.EmployeeRequestPassword) error

	GetModuleName() string
}

type EmployeeService struct {
	repo          IEmployeeRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository
}

func NewEmployeeService(repo IEmployeeRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *EmployeeService {
	return &EmployeeService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}
}

func (svc EmployeeService) Login(shopID string, loginReq models.EmployeeRequestLogin) (*models.EmployeeInfo, error) {

	loginReq.Username = strings.TrimSpace(loginReq.Username)

	findUser, err := svc.repo.FindEmployeeByUsername(shopID, loginReq.Username)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, errors.New("auth: database connect error")
	}

	if len(findUser.Username) < 1 {
		return nil, errors.New("username is exists")
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
		return "", errors.New("code is exists")
	}

	userFind, err := svc.repo.FindEmployeeByUsername(shopID, emp.Username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return "", err
	}

	if len(userFind.Username) > 0 {
		return "", errors.New("username is exists")
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

func (svc EmployeeService) Get(shopID string, username string) (models.EmployeeInfo, error) {
	doc, err := svc.repo.FindEmployeeByUsername(shopID, username)

	if err != nil {
		return models.EmployeeInfo{}, err
	}

	return doc.EmployeeInfo, nil
}

func (svc EmployeeService) List(shopID string, q string, page int, limit int) ([]models.EmployeeInfo, mongopagination.PaginationData, error) {

	docList, pagination, err := svc.repo.FindEmployeeByShopIDPage(shopID, q, page, limit)

	if err != nil {
		return []models.EmployeeInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc EmployeeService) Update(shopID string, authUsername string, emp models.EmployeeRequestUpdate) error {

	userFind, err := svc.repo.FindEmployeeByUsername(shopID, emp.Username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.Username) < 1 {
		return errors.New("username is exists")
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

	userFind, err := svc.repo.FindEmployeeByUsername(shopID, emp.Username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.Username) < 1 {
		return errors.New("username is exists")
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

func (svc EmployeeService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, mongopagination.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.EmployeeDeleteActivity
	var pagination1 mongopagination.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.repo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.EmployeeActivity
	var pagination2 mongopagination.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.repo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return common.LastActivity{}, pagination1, err1
	}

	if err2 != nil {
		return common.LastActivity{}, pagination2, err2
	}

	lastActivity := common.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	pagination := pagination1

	if pagination.Total < pagination2.Total {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc EmployeeService) LastActivityOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) (common.LastActivity, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.EmployeeDeleteActivity
	var err1 error

	go func() {
		deleteDocList, err1 = svc.repo.FindDeletedOffset(shopID, lastUpdatedDate, skip, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.EmployeeActivity

	var err2 error

	go func() {
		createAndUpdateDocList, err2 = svc.repo.FindCreatedOrUpdatedOffset(shopID, lastUpdatedDate, skip, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return common.LastActivity{}, err1
	}

	if err2 != nil {
		return common.LastActivity{}, err2
	}

	lastActivity := common.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	return lastActivity, nil
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
