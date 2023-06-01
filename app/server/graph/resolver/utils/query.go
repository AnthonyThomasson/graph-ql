package utils

import (
	"github.com/AnthonyThomasson/graph-ql/graph/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Order(query *gorm.DB, order *model.Order, defaultOrder string) (*gorm.DB, string) {
	orderField := defaultOrder
	if order != nil && order.Field != "" {
		orderField = order.Field
	}

	query.Order(clause.OrderByColumn{
		Column: clause.Column{Name: orderField},
		Desc:   order.Direction == "DESC",
	})
	return query, orderField
}
