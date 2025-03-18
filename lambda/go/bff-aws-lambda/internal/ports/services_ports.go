package ports

import (
	"context"
	"github.com/dcastellini/bff-lambda-service/internal/core/domain"
)

type ProductsBFFService interface {
	CreateProduct(ctx context.Context, request *domain.CreateProductRequest) (*domain.CreateProductResponse, error)
	EditProduct(ctx context.Context, request *domain.EditProductRequest) (*domain.EditProductResponse, error)
	DeleteProduct(ctx context.Context, request *domain.DeleteProductRequest) (*domain.DeleteProductResponse, error)
	GetProducts(ctx context.Context, request *domain.GetProductsRequest) (*domain.GetProductsResponse, error)
}
