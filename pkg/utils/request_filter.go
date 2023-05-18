package utils

import "strconv"

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

		filterValue := getParamFunc(field.Param)
		if filterValue == "" {
			continue
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
		default:
			filters[field.Field] = filterValue
		}
	}

	return filters
}
