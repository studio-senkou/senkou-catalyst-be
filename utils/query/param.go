package query

import (
	"slices"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type QueryParams struct {
	Page    int               `json:"page"`    // current page number
	Limit   int               `json:"limit"`   // number of items per page
	Sort    string            `json:"sort"`    // sort field
	Order   string            `json:"order"`   // order `asc` for ascending and `desc` for descending
	Search  string            `json:"search"`  // search keyword
	Filters map[string]string `json:"filters"` // additional filters
}

type PaginationResponse struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int   `json:"items_per_page"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

func ParseQueryParams(c *fiber.Ctx) *QueryParams {
	params := &QueryParams{
		Page:    1,
		Limit:   10,
		Sort:    "created_at",
		Order:   "desc",
		Filters: make(map[string]string),
	}

	// pagination params
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if p, err := strconv.Atoi(limit); err == nil && p > 0 {
			params.Limit = p
		}
	}

	// sorting params
	// currently only supports for single sorting on `sort` query param
	// there is also support for order sorting on `order` query param
	if sort := c.Query("sort"); sort != "" {
		params.Sort = sort
	}
	if order := c.Query("order"); order != "" {
		params.Order = order
	}

	// search params
	// currently only supports for single search on `search` query param
	params.Search = c.Query("search")

	// filters params
	for key, value := range c.Queries() {
		if !isReservedParam(key) && value != "" {
			params.Filters[key] = value
		}
	}

	return params
}

func isReservedParam(key string) bool {
	reserved := []string{"page", "limit", "sort", "order", "search"}
	return slices.Contains(reserved, key)
}

func CalculatePagination(page, limit int, total int64) *PaginationResponse {
	pages := int((total + int64(limit) - 1) / int64(limit))

	return &PaginationResponse{
		CurrentPage:  page,
		TotalPages:   pages,
		TotalItems:   total,
		ItemsPerPage: limit,
		HasNext:      page < pages,
		HasPrev:      page > 1,
	}
}
