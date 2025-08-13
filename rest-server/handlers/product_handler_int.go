package handlers

import (
	"github.com/gin-gonic/gin"
)

type ProductHandlerInterface interface {
	CreateProduct(c *gin.Context)
	GetProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	ListProducts(c *gin.Context)
	GetProductsByCategory(c *gin.Context)
}
