package container

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/dcastellini/bff-lambda-service/internal/core/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func loggerMiddlewareFunc(next lambdaFunc, lambdaLogger Logger) lambdaFunc {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		ctxLogger := NewContextLogger(lambdaLogger)

		lCtx, log := ctxLogger.WithLogger(
			ctx,
			zap.String("user.username", getUserName(req)),
			zap.String("user.ua", req.RequestContext.Identity.UserAgent),
		)

		timeNow := time.Now()

		log.Info("start request")

		resp, err := next(lCtx, req)

		log.Info("end request", log.AnyField("object.requestTime", time.Since(timeNow).String()))

		return resp, err
	}
}

func getUserName(req events.APIGatewayProxyRequest) string {
	return req.RequestContext.Identity.UserAgent
}

func errorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		lastErr := ctx.Errors.Last()

		if lastErr != nil {
			ctxLogger := GetLogger(ctx.Request.Context())

			ctxLogger.Errorf("handler return error", lastErr)

			var customError domain.CustomError

			if !errors.As(lastErr, &customError) {
				customError = domain.BuildCustomError(lastErr)
			}

			ctx.JSON(customError.HTTPStatusCode, customError)
		}
	}
}
