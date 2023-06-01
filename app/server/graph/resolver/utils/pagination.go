package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/AnthonyThomasson/graph-ql/graph/model"
	"gorm.io/gorm"
)

type Connection[T any] struct {
	Total    int             `json:"Total"`
	Edges    []Edge[T]       `json:"Edges"`
	PageInfo *model.PageInfo `json:"PageInfo"`
}

type Edge[T any] struct {
	Cursor string `json:"Cursor"`
	Node   *T     `json:"Node"`
}

func ExecutePagination[T any](query *gorm.DB, page *model.PaginationInput, order *model.Order, obj T) (*Connection[T], error) {
	query = query.Debug()

	if page.First != nil && page.Before != nil {
		return nil, errors.New("cannot use first and before at the same time")
	}

	if page.Last != nil && page.After != nil {
		return nil, errors.New("cannot use last and after at the same time")
	}

	if order.Field == "" {
		order.Field = "created_at"
	}
	if order.Direction == "" {
		order.Direction = "ASC"
	}

	var models *[]*T
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	hasNextPage := false
	hasPreviousPage := false
	if page.Last != nil {
		query, hasNextPage, hasPreviousPage = executeOffsetPagination(query, total, page)
	} else if page.Before != nil || page.After != nil {
		query, _ = executeCursorPagination[T](query, page, order, obj)
	}

	query = query.Limit(getLimit(page.First))

	query.Find(&models)
	Connection := &Connection[T]{
		Total: int(total),
		Edges: toEdges(*models, order.Field),
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     encodeCursor(getStartCursor[T](*models, order.Field)),
			EndCursor:       encodeCursor(getEndCursor[T](*models, order.Field)),
		},
	}

	return Connection, nil
}

func toEdges[T any](models []*T, orderedByField string) []Edge[T] {
	var edges []Edge[T]
	for _, model := range models {
		edges = append(edges, Edge[T]{
			Cursor: encodeCursor(getCursorValue(model, orderedByField)),
			Node:   model,
		})
	}
	return edges
}

func executeOffsetPagination(query *gorm.DB, total int64, page *model.PaginationInput) (*gorm.DB, bool, bool) {
	offset := getOffsetFromLast(int(total), page.Last)
	query.Offset(offset)

	hasNextPage := offset+getLimit(page.Last) < int(total)

	hasPreviousPage := offset > 0

	return query, hasNextPage, hasPreviousPage
}

func executeCursorPagination[T any](query *gorm.DB, page *model.PaginationInput, order *model.Order, obj T) (*gorm.DB, error) {
	field, _ := getFieldInfoByTag(obj, order.Field)
	stmt := &gorm.Statement{DB: query}
	stmt.Parse(&obj)
	tableName := stmt.Schema.Table
	if page.After != nil {
		cursor, err := decodeCursor(*page.After)
		if err != nil {
			return nil, err
		}
		if order.Direction == "DESC" {
			query = query.Where(fmt.Sprintf("%s.%s < ?", tableName, field), cursor)
		} else {
			query = query.Where(fmt.Sprintf("%s.%s > ?", tableName, field), cursor)
		}

	} else if page.Before != nil {
		cursor, err := decodeCursor(*page.Before)
		if err != nil {
			return nil, err
		}
		if order.Direction == "DESC" {
			query = query.Where(fmt.Sprintf("%s.%s > ?", tableName, field), cursor)
		} else {
			query = query.Where(fmt.Sprintf("%s.%s < ?", tableName, field), cursor)
		}
	}
	return query, nil
}

func getOffsetFromLast(count int, last *int) int {
	if last != nil {
		return count - *last
	}
	return 0
}

func getLimit(first *int) int {
	if first != nil {
		return *first
	}
	return 10
}

func decodeCursor(cursor string) (string, error) {
	val, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func encodeCursor(cursor string) string {
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}

func getStartCursor[T any](model []*T, orderedByField string) string {
	if len(model) == 0 {
		return ""
	}
	return getCursorValue(model[0], orderedByField)
}

func getCursorValue[T any](model T, cursorName string) string {
	_, value := getFieldInfoByTag(model, cursorName)
	return value
}

func getEndCursor[T any](models []*T, orderedByField string) string {
	if len(models) == 0 {
		return ""
	}
	return getCursorValue(models[len(models)-1], orderedByField)
}

func getFieldInfoByTag(structValue any, fieldName string) (string, string) {
	value := reflect.ValueOf(structValue)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return "", ""
	}
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		jsonFieldName := strings.Split(tag, ",")[0]
		if strings.EqualFold(jsonFieldName, fieldName) {
			field := value.Field(i)
			return jsonFieldName, fmt.Sprintf("%v", field)
		}
	}

	return "", ""
}
