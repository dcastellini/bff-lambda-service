package arg

import (
	"github.com/dcastellini/bff-lambda-service/internal/core/domain"
	"github.com/dcastellini/bff-lambda-service/internal/ports"
)

type ProductBffArgService struct {
	productApi ports.ProductServicePort
}

func NewProductBffArgService(productApi ports.ProductServicePort) ports.ProductServicePort {
	return &ProductBffArgService{
		productApi: productApi,
	}
}

func (p ProductBffArgService) CreateProduct(ctx context.Context, request *domain.CreateProductRequest) (domain.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductBffArgService) EditProduct(ctx context.Context, request *domain.EditProductRequest) (domain.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductBffArgService) DeleteProduct(ctx context.Context, billReminderID string) error {
	//TODO implement me
	panic("implement me")
}

func (p ProductBffArgService) GetProducts(ctx context.Context, clientID string) ([]domain.Product, error) {
	//TODO implement me
	panic("implement me")
}
