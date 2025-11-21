package requestutil

import (
	"net/http"
	"strconv"
)

const (
	defaultItemsPerPage int64 = 10
	maxItemsPerPage     int64 = 100
)

func ParseQueryParamLimitOffset(r *http.Request) (int64, int64) {
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 0, 64)
	if err != nil || limit < 1 {
		limit = defaultItemsPerPage
	} else if limit > maxItemsPerPage {
		limit = maxItemsPerPage
	}

	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 0, 64)
	if err != nil || offset < 1 {
		offset = 0
	}

	return limit, offset
}
