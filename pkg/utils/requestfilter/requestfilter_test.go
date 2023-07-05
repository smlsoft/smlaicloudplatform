package requestfilter_test

import (
	"smlcloudplatform/pkg/utils/requestfilter"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestGenerateFilters(t *testing.T) {
	filterFields := []requestfilter.FilterRequest{
		{Param: "strparam", Field: "strfield", Type: requestfilter.FieldTypeString},
		{Param: "strparamarr", Field: "strfieldarr", Type: requestfilter.FieldTypeString},
		{Param: "intparam", Field: "intfield", Type: requestfilter.FieldTypeInt},
		{Param: "floatparam", Field: "floatfield", Type: requestfilter.FieldTypeFloat64},
		{Param: "boolparam", Field: "boolfield", Type: requestfilter.FieldTypeBoolean},
		{Param: "dateparam", Field: "datefield", Type: requestfilter.FieldTypeRangeDate},
	}

	getParamFunc := func(param string) string {
		switch param {
		case "strparam":
			return "strvalue"
		case "strparamarr":
			return "strvalue1,strvalue2,strvalue3"
		case "intparam":
			return "42"
		case "floatparam":
			return "42.42"
		case "boolparam":
			return "true"
		case "dateparam":
			return "2023-07-05"
		case "fromdate":
			return "2023-07-05"
		case "todate":
			return "2023-07-06"
		default:
			return ""
		}
	}

	filters := requestfilter.GenerateFilters(getParamFunc, filterFields)

	if filters["strfield"] != "strvalue" {
		t.Errorf("Unexpected value for strfield: got %v", filters["strfield"])
	}

	if len(filters["strfieldarr"].(bson.M)["$in"].([]string)) != 3 {
		t.Errorf("Unexpected value for strfieldarr: got %v", filters["strfieldarr"])
	}

	if filters["intfield"] != 42 {
		t.Errorf("Unexpected value for intfield: got %v", filters["intfield"])
	}

	if filters["floatfield"] != 42.42 {
		t.Errorf("Unexpected value for floatfield: got %v", filters["floatfield"])
	}

	if filters["boolfield"] != true {
		t.Errorf("Unexpected value for boolfield: got %v", filters["boolfield"])
	}

	fromDate, _ := time.Parse(requestfilter.DateFormate, "2023-07-05")
	toDate, _ := time.Parse(requestfilter.DateFormate, "2023-07-06")

	expectedDateRange := bson.M{
		"$gte": fromDate,
		"$lt":  toDate.AddDate(0, 0, 1),
	}

	if dateRange, ok := filters["datefield"].(bson.M); ok {
		if !dateRange["$gte"].(time.Time).Equal(expectedDateRange["$gte"].(time.Time)) || !dateRange["$lt"].(time.Time).Equal(expectedDateRange["$lt"].(time.Time)) {
			t.Errorf("Unexpected value for datefield: got %v", filters["datefield"])
		}
	} else {
		t.Errorf("Unexpected type for datefield: got %T", filters["datefield"])
	}
}

func TestGenerateFiltersComma(t *testing.T) {
	// Adding test for parameters with comma-separated values
	filterFields := []requestfilter.FilterRequest{
		{Param: "strparam", Field: "strfield", Type: requestfilter.FieldTypeString},
		{Param: "intparam", Field: "intfield", Type: requestfilter.FieldTypeInt},
		{Param: "floatparam", Field: "floatfield", Type: requestfilter.FieldTypeFloat64},
	}

	getParamFunc := func(param string) string {
		switch param {
		case "strparam":
			return "value1,value2"
		case "intparam":
			return "42,43"
		case "floatparam":
			return "42.42,43.43"
		default:
			return ""
		}
	}

	filters := requestfilter.GenerateFilters(getParamFunc, filterFields)

	// Test for FieldTypeString
	if values, ok := filters["strfield"].(bson.M); ok {
		if valSlice, ok := values["$in"].([]interface{}); ok {
			if !contains(valSlice, "value1") || !contains(valSlice, "value2") {
				t.Errorf("Unexpected value for strfield: got %v", values["$in"])
			}
		} else {
			t.Errorf("Unexpected type for strfield")
		}
	} else {
		t.Errorf("Unexpected type for strfield: got")
	}

	// Test for FieldTypeInt
	if values, ok := filters["intfield"].(bson.M); ok {
		if valSlice, ok := values["$in"].([]interface{}); ok {
			if !contains(valSlice, 42) || !contains(valSlice, 43) {
				t.Errorf("Unexpected value for intfield: got %v", values["$in"])
			}
		} else {
			t.Errorf("Unexpected type for intfield")
		}
	} else {
		t.Errorf("Unexpected type for intfield: got %T", filters["intfield"])
	}

	// Test for FieldTypeFloat64
	if values, ok := filters["floatfield"].(bson.M); ok {
		if valSlice, ok := values["$in"].([]interface{}); ok {
			if !contains(valSlice, 42.42) || !contains(valSlice, 43.43) {
				t.Errorf("Unexpected value for floatfield: got %v", values["$in"])
			}
		} else {
			t.Errorf("Unexpected type for floatfield")
		}
	} else {
		t.Errorf("Unexpected type for floatfield: got %T", filters["floatfield"])
	}
}

// Helper function to check if a slice contains a certain value
func contains(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
