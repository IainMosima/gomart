package routes

import (
	"github.com/IainMosima/gomart/rest-server/handlers"
	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(router *gin.Engine, categoryHandler handlers.CategoryHandlerInterface) {
	categories := router.Group("/categories")
	{
		categories.POST("", categoryHandler.CreateCategory)
		categories.GET("", categoryHandler.ListCategories)
		categories.GET("/:id", categoryHandler.GetCategory)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
		categories.GET("/:id/children", categoryHandler.GetCategoryChildren)
		categories.GET("/:id/average-price", categoryHandler.GetCategoryAverageProductPrice)
	}
}
