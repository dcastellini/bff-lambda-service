package container

import (
	"github.com/dcastellini/bff-lambda-service/internal/adapters"
	"github.com/dcastellini/bff-lambda-service/internal/config"
	"github.com/dcastellini/bff-lambda-service/internal/core/domain"
	"os"
	"strings"
)

const (
	argCountryCode = "ARG"
)

func startHTTPHandlerByCountry(cfg *config.LambdaConfiguration) adapters.HTTPHandler {
	switch strings.ToUpper(cfg.General.Country) {
	case argCountryCode:
		return newArgHTTPHandler(cfg)
	default:
		customError := domain.ErrNotImplemented
		os.Exit(customError.MessageCode)
		return adapters.HTTPHandler{}
	}
}

func newArgHTTPHandler(cfg *config.LambdaConfiguration) adapters.HTTPHandler {

	productServiceAdapter := argAdapters.NewProductServiceAdapter()
	return adapters.NewHTTPHandler(productServiceAdapter)
}
