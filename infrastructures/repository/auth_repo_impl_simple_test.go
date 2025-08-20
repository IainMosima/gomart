package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/IainMosima/gomart/domains/auth/entity"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockStore struct {
	createUserFunc     func(ctx context.Context, arg db.CreateUserParams) (db.Customer, error)
	getUserFunc        func(ctx context.Context, userID uuid.UUID) (db.Customer, error)
	getUserByEmailFunc func(ctx context.Context, email string) (db.Customer, error)
}

func (m *MockStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.Customer, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, arg)
	}
	return db.Customer{}, nil
}

func (m *MockStore) GetUser(ctx context.Context, userID uuid.UUID) (db.Customer, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, userID)
	}
	return db.Customer{}, nil
}

func (m *MockStore) GetUserByEmail(ctx context.Context, email string) (db.Customer, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	return db.Customer{}, nil
}

func (m *MockStore) CreateCategory(ctx context.Context, arg db.CreateCategoryParams) (db.Category, error) {
	return db.Category{}, errors.New("not implemented")
}

func (m *MockStore) CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
	return db.Order{}, errors.New("not implemented")
}

func (m *MockStore) CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItem, error) {
	return db.OrderItem{}, errors.New("not implemented")
}

func (m *MockStore) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
	return db.Product{}, errors.New("not implemented")
}

func (m *MockStore) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *MockStore) GetCategory(ctx context.Context, categoryID uuid.UUID) (db.Category, error) {
	return db.Category{}, errors.New("not implemented")
}

func (m *MockStore) GetCategoryAverageProductPrice(ctx context.Context, categoryID uuid.UUID) (pgtype.Numeric, error) {
	return pgtype.Numeric{}, errors.New("not implemented")
}

func (m *MockStore) GetCategoryChildren(ctx context.Context, parentID pgtype.UUID) ([]db.Category, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStore) GetOrder(ctx context.Context, orderID uuid.UUID) (db.Order, error) {
	return db.Order{}, errors.New("not implemented")
}

func (m *MockStore) GetProduct(ctx context.Context, productID uuid.UUID) (db.Product, error) {
	return db.Product{}, errors.New("not implemented")
}

func (m *MockStore) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]db.Product, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStore) GetRootCategories(ctx context.Context) ([]db.Category, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStore) ListCategories(ctx context.Context) ([]db.Category, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStore) ListProducts(ctx context.Context) ([]db.Product, error) {
	return nil, errors.New("not implemented")
}

func (m *MockStore) SoftDeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *MockStore) UpdateCategory(ctx context.Context, arg db.UpdateCategoryParams) (db.Category, error) {
	return db.Category{}, errors.New("not implemented")
}

func (m *MockStore) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
	return db.Product{}, errors.New("not implemented")
}

func TestNewAuthRepository_Simple(t *testing.T) {
	mockStore := &MockStore{}
	repo := NewAuthRepository(mockStore)

	assert.NotNil(t, repo)
	_, ok := repo.(*AuthRepositoryImpl)
	assert.True(t, ok)
}

func TestAuthRepositoryImpl_CreateUser_Success_Simple(t *testing.T) {
	userID := uuid.New()
	customer := &entity.Customer{
		UserID:      userID,
		PhoneNumber: "+1234567890",
		UserName:    "testuser",
		Email:       "test@example.com",
	}

	expectedDBCustomer := db.Customer{
		UserID:      userID,
		PhoneNumber: "+1234567890",
		UserName:    "testuser",
		Email:       "test@example.com",
	}

	mockStore := &MockStore{
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.Customer, error) {
			assert.Equal(t, userID, arg.UserID)
			assert.Equal(t, "+1234567890", arg.PhoneNumber)
			assert.Equal(t, "testuser", arg.UserName)
			assert.Equal(t, "test@example.com", arg.Email)
			return expectedDBCustomer, nil
		},
	}

	repo := &AuthRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.CreateUser(ctx, customer)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, customer.UserID, result.UserID)
	assert.Equal(t, customer.PhoneNumber, result.PhoneNumber)
	assert.Equal(t, customer.UserName, result.UserName)
	assert.Equal(t, customer.Email, result.Email)
}

func TestAuthRepositoryImpl_CreateUser_DatabaseError_Simple(t *testing.T) {
	customer := &entity.Customer{
		UserID:      uuid.New(),
		PhoneNumber: "+1234567890",
		UserName:    "testuser",
		Email:       "test@example.com",
	}

	expectedError := errors.New("database connection failed")

	mockStore := &MockStore{
		createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.Customer, error) {
			return db.Customer{}, expectedError
		},
	}

	repo := &AuthRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.CreateUser(ctx, customer)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestAuthRepositoryImpl_GetUserByID_Success_Simple(t *testing.T) {
	userID := uuid.New()
	expectedDBCustomer := db.Customer{
		UserID:      userID,
		PhoneNumber: "+1234567890",
		UserName:    "testuser",
		Email:       "test@example.com",
	}

	mockStore := &MockStore{
		getUserFunc: func(ctx context.Context, id uuid.UUID) (db.Customer, error) {
			assert.Equal(t, userID, id)
			return expectedDBCustomer, nil
		},
	}

	repo := &AuthRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetUserByID(ctx, userID)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "+1234567890", result.PhoneNumber)
	assert.Equal(t, "testuser", result.UserName)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestAuthRepositoryImpl_GetUserByID_NotFound_Simple(t *testing.T) {
	userID := uuid.New()
	notFoundError := errors.New("no rows in result set")

	mockStore := &MockStore{
		getUserFunc: func(ctx context.Context, id uuid.UUID) (db.Customer, error) {
			return db.Customer{}, notFoundError
		},
	}

	repo := &AuthRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, notFoundError, err)
}

func TestAuthRepositoryImpl_GetUserByEmail_Success_Simple(t *testing.T) {
	email := "test@example.com"
	userID := uuid.New()
	expectedDBCustomer := db.Customer{
		UserID:      userID,
		PhoneNumber: "+1234567890",
		UserName:    "testuser",
		Email:       email,
	}

	mockStore := &MockStore{
		getUserByEmailFunc: func(ctx context.Context, e string) (db.Customer, error) {
			assert.Equal(t, email, e)
			return expectedDBCustomer, nil
		},
	}

	repo := &AuthRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetUserByEmail(ctx, email)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "+1234567890", result.PhoneNumber)
	assert.Equal(t, "testuser", result.UserName)
	assert.Equal(t, email, result.Email)
}

func TestAuthRepositoryImpl_GetUserByEmail_NotFound_Simple(t *testing.T) {
	email := "nonexistent@example.com"
	notFoundError := errors.New("no rows in result set")

	mockStore := &MockStore{
		getUserByEmailFunc: func(ctx context.Context, e string) (db.Customer, error) {
			return db.Customer{}, notFoundError
		},
	}

	repo := &AuthRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetUserByEmail(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, notFoundError, err)
}

func TestAuthRepositoryImpl_convertToEntity_Simple(t *testing.T) {
	repo := &AuthRepositoryImpl{}
	userID := uuid.New()

	dbCustomer := db.Customer{
		UserID:      userID,
		PhoneNumber: "+1234567890",
		UserName:    "testuser",
		Email:       "test@example.com",
	}

	result := repo.convertToEntity(dbCustomer)

	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "+1234567890", result.PhoneNumber)
	assert.Equal(t, "testuser", result.UserName)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestAuthRepositoryImpl_ErrorScenarios_Simple(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() db.Store
		testFunc    func(repo *AuthRepositoryImpl) error
		expectError string
	}{
		{
			name: "CreateUser duplicate key error",
			setupMock: func() db.Store {
				return &MockStore{
					createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.Customer, error) {
						return db.Customer{}, errors.New("duplicate key value violates unique constraint")
					},
				}
			},
			testFunc: func(repo *AuthRepositoryImpl) error {
				customer := &entity.Customer{
					UserID:      uuid.New(),
					PhoneNumber: "+1234567890",
					UserName:    "testuser",
					Email:       "test@example.com",
				}
				_, err := repo.CreateUser(context.Background(), customer)
				return err
			},
			expectError: "duplicate key",
		},
		{
			name: "GetUserByID database timeout",
			setupMock: func() db.Store {
				return &MockStore{
					getUserFunc: func(ctx context.Context, userID uuid.UUID) (db.Customer, error) {
						return db.Customer{}, errors.New("database connection timeout")
					},
				}
			},
			testFunc: func(repo *AuthRepositoryImpl) error {
				_, err := repo.GetUserByID(context.Background(), uuid.New())
				return err
			},
			expectError: "timeout",
		},
		{
			name: "GetUserByEmail connection error",
			setupMock: func() db.Store {
				return &MockStore{
					getUserByEmailFunc: func(ctx context.Context, email string) (db.Customer, error) {
						return db.Customer{}, errors.New("connection refused")
					},
				}
			},
			testFunc: func(repo *AuthRepositoryImpl) error {
				_, err := repo.GetUserByEmail(context.Background(), "test@example.com")
				return err
			},
			expectError: "connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := tt.setupMock()
			repo := &AuthRepositoryImpl{store: mockStore}

			err := tt.testFunc(repo)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}
}
