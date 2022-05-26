package utils

import "strconv"

const RequestSearchMaxPage = 1
const RequestSearchMinPage = 1
const RequestSearchMaxLimit = 1000
const RequestSearchMinLimit = 1

const RequestSearchDefaultPage = 1
const RequestSearchDefaultLimit = 20

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
