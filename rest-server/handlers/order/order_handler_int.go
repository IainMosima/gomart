package order

import (
	"github.com/gin-gonic/gin"
)

type OrderHandlerInterface interface {
	CreateOrder(c *gin.Context)
	GetOrderStatus(c *gin.Context)
}
