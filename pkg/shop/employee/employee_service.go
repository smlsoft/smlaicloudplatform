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
}

type EmployeeService struct {
	empRepo   IEmployeeRepository
	cacheRepo mastersync.IMasterSyncCacheRepository
}

func NewEmployeeService(empRepo IEmployeeRepository, cacheRepo mastersync.IMasterSyncCacheRepository) *EmployeeService {
	return &EmployeeService{
		empRepo:   empRepo,
		cacheRepo: cacheRepo,
	}
}

func (svc EmployeeService) Login(shopID string, loginReq models.EmployeeRequestLogin) (*models.EmployeeInfo, error) {

	loginReq.Username = strings.TrimSpace(loginReq.Username)

	findUser, err := svc.empRepo.FindEmployeeByUsername(shopID, loginReq.Username)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, errors.New("auth: database connect error")
	}

	if len(findUser.Username) < 1 {
		return nil, errors.New("username is not exists")
	}

	passwordInvalid := !utils.CheckHashPassword(loginReq.Password, findUser.Password)

	if passwordInvalid {
		return nil, errors.New("password is not invalid")
	}

	return &findUser.EmployeeInfo, nil
}

func (svc EmployeeService) Register(shopID string, authUsername string, emp models.EmployeeRequestRegister) (string, error) {

	userFind, err := svc.empRepo.FindEmployeeByUsername(shopID, emp.Username)
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

	_, err = svc.empRepo.Create(empDoc)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuid, nil
}

func (svc EmployeeService) Get(shopID string, username string) (models.EmployeeInfo, error) {
	doc, err := svc.empRepo.FindEmployeeByUsername(shopID, username)

	if err != nil {
		return models.EmployeeInfo{}, err
	}

	return doc.EmployeeInfo, nil
}

func (svc EmployeeService) List(shopID string, q string, page int, limit int) ([]models.EmployeeInfo, mongopagination.PaginationData, error) {

	docList, pagination, err := svc.empRepo.FindEmployeeByShopIDPage(shopID, q, page, limit)

	if err != nil {
		return []models.EmployeeInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc EmployeeService) Update(shopID string, authUsername string, emp models.EmployeeRequestUpdate) error {

	userFind, err := svc.empRepo.FindEmployeeByUsername(shopID, emp.Username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.Username) < 1 {
		return errors.New("username is not exists")
	}

	userFind.Name = emp.Name
	userFind.Roles = *emp.Roles

	userFind.UpdatedBy = authUsername
	userFind.UpdatedAt = time.Now()

	err = svc.empRepo.Update(shopID, userFind.GuidFixed, userFind)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc EmployeeService) UpdatePassword(shopID string, authUsername string, emp models.EmployeeRequestPassword) error {

	userFind, err := svc.empRepo.FindEmployeeByUsername(shopID, emp.Username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.Username) < 1 {
		return errors.New("username is not exists")
	}

	hashPassword, err := utils.HashPassword(emp.Password)

	if err != nil {
		return err
	}

	userFind.Password = hashPassword

	userFind.UpdatedBy = authUsername
	userFind.UpdatedAt = time.Now()

	err = svc.empRepo.Update(shopID, userFind.GuidFixed, userFind)

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
		deleteDocList, pagination1, err1 = svc.empRepo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.EmployeeActivity
	var pagination2 mongopagination.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.empRepo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
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

func (svc EmployeeService) saveMasterSync(shopID string) {
	err := svc.cacheRepo.Save(shopID)

	if err != nil {
		fmt.Println("save category master cache error :: " + err.Error())
	}
}
