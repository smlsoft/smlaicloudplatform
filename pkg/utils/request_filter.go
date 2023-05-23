package utils

import (
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type FilterRequest struct {
	Param string
	Field string
	Type  string
}

func GetFilters(getParamFunc func(string) string, filterFields []FilterRequest) map[string]interface{} {
	filters := make(map[string]interface{})

	for _, field := range filterFields {
		if field.Param == "" {
			continue
		}

		filterValue := ""
		if field.Param != "-" {
			filterValue = getParamFunc(field.Param)
			if filterValue == "" {
				continue
			}
		}

		if field.Field == "" {
			field.Field = field.Param
		}

		switch field.Type {
		case "string":
			filters[field.Field] = filterValue
		case "int":
			if val, err := strconv.Atoi(filterValue); err == nil {
				filters[field.Field] = val
			}
		case "float64":
			if val, err := strconv.ParseFloat(filterValue, 64); err == nil {
				filters[field.Field] = val
			}
		case "bool":
			if val, err := strconv.ParseBool(filterValue); err == nil {
				filters[field.Field] = val
			}

		case "rangeDate":

			fromDateStr := strings.TrimSpace(getParamFunc("fromdate"))
			toDateStr := strings.TrimSpace(getParamFunc("todate"))

			if len(fromDateStr) > 0 && len(toDateStr) > 0 {
				fromDate, err1 := time.Parse("2006-01-02", fromDateStr)
				toDate, err2 := time.Parse("2006-01-02", toDateStr)

				if err1 == nil && err2 == nil {
					filters[field.Field] = bson.M{
						"$gte": fromDate,
						"$lt":  toDate.AddDate(0, 0, 1),
					}
				}
			}

		default:
			filters[field.Field] = filterValue
		}
	}

	return filters
}
