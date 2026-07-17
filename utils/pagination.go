package utils

import "strconv"

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Page is a generic paginated result: the items on this page plus the window
// they came from. HasMore reports whether another page exists.
type Page[T any] struct {
	Items   []T  `json:"items"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Count   int  `json:"count"`
	HasMore bool `json:"has_more"`
}

// NewPage builds a Page from the rows fetched with a one-extra-row lookahead:
// pass the rows queried with limit+1, and it trims to the page and sets HasMore.
func NewPage[T any](rows []T, limit, offset int) Page[T] {
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	if rows == nil {
		rows = []T{}
	}
	return Page[T]{
		Items:   rows,
		Limit:   limit,
		Offset:  offset,
		Count:   len(rows),
		HasMore: hasMore,
	}
}

// PageParams reads limit and offset from query params, applying defaults and a
// maximum page size.
func PageParams(params map[string][]string) (limit, offset int) {
	limit, offset = DefaultPageSize, 0
	if n, err := strconv.Atoi(firstValue(params, "limit")); err == nil && n > 0 {
		limit = n
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	if n, err := strconv.Atoi(firstValue(params, "offset")); err == nil && n > 0 {
		offset = n
	}
	return limit, offset
}

func firstValue(params map[string][]string, key string) string {
	if v := params[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}
