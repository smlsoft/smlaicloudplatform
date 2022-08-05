package mogoutil

import (
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
)

func AggregatePageDecode[T any](paginatedData *mongopagination.PaginatedData) ([]T, error) {
	var aggList []T
	for _, raw := range paginatedData.Data {
		var product *T
		marshallErr := bson.Unmarshal(raw, &product)

		if marshallErr != nil {
			return aggList, marshallErr
		}
		aggList = append(aggList, *product)
	}

	return aggList, nil
}
