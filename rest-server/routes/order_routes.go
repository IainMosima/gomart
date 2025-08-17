package routes

import (
	"github.com/IainMosima/gomart/rest-server/handlers/order"
	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.Engine, orderHandler order.OrderHandlerInterface) {
	orders := router.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("/:id/status", orderHandler.GetOrderStatus)
	}
}
