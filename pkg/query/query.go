package query

import (
	"net/http"
	"strconv"
)

var (
	paramLimit  = "limit"
	paramOffset = "offset"

	defaultLimit = 10
	maxLimit     = 100
)

func ParseLimit(r *http.Request) int {
	limit, err := strconv.Atoi(r.URL.Query().Get(paramLimit))
	if err != nil || limit == 0 {
		limit = defaultLimit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	return limit
}

func ParseOffset(r *http.Request) int {
	offset, err := strconv.Atoi(r.URL.Query().Get(paramOffset))
	if err != nil {
		return 0
	}

	return offset
}
