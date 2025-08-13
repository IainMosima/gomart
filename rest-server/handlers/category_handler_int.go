package handlers

import (
	"github.com/gin-gonic/gin"
)

type CategoryHandlerInterface interface {
	CreateCategory(c *gin.Context)
	GetCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	ListCategories(c *gin.Context)
	GetCategoryChildren(c *gin.Context)
	GetCategoryAverageProductPrice(c *gin.Context)
}
