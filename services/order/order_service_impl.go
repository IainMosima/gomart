package order

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	authRepo "github.com/IainMosima/gomart/domains/auth/repository"
	"github.com/IainMosima/gomart/domains/order/entity"
	orderRepo "github.com/IainMosima/gomart/domains/order/repository"
	"github.com/IainMosima/gomart/domains/order/schema"
	"github.com/IainMosima/gomart/domains/order/service"
	productRepo "github.com/IainMosima/gomart/domains/product/repository"
	"github.com/google/uuid"
)

type OrderServiceImpl struct {
	orderRepo    orderRepo.OrderRepository
	productsRepo productRepo.ProductRepository
	authRepo     authRepo.AuthRepository
	notification service.NotificationService
}

func NewOrderServiceImpl(orderRepository orderRepo.OrderRepository, productsRepo productRepo.ProductRepository, authRepository authRepo.AuthRepository, notificationService service.NotificationService) service.OrderService {
	return &OrderServiceImpl{
		orderRepo:    orderRepository,
		productsRepo: productsRepo,
		authRepo:     authRepository,
		notification: notificationService,
	}
}

func (o *OrderServiceImpl) CreateOrder(ctx context.Context, req *schema.CreateOrderRequest) (*schema.OrderResponse, error) {
	_, err := o.authRepo.GetUserByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	var totalAmount float64
	orderItems := make([]*entity.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		product, err := o.productsRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found: %w", item.ProductID, err)
		}

		unitPrice := product.Price
		itemTotal := unitPrice * float64(item.Quantity)
		totalAmount += itemTotal

		orderItem := &entity.OrderItem{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  unitPrice,
			TotalPrice: itemTotal,
		}
		orderItems = append(orderItems, orderItem)
	}

	orderNumber := o.generateOrderNumber()

	order := &entity.Order{
		OrderID:     uuid.New(),
		CustomerID:  req.CustomerID,
		OrderNumber: orderNumber,
		Status:      entity.OrderStatusConfirmed,
		TotalAmount: totalAmount,
	}

	createdOrder, err := o.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range orderItems {
		item.OrderID = createdOrder.OrderID
		item.OrderItemID = uuid.New()
		_, err := o.orderRepo.CreateOrderItem(ctx, item)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	orderResponse := o.orderToResponse(createdOrder)

	if err := o.notification.SendEmail(ctx, orderResponse); err != nil {
		fmt.Printf("Failed to send email notification: %v\n", err)
	}

	if err := o.notification.SendSMS(ctx, orderResponse); err != nil {
		fmt.Printf("Failed to send SMS notification: %v\n", err)
	}

	return orderResponse, nil
}

func (o *OrderServiceImpl) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (*schema.OrderStatusResponse, error) {
	order, err := o.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	return &schema.OrderStatusResponse{
		OrderID:     order.OrderID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
	}, nil
}

func (o *OrderServiceImpl) generateOrderNumber() string {
	timestamp := time.Now().Unix()
	random := rand.Intn(999) + 1
	return fmt.Sprintf("ORD-%03d", (timestamp%1000000+int64(random))%1000)
}

func (o *OrderServiceImpl) orderToResponse(order *entity.Order) *schema.OrderResponse {
	return &schema.OrderResponse{
		OrderID:     order.OrderID,
		CustomerID:  order.CustomerID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
	}
}
