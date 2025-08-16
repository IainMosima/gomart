//go:generate mockgen -source=auth_repo_int.go -destination=auth_repo_mock.go -package=repository

package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/auth/entity"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, customer *entity.Customer) (*entity.Customer, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.Customer, error)
	GetUserByID(ctx context.Context, id string) (*entity.Customer, error)
}
