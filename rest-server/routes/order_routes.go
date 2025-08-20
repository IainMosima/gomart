package routes

import (
	"github.com/IainMosima/gomart/rest-server/handlers/order"
	"github.com/IainMosima/gomart/rest-server/middleware"
	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.Engine, orderHandler order.OrderHandlerInterface, authMiddleware *middleware.AuthMiddleware) {
	orders := router.Group("/orders")
	{
		orders.POST("", authMiddleware.RequireAuth(), orderHandler.CreateOrder)
		orders.GET("/:id/status", orderHandler.GetOrderStatus)
	}
}
