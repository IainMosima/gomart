package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/IainMosima/gomart/domains/order/entity"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockOrderStore struct {
	createOrderFunc     func(ctx context.Context, arg db.CreateOrderParams) (db.Order, error)
	createOrderItemFunc func(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItem, error)
	getOrderFunc        func(ctx context.Context, orderID uuid.UUID) (db.Order, error)
}

func (m *MockOrderStore) CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
	if m.createOrderFunc != nil {
		return m.createOrderFunc(ctx, arg)
	}
	return db.Order{}, nil
}

func (m *MockOrderStore) CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItem, error) {
	if m.createOrderItemFunc != nil {
		return m.createOrderItemFunc(ctx, arg)
	}
	return db.OrderItem{}, nil
}

func (m *MockOrderStore) GetOrder(ctx context.Context, orderID uuid.UUID) (db.Order, error) {
	if m.getOrderFunc != nil {
		return m.getOrderFunc(ctx, orderID)
	}
	return db.Order{}, nil
}

func (m *MockOrderStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.Customer, error) {
	return db.Customer{}, errors.New("not implemented")
}

func (m *MockOrderStore) GetUser(ctx context.Context, userID uuid.UUID) (db.Customer, error) {
	return db.Customer{}, errors.New("not implemented")
}

func (m *MockOrderStore) GetUserByEmail(ctx context.Context, email string) (db.Customer, error) {
	return db.Customer{}, errors.New("not implemented")
}

func (m *MockOrderStore) CreateCategory(ctx context.Context, arg db.CreateCategoryParams) (db.Category, error) {
	return db.Category{}, errors.New("not implemented")
}

func (m *MockOrderStore) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
	return db.Product{}, errors.New("not implemented")
}

func (m *MockOrderStore) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *MockOrderStore) GetCategory(ctx context.Context, categoryID uuid.UUID) (db.Category, error) {
	return db.Category{}, errors.New("not implemented")
}

func (m *MockOrderStore) GetCategoryAverageProductPrice(ctx context.Context, categoryID uuid.UUID) (pgtype.Numeric, error) {
	return pgtype.Numeric{}, errors.New("not implemented")
}

func (m *MockOrderStore) GetCategoryChildren(ctx context.Context, parentID pgtype.UUID) ([]db.Category, error) {
	return nil, errors.New("not implemented")
}

func (m *MockOrderStore) GetProduct(ctx context.Context, productID uuid.UUID) (db.Product, error) {
	return db.Product{}, errors.New("not implemented")
}

func (m *MockOrderStore) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]db.Product, error) {
	return nil, errors.New("not implemented")
}

func (m *MockOrderStore) GetRootCategories(ctx context.Context) ([]db.Category, error) {
	return nil, errors.New("not implemented")
}

func (m *MockOrderStore) ListCategories(ctx context.Context) ([]db.Category, error) {
	return nil, errors.New("not implemented")
}

func (m *MockOrderStore) ListProducts(ctx context.Context) ([]db.Product, error) {
	return nil, errors.New("not implemented")
}

func (m *MockOrderStore) SoftDeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *MockOrderStore) UpdateCategory(ctx context.Context, arg db.UpdateCategoryParams) (db.Category, error) {
	return db.Category{}, errors.New("not implemented")
}

func (m *MockOrderStore) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
	return db.Product{}, errors.New("not implemented")
}

func TestNewOrderRepository(t *testing.T) {
	mockStore := &MockOrderStore{}
	repo := NewOrderRepository(mockStore)

	assert.NotNil(t, repo)
	_, ok := repo.(*OrderRepositoryImpl)
	assert.True(t, ok)
}

func TestOrderRepositoryImpl_CreateOrder_Success(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()

	order := &entity.Order{
		CustomerID:  customerID,
		OrderNumber: "ORD-001",
		Status:      entity.OrderStatusPending,
		TotalAmount: 99.99,
	}

	expectedDBOrder := db.Order{
		OrderID:     orderID,
		CustomerID:  customerID,
		OrderNumber: "ORD-001",
		Status:      pgtype.Text{String: string(entity.OrderStatusPending), Valid: true},
	}
	expectedDBOrder.TotalAmount.Scan("99.99")

	mockStore := &MockOrderStore{
		createOrderFunc: func(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
			assert.Equal(t, customerID, arg.CustomerID)
			assert.Equal(t, "ORD-001", arg.OrderNumber)
			assert.Equal(t, string(entity.OrderStatusPending), arg.Status.String)
			assert.True(t, arg.Status.Valid)
			return expectedDBOrder, nil
		},
	}

	repo := &OrderRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.CreateOrder(ctx, order)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, orderID, result.OrderID)
	assert.Equal(t, customerID, result.CustomerID)
	assert.Equal(t, "ORD-001", result.OrderNumber)
	assert.Equal(t, entity.OrderStatusPending, result.Status)
	assert.Equal(t, 99.99, result.TotalAmount)
}

func TestOrderRepositoryImpl_CreateOrder_DatabaseError(t *testing.T) {
	order := &entity.Order{
		CustomerID:  uuid.New(),
		OrderNumber: "ORD-001",
		Status:      entity.OrderStatusPending,
		TotalAmount: 99.99,
	}

	mockStore := &MockOrderStore{
		createOrderFunc: func(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
			return db.Order{}, errors.New("database error")
		},
	}

	repo := &OrderRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.CreateOrder(ctx, order)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestOrderRepositoryImpl_CreateOrderItem_Success(t *testing.T) {
	orderItemID := uuid.New()
	orderID := uuid.New()
	productID := uuid.New()

	orderItem := &entity.OrderItem{
		OrderID:    orderID,
		ProductID:  productID,
		Quantity:   2,
		UnitPrice:  25.50,
		TotalPrice: 51.00,
	}

	expectedDBOrderItem := db.OrderItem{
		OrderItemID: orderItemID,
		OrderID:     orderID,
		ProductID:   productID,
		Quantity:    2,
	}
	expectedDBOrderItem.UnitPrice.Scan("25.50")
	expectedDBOrderItem.TotalPrice.Scan("51.00")

	mockStore := &MockOrderStore{
		createOrderItemFunc: func(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItem, error) {
			assert.Equal(t, orderID, arg.OrderID)
			assert.Equal(t, productID, arg.ProductID)
			assert.Equal(t, int32(2), arg.Quantity)
			return expectedDBOrderItem, nil
		},
	}

	repo := &OrderRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.CreateOrderItem(ctx, orderItem)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, orderItemID, result.OrderItemID)
	assert.Equal(t, orderID, result.OrderID)
	assert.Equal(t, productID, result.ProductID)
	assert.Equal(t, int32(2), result.Quantity)
	assert.Equal(t, 25.50, result.UnitPrice)
	assert.Equal(t, 51.00, result.TotalPrice)
}

func TestOrderRepositoryImpl_CreateOrderItem_DatabaseError(t *testing.T) {
	orderItem := &entity.OrderItem{
		OrderID:    uuid.New(),
		ProductID:  uuid.New(),
		Quantity:   1,
		UnitPrice:  10.00,
		TotalPrice: 10.00,
	}

	mockStore := &MockOrderStore{
		createOrderItemFunc: func(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItem, error) {
			return db.OrderItem{}, errors.New("database error")
		},
	}

	repo := &OrderRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.CreateOrderItem(ctx, orderItem)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestOrderRepositoryImpl_GetOrder_Success(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()

	expectedDBOrder := db.Order{
		OrderID:     orderID,
		CustomerID:  customerID,
		OrderNumber: "ORD-002",
		Status:      pgtype.Text{String: string(entity.OrderStatusConfirmed), Valid: true},
		IsDeleted:   pgtype.Bool{Bool: false, Valid: true},
	}
	expectedDBOrder.TotalAmount.Scan("149.99")

	mockStore := &MockOrderStore{
		getOrderFunc: func(ctx context.Context, id uuid.UUID) (db.Order, error) {
			assert.Equal(t, orderID, id)
			return expectedDBOrder, nil
		},
	}

	repo := &OrderRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetOrder(ctx, orderID)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, orderID, result.OrderID)
	assert.Equal(t, customerID, result.CustomerID)
	assert.Equal(t, "ORD-002", result.OrderNumber)
	assert.Equal(t, entity.OrderStatusConfirmed, result.Status)
	assert.Equal(t, 149.99, result.TotalAmount)
	assert.False(t, result.IsDeleted)
}

func TestOrderRepositoryImpl_GetOrder_NotFound(t *testing.T) {
	orderID := uuid.New()

	mockStore := &MockOrderStore{
		getOrderFunc: func(ctx context.Context, id uuid.UUID) (db.Order, error) {
			return db.Order{}, errors.New("no rows in result set")
		},
	}

	repo := &OrderRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetOrder(ctx, orderID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestOrderRepositoryImpl_convertOrderToEntity(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()

	dbOrder := db.Order{
		OrderID:     orderID,
		CustomerID:  customerID,
		OrderNumber: "ORD-003",
		Status:      pgtype.Text{String: string(entity.OrderStatusConfirmed), Valid: true},
		IsDeleted:   pgtype.Bool{Bool: false, Valid: true},
	}
	dbOrder.TotalAmount.Scan("75.25")

	repo := &OrderRepositoryImpl{}
	result := repo.convertOrderToEntity(dbOrder)

	assert.NotNil(t, result)
	assert.Equal(t, orderID, result.OrderID)
	assert.Equal(t, customerID, result.CustomerID)
	assert.Equal(t, "ORD-003", result.OrderNumber)
	assert.Equal(t, entity.OrderStatusConfirmed, result.Status)
	assert.Equal(t, 75.25, result.TotalAmount)
	assert.False(t, result.IsDeleted)
}

func TestOrderRepositoryImpl_convertOrderToEntity_NullStatus(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()

	dbOrder := db.Order{
		OrderID:     orderID,
		CustomerID:  customerID,
		OrderNumber: "ORD-004",
		Status:      pgtype.Text{Valid: false},
		TotalAmount: pgtype.Numeric{Valid: false},
	}

	repo := &OrderRepositoryImpl{}
	result := repo.convertOrderToEntity(dbOrder)

	assert.NotNil(t, result)
	assert.Equal(t, orderID, result.OrderID)
	assert.Equal(t, customerID, result.CustomerID)
	assert.Equal(t, "ORD-004", result.OrderNumber)
	assert.Empty(t, result.Status)
	assert.Equal(t, 0.0, result.TotalAmount)
}

func TestOrderRepositoryImpl_convertOrderItemToEntity(t *testing.T) {
	orderItemID := uuid.New()
	orderID := uuid.New()
	productID := uuid.New()

	dbOrderItem := db.OrderItem{
		OrderItemID: orderItemID,
		OrderID:     orderID,
		ProductID:   productID,
		Quantity:    3,
		IsDeleted:   pgtype.Bool{Bool: false, Valid: true},
	}
	dbOrderItem.UnitPrice.Scan("15.99")
	dbOrderItem.TotalPrice.Scan("47.97")

	repo := &OrderRepositoryImpl{}
	result := repo.convertOrderItemToEntity(dbOrderItem)

	assert.NotNil(t, result)
	assert.Equal(t, orderItemID, result.OrderItemID)
	assert.Equal(t, orderID, result.OrderID)
	assert.Equal(t, productID, result.ProductID)
	assert.Equal(t, int32(3), result.Quantity)
	assert.Equal(t, 15.99, result.UnitPrice)
	assert.Equal(t, 47.97, result.TotalPrice)
	assert.False(t, result.IsDeleted)
}

func TestOrderRepositoryImpl_convertOrderItemToEntity_NullPrices(t *testing.T) {
	orderItemID := uuid.New()
	orderID := uuid.New()
	productID := uuid.New()

	dbOrderItem := db.OrderItem{
		OrderItemID: orderItemID,
		OrderID:     orderID,
		ProductID:   productID,
		Quantity:    1,
		UnitPrice:   pgtype.Numeric{Valid: false},
		TotalPrice:  pgtype.Numeric{Valid: false},
	}

	repo := &OrderRepositoryImpl{}
	result := repo.convertOrderItemToEntity(dbOrderItem)

	assert.NotNil(t, result)
	assert.Equal(t, orderItemID, result.OrderItemID)
	assert.Equal(t, orderID, result.OrderID)
	assert.Equal(t, productID, result.ProductID)
	assert.Equal(t, int32(1), result.Quantity)
	assert.Equal(t, 0.0, result.UnitPrice)
	assert.Equal(t, 0.0, result.TotalPrice)
}
