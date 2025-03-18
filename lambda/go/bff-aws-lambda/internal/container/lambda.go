package container

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/dcastellini/bff-lambda-service/internal/config"
)

type lambdaFunc func(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error)

type Lambda struct {
	httpHandlerFunc lambdaFunc
}

func Start() {
	initializedLambda := initialize()
	lambda.Start(initializedLambda.Handler)
}

func initialize() Lambda {
	cfgLambda := loadLambdaConfiguration()

	httpHandler := startHTTPHandlerByCountry(cfgLambda)
	ginEngine := startRouter(httpHandler)
	ginAdapter := ginAdapter.New(ginEngine)

	return Lambda{
		httpHandlerFunc: loggerMiddlewareFunc(
			ginAdapter.ProxyWithContext,
		),
	}
}

func loadLambdaConfiguration() *config.LambdaConfiguration {
	cfgLambda := config.NewConfigLambda()
	cfgLambda.LoadFromEnvs()
	return cfgLambda
}

func (h *Lambda) Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return h.httpHandlerFunc(ctx, req)
}
