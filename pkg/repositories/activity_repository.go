package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/restaurant"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IActivityDeleteRepo interface {
	models.CategoryDeleteActivity | restaurant.ShopZoneDeleteActivity | restaurant.ShopTableDeleteActivity | restaurant.PrinterTerminalDeleteActivity | restaurant.KitchenDeleteActivity
}

type IActivityCreateUpdateReo interface {
	models.CategoryActivity | restaurant.ShopZoneActivity | restaurant.ShopTableActivity | restaurant.PrinterTerminalActivity | restaurant.KitchenActivity
}

type ActivityRepository[TCU IActivityCreateUpdateReo, TDEL IActivityDeleteRepo] struct {
	pst microservice.IPersisterMongo
}

func NewActivityRepository[TCU IActivityCreateUpdateReo, TDEL IActivityDeleteRepo](pst microservice.IPersisterMongo) ActivityRepository[TCU, TDEL] {
	return ActivityRepository[TCU, TDEL]{
		pst: pst,
	}
}

func (repo ActivityRepository[TCU, TDEL]) FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]TDEL, mongopagination.PaginationData, error) {

	docList := []TDEL{}
	pagination, err := repo.pst.FindPage(new(TDEL), limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}, &docList)

	if err != nil {
		return []TDEL{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]TCU, mongopagination.PaginationData, error) {

	docList := []TCU{}
	pagination, err := repo.pst.FindPage(new(TCU), limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}, &docList)

	if err != nil {
		return []TCU{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
