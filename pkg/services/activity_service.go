package services

import (
	"smlcloudplatform/pkg/repositories"
	"sync"
	"time"

	common "smlcloudplatform/pkg/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type ActivityService[TCU any, TDEL any] struct {
	repo repositories.IActivityRepository[TCU, TDEL]
}

func NewActivityService[TCU any, TDEL any](repo repositories.IActivityRepository[TCU, TDEL]) ActivityService[TCU, TDEL] {
	return ActivityService[TCU, TDEL]{
		repo: repo,
	}
}

func (svc ActivityService[TCU, TDEL]) LastActivity(shopID string, action string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, mongopagination.PaginationData, error) {
	lastActivity := common.LastActivity{}
	var wg sync.WaitGroup

	isActionRemove := len(action) == 0 || action == "all" || action == "remove"
	isActionNew := len(action) == 0 || action == "all" || action == "new"

	var deleteDocList []TDEL
	pagination1 := mongopagination.PaginationData{}
	var errFindDel error

	if isActionRemove {
		wg.Add(1)
		go func() {
			deleteDocList, pagination1, errFindDel = svc.repo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
			wg.Done()
		}()
	}

	var createAndUpdateDocList []TCU
	pagination2 := mongopagination.PaginationData{}
	var errFindCreateUpdate error

	if isActionNew {
		wg.Add(1)
		go func() {
			createAndUpdateDocList, pagination2, errFindCreateUpdate = svc.repo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
			wg.Done()
		}()
	}

	wg.Wait()

	if isActionRemove {
		if errFindDel != nil {
			return common.LastActivity{}, pagination1, errFindDel
		}
		lastActivity.Remove = &deleteDocList
	}

	if isActionNew {
		if errFindCreateUpdate != nil {
			return common.LastActivity{}, pagination2, errFindCreateUpdate
		}
		lastActivity.New = &createAndUpdateDocList
	}

	pagination := pagination1

	if pagination.Total < pagination2.Total {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc ActivityService[TCU, TDEL]) LastActivityOffset(shopID string, action string, lastUpdatedDate time.Time, skip int, limit int) (common.LastActivity, error) {
	lastActivity := common.LastActivity{}

	var wg sync.WaitGroup

	isActionRemove := len(action) == 0 || action == "all" || action == "remove"
	isActionNew := len(action) == 0 || action == "all" || action == "new"

	var errFindDel error
	var deleteDocList []TDEL
	if isActionRemove {
		wg.Add(1)
		go func() {
			deleteDocList, errFindDel = svc.repo.FindDeletedOffset(shopID, lastUpdatedDate, skip, limit)
			wg.Done()
		}()
	}

	var createAndUpdateDocList []TCU
	var errFindCreateUpdate error
	if isActionNew {
		wg.Add(1)
		go func() {
			createAndUpdateDocList, errFindCreateUpdate = svc.repo.FindCreatedOrUpdatedOffset(shopID, lastUpdatedDate, skip, limit)
			wg.Done()
		}()
	}

	wg.Wait()

	if isActionRemove {
		if errFindDel != nil {
			return common.LastActivity{}, errFindDel
		}
		lastActivity.Remove = &deleteDocList
	}

	if isActionNew {
		if errFindCreateUpdate != nil {
			return common.LastActivity{}, errFindCreateUpdate
		}
		lastActivity.New = &createAndUpdateDocList
	}

	return lastActivity, nil
}
