package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/auth/entity"
	"github.com/google/uuid"
)

type AuthRepository interface {
	// Customer authentication methods
	GetCustomerByOpenIDSub(ctx context.Context, openidSub string) (*entity.CustomerAuth, error)
	GetCustomerByEmail(ctx context.Context, email string) (*entity.CustomerAuth, error)
	GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*entity.CustomerAuth, error)

	// Customer creation/update for first-time auth
	CreateCustomer(ctx context.Context, customer *entity.CustomerAuth) (*entity.CustomerAuth, error)
	UpdateCustomerSetupStatus(ctx context.Context, customerID uuid.UUID, setupCompleted bool) (*entity.CustomerAuth, error)
	UpdateCustomer(ctx context.Context, customer *entity.CustomerAuth) (*entity.CustomerAuth, error)

	// Customer validation
	CustomerExists(ctx context.Context, openidSub string) (bool, error)
}
