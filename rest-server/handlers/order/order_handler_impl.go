package order

import (
	"net/http"

	"github.com/IainMosima/gomart/domains/order/schema"
	"github.com/IainMosima/gomart/domains/order/service"
	"github.com/IainMosima/gomart/rest-server/dtos"
	"github.com/IainMosima/gomart/rest-server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandlerImpl struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) OrderHandlerInterface {
	return &OrderHandlerImpl{
		orderService: orderService,
	}
}

func (h *OrderHandlerImpl) CreateOrder(c *gin.Context) {
	var req dtos.CreateOrderRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if req.CustomerID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot create order for another customer"})
		return
	}

	items := make([]schema.CreateOrderItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = schema.CreateOrderItemRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	schemaReq := &schema.CreateOrderRequest{
		CustomerID: req.CustomerID,
		Items:      items,
	}

	result, err := h.orderService.CreateOrder(c.Request.Context(), schemaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &dtos.OrderResponseDTO{
		OrderID:     result.OrderID,
		CustomerID:  result.CustomerID,
		OrderNumber: result.OrderNumber,
		Status:      string(result.Status),
		TotalAmount: result.TotalAmount,
		CreatedAt:   result.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *OrderHandlerImpl) GetOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	result, err := h.orderService.GetOrderStatus(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := &dtos.OrderStatusResponseDTO{
		OrderID:     result.OrderID,
		OrderNumber: result.OrderNumber,
		Status:      string(result.Status),
		CreatedAt:   result.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}
