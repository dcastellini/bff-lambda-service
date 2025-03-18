package ports

import (
	"context"
	"github.com/dcastellini/bff-lambda-service/internal/core/domain"
)

type ProductServicePort interface {
	CreateProduct(ctx context.Context, request *domain.CreateProductRequest) (domain.Product, error)
	EditProduct(ctx context.Context, request *domain.EditProductRequest) (domain.Product, error)
	DeleteProduct(ctx context.Context, billReminderID string) error
	GetProducts(ctx context.Context, clientID string) ([]domain.Product, error)
}
