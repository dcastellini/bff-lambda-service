package container

import (
	"github.com/dcastellini/bff-lambda-service/internal/adapters"
	"github.com/gin-gonic/gin"
)

func startRouter(handler adapters.HTTPHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use()

	v1 := router.Group("/api/v1")
	{
		reminders := v1.Group("/products", errorHandlerMiddleware())
		{
			reminders.POST("", handler.CreateProductHandler())
			reminders.GET("", handler.GetProductsHandler())
			reminders.PUT("/:uid", handler.EditProductHandler())
			reminders.DELETE("/:uid", handler.DeleteProductHandler())
		}
	}
	return router
}
