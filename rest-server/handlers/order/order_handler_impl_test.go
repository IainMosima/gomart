package order

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authSchema "github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/IainMosima/gomart/domains/order/entity"
	"github.com/IainMosima/gomart/domains/order/schema"
	"github.com/IainMosima/gomart/domains/order/service"
	"github.com/IainMosima/gomart/rest-server/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Helper function to set user in gin context
func setUserInContext(c *gin.Context, userID uuid.UUID) {
	user := &authSchema.UserInfoResponse{
		UserID:        userID,
		UserName:      "testuser",
		Email:         "test@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}
	c.Set("user", user)
}

func TestNewOrderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockOrderService(ctrl)
	handler := NewOrderHandler(mockService)

	assert.NotNil(t, handler)
	assert.IsType(t, &OrderHandlerImpl{}, handler)
}

func TestOrderHandlerImpl_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockOrderService(ctrl)
	handler := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)

	customerID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		requestBody    interface{}
		setup          func(*gin.Context)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful order creation",
			requestBody: dtos.CreateOrderRequestDTO{
				CustomerID: customerID,
				Items: []dtos.CreateOrderItemRequestDTO{
					{
						ProductID: productID,
						Quantity:  2,
					},
				},
			},
			setup: func(c *gin.Context) {
				// Set authenticated user in context
				setUserInContext(c, customerID)

				expectedReq := &schema.CreateOrderRequest{
					CustomerID: customerID,
					Items: []schema.CreateOrderItemRequest{
						{
							ProductID: productID,
							Quantity:  2,
						},
					},
				}

				mockService.EXPECT().
					CreateOrder(gomock.Any(), expectedReq).
					Return(&schema.OrderResponse{
						OrderID:     orderID,
						CustomerID:  customerID,
						OrderNumber: "ORD-001",
						Status:      entity.OrderStatusConfirmed,
						TotalAmount: 200.00,
						CreatedAt:   now,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &dtos.OrderResponseDTO{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-001",
				Status:      string(entity.OrderStatusConfirmed),
				TotalAmount: 200.00,
				CreatedAt:   now,
			},
		},
		{
			name: "user not authenticated",
			requestBody: dtos.CreateOrderRequestDTO{
				CustomerID: customerID,
				Items: []dtos.CreateOrderItemRequestDTO{
					{
						ProductID: productID,
						Quantity:  2,
					},
				},
			},
			setup: func(c *gin.Context) {
				// Don't set user in context to simulate unauthenticated request
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   gin.H{"error": "User not authenticated"},
		},
		{
			name: "customer ID mismatch with authenticated user",
			requestBody: dtos.CreateOrderRequestDTO{
				CustomerID: customerID,
				Items: []dtos.CreateOrderItemRequestDTO{
					{
						ProductID: productID,
						Quantity:  2,
					},
				},
			},
			setup: func(c *gin.Context) {
				// Set different user ID to simulate unauthorized access
				differentUserID := uuid.New()
				setUserInContext(c, differentUserID)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   gin.H{"error": "Cannot create order for another customer"},
		},
		{
			name:           "invalid JSON request body",
			requestBody:    "invalid-json",
			setup:          func(c *gin.Context) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid character 'i' looking for beginning of value"},
		},
		{
			name: "service error",
			requestBody: dtos.CreateOrderRequestDTO{
				CustomerID: customerID,
				Items: []dtos.CreateOrderItemRequestDTO{
					{
						ProductID: productID,
						Quantity:  1,
					},
				},
			},
			setup: func(c *gin.Context) {
				setUserInContext(c, customerID)
				mockService.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("customer not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "customer not found"},
		},
		{
			name: "empty items array",
			requestBody: dtos.CreateOrderRequestDTO{
				CustomerID: customerID,
				Items:      []dtos.CreateOrderItemRequestDTO{},
			},
			setup: func(c *gin.Context) {
				setUserInContext(c, customerID)
				mockService.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("order must contain at least one item"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "order must contain at least one item"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var bodyBytes []byte
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			tt.setup(c)

			handler.(*OrderHandlerImpl).CreateOrder(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusCreated {
				expectedDTO := tt.expectedBody.(*dtos.OrderResponseDTO)
				actualResponse := response.(map[string]interface{})
				assert.Equal(t, expectedDTO.OrderNumber, actualResponse["order_number"])
				assert.Equal(t, expectedDTO.Status, actualResponse["status"])
				assert.Equal(t, expectedDTO.TotalAmount, actualResponse["total_amount"])
				assert.NotEmpty(t, actualResponse["order_id"])
			} else {
				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}

func TestOrderHandlerImpl_GetOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockOrderService(ctrl)
	handler := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)

	orderID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		orderID        string
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:    "successful status retrieval with valid UUID",
			orderID: orderID.String(),
			setup: func() {
				mockService.EXPECT().
					GetOrderStatus(gomock.Any(), orderID).
					Return(&schema.OrderStatusResponse{
						OrderID:     orderID,
						OrderNumber: "ORD-001",
						Status:      entity.OrderStatusConfirmed,
						CreatedAt:   now,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &dtos.OrderStatusResponseDTO{
				OrderID:     orderID,
				OrderNumber: "ORD-001",
				Status:      string(entity.OrderStatusConfirmed),
				CreatedAt:   now,
			},
		},
		{
			name:           "invalid order ID format",
			orderID:        "invalid-uuid",
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Invalid order ID"},
		},
		{
			name:    "order not found",
			orderID: orderID.String(),
			setup: func() {
				mockService.EXPECT().
					GetOrderStatus(gomock.Any(), orderID).
					Return(nil, errors.New("order not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "order not found"},
		},
		{
			name:    "service internal error",
			orderID: orderID.String(),
			setup: func() {
				mockService.EXPECT().
					GetOrderStatus(gomock.Any(), orderID).
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "database connection failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/orders/"+tt.orderID+"/status", nil)
			c.Params = gin.Params{{Key: "id", Value: tt.orderID}}

			handler.(*OrderHandlerImpl).GetOrderStatus(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusOK {
				expectedDTO := tt.expectedBody.(*dtos.OrderStatusResponseDTO)
				actualResponse := response.(map[string]interface{})
				assert.Equal(t, expectedDTO.OrderNumber, actualResponse["order_number"])
				assert.Equal(t, expectedDTO.Status, actualResponse["status"])
				assert.NotEmpty(t, actualResponse["order_id"])
			} else {
				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}

func TestOrderHandlerImpl_CreateOrder_ValidationCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockOrderService(ctrl)
	handler := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)

	customerID := uuid.New()
	productID := uuid.New()

	tests := []struct {
		name           string
		requestBody    interface{}
		setup          func(*gin.Context)
		expectedStatus int
	}{
		{
			name: "service validation error",
			requestBody: dtos.CreateOrderRequestDTO{
				CustomerID: customerID,
				Items: []dtos.CreateOrderItemRequestDTO{
					{
						ProductID: productID,
						Quantity:  0,
					},
				},
			},
			setup: func(c *gin.Context) {
				setUserInContext(c, customerID)
				mockService.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("invalid quantity"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			tt.setup(c)

			handler.(*OrderHandlerImpl).CreateOrder(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
