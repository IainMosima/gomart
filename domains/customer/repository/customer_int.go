package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/customer/entity"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	// Basic CRUD operations
	CreateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error)
	GetCustomer(ctx context.Context, customerID uuid.UUID) (*entity.Customer, error)
	GetCustomerByEmail(ctx context.Context, email string) (*entity.Customer, error)
	GetCustomerByOpenIDSub(ctx context.Context, openidSub string) (*entity.Customer, error)
	UpdateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error)
	UpdateCustomerSetupStatus(ctx context.Context, customerID uuid.UUID, setupCompleted bool) (*entity.Customer, error)
	SoftDeleteCustomer(ctx context.Context, customerID uuid.UUID) error

	// Customer listing and filtering
	ListCustomers(ctx context.Context) ([]*entity.Customer, error)
	ListCustomersBySetupStatus(ctx context.Context, setupCompleted bool) ([]*entity.Customer, error)
	GetCustomersWithPendingSetup(ctx context.Context) ([]*entity.Customer, error)

	// Customer search
	SearchCustomers(ctx context.Context, query string) ([]*entity.Customer, error)
	SearchCustomersByName(ctx context.Context, firstName, lastName string) ([]*entity.Customer, error)

	// Statistics and analytics
	CountCustomers(ctx context.Context) (int64, error)
	CountCustomersBySetupStatus(ctx context.Context, setupCompleted bool) (int64, error)
	GetCustomerStats(ctx context.Context) (*entity.CustomerStats, error)

	// Customer profile with order data
	GetCustomerProfile(ctx context.Context, customerID uuid.UUID) (*entity.CustomerProfile, error)

	// Validation helpers
	CustomerExists(ctx context.Context, customerID uuid.UUID) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	OpenIDSubExists(ctx context.Context, openidSub string) (bool, error)
}
