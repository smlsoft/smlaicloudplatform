package mastersync

import (
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/models"
	"time"

	"github.com/userplant/mongopagination"
)

type ActivityModuleManager struct {
	activityModuleList map[string]ActivityModule
}

func NewActivityModuleManager() *ActivityModuleManager {
	return &ActivityModuleManager{
		activityModuleList: map[string]ActivityModule{},
	}
}

func (m *ActivityModuleManager) Add(activityModule ActivityModule) *ActivityModuleManager {
	m.activityModuleList[activityModule.GetModuleName()] = activityModule
	return m
}

func (m ActivityModuleManager) GetList() map[string]ActivityModule {
	return m.activityModuleList
}

func (m ActivityModuleManager) GetModules() []string {
	modules := []string{}
	for module := range m.activityModuleList {
		modules = append(modules, module)
	}
	return modules
}

func (m ActivityModuleManager) GetPage(moduleSelectList map[string]struct{}, activityParam ActivityParamPage) (map[string]interface{}, mongopagination.PaginationData, error) {
	moduleList := map[string]ActivityModule{}

	for _, activityModule := range m.activityModuleList {
		moduleList[activityModule.GetModuleName()] = activityModule
	}

	return listDataModulePage(moduleList, moduleSelectList, activityParam)
}

type ActivityModule interface {
	LastActivity(string, string, time.Time, micromodels.Pageable) (models.LastActivity, mongopagination.PaginationData, error)
	LastActivityStep(string, string, time.Time, micromodels.PageableStep) (models.LastActivity, error)
	GetModuleName() string
}

type ActivityParamPage struct {
	ShopID     string
	Action     string
	LastUpdate time.Time
	Pageable   micromodels.Pageable
}

type ActivityParamOffset struct {
	ShopID       string
	Action       string
	LastUpdate   time.Time
	PageableStep micromodels.PageableStep
}

func listDataModulePage(appModules map[string]ActivityModule, moduleSelectList map[string]struct{}, param ActivityParamPage) (map[string]interface{}, mongopagination.PaginationData, error) {

	result := map[string]interface{}{}

	resultPagination := mongopagination.PaginationData{}
	for moduleName, appModule := range appModules {
		if len(moduleSelectList) == 0 || isSelectModule(moduleSelectList, moduleName) {
			docList, pagination, err := appModule.LastActivity(param.ShopID, param.Action, param.LastUpdate, param.Pageable)

			if err != nil {
				return map[string]interface{}{}, mongopagination.PaginationData{}, err
			}

			result[moduleName] = docList

			if pagination.Total > resultPagination.Total {
				resultPagination = pagination
			}
		}
	}

	return result, resultPagination, nil
}

func listDataModuleOffset(appModules map[string]ActivityModule, moduleSelectList map[string]struct{}, param ActivityParamOffset) (map[string]interface{}, error) {

	result := map[string]interface{}{}

	for moduleName, appModule := range appModules {
		if len(moduleSelectList) == 0 || isSelectModule(moduleSelectList, moduleName) {

			docList, err := appModule.LastActivityStep(param.ShopID, param.Action, param.LastUpdate, param.PageableStep)

			if err != nil {
				return map[string]interface{}{}, err
			}

			result[moduleName] = docList

		}
	}

	return result, nil
}

func isSelectModule(moduleList map[string]struct{}, moduleKey string) bool {
	_, ok := moduleList[moduleKey]
	return ok
}
