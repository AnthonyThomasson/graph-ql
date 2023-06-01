package resolver

import (
	"context"

	"github.com/AnthonyThomasson/graph-ql/graph/model"
	"github.com/AnthonyThomasson/graph-ql/graph/resolver/utils"
)

func (r *mutationResolver) createProduct(ctx context.Context, input model.NewProduct) (*model.Product, error) {
	product := &model.Product{
		Name:  input.Name,
		Price: input.Price,
		Rank:  int32(input.Rank),
	}
	err := r.Resolver.DB.Create(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *Resolver) products(ctx context.Context, page *model.PaginationInput, order *model.Order) (*model.ProductConnection, error) {
	query := r.DB.Model(&model.Product{})
	query, orderByField := utils.Order(query, order, "created_at")

	connection, err := utils.ExecutePagination[model.Product](query, page, order, model.Product{})
	if err != nil {
		return nil, err
	}

	response := &model.ProductConnection{
		Total:    connection.Total,
		Edges:    toProductEdges(connection.Edges, orderByField),
		PageInfo: connection.PageInfo,
	}

	return response, nil
}

func toProductEdges(products []utils.Edge[model.Product], orderByField string) []*model.ProductEdge {
	var edges []*model.ProductEdge
	for _, product := range products {
		edges = append(edges, &model.ProductEdge{
			Cursor: product.Cursor,
			Node:   product.Node,
		})
	}
	return edges
}
