package employee

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
)

type IEmployeeService interface {
	Login(shopID string, loginReq models.EmployeeRequestLogin) (*models.EmployeeInfo, error)
	Register(shopID string, authUsername string, emp models.EmployeeRequestRegister) (string, error)
	Get(shopID string, username string) (models.EmployeeInfo, error)
	List(shopID string, q string, page int, limit int) ([]models.EmployeeInfo, paginate.PaginationData, error)
	Update(shopID string, authUsername string, emp models.EmployeeRequestUpdate) error
	UpdatePassword(shopID string, authUsername string, emp models.EmployeeRequestPassword) error
}

type EmployeeService struct {
	empRepo IEmployeeRepository
}

func NewEmployeeService(empRepo IEmployeeRepository) EmployeeService {
	return EmployeeService{
		empRepo: empRepo,
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

	return newGuid, nil
}

func (svc EmployeeService) Get(shopID string, username string) (models.EmployeeInfo, error) {
	doc, err := svc.empRepo.FindEmployeeByUsername(shopID, username)

	if err != nil {
		return models.EmployeeInfo{}, err
	}

	return doc.EmployeeInfo, nil
}

func (svc EmployeeService) List(shopID string, q string, page int, limit int) ([]models.EmployeeInfo, paginate.PaginationData, error) {

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
	userFind.Role = *emp.Role

	userFind.UpdatedBy = authUsername
	userFind.UpdatedAt = time.Now()

	err = svc.empRepo.Update(shopID, userFind.GuidFixed, userFind)

	if err != nil {
		return err
	}

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

	return nil
}
