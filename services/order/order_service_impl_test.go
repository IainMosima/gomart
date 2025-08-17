package order

import (
	"context"
	"errors"
	"testing"
	"time"

	authEntity "github.com/IainMosima/gomart/domains/auth/entity"
	authRepo "github.com/IainMosima/gomart/domains/auth/repository"
	"github.com/IainMosima/gomart/domains/order/entity"
	orderRepo "github.com/IainMosima/gomart/domains/order/repository"
	"github.com/IainMosima/gomart/domains/order/schema"
	"github.com/IainMosima/gomart/domains/order/service"
	productEntity "github.com/IainMosima/gomart/domains/product/entity"
	productRepo "github.com/IainMosima/gomart/domains/product/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderServiceImpl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := orderRepo.NewMockOrderRepository(ctrl)
	mockProductRepo := productRepo.NewMockProductRepository(ctrl)
	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)
	mockNotificationService := service.NewMockNotificationService(ctrl)

	orderService := NewOrderServiceImpl(mockOrderRepo, mockProductRepo, mockAuthRepo, mockNotificationService)

	assert.NotNil(t, orderService)
	assert.IsType(t, &OrderServiceImpl{}, orderService)
}

func TestOrderServiceImpl_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := orderRepo.NewMockOrderRepository(ctrl)
	mockProductRepo := productRepo.NewMockProductRepository(ctrl)
	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)
	mockNotificationService := service.NewMockNotificationService(ctrl)
	orderService := NewOrderServiceImpl(mockOrderRepo, mockProductRepo, mockAuthRepo, mockNotificationService)

	ctx := context.Background()
	now := time.Now()
	customerID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()

	tests := []struct {
		name    string
		request *schema.CreateOrderRequest
		setup   func()
		want    *schema.OrderResponse
		wantErr bool
	}{
		{
			name: "successful order creation",
			request: &schema.CreateOrderRequest{
				CustomerID: customerID,
				Items: []schema.CreateOrderItemRequest{
					{
						ProductID: productID,
						Quantity:  2,
					},
				},
			},
			setup: func() {
				// Mock customer validation
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(&authEntity.Customer{
						UserID: customerID,
						Email:  "test@example.com",
					}, nil)

				// Mock product validation
				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(&productEntity.Product{
						ProductID: productID,
						Price:     100.00,
					}, nil)

				// Mock order creation
				mockOrderRepo.EXPECT().
					CreateOrder(ctx, gomock.Any()).
					Return(&entity.Order{
						OrderID:     orderID,
						CustomerID:  customerID,
						OrderNumber: "ORD-001",
						Status:      entity.OrderStatusConfirmed,
						TotalAmount: 200.00,
						CreatedAt:   now,
					}, nil)

				// Mock order item creation
				mockOrderRepo.EXPECT().
					CreateOrderItem(ctx, gomock.Any()).
					Return(&entity.OrderItem{
						OrderItemID: uuid.New(),
						OrderID:     orderID,
						ProductID:   productID,
						Quantity:    2,
						UnitPrice:   100.00,
						TotalPrice:  200.00,
					}, nil)

				// Mock notifications (should not fail order creation)
				mockNotificationService.EXPECT().
					SendEmail(ctx, gomock.Any()).
					Return(nil)

				mockNotificationService.EXPECT().
					SendSMS(ctx, gomock.Any()).
					Return(nil)
			},
			want: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-001",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 200.00,
				CreatedAt:   now,
			},
			wantErr: false,
		},
		{
			name: "customer not found",
			request: &schema.CreateOrderRequest{
				CustomerID: customerID,
				Items: []schema.CreateOrderItemRequest{
					{
						ProductID: productID,
						Quantity:  1,
					},
				},
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(nil, errors.New("customer not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "product not found",
			request: &schema.CreateOrderRequest{
				CustomerID: customerID,
				Items: []schema.CreateOrderItemRequest{
					{
						ProductID: productID,
						Quantity:  1,
					},
				},
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(&authEntity.Customer{
						UserID: customerID,
						Email:  "test@example.com",
					}, nil)

				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(nil, errors.New("product not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "order creation fails",
			request: &schema.CreateOrderRequest{
				CustomerID: customerID,
				Items: []schema.CreateOrderItemRequest{
					{
						ProductID: productID,
						Quantity:  1,
					},
				},
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(&authEntity.Customer{
						UserID: customerID,
						Email:  "test@example.com",
					}, nil)

				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(&productEntity.Product{
						ProductID: productID,
						Price:     50.00,
					}, nil)

				mockOrderRepo.EXPECT().
					CreateOrder(ctx, gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "order item creation fails",
			request: &schema.CreateOrderRequest{
				CustomerID: customerID,
				Items: []schema.CreateOrderItemRequest{
					{
						ProductID: productID,
						Quantity:  1,
					},
				},
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(&authEntity.Customer{
						UserID: customerID,
						Email:  "test@example.com",
					}, nil)

				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(&productEntity.Product{
						ProductID: productID,
						Price:     50.00,
					}, nil)

				mockOrderRepo.EXPECT().
					CreateOrder(ctx, gomock.Any()).
					Return(&entity.Order{
						OrderID:     orderID,
						CustomerID:  customerID,
						OrderNumber: "ORD-002",
						Status:      entity.OrderStatusConfirmed,
						TotalAmount: 50.00,
						CreatedAt:   now,
					}, nil)

				mockOrderRepo.EXPECT().
					CreateOrderItem(ctx, gomock.Any()).
					Return(nil, errors.New("order item creation failed"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "successful order with notification failures",
			request: &schema.CreateOrderRequest{
				CustomerID: customerID,
				Items: []schema.CreateOrderItemRequest{
					{
						ProductID: productID,
						Quantity:  1,
					},
				},
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(&authEntity.Customer{
						UserID: customerID,
						Email:  "test@example.com",
					}, nil)

				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(&productEntity.Product{
						ProductID: productID,
						Price:     75.00,
					}, nil)

				mockOrderRepo.EXPECT().
					CreateOrder(ctx, gomock.Any()).
					Return(&entity.Order{
						OrderID:     orderID,
						CustomerID:  customerID,
						OrderNumber: "ORD-003",
						Status:      entity.OrderStatusConfirmed,
						TotalAmount: 75.00,
						CreatedAt:   now,
					}, nil)

				mockOrderRepo.EXPECT().
					CreateOrderItem(ctx, gomock.Any()).
					Return(&entity.OrderItem{
						OrderItemID: uuid.New(),
						OrderID:     orderID,
						ProductID:   productID,
						Quantity:    1,
						UnitPrice:   75.00,
						TotalPrice:  75.00,
					}, nil)

				// Notifications fail but should not affect order creation
				mockNotificationService.EXPECT().
					SendEmail(ctx, gomock.Any()).
					Return(errors.New("email service unavailable"))

				mockNotificationService.EXPECT().
					SendSMS(ctx, gomock.Any()).
					Return(errors.New("SMS service unavailable"))
			},
			want: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-003",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 75.00,
				CreatedAt:   now,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := orderService.CreateOrder(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestOrderServiceImpl_GetOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := orderRepo.NewMockOrderRepository(ctrl)
	mockProductRepo := productRepo.NewMockProductRepository(ctrl)
	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)
	mockNotificationService := service.NewMockNotificationService(ctrl)
	orderService := NewOrderServiceImpl(mockOrderRepo, mockProductRepo, mockAuthRepo, mockNotificationService)

	ctx := context.Background()
	orderID := uuid.New()
	now := time.Now()

	tests := []struct {
		name    string
		orderID uuid.UUID
		setup   func()
		want    *schema.OrderStatusResponse
		wantErr bool
	}{
		{
			name:    "successful status retrieval",
			orderID: orderID,
			setup: func() {
				mockOrderRepo.EXPECT().
					GetOrder(ctx, orderID).
					Return(&entity.Order{
						OrderID:     orderID,
						OrderNumber: "ORD-001",
						Status:      entity.OrderStatusConfirmed,
						CreatedAt:   now,
					}, nil)
			},
			want: &schema.OrderStatusResponse{
				OrderID:     orderID,
				OrderNumber: "ORD-001",
				Status:      entity.OrderStatusConfirmed,
				CreatedAt:   now,
			},
			wantErr: false,
		},
		{
			name:    "order not found",
			orderID: orderID,
			setup: func() {
				mockOrderRepo.EXPECT().
					GetOrder(ctx, orderID).
					Return(nil, errors.New("order not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := orderService.GetOrderStatus(ctx, tt.orderID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestOrderServiceImpl_generateOrderNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := orderRepo.NewMockOrderRepository(ctrl)
	mockProductRepo := productRepo.NewMockProductRepository(ctrl)
	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)
	mockNotificationService := service.NewMockNotificationService(ctrl)
	orderService := NewOrderServiceImpl(mockOrderRepo, mockProductRepo, mockAuthRepo, mockNotificationService).(*OrderServiceImpl)

	orderNumber := orderService.generateOrderNumber()

	assert.NotEmpty(t, orderNumber)
	assert.Contains(t, orderNumber, "ORD-")
	assert.Len(t, orderNumber, 7) // "ORD-" + 3 digits
}

func TestOrderServiceImpl_orderToResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := orderRepo.NewMockOrderRepository(ctrl)
	mockProductRepo := productRepo.NewMockProductRepository(ctrl)
	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)
	mockNotificationService := service.NewMockNotificationService(ctrl)
	orderService := NewOrderServiceImpl(mockOrderRepo, mockProductRepo, mockAuthRepo, mockNotificationService).(*OrderServiceImpl)

	orderID := uuid.New()
	customerID := uuid.New()
	now := time.Now()

	order := &entity.Order{
		OrderID:     orderID,
		CustomerID:  customerID,
		OrderNumber: "ORD-123",
		Status:      entity.OrderStatusConfirmed,
		TotalAmount: 150.00,
		CreatedAt:   now,
	}

	response := orderService.orderToResponse(order)

	expected := &schema.OrderResponse{
		OrderID:     orderID,
		CustomerID:  customerID,
		OrderNumber: "ORD-123",
		Status:      entity.OrderStatusConfirmed,
		TotalAmount: 150.00,
		CreatedAt:   now,
	}

	assert.Equal(t, expected, response)
}
