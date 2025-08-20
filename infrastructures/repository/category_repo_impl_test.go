package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/IainMosima/gomart/domains/category/entity"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockCategoryStore struct {
	createCategoryFunc                 func(ctx context.Context, arg db.CreateCategoryParams) (db.Category, error)
	getCategoryFunc                    func(ctx context.Context, categoryID uuid.UUID) (db.Category, error)
	getCategoryChildrenFunc            func(ctx context.Context, parentID pgtype.UUID) ([]db.Category, error)
	getRootCategoriesFunc              func(ctx context.Context) ([]db.Category, error)
	listCategoriesFunc                 func(ctx context.Context) ([]db.Category, error)
	updateCategoryFunc                 func(ctx context.Context, arg db.UpdateCategoryParams) (db.Category, error)
	softDeleteCategoryFunc             func(ctx context.Context, categoryID uuid.UUID) error
	getCategoryAverageProductPriceFunc func(ctx context.Context, categoryID uuid.UUID) (pgtype.Numeric, error)
}

func (m *MockCategoryStore) CreateCategory(ctx context.Context, arg db.CreateCategoryParams) (db.Category, error) {
	if m.createCategoryFunc != nil {
		return m.createCategoryFunc(ctx, arg)
	}
	return db.Category{}, nil
}

func (m *MockCategoryStore) GetCategory(ctx context.Context, categoryID uuid.UUID) (db.Category, error) {
	if m.getCategoryFunc != nil {
		return m.getCategoryFunc(ctx, categoryID)
	}
	return db.Category{}, nil
}

func (m *MockCategoryStore) GetCategoryChildren(ctx context.Context, parentID pgtype.UUID) ([]db.Category, error) {
	if m.getCategoryChildrenFunc != nil {
		return m.getCategoryChildrenFunc(ctx, parentID)
	}
	return nil, nil
}

func (m *MockCategoryStore) GetRootCategories(ctx context.Context) ([]db.Category, error) {
	if m.getRootCategoriesFunc != nil {
		return m.getRootCategoriesFunc(ctx)
	}
	return nil, nil
}

func (m *MockCategoryStore) ListCategories(ctx context.Context) ([]db.Category, error) {
	if m.listCategoriesFunc != nil {
		return m.listCategoriesFunc(ctx)
	}
	return nil, nil
}

func (m *MockCategoryStore) UpdateCategory(ctx context.Context, arg db.UpdateCategoryParams) (db.Category, error) {
	if m.updateCategoryFunc != nil {
		return m.updateCategoryFunc(ctx, arg)
	}
	return db.Category{}, nil
}

func (m *MockCategoryStore) SoftDeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	if m.softDeleteCategoryFunc != nil {
		return m.softDeleteCategoryFunc(ctx, categoryID)
	}
	return nil
}

func (m *MockCategoryStore) GetCategoryAverageProductPrice(ctx context.Context, categoryID uuid.UUID) (pgtype.Numeric, error) {
	if m.getCategoryAverageProductPriceFunc != nil {
		return m.getCategoryAverageProductPriceFunc(ctx, categoryID)
	}
	return pgtype.Numeric{}, nil
}

func (m *MockCategoryStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.Customer, error) {
	return db.Customer{}, errors.New("not implemented")
}

func (m *MockCategoryStore) GetUser(ctx context.Context, userID uuid.UUID) (db.Customer, error) {
	return db.Customer{}, errors.New("not implemented")
}

func (m *MockCategoryStore) GetUserByEmail(ctx context.Context, email string) (db.Customer, error) {
	return db.Customer{}, errors.New("not implemented")
}

func (m *MockCategoryStore) CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
	return db.Order{}, errors.New("not implemented")
}

func (m *MockCategoryStore) CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItem, error) {
	return db.OrderItem{}, errors.New("not implemented")
}

func (m *MockCategoryStore) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
	return db.Product{}, errors.New("not implemented")
}

func (m *MockCategoryStore) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *MockCategoryStore) GetOrder(ctx context.Context, orderID uuid.UUID) (db.Order, error) {
	return db.Order{}, errors.New("not implemented")
}

func (m *MockCategoryStore) GetProduct(ctx context.Context, productID uuid.UUID) (db.Product, error) {
	return db.Product{}, errors.New("not implemented")
}

func (m *MockCategoryStore) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]db.Product, error) {
	return nil, errors.New("not implemented")
}

func (m *MockCategoryStore) ListProducts(ctx context.Context) ([]db.Product, error) {
	return nil, errors.New("not implemented")
}

func (m *MockCategoryStore) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
	return db.Product{}, errors.New("not implemented")
}

func TestNewCategoryRepository(t *testing.T) {
	mockStore := &MockCategoryStore{}
	repo := NewCategoryRepository(mockStore)

	assert.NotNil(t, repo)
	_, ok := repo.(*CategoryRepositoryImpl)
	assert.True(t, ok)
}

func TestCategoryRepositoryImpl_Create_Success(t *testing.T) {
	categoryID := uuid.New()
	parentID := uuid.New()

	category := &entity.Category{
		CategoryName: "Electronics",
		ParentID:     &parentID,
	}

	expectedDBCategory := db.Category{
		CategoryID:   categoryID,
		CategoryName: "Electronics",
		ParentID: pgtype.UUID{
			Bytes: parentID,
			Valid: true,
		},
	}

	mockStore := &MockCategoryStore{
		createCategoryFunc: func(ctx context.Context, arg db.CreateCategoryParams) (db.Category, error) {
			assert.Equal(t, "Electronics", arg.CategoryName)
			assert.Equal(t, parentID, uuid.UUID(arg.ParentID.Bytes))
			assert.True(t, arg.ParentID.Valid)
			return expectedDBCategory, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.Create(ctx, category)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, categoryID, result.CategoryID)
	assert.Equal(t, "Electronics", result.CategoryName)
	assert.Equal(t, parentID, *result.ParentID)
}

func TestCategoryRepositoryImpl_Create_NoParent(t *testing.T) {
	categoryID := uuid.New()

	category := &entity.Category{
		CategoryName: "Electronics",
		ParentID:     nil,
	}

	expectedDBCategory := db.Category{
		CategoryID:   categoryID,
		CategoryName: "Electronics",
		ParentID: pgtype.UUID{
			Valid: false,
		},
	}

	mockStore := &MockCategoryStore{
		createCategoryFunc: func(ctx context.Context, arg db.CreateCategoryParams) (db.Category, error) {
			assert.Equal(t, "Electronics", arg.CategoryName)
			assert.False(t, arg.ParentID.Valid)
			return expectedDBCategory, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.Create(ctx, category)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, categoryID, result.CategoryID)
	assert.Equal(t, "Electronics", result.CategoryName)
	assert.Nil(t, result.ParentID)
}

func TestCategoryRepositoryImpl_Create_DatabaseError(t *testing.T) {
	category := &entity.Category{
		CategoryName: "Electronics",
	}

	mockStore := &MockCategoryStore{
		createCategoryFunc: func(ctx context.Context, arg db.CreateCategoryParams) (db.Category, error) {
			return db.Category{}, errors.New("database error")
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.Create(ctx, category)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestCategoryRepositoryImpl_GetByID_Success(t *testing.T) {
	categoryID := uuid.New()
	parentID := uuid.New()

	expectedDBCategory := db.Category{
		CategoryID:   categoryID,
		CategoryName: "Electronics",
		ParentID: pgtype.UUID{
			Bytes: parentID,
			Valid: true,
		},
	}

	mockStore := &MockCategoryStore{
		getCategoryFunc: func(ctx context.Context, id uuid.UUID) (db.Category, error) {
			assert.Equal(t, categoryID, id)
			return expectedDBCategory, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetByID(ctx, categoryID)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, categoryID, result.CategoryID)
	assert.Equal(t, "Electronics", result.CategoryName)
	assert.Equal(t, parentID, *result.ParentID)
}

func TestCategoryRepositoryImpl_GetByID_NotFound(t *testing.T) {
	categoryID := uuid.New()

	mockStore := &MockCategoryStore{
		getCategoryFunc: func(ctx context.Context, id uuid.UUID) (db.Category, error) {
			return db.Category{}, errors.New("no rows in result set")
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetByID(ctx, categoryID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCategoryRepositoryImpl_GetByParent_Success(t *testing.T) {
	parentID := uuid.New()
	categoryID1 := uuid.New()
	categoryID2 := uuid.New()

	expectedDBCategories := []db.Category{
		{
			CategoryID:   categoryID1,
			CategoryName: "Phones",
			ParentID: pgtype.UUID{
				Bytes: parentID,
				Valid: true,
			},
		},
		{
			CategoryID:   categoryID2,
			CategoryName: "Laptops",
			ParentID: pgtype.UUID{
				Bytes: parentID,
				Valid: true,
			},
		},
	}

	mockStore := &MockCategoryStore{
		getCategoryChildrenFunc: func(ctx context.Context, pgParentID pgtype.UUID) ([]db.Category, error) {
			assert.Equal(t, parentID, uuid.UUID(pgParentID.Bytes))
			assert.True(t, pgParentID.Valid)
			return expectedDBCategories, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetByParent(ctx, &parentID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, categoryID1, result[0].CategoryID)
	assert.Equal(t, "Phones", result[0].CategoryName)
	assert.Equal(t, categoryID2, result[1].CategoryID)
	assert.Equal(t, "Laptops", result[1].CategoryName)
}

func TestCategoryRepositoryImpl_GetByParent_NilParent(t *testing.T) {
	mockStore := &MockCategoryStore{
		getCategoryChildrenFunc: func(ctx context.Context, pgParentID pgtype.UUID) ([]db.Category, error) {
			assert.False(t, pgParentID.Valid)
			return []db.Category{}, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetByParent(ctx, nil)

	require.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestCategoryRepositoryImpl_GetRootCategories_Success(t *testing.T) {
	categoryID1 := uuid.New()
	categoryID2 := uuid.New()

	expectedDBCategories := []db.Category{
		{
			CategoryID:   categoryID1,
			CategoryName: "Electronics",
		},
		{
			CategoryID:   categoryID2,
			CategoryName: "Clothing",
		},
	}

	mockStore := &MockCategoryStore{
		getRootCategoriesFunc: func(ctx context.Context) ([]db.Category, error) {
			return expectedDBCategories, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetRootCategories(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, categoryID1, result[0].CategoryID)
	assert.Equal(t, "Electronics", result[0].CategoryName)
	assert.Equal(t, categoryID2, result[1].CategoryID)
	assert.Equal(t, "Clothing", result[1].CategoryName)
}

func TestCategoryRepositoryImpl_GetAll_Success(t *testing.T) {
	categoryID := uuid.New()

	expectedDBCategories := []db.Category{
		{
			CategoryID:   categoryID,
			CategoryName: "Electronics",
		},
	}

	mockStore := &MockCategoryStore{
		listCategoriesFunc: func(ctx context.Context) ([]db.Category, error) {
			return expectedDBCategories, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetAll(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, categoryID, result[0].CategoryID)
	assert.Equal(t, "Electronics", result[0].CategoryName)
}

func TestCategoryRepositoryImpl_Update_Success(t *testing.T) {
	categoryID := uuid.New()

	category := &entity.Category{
		CategoryID:   categoryID,
		CategoryName: "Updated Electronics",
	}

	expectedDBCategory := db.Category{
		CategoryID:   categoryID,
		CategoryName: "Updated Electronics",
	}

	mockStore := &MockCategoryStore{
		updateCategoryFunc: func(ctx context.Context, arg db.UpdateCategoryParams) (db.Category, error) {
			assert.Equal(t, categoryID, arg.CategoryID)
			assert.Equal(t, "Updated Electronics", arg.CategoryName)
			return expectedDBCategory, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.Update(ctx, category)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, categoryID, result.CategoryID)
	assert.Equal(t, "Updated Electronics", result.CategoryName)
}

func TestCategoryRepositoryImpl_SoftDelete_Success(t *testing.T) {
	categoryID := uuid.New()

	mockStore := &MockCategoryStore{
		softDeleteCategoryFunc: func(ctx context.Context, id uuid.UUID) error {
			assert.Equal(t, categoryID, id)
			return nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	err := repo.SoftDelete(ctx, categoryID)

	assert.NoError(t, err)
}

func TestCategoryRepositoryImpl_GetAverageProductPrice_Success(t *testing.T) {
	categoryID := uuid.New()

	numeric := pgtype.Numeric{}
	numeric.Scan("25.99")

	mockStore := &MockCategoryStore{
		getCategoryAverageProductPriceFunc: func(ctx context.Context, id uuid.UUID) (pgtype.Numeric, error) {
			assert.Equal(t, categoryID, id)
			return numeric, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetAverageProductPrice(ctx, categoryID)

	require.NoError(t, err)
	assert.Equal(t, 25.99, result)
}

func TestCategoryRepositoryImpl_GetAverageProductPrice_NoData(t *testing.T) {
	categoryID := uuid.New()

	numeric := pgtype.Numeric{Valid: false}

	mockStore := &MockCategoryStore{
		getCategoryAverageProductPriceFunc: func(ctx context.Context, id uuid.UUID) (pgtype.Numeric, error) {
			return numeric, nil
		},
	}

	repo := &CategoryRepositoryImpl{store: mockStore}
	ctx := context.Background()

	result, err := repo.GetAverageProductPrice(ctx, categoryID)

	require.NoError(t, err)
	assert.Equal(t, 0.0, result)
}

func TestCategoryRepositoryImpl_convertToEntity(t *testing.T) {
	categoryID := uuid.New()
	parentID := uuid.New()

	dbCategory := db.Category{
		CategoryID:   categoryID,
		CategoryName: "Electronics",
		ParentID: pgtype.UUID{
			Bytes: parentID,
			Valid: true,
		},
		IsDeleted: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	}

	repo := &CategoryRepositoryImpl{}
	result := repo.convertToEntity(dbCategory)

	assert.NotNil(t, result)
	assert.Equal(t, categoryID, result.CategoryID)
	assert.Equal(t, "Electronics", result.CategoryName)
	assert.Equal(t, parentID, *result.ParentID)
	assert.False(t, result.IsDeleted)
}
