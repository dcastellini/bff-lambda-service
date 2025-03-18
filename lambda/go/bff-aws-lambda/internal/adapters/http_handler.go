package adapters

import "github.com/dcastellini/bff-lambda-service/internal/ports"

type HTTPHandler struct {
	ProductsBFFService ports.ProductsBFFService
}

func NewHTTPHandler(productsService ports.ProductsBFFService) HTTPHandler {
	return HTTPHandler{
		ProductsBFFService: productsService,
	}

}
