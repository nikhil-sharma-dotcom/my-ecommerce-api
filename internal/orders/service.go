package orders

import (
	"context"
	"errors"

	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/adapters/postgresql/sqlc"
)

type Service struct {
	queries *sqlc.Queries
}

func NewService(queries *sqlc.Queries) *Service {
	return &Service{queries: queries}
}

type OrderItemInput struct {
	ProductID int64
	Quantity  int32
}

func (s *Service) Create(ctx context.Context, customerID int64, items []OrderItemInput) (sqlc.Order, error) {
	if len(items) == 0 {
		return sqlc.Order{}, errors.New("order must contain at least one item")
	}

	var totalCents int64
	for _, item := range items {
		product, err := s.queries.GetProduct(ctx, item.ProductID)
		if err != nil {
			return sqlc.Order{}, errors.New("product not found")
		}
		if product.Quantity < item.Quantity {
			return sqlc.Order{}, errors.New("insufficient stock")
		}

		err = s.queries.UpdateProductQuantity(ctx, sqlc.UpdateProductQuantityParams{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		})
		if err != nil {
			return sqlc.Order{}, err
		}

		totalCents += int64(product.PriceInCents) * int64(item.Quantity)
	}

	order, err := s.queries.CreateOrder(ctx, customerID)
	if err != nil {
		return sqlc.Order{}, err
	}

	for _, item := range items {
		product, _ := s.queries.GetProduct(ctx, item.ProductID)
		_, err = s.queries.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
		})
		if err != nil {
			return sqlc.Order{}, err
		}
	}

	return order, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (sqlc.Order, error) {
	return s.queries.GetOrder(ctx, id)
}

func (s *Service) GetByCustomerID(ctx context.Context, customerID int64) ([]sqlc.Order, error) {
	return s.queries.ListOrdersByCustomer(ctx, customerID)
}

func (s *Service) GetOrderItems(ctx context.Context, orderID int64) ([]sqlc.OrderItem, error) {
	return s.queries.GetOrderItems(ctx, orderID)
}
