package query

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type QueryBuilder struct {
	DB           *gorm.DB
	AllowedSorts map[string]string
	SearchFields []string
}

func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{
		DB:           db,
		AllowedSorts: make(map[string]string),
		SearchFields: []string{},
	}
}

func (qb *QueryBuilder) SetAllowedSorts(sorts map[string]string) *QueryBuilder {
	qb.AllowedSorts = sorts
	return qb
}

func (qb *QueryBuilder) SetSearchFields(fields []string) *QueryBuilder {
	qb.SearchFields = fields
	return qb
}

func (qb *QueryBuilder) ApplyFiltering(params *QueryParams) *gorm.DB {
	query := qb.DB

	if params.Search != "" && len(qb.SearchFields) > 0 {

		searchQuery := ""
		searchValues := []any{}

		for i, field := range qb.SearchFields {
			if i > 0 {
				searchQuery += " OR "
			}

			searchQuery += field + " ILIKE ?"
			searchValues = append(searchValues, "%"+params.Search+"%")
		}

		query = query.Where(searchQuery, searchValues...)
	}

	for key, value := range params.Filters {
		if value != "" {
			switch {
			case strings.HasSuffix(key, "_like"):
				query = query.Where(fmt.Sprintf("%s ILIKE ?", key), "%"+value+"%")
			default:
				query = query.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}

	return query
}

func (qb *QueryBuilder) ApplySorting(query *gorm.DB, params *QueryParams) *gorm.DB {
	sortField := params.Sort

	if dbField, exists := qb.AllowedSorts[sortField]; exists {
		sortField = dbField
	} else if len(qb.AllowedSorts) > 0 {
		for _, allowedField := range qb.AllowedSorts {
			sortField = allowedField
			break
		}
	}

	return query.Order(fmt.Sprintf("%s %s", sortField, strings.ToUpper(params.Order)))
}

func (qb *QueryBuilder) ApplyPagination(query *gorm.DB, params *QueryParams) *gorm.DB {
	offset := (params.Page - 1) * params.Limit
	return query.Offset(offset).Limit(params.Limit)
}

func (qb *QueryBuilder) BuildQuery(params *QueryParams) *gorm.DB {
	query := qb.ApplyFiltering(params)
	query = qb.ApplySorting(query, params)
	return query
}
