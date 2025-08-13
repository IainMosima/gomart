package routes

import (
	"github.com/IainMosima/gomart/rest-server/handlers"
	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(router *gin.Engine, productHandler handlers.ProductHandlerInterface) {
	products := router.Group("/products")
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.ListProducts)
		products.GET("/:id", productHandler.GetProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
		products.GET("/category/:categoryId", productHandler.GetProductsByCategory)
	}
}
