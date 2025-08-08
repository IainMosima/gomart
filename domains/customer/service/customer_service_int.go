package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/customer/schema"
	"github.com/google/uuid"
)

type CustomerService interface {
	// Basic CRUD operations
	CreateCustomer(ctx context.Context, req *schema.CreateCustomerRequest) (*schema.CustomerResponse, error)
	GetCustomer(ctx context.Context, customerID uuid.UUID) (*schema.CustomerResponse, error)
	GetCustomerByEmail(ctx context.Context, email string) (*schema.CustomerResponse, error)
	UpdateCustomer(ctx context.Context, customerID uuid.UUID, req *schema.UpdateCustomerRequest) (*schema.CustomerResponse, error)
	DeleteCustomer(ctx context.Context, customerID uuid.UUID) error

	// Customer profile management
	GetCustomerProfile(ctx context.Context, customerID uuid.UUID) (*schema.CustomerProfileResponse, error)
	UpdateSetupStatus(ctx context.Context, customerID uuid.UUID, req *schema.UpdateSetupStatusRequest) (*schema.CustomerResponse, error)
	CompleteCustomerSetup(ctx context.Context, customerID uuid.UUID, req *schema.UpdateCustomerRequest) (*schema.CustomerResponse, error)

	// Customer listing and search
	ListCustomers(ctx context.Context, req *schema.CustomerListRequest) (*schema.CustomerListResponse, error)
	SearchCustomers(ctx context.Context, req *schema.CustomerSearchRequest) (*schema.CustomerSearchResponse, error)
	GetCustomersWithPendingSetup(ctx context.Context) (*schema.CustomerListResponse, error)

	// Customer dashboard
	GetCustomerDashboard(ctx context.Context, customerID uuid.UUID) (*schema.CustomerDashboardResponse, error)

	// Statistics and analytics
	GetCustomerStats(ctx context.Context) (*schema.CustomerStatsResponse, error)

	// Validation and utilities
	ValidateCustomer(ctx context.Context, req *schema.CreateCustomerRequest) error
	CheckEmailAvailability(ctx context.Context, email string) (bool, error)
	GetCustomerByAuthInfo(ctx context.Context, openidSub string) (*schema.CustomerResponse, error)
}
