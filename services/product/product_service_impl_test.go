package product

import (
	"context"
	"errors"
	"testing"
	"time"

	categoryEntity "github.com/IainMosima/gomart/domains/category/entity"
	categoryRepo "github.com/IainMosima/gomart/domains/category/repository"
	"github.com/IainMosima/gomart/domains/product/entity"
	"github.com/IainMosima/gomart/domains/product/repository"
	"github.com/IainMosima/gomart/domains/product/schema"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewProductService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := repository.NewMockProductRepository(ctrl)
	mockCategoryRepo := categoryRepo.NewMockCategoryRepository(ctrl)
	service := NewProductService(mockProductRepo, mockCategoryRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &ProductServiceImpl{}, service)
}

func TestProductServiceImpl_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := repository.NewMockProductRepository(ctrl)
	mockCategoryRepo := categoryRepo.NewMockCategoryRepository(ctrl)
	service := NewProductService(mockProductRepo, mockCategoryRepo)

	ctx := context.Background()
	now := time.Now()
	productID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name    string
		request *schema.CreateProductRequest
		setup   func()
		want    *schema.ProductResponse
		wantErr bool
	}{
		{
			name: "successful creation",
			request: &schema.CreateProductRequest{
				ProductName:   "iPhone 15",
				Description:   stringPtr("Latest iPhone"),
				Price:         999.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
				IsActive:      boolPtr(true),
			},
			setup: func() {
				// Mock category validation
				mockCategoryRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(&categoryEntity.Category{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
					}, nil)

				// Mock product creation
				mockProductRepo.EXPECT().
					Create(ctx, gomock.Any()).
					Return(&entity.Product{
						ProductID:     productID,
						ProductName:   "iPhone 15",
						Description:   stringPtr("Latest iPhone"),
						Price:         999.99,
						SKU:           "IPHONE15-001",
						StockQuantity: 100,
						CategoryID:    categoryID,
						IsActive:      true,
						CreatedAt:     now,
					}, nil)
			},
			want: &schema.ProductResponse{
				ProductID:     productID,
				ProductName:   "iPhone 15",
				Description:   stringPtr("Latest iPhone"),
				Price:         999.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
				IsActive:      true,
				CreatedAt:     now,
			},
			wantErr: false,
		},
		{
			name: "category not found",
			request: &schema.CreateProductRequest{
				ProductName:   "iPhone 15",
				Price:         999.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
			},
			setup: func() {
				mockCategoryRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(nil, errors.New("category not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repository error",
			request: &schema.CreateProductRequest{
				ProductName:   "iPhone 15",
				Price:         999.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
			},
			setup: func() {
				mockCategoryRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(&categoryEntity.Category{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
					}, nil)

				mockProductRepo.EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.CreateProduct(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProductServiceImpl_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := repository.NewMockProductRepository(ctrl)
	mockCategoryRepo := categoryRepo.NewMockCategoryRepository(ctrl)
	service := NewProductService(mockProductRepo, mockCategoryRepo)

	ctx := context.Background()
	productID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		productID uuid.UUID
		setup     func()
		want      *schema.ProductResponse
		wantErr   bool
	}{
		{
			name:      "successful get",
			productID: productID,
			setup: func() {
				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(&entity.Product{
						ProductID:     productID,
						ProductName:   "iPhone 15",
						Description:   stringPtr("Latest iPhone"),
						Price:         999.99,
						SKU:           "IPHONE15-001",
						StockQuantity: 100,
						CategoryID:    categoryID,
						IsActive:      true,
						CreatedAt:     now,
					}, nil)
			},
			want: &schema.ProductResponse{
				ProductID:     productID,
				ProductName:   "iPhone 15",
				Description:   stringPtr("Latest iPhone"),
				Price:         999.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
				IsActive:      true,
				CreatedAt:     now,
			},
			wantErr: false,
		},
		{
			name:      "product not found",
			productID: productID,
			setup: func() {
				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(nil, errors.New("product not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.GetProduct(ctx, tt.productID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProductServiceImpl_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := repository.NewMockProductRepository(ctrl)
	mockCategoryRepo := categoryRepo.NewMockCategoryRepository(ctrl)
	service := NewProductService(mockProductRepo, mockCategoryRepo)

	ctx := context.Background()
	productID := uuid.New()
	categoryID := uuid.New()
	newCategoryID := uuid.New()
	now := time.Now()

	existingProduct := &entity.Product{
		ProductID:     productID,
		ProductName:   "iPhone 14",
		Description:   stringPtr("Old iPhone"),
		Price:         899.99,
		SKU:           "IPHONE14-001",
		StockQuantity: 50,
		CategoryID:    categoryID,
		IsActive:      true,
		CreatedAt:     now,
	}

	tests := []struct {
		name      string
		productID uuid.UUID
		request   *schema.UpdateProductRequest
		setup     func()
		want      *schema.ProductResponse
		wantErr   bool
	}{
		{
			name:      "successful update with category change",
			productID: productID,
			request: &schema.UpdateProductRequest{
				ProductName: stringPtr("iPhone 15"),
				Price:       float64Ptr(999.99),
				CategoryID:  &newCategoryID,
			},
			setup: func() {
				// Get existing product
				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(existingProduct, nil)

				// Validate new category
				mockCategoryRepo.EXPECT().
					GetByID(ctx, newCategoryID).
					Return(&categoryEntity.Category{
						CategoryID:   newCategoryID,
						CategoryName: "Mobile Phones",
					}, nil)

				// Update product
				mockProductRepo.EXPECT().
					Update(ctx, gomock.Any()).
					Return(&entity.Product{
						ProductID:     productID,
						ProductName:   "iPhone 15",
						Description:   stringPtr("Old iPhone"),
						Price:         999.99,
						SKU:           "IPHONE14-001",
						StockQuantity: 50,
						CategoryID:    newCategoryID,
						IsActive:      true,
						CreatedAt:     now,
					}, nil)
			},
			want: &schema.ProductResponse{
				ProductID:     productID,
				ProductName:   "iPhone 15",
				Description:   stringPtr("Old iPhone"),
				Price:         999.99,
				SKU:           "IPHONE14-001",
				StockQuantity: 50,
				CategoryID:    newCategoryID,
				IsActive:      true,
				CreatedAt:     now,
			},
			wantErr: false,
		},
		{
			name:      "product not found",
			productID: productID,
			request: &schema.UpdateProductRequest{
				ProductName: stringPtr("iPhone 15"),
			},
			setup: func() {
				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(nil, errors.New("product not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.UpdateProduct(ctx, tt.productID, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProductServiceImpl_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := repository.NewMockProductRepository(ctrl)
	mockCategoryRepo := categoryRepo.NewMockCategoryRepository(ctrl)
	service := NewProductService(mockProductRepo, mockCategoryRepo)

	ctx := context.Background()
	productID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name      string
		productID uuid.UUID
		setup     func()
		wantErr   bool
	}{
		{
			name:      "successful delete",
			productID: productID,
			setup: func() {
				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(&entity.Product{
						ProductID:  productID,
						CategoryID: categoryID,
					}, nil)

				mockProductRepo.EXPECT().
					SoftDelete(ctx, productID).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "product not found",
			productID: productID,
			setup: func() {
				mockProductRepo.EXPECT().
					GetByID(ctx, productID).
					Return(nil, errors.New("product not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := service.DeleteProduct(ctx, tt.productID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestProductServiceImpl_ListProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := repository.NewMockProductRepository(ctrl)
	mockCategoryRepo := categoryRepo.NewMockCategoryRepository(ctrl)
	service := NewProductService(mockProductRepo, mockCategoryRepo)

	ctx := context.Background()
	productID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	products := []*entity.Product{
		{
			ProductID:     productID,
			ProductName:   "iPhone 15",
			Price:         999.99,
			SKU:           "IPHONE15-001",
			StockQuantity: 100,
			CategoryID:    categoryID,
			IsActive:      true,
			CreatedAt:     now,
		},
	}

	tests := []struct {
		name    string
		request *schema.ProductSearchRequest
		setup   func()
		want    *schema.ProductListResponse
		wantErr bool
	}{
		{
			name:    "list all products",
			request: nil,
			setup: func() {
				mockProductRepo.EXPECT().
					GetAll(ctx).
					Return(products, nil)
			},
			want: &schema.ProductListResponse{
				Products: []*schema.ProductResponse{{
					ProductID:     productID,
					ProductName:   "iPhone 15",
					Price:         999.99,
					SKU:           "IPHONE15-001",
					StockQuantity: 100,
					CategoryID:    categoryID,
					IsActive:      true,
					CreatedAt:     now,
				}},
				Total:   1,
				Page:    1,
				Limit:   1,
				HasNext: false,
			},
			wantErr: false,
		},
		{
			name: "list products by category",
			request: &schema.ProductSearchRequest{
				CategoryID: &categoryID,
			},
			setup: func() {
				mockProductRepo.EXPECT().
					GetByCategory(ctx, categoryID).
					Return(products, nil)
			},
			want: &schema.ProductListResponse{
				Products: []*schema.ProductResponse{{
					ProductID:     productID,
					ProductName:   "iPhone 15",
					Price:         999.99,
					SKU:           "IPHONE15-001",
					StockQuantity: 100,
					CategoryID:    categoryID,
					IsActive:      true,
					CreatedAt:     now,
				}},
				Total:   1,
				Page:    1,
				Limit:   1,
				HasNext: false,
			},
			wantErr: false,
		},
		{
			name: "list active products",
			request: &schema.ProductSearchRequest{
				IsActive: boolPtr(true),
			},
			setup: func() {
				mockProductRepo.EXPECT().
					GetActiveProducts(ctx).
					Return(products, nil)
			},
			want: &schema.ProductListResponse{
				Products: []*schema.ProductResponse{{
					ProductID:     productID,
					ProductName:   "iPhone 15",
					Price:         999.99,
					SKU:           "IPHONE15-001",
					StockQuantity: 100,
					CategoryID:    categoryID,
					IsActive:      true,
					CreatedAt:     now,
				}},
				Total:   1,
				Page:    1,
				Limit:   1,
				HasNext: false,
			},
			wantErr: false,
		},
		{
			name: "repository error",
			request: &schema.ProductSearchRequest{
				CategoryID: &categoryID,
			},
			setup: func() {
				mockProductRepo.EXPECT().
					GetByCategory(ctx, categoryID).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.ListProducts(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProductServiceImpl_GetProductsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := repository.NewMockProductRepository(ctrl)
	mockCategoryRepo := categoryRepo.NewMockCategoryRepository(ctrl)
	service := NewProductService(mockProductRepo, mockCategoryRepo)

	ctx := context.Background()
	productID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		categoryID uuid.UUID
		setup      func()
		want       *schema.ProductListResponse
		wantErr    bool
	}{
		{
			name:       "successful get products by category",
			categoryID: categoryID,
			setup: func() {
				mockCategoryRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(&categoryEntity.Category{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
					}, nil)

				mockProductRepo.EXPECT().
					GetByCategory(ctx, categoryID).
					Return([]*entity.Product{{
						ProductID:     productID,
						ProductName:   "iPhone 15",
						Price:         999.99,
						SKU:           "IPHONE15-001",
						StockQuantity: 100,
						CategoryID:    categoryID,
						IsActive:      true,
						CreatedAt:     now,
					}}, nil)
			},
			want: &schema.ProductListResponse{
				Products: []*schema.ProductResponse{{
					ProductID:     productID,
					ProductName:   "iPhone 15",
					Price:         999.99,
					SKU:           "IPHONE15-001",
					StockQuantity: 100,
					CategoryID:    categoryID,
					IsActive:      true,
					CreatedAt:     now,
				}},
				Total:   1,
				Page:    1,
				Limit:   1,
				HasNext: false,
			},
			wantErr: false,
		},
		{
			name:       "category not found",
			categoryID: categoryID,
			setup: func() {
				mockCategoryRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(nil, errors.New("category not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.GetProductsByCategory(ctx, tt.categoryID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
