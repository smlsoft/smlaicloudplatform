package utils

import (
	"smlcloudplatform/pkg/models"
	"strconv"
	"strings"
)

const RequestSearchMaxPage = 2147483647
const RequestSearchMinPage = 1
const RequestSearchMaxLimit = 1000
const RequestSearchMinLimit = 1

const RequestSearchDefaultPage = 1
const RequestSearchDefaultLimit = 20

const RequestMaxOffset = 2147483647
const RequestMinOffset = 0

func GetSearchQueryParam(fnGetParam func(string) string) string {
	q := strings.Trim(fnGetParam("q"), " ")
	return q
}

func GetPaginationParam(fnGetParam func(string) string) (int, int) {

	pageRawText := fnGetParam("page")
	limitRawText := fnGetParam("limit")

	page, err := strconv.Atoi(pageRawText) //strconv.ParseUint(pageRawText, 10, 32)
	if err != nil {
		page = RequestSearchDefaultPage
	}

	limit, err := strconv.Atoi(limitRawText) //strconv.ParseUint(limitRawText, 10, 32)

	if err != nil {
		limit = RequestSearchDefaultLimit
	}

	if page < 0 {
		page = RequestSearchMinPage
	}

	if page > RequestSearchMaxPage {
		page = RequestSearchMaxPage
	}

	if limit < RequestSearchMinLimit {
		limit = RequestSearchMinLimit
	}

	if limit > RequestSearchMaxLimit {
		limit = RequestSearchMaxLimit
	}

	return page, limit
}

func GetParamOffsetLimit(fnGetParam func(string) string) (int, int) {

	pageRawText := fnGetParam("offset")
	limitRawText := fnGetParam("limit")

	offset, err := strconv.Atoi(pageRawText) //strconv.ParseUint(pageRawText, 10, 32)
	if err != nil {
		offset = RequestSearchDefaultPage
	}

	limit, err := strconv.Atoi(limitRawText) //strconv.ParseUint(limitRawText, 10, 32)

	if err != nil {
		limit = RequestSearchDefaultLimit
	}

	if offset < 0 {
		offset = RequestMinOffset
	}

	if offset > RequestSearchMaxPage {
		offset = RequestMaxOffset
	}

	if limit < RequestSearchMinLimit {
		limit = RequestSearchMinLimit
	}

	if limit > RequestSearchMaxLimit {
		limit = RequestSearchMaxLimit
	}

	return offset, limit
}

func GetSortParam(fnGetParam func(string) string) map[string]int {
	tempSort := make(map[string]int)

	sortRawText := strings.Trim(fnGetParam("sort"), " ")

	if sortRawText == "" {
		return tempSort
	}

	sortSplitArr := strings.Split(sortRawText, ",")

	for _, sortFieldRaw := range sortSplitArr {
		sortFieldRawArr := strings.Split(sortFieldRaw, ":")
		if len(sortFieldRawArr) != 2 {
			continue
		}

		sortKey := sortFieldRawArr[0]
		sortValRaw := sortFieldRawArr[1]

		sortVal, err := strconv.Atoi(sortValRaw)
		if err != nil {
			sortVal = 1
		}

		tempSort[sortKey] = sortVal
	}

	return tempSort
}

func GetSearchParam(fnGetParam func(string) string) (string, int, int, map[string]int) {
	q := GetSearchQueryParam(fnGetParam)
	page, limit := GetPaginationParam(fnGetParam)
	sorts := GetSortParam(fnGetParam)

	return q, page, limit, sorts
}

func GetSearchPageable(fnGetParam func(string) string) models.Pageable {
	q := GetSearchQueryParam(fnGetParam)
	page, limit := GetPaginationParam(fnGetParam)
	sorts := GetSortParam(fnGetParam)

	return models.Pageable{
		Q:     q,
		Page:  page,
		Limit: limit,
		Sorts: sorts,
	}
}
