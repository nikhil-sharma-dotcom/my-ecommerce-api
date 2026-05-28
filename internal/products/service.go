package products

import (
	"context"

	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/adapters/postgresql/sqlc"
)

type Service struct {
	queries *sqlc.Queries
}

func NewService(queries *sqlc.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) Create(ctx context.Context, name string, priceInCents int32, quantity int32) (sqlc.Product, error) {
	return s.queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Name:         name,
		PriceInCents: priceInCents,
		Quantity:     quantity,
	})
}

func (s *Service) GetByID(ctx context.Context, id int64) (sqlc.Product, error) {
	return s.queries.GetProduct(ctx, id)
}

func (s *Service) GetAll(ctx context.Context) ([]sqlc.Product, error) {
	return s.queries.ListProducts(ctx)
}

func (s *Service) Update(ctx context.Context, id int64, name string, priceInCents int32, quantity int32) error {
	return s.queries.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:           id,
		Name:         name,
		PriceInCents: priceInCents,
		Quantity:     quantity,
	})
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.queries.DeleteProduct(ctx, id)
}
