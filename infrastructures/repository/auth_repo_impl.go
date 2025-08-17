package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/auth/entity"
	domainRepo "github.com/IainMosima/gomart/domains/auth/repository"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type AuthRepositoryImpl struct {
	store db.Store
}

func NewAuthRepository(store db.Store) domainRepo.AuthRepository {
	return &AuthRepositoryImpl{
		store: store,
	}
}

func (r *AuthRepositoryImpl) CreateUser(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	params := db.CreateUserParams{
		UserID:      customer.UserID,
		PhoneNumber: customer.PhoneNumber,
		UserName:    customer.UserName,
		Email:       customer.Email,
	}

	result, err := r.store.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(result), nil
}

func (r *AuthRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	result, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(result), nil
}

func (r *AuthRepositoryImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.Customer, error) {
	result, err := r.store.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(result), nil
}

func (r *AuthRepositoryImpl) convertToEntity(dbCustomer db.Customer) *entity.Customer {
	return &entity.Customer{
		UserID:      dbCustomer.UserID,
		PhoneNumber: dbCustomer.PhoneNumber,
		UserName:    dbCustomer.UserName,
		Email:       dbCustomer.Email,
	}
}
