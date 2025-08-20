package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/IainMosima/gomart/domains/product/entity"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockProductStore struct {
	createProductFunc         func(ctx context.Context, arg db.CreateProductParams) (db.Product, error)
	getProductFunc            func(ctx context.Context, productID uuid.UUID) (db.Product, error)
	updateProductFunc         func(ctx context.Context, arg db.UpdateProductParams) (db.Product, error)
	deleteProductFunc         func(ctx context.Context, productID uuid.UUID) error
	listProductsFunc          func(ctx context.Context) ([]db.Product, error)
	getProductsByCategoryFunc func(ctx context.Context, categoryID uuid.UUID) ([]db.Product, error)
}

func (m *MockProductStore) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
	if m.createProductFunc != nil {
		return m.createProductFunc(ctx, arg)
	}
	return db.Product{}, nil
}

func (m *MockProductStore) GetProduct(ctx context.Context, productID uuid.UUID) (db.Product, error) {
	if m.getProductFunc != nil {
		return m.getProductFunc(ctx, productID)
	}
	return db.Product{}, nil
}

func (m *MockProductStore) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
	if m.updateProductFunc != nil {
		return m.updateProductFunc(ctx, arg)
	}
	return db.Product{}, nil
}

func (m *MockProductStore) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	if m.deleteProductFunc != nil {
		return m.deleteProductFunc(ctx, productID)
	}
	return nil
}

func (m *MockProductStore) ListProducts(ctx context.Context) ([]db.Product, error) {
	if m.listProductsFunc != nil {
		return m.listProductsFunc(ctx)
	}
	return nil, nil
}

func (m *MockProductStore) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]db.Product, error) {
	if m.getProductsByCategoryFunc != nil {
		return m.getProductsByCategoryFunc(ctx, categoryID)
	}
	return nil, nil
}

func (m *MockProductStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.Customer, error) {
	return db.Customer{}, errors.New("not implemented")
}

func (m *MockProductStore) GetUser(ctx context.Context, userID uuid.UUID) (db.Customer, error) {
	return db.Customer{}, errors.New("not implemented")
}

func (m *MockProductStore) GetUserByEmail(ctx context.Context, email string) (db.Customer, error) {
	return db.Customer{}, errors.New("not implemented")
}

func (m *MockProductStore) CreateCategory(ctx context.Context, arg db.CreateCategoryParams) (db.Category, error) {
	return db.Category{}, errors.New("not implemented")
}

func (m *MockProductStore) CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
	return db.Order{}, errors.New("not implemented")
}

func (m *MockProductStore) CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItem, error) {
	return db.OrderItem{}, errors.New("not implemented")
}

func (m *MockProductStore) GetCategory(ctx context.Context, categoryID uuid.UUID) (db.Category, error) {
	return db.Category{}, errors.New("not implemented")
}

func (m *MockProductStore) GetCategoryAverageProductPrice(ctx context.Context, categoryID uuid.UUID) (pgtype.Numeric, error) {
	return pgtype.Numeric{}, errors.New("not implemented")
}

func (m *MockProductStore) GetCategoryChildren(ctx context.Context, parentID pgtype.UUID) ([]db.Category, error) {
	return nil, errors.New("not implemented")
}

func (m *MockProductStore) GetOrder(ctx context.Context, orderID uuid.UUID) (db.Order, error) {
	return db.Order{}, errors.New("not implemented")
}

func (m *MockProductStore) GetRootCategories(ctx context.Context) ([]db.Category, error) {
	return nil, errors.New("not implemented")
}

func (m *MockProductStore) ListCategories(ctx context.Context) ([]db.Category, error) {
	return nil, errors.New("not implemented")
}

func (m *MockProductStore) SoftDeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *MockProductStore) UpdateCategory(ctx context.Context, arg db.UpdateCategoryParams) (db.Category, error) {
	return db.Category{}, errors.New("not implemented")
}

func TestNewProductRepository(t *testing.T) {
	mockStore := &MockProductStore{}
	repo := NewProductRepository(mockStore)

	assert.NotNil(t, repo)
	_, ok := repo.(*ProductRepositoryImpl)
	assert.True(t, ok)
}

func TestProductRepositoryImpl_Create_Success(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()
	description := "Test product description"

	product := &entity.Product{
		ProductName:   "Test Product",
		Description:   &description,
		Price:         29.99,
		SKU:           "TEST-001",
		StockQuantity: 100,
		CategoryID:    categoryID,
		IsActive:      true,
	}

	expectedDBProduct := db.Product{
		ProductID:     productID,
		ProductName:   "Test Product",
		Description:   pgtype.Text{String: description, Valid: true},
		Sku:           "TEST-001",
		StockQuantity: 100,
		CategoryID:    categoryID,
		IsActive:      pgtype.Bool{Bool: true, Valid: true},
	}
	expectedDBProduct.Price.Scan("29.99")

	mockStore := &MockProductStore{
		createProductFunc: func(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
			assert.Equal(t, "Test Product", arg.ProductName)
			assert.Equal(t, description, arg.Description.String)
			assert.True(t, arg.Description.Valid)
			assert.Equal(t, "TEST-001", arg.Sku)
			assert.Equal(t, int32(100), arg.StockQuantity)
			assert.Equal(t, categoryID, arg.CategoryID)
			assert.True(t, arg.IsActive.Bool)
			return expectedDBProduct, nil
		},
	}

	repo := &ProductRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.Create(ctx, product)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, productID, result.ProductID)
	assert.Equal(t, "Test Product", result.ProductName)
	assert.Equal(t, description, *result.Description)
	assert.Equal(t, "TEST-001", result.SKU)
	assert.Equal(t, int32(100), result.StockQuantity)
	assert.Equal(t, categoryID, result.CategoryID)
	assert.True(t, result.IsActive)
}

func TestProductRepositoryImpl_Create_NoDescription(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()

	product := &entity.Product{
		ProductName:   "Test Product",
		Description:   nil,
		Price:         29.99,
		SKU:           "TEST-001",
		StockQuantity: 100,
		CategoryID:    categoryID,
		IsActive:      true,
	}

	expectedDBProduct := db.Product{
		ProductID:     productID,
		ProductName:   "Test Product",
		Description:   pgtype.Text{Valid: false},
		Sku:           "TEST-001",
		StockQuantity: 100,
		CategoryID:    categoryID,
	}

	mockStore := &MockProductStore{
		createProductFunc: func(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
			assert.False(t, arg.Description.Valid)
			return expectedDBProduct, nil
		},
	}

	repo := &ProductRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.Create(ctx, product)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.Description)
}

func TestProductRepositoryImpl_Create_DatabaseError(t *testing.T) {
	product := &entity.Product{
		ProductName: "Test Product",
		Price:       29.99,
	}

	mockStore := &MockProductStore{
		createProductFunc: func(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
			return db.Product{}, errors.New("database error")
		},
	}

	repo := &ProductRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.Create(ctx, product)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestProductRepositoryImpl_GetByID_Success(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()
	description := "Test description"

	expectedDBProduct := db.Product{
		ProductID:     productID,
		ProductName:   "Test Product",
		Description:   pgtype.Text{String: description, Valid: true},
		Sku:           "TEST-001",
		StockQuantity: 100,
		CategoryID:    categoryID,
		IsActive:      pgtype.Bool{Bool: true, Valid: true},
	}
	expectedDBProduct.Price.Scan("29.99")

	mockStore := &MockProductStore{
		getProductFunc: func(ctx context.Context, id uuid.UUID) (db.Product, error) {
			assert.Equal(t, productID, id)
			return expectedDBProduct, nil
		},
	}

	repo := &ProductRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetByID(ctx, productID)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, productID, result.ProductID)
	assert.Equal(t, "Test Product", result.ProductName)
	assert.Equal(t, description, *result.Description)
	assert.Equal(t, 29.99, result.Price)
}

func TestProductRepositoryImpl_GetByID_NotFound(t *testing.T) {
	productID := uuid.New()

	mockStore := &MockProductStore{
		getProductFunc: func(ctx context.Context, id uuid.UUID) (db.Product, error) {
			return db.Product{}, errors.New("no rows in result set")
		},
	}

	repo := &ProductRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetByID(ctx, productID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestProductRepositoryImpl_Update_Success(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()
	description := "Updated description"

	product := &entity.Product{
		ProductID:     productID,
		ProductName:   "Updated Product",
		Description:   &description,
		Price:         39.99,
		StockQuantity: 50,
		CategoryID:    categoryID,
		IsActive:      false,
	}

	expectedDBProduct := db.Product{
		ProductID:     productID,
		ProductName:   "Updated Product",
		Description:   pgtype.Text{String: description, Valid: true},
		StockQuantity: 50,
		CategoryID:    categoryID,
		IsActive:      pgtype.Bool{Bool: false, Valid: true},
	}
	expectedDBProduct.Price.Scan("39.99")

	mockStore := &MockProductStore{
		updateProductFunc: func(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
			assert.Equal(t, productID, arg.ProductID)
			assert.Equal(t, "Updated Product", arg.ProductName)
			assert.Equal(t, description, arg.Description.String)
			assert.Equal(t, int32(50), arg.StockQuantity)
			assert.Equal(t, categoryID, arg.CategoryID)
			assert.False(t, arg.IsActive.Bool)
			return expectedDBProduct, nil
		},
	}

	repo := &ProductRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.Update(ctx, product)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, productID, result.ProductID)
	assert.Equal(t, "Updated Product", result.ProductName)
	assert.Equal(t, description, *result.Description)
	assert.Equal(t, 39.99, result.Price)
	assert.False(t, result.IsActive)
}

func TestProductRepositoryImpl_Delete_Success(t *testing.T) {
	productID := uuid.New()

	mockStore := &MockProductStore{
		deleteProductFunc: func(ctx context.Context, id uuid.UUID) error {
			assert.Equal(t, productID, id)
			return nil
		},
	}

	repo := &ProductRepositoryImpl{store: mockStore}
	ctx := context.Background()

	err := repo.Delete(ctx, productID)

	assert.NoError(t, err)
}

func TestProductRepositoryImpl_GetAll_Success(t *testing.T) {
	productID1 := uuid.New()
	productID2 := uuid.New()
	categoryID := uuid.New()

	expectedDBProducts := []db.Product{
		{
			ProductID:     productID1,
			ProductName:   "Product 1",
			StockQuantity: 100,
			CategoryID:    categoryID,
		},
		{
			ProductID:     productID2,
			ProductName:   "Product 2",
			StockQuantity: 50,
			CategoryID:    categoryID,
		},
	}

	mockStore := &MockProductStore{
		listProductsFunc: func(ctx context.Context) ([]db.Product, error) {
			return expectedDBProducts, nil
		},
	}

	repo := &ProductRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetAll(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, productID1, result[0].ProductID)
	assert.Equal(t, "Product 1", result[0].ProductName)
	assert.Equal(t, productID2, result[1].ProductID)
	assert.Equal(t, "Product 2", result[1].ProductName)
}

func TestProductRepositoryImpl_GetByCategory_Success(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()

	expectedDBProducts := []db.Product{
		{
			ProductID:     productID,
			ProductName:   "Category Product",
			StockQuantity: 25,
			CategoryID:    categoryID,
		},
	}

	mockStore := &MockProductStore{
		getProductsByCategoryFunc: func(ctx context.Context, catID uuid.UUID) ([]db.Product, error) {
			assert.Equal(t, categoryID, catID)
			return expectedDBProducts, nil
		},
	}

	repo := &ProductRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetByCategory(ctx, categoryID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, productID, result[0].ProductID)
	assert.Equal(t, "Category Product", result[0].ProductName)
	assert.Equal(t, categoryID, result[0].CategoryID)
}

func TestProductRepositoryImpl_convertToEntity(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()
	description := "Test description"

	dbProduct := db.Product{
		ProductID:     productID,
		ProductName:   "Test Product",
		Description:   pgtype.Text{String: description, Valid: true},
		Sku:           "TEST-001",
		StockQuantity: 100,
		CategoryID:    categoryID,
		IsActive:      pgtype.Bool{Bool: true, Valid: true},
		IsDeleted:     pgtype.Bool{Bool: false, Valid: true},
	}
	dbProduct.Price.Scan("29.99")

	repo := &ProductRepositoryImpl{}
	result := repo.convertToEntity(dbProduct)

	assert.NotNil(t, result)
	assert.Equal(t, productID, result.ProductID)
	assert.Equal(t, "Test Product", result.ProductName)
	assert.Equal(t, description, *result.Description)
	assert.Equal(t, 29.99, result.Price)
	assert.Equal(t, "TEST-001", result.SKU)
	assert.Equal(t, int32(100), result.StockQuantity)
	assert.Equal(t, categoryID, result.CategoryID)
	assert.True(t, result.IsActive)
	assert.False(t, result.IsDeleted)
}

func TestProductRepositoryImpl_convertToEntity_NullFields(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()

	dbProduct := db.Product{
		ProductID:     productID,
		ProductName:   "Test Product",
		Description:   pgtype.Text{Valid: false},
		Price:         pgtype.Numeric{Valid: false},
		Sku:           "TEST-001",
		StockQuantity: 100,
		CategoryID:    categoryID,
		IsActive:      pgtype.Bool{Valid: false},
	}

	repo := &ProductRepositoryImpl{}
	result := repo.convertToEntity(dbProduct)

	assert.NotNil(t, result)
	assert.Equal(t, productID, result.ProductID)
	assert.Equal(t, "Test Product", result.ProductName)
	assert.Nil(t, result.Description)
	assert.Equal(t, 0.0, result.Price)
	assert.False(t, result.IsActive)
}
