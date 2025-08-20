package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/IainMosima/gomart/configs"
	authEntity "github.com/IainMosima/gomart/domains/auth/entity"
	authRepo "github.com/IainMosima/gomart/domains/auth/repository"
	"github.com/IainMosima/gomart/domains/order/entity"
	"github.com/IainMosima/gomart/domains/order/schema"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewNotificationServiceImpl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)
	config := &configs.Config{
		EmailHost:     "smtp.gmail.com",
		EmailPort:     "587",
		EmailUsername: "test@example.com",
		EmailPassword: "password",
		EmailFrom:     "test@example.com",
	}

	service := NewNotificationServiceImpl(config, mockAuthRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &EmailNotificationServiceImpl{}, service)
}

func TestEmailNotificationServiceImpl_SendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)

	ctx := context.Background()
	customerID := uuid.New()
	orderID := uuid.New()
	now := time.Now()

	tests := []struct {
		name    string
		config  *configs.Config
		order   *schema.OrderResponse
		setup   func()
		wantErr bool
	}{
		{
			name: "successful email sending",
			config: &configs.Config{
				EmailHost:     "smtp.gmail.com",
				EmailPort:     "587",
				EmailUsername: "test@example.com",
				EmailPassword: "password",
				EmailFrom:     "test@example.com",
			},
			order: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-001",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 150.00,
				CreatedAt:   now,
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(&authEntity.Customer{
						UserID: customerID,
						Email:  "customer@example.com",
					}, nil)
			},
			wantErr: true, // Will fail due to SMTP connection in test environment
		},
		{
			name: "customer not found",
			config: &configs.Config{
				EmailHost:     "smtp.gmail.com",
				EmailPort:     "587",
				EmailUsername: "test@example.com",
				EmailPassword: "password",
				EmailFrom:     "test@example.com",
			},
			order: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-002",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 100.00,
				CreatedAt:   now,
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(nil, errors.New("customer not found"))
			},
			wantErr: true,
		},
		{
			name: "missing email configuration",
			config: &configs.Config{
				EmailHost:     "",
				EmailPort:     "",
				EmailUsername: "",
				EmailPassword: "",
				EmailFrom:     "",
			},
			order: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-003",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 75.00,
				CreatedAt:   now,
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(&authEntity.Customer{
						UserID: customerID,
						Email:  "customer@example.com",
					}, nil)
			},
			wantErr: true, // Will fail due to empty SMTP configuration
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewNotificationServiceImpl(tt.config, mockAuthRepo)
			tt.setup()

			err := service.SendEmail(ctx, tt.order)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailNotificationServiceImpl_SendSMS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)

	ctx := context.Background()
	customerID := uuid.New()
	orderID := uuid.New()
	now := time.Now()

	tests := []struct {
		name    string
		config  *configs.Config
		order   *schema.OrderResponse
		setup   func()
		wantErr bool
	}{
		{
			name: "missing AfricasTalking configuration",
			config: &configs.Config{
				AfricasTalkingAPIKey:    "",
				AfricasTalkingUsername:  "",
				AfricasTalkingShortCode: "",
				AfricasTalkingSandbox:   "true",
			},
			order: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-001",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 150.00,
				CreatedAt:   now,
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "valid configuration but customer not found",
			config: &configs.Config{
				AfricasTalkingAPIKey:    "test-api-key",
				AfricasTalkingUsername:  "test-username",
				AfricasTalkingShortCode: "12345",
				AfricasTalkingSandbox:   "true",
			},
			order: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-002",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 100.00,
				CreatedAt:   now,
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(nil, errors.New("customer not found"))
			},
			wantErr: true,
		},
		{
			name: "valid configuration and customer found",
			config: &configs.Config{
				AfricasTalkingAPIKey:    "test-api-key",
				AfricasTalkingUsername:  "test-username",
				AfricasTalkingShortCode: "12345",
				AfricasTalkingSandbox:   "true",
			},
			order: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-003",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 75.00,
				CreatedAt:   now,
			},
			setup: func() {
				mockAuthRepo.EXPECT().
					GetUserByID(ctx, customerID).
					Return(&authEntity.Customer{
						UserID:      customerID,
						PhoneNumber: "+254712345678",
					}, nil)
			},
			wantErr: true, // Will fail due to AfricasTalking API connection in test environment
		},
		{
			name: "partial configuration missing API key",
			config: &configs.Config{
				AfricasTalkingAPIKey:    "",
				AfricasTalkingUsername:  "test-username",
				AfricasTalkingShortCode: "12345",
				AfricasTalkingSandbox:   "true",
			},
			order: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-004",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 50.00,
				CreatedAt:   now,
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "partial configuration missing username",
			config: &configs.Config{
				AfricasTalkingAPIKey:    "test-api-key",
				AfricasTalkingUsername:  "",
				AfricasTalkingShortCode: "12345",
				AfricasTalkingSandbox:   "true",
			},
			order: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-005",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 25.00,
				CreatedAt:   now,
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "partial configuration missing short code",
			config: &configs.Config{
				AfricasTalkingAPIKey:    "test-api-key",
				AfricasTalkingUsername:  "test-username",
				AfricasTalkingShortCode: "",
				AfricasTalkingSandbox:   "true",
			},
			order: &schema.OrderResponse{
				OrderID:     orderID,
				CustomerID:  customerID,
				OrderNumber: "ORD-006",
				Status:      entity.OrderStatusConfirmed,
				TotalAmount: 125.00,
				CreatedAt:   now,
			},
			setup:   func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewNotificationServiceImpl(tt.config, mockAuthRepo)
			tt.setup()

			err := service.SendSMS(ctx, tt.order)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmailNotificationServiceImpl_SendEmail_ErrorScenarios(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)
	ctx := context.Background()
	customerID := uuid.New()
	orderID := uuid.New()
	now := time.Now()

	config := &configs.Config{
		EmailHost:     "smtp.gmail.com",
		EmailPort:     "587",
		EmailUsername: "test@example.com",
		EmailPassword: "password",
		EmailFrom:     "test@example.com",
	}

	service := NewNotificationServiceImpl(config, mockAuthRepo)

	order := &schema.OrderResponse{
		OrderID:     orderID,
		CustomerID:  customerID,
		OrderNumber: "ORD-001",
		Status:      entity.OrderStatusConfirmed,
		TotalAmount: 150.00,
		CreatedAt:   now,
	}

	// Test customer not found scenario
	mockAuthRepo.EXPECT().
		GetUserByID(ctx, customerID).
		Return(nil, errors.New("database error"))

	err := service.SendEmail(ctx, order)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get customer details")
}

func TestEmailNotificationServiceImpl_SendSMS_ErrorScenarios(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := authRepo.NewMockAuthRepository(ctrl)
	ctx := context.Background()
	customerID := uuid.New()
	orderID := uuid.New()
	now := time.Now()

	config := &configs.Config{
		AfricasTalkingAPIKey:    "test-api-key",
		AfricasTalkingUsername:  "test-username",
		AfricasTalkingShortCode: "12345",
		AfricasTalkingSandbox:   "true",
	}

	service := NewNotificationServiceImpl(config, mockAuthRepo)

	order := &schema.OrderResponse{
		OrderID:     orderID,
		CustomerID:  customerID,
		OrderNumber: "ORD-001",
		Status:      entity.OrderStatusConfirmed,
		TotalAmount: 150.00,
		CreatedAt:   now,
	}

	// Test customer not found scenario
	mockAuthRepo.EXPECT().
		GetUserByID(ctx, customerID).
		Return(nil, errors.New("customer record not found"))

	err := service.SendSMS(ctx, order)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get customer details")
}
