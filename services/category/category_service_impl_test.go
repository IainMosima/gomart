package category

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/IainMosima/gomart/domains/category/entity"
	"github.com/IainMosima/gomart/domains/category/repository"
	"github.com/IainMosima/gomart/domains/category/schema"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewCategoryService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCategoryRepository(ctrl)
	service := NewCategoryService(mockRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &CategoryServiceImpl{}, service)
}

func TestCategoryServiceImpl_CreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCategoryRepository(ctrl)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	now := time.Now()
	categoryID := uuid.New()
	parentID := uuid.New()

	tests := []struct {
		name    string
		request *schema.CreateCategoryRequest
		setup   func()
		want    *schema.CategoryResponse
		wantErr bool
	}{
		{
			name: "successful creation with parent",
			request: &schema.CreateCategoryRequest{
				CategoryName: "Electronics",
				ParentID:     &parentID,
			},
			setup: func() {
				mockRepo.EXPECT().
					Create(ctx, gomock.Any()).
					Return(&entity.Category{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
						ParentID:     &parentID,
						CreatedAt:    now,
					}, nil)
			},
			want: &schema.CategoryResponse{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     &parentID,
				CreatedAt:    now,
			},
			wantErr: false,
		},
		{
			name: "successful creation without parent (root category)",
			request: &schema.CreateCategoryRequest{
				CategoryName: "Electronics",
				ParentID:     nil,
			},
			setup: func() {
				mockRepo.EXPECT().
					Create(ctx, gomock.Any()).
					Return(&entity.Category{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
						ParentID:     nil,
						CreatedAt:    now,
					}, nil)
			},
			want: &schema.CategoryResponse{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     nil,
				CreatedAt:    now,
			},
			wantErr: false,
		},
		{
			name: "repository error",
			request: &schema.CreateCategoryRequest{
				CategoryName: "Electronics",
				ParentID:     &parentID,
			},
			setup: func() {
				mockRepo.EXPECT().
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

			result, err := service.CreateCategory(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.CategoryName, result.CategoryName)
				assert.Equal(t, tt.want.ParentID, result.ParentID)
				assert.NotEmpty(t, result.CategoryID)
			}
		})
	}
}

func TestCategoryServiceImpl_GetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCategoryRepository(ctrl)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		categoryID uuid.UUID
		setup      func()
		want       *schema.CategoryResponse
		wantErr    bool
	}{
		{
			name:       "successful get",
			categoryID: categoryID,
			setup: func() {
				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(&entity.Category{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
						ParentID:     nil,
						CreatedAt:    now,
					}, nil)
			},
			want: &schema.CategoryResponse{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     nil,
				CreatedAt:    now,
			},
			wantErr: false,
		},
		{
			name:       "category not found",
			categoryID: categoryID,
			setup: func() {
				mockRepo.EXPECT().
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

			result, err := service.GetCategory(ctx, tt.categoryID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.CategoryID, result.CategoryID)
				assert.Equal(t, tt.want.CategoryName, result.CategoryName)
				assert.Equal(t, tt.want.ParentID, result.ParentID)
			}
		})
	}
}

func TestCategoryServiceImpl_UpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCategoryRepository(ctrl)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	categoryID := uuid.New()
	now := time.Now()
	updated := now.Add(time.Hour)

	tests := []struct {
		name       string
		categoryID uuid.UUID
		request    *schema.UpdateCategoryRequest
		setup      func()
		want       *schema.CategoryResponse
		wantErr    bool
	}{
		{
			name:       "successful update",
			categoryID: categoryID,
			request: &schema.UpdateCategoryRequest{
				CategoryName: "Updated Electronics",
			},
			setup: func() {
				existing := &entity.Category{
					CategoryID:   categoryID,
					CategoryName: "Electronics",
					ParentID:     nil,
					CreatedAt:    now,
				}
				updated := &entity.Category{
					CategoryID:   categoryID,
					CategoryName: "Updated Electronics",
					ParentID:     nil,
					CreatedAt:    now,
					UpdatedAt:    &updated,
				}

				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(existing, nil)

				mockRepo.EXPECT().
					Update(ctx, gomock.Any()).
					Return(updated, nil)
			},
			want: &schema.CategoryResponse{
				CategoryID:   categoryID,
				CategoryName: "Updated Electronics",
				ParentID:     nil,
				CreatedAt:    now,
				UpdatedAt:    &updated,
			},
			wantErr: false,
		},
		{
			name:       "category not found for update",
			categoryID: categoryID,
			request: &schema.UpdateCategoryRequest{
				CategoryName: "Updated Electronics",
			},
			setup: func() {
				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(nil, errors.New("category not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:       "update repository error",
			categoryID: categoryID,
			request: &schema.UpdateCategoryRequest{
				CategoryName: "Updated Electronics",
			},
			setup: func() {
				existing := &entity.Category{
					CategoryID:   categoryID,
					CategoryName: "Electronics",
					ParentID:     nil,
					CreatedAt:    now,
				}

				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(existing, nil)

				mockRepo.EXPECT().
					Update(ctx, gomock.Any()).
					Return(nil, errors.New("update failed"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := service.UpdateCategory(ctx, tt.categoryID, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.CategoryID, result.CategoryID)
				assert.Equal(t, tt.want.CategoryName, result.CategoryName)
			}
		})
	}
}

func TestCategoryServiceImpl_ListCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCategoryRepository(ctrl)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	now := time.Now()
	parentID := uuid.New()

	tests := []struct {
		name    string
		request *schema.ListCategoriesRequest
		setup   func()
		want    *schema.CategoryListResponse
		wantErr bool
	}{
		{
			name:    "successful list all categories (nil request)",
			request: nil,
			setup: func() {
				categories := []entity.Category{
					{
						CategoryID:   uuid.New(),
						CategoryName: "Electronics",
						ParentID:     nil,
						CreatedAt:    now,
					},
					{
						CategoryID:   uuid.New(),
						CategoryName: "Books",
						ParentID:     nil,
						CreatedAt:    now,
					},
				}
				mockRepo.EXPECT().
					GetAll(ctx).
					Return(categories, nil)
			},
			want: &schema.CategoryListResponse{
				Total: 2,
			},
			wantErr: false,
		},
		{
			name:    "successful list all categories (empty request)",
			request: &schema.ListCategoriesRequest{},
			setup: func() {
				categories := []entity.Category{
					{
						CategoryID:   uuid.New(),
						CategoryName: "Electronics",
						ParentID:     nil,
						CreatedAt:    now,
					},
				}
				mockRepo.EXPECT().
					GetAll(ctx).
					Return(categories, nil)
			},
			want: &schema.CategoryListResponse{
				Total: 1,
			},
			wantErr: false,
		},
		{
			name: "successful list root categories only",
			request: &schema.ListCategoriesRequest{
				RootOnly: true,
			},
			setup: func() {
				categories := []entity.Category{
					{
						CategoryID:   uuid.New(),
						CategoryName: "Electronics",
						ParentID:     nil,
						CreatedAt:    now,
					},
				}
				mockRepo.EXPECT().
					GetRootCategories(ctx).
					Return(categories, nil)
			},
			want: &schema.CategoryListResponse{
				Total: 1,
			},
			wantErr: false,
		},
		{
			name: "successful list by parent ID",
			request: &schema.ListCategoriesRequest{
				ParentID: &parentID,
			},
			setup: func() {
				categories := []entity.Category{
					{
						CategoryID:   uuid.New(),
						CategoryName: "Smartphones",
						ParentID:     &parentID,
						CreatedAt:    now,
					},
					{
						CategoryID:   uuid.New(),
						CategoryName: "Laptops",
						ParentID:     &parentID,
						CreatedAt:    now,
					},
				}
				mockRepo.EXPECT().
					GetByParent(ctx, &parentID).
					Return(categories, nil)
			},
			want: &schema.CategoryListResponse{
				Total: 2,
			},
			wantErr: false,
		},
		{
			name:    "repository error",
			request: &schema.ListCategoriesRequest{},
			setup: func() {
				mockRepo.EXPECT().
					GetAll(ctx).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repository error - root only",
			request: &schema.ListCategoriesRequest{
				RootOnly: true,
			},
			setup: func() {
				mockRepo.EXPECT().
					GetRootCategories(ctx).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repository error - by parent",
			request: &schema.ListCategoriesRequest{
				ParentID: &parentID,
			},
			setup: func() {
				mockRepo.EXPECT().
					GetByParent(ctx, &parentID).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := service.ListCategories(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.Total, result.Total)
				if tt.want.Total > 0 {
					assert.Len(t, result.Categories, int(tt.want.Total))
				}
			}
		})
	}
}

func TestCategoryServiceImpl_GetCategoryChildren(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCategoryRepository(ctrl)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	parentID := uuid.New()
	now := time.Now()

	tests := []struct {
		name     string
		parentID *uuid.UUID
		setup    func()
		want     *schema.CategoryListResponse
		wantErr  bool
	}{
		{
			name:     "successful get children with parent ID",
			parentID: &parentID,
			setup: func() {
				children := []entity.Category{
					{
						CategoryID:   uuid.New(),
						CategoryName: "Smartphones",
						ParentID:     &parentID,
						CreatedAt:    now,
					},
					{
						CategoryID:   uuid.New(),
						CategoryName: "Laptops",
						ParentID:     &parentID,
						CreatedAt:    now,
					},
				}
				mockRepo.EXPECT().
					GetByParent(ctx, &parentID).
					Return(children, nil)
			},
			want: &schema.CategoryListResponse{
				Total: 2,
			},
			wantErr: false,
		},
		{
			name:     "successful get root categories (nil parent)",
			parentID: nil,
			setup: func() {
				rootCategories := []entity.Category{
					{
						CategoryID:   uuid.New(),
						CategoryName: "Electronics",
						ParentID:     nil,
						CreatedAt:    now,
					},
				}
				mockRepo.EXPECT().
					GetByParent(ctx, (*uuid.UUID)(nil)).
					Return(rootCategories, nil)
			},
			want: &schema.CategoryListResponse{
				Total: 1,
			},
			wantErr: false,
		},
		{
			name:     "no children found",
			parentID: &parentID,
			setup: func() {
				mockRepo.EXPECT().
					GetByParent(ctx, &parentID).
					Return([]entity.Category{}, nil)
			},
			want: &schema.CategoryListResponse{
				Categories: []*schema.CategoryResponse{},
				Total:      0,
			},
			wantErr: false,
		},
		{
			name:     "repository error",
			parentID: &parentID,
			setup: func() {
				mockRepo.EXPECT().
					GetByParent(ctx, &parentID).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := service.GetCategoryChildren(ctx, tt.parentID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.Total, result.Total)
				if tt.want.Total > 0 {
					assert.Len(t, result.Categories, int(tt.want.Total))
				}
			}
		})
	}
}

func TestCategoryServiceImpl_entityToResponse(t *testing.T) {
	service := &CategoryServiceImpl{}
	categoryID := uuid.New()
	parentID := uuid.New()
	now := time.Now()
	updated := now.Add(time.Hour)

	tests := []struct {
		name     string
		category *entity.Category
		want     *schema.CategoryResponse
	}{
		{
			name: "convert entity with all fields",
			category: &entity.Category{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     &parentID,
				CreatedAt:    now,
				UpdatedAt:    &updated,
			},
			want: &schema.CategoryResponse{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     &parentID,
				CreatedAt:    now,
				UpdatedAt:    &updated,
			},
		},
		{
			name: "convert root category (no parent)",
			category: &entity.Category{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     nil,
				CreatedAt:    now,
				UpdatedAt:    nil,
			},
			want: &schema.CategoryResponse{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     nil,
				CreatedAt:    now,
				UpdatedAt:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.entityToResponse(tt.category)

			assert.Equal(t, tt.want.CategoryID, result.CategoryID)
			assert.Equal(t, tt.want.CategoryName, result.CategoryName)
			assert.Equal(t, tt.want.ParentID, result.ParentID)
			assert.Equal(t, tt.want.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.want.UpdatedAt, result.UpdatedAt)
		})
	}
}
func TestCategoryServiceImpl_GetCategoryAverageProductPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCategoryRepository(ctrl)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		categoryID uuid.UUID
		setup      func()
		want       *schema.CategoryAverageProductPriceResponse
		wantErr    bool
	}{
		{
			name:       "successful get average price",
			categoryID: categoryID,
			setup: func() {
				category := &entity.Category{
					CategoryID:   categoryID,
					CategoryName: "Electronics",
					ParentID:     nil,
					CreatedAt:    now,
				}

				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(category, nil)

				mockRepo.EXPECT().
					GetAverageProductPrice(ctx, categoryID).
					Return(299.99, nil)
			},
			want: &schema.CategoryAverageProductPriceResponse{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				AveragePrice: 299.99,
			},
			wantErr: false,
		},
		{
			name:       "successful get average price with zero (no products)",
			categoryID: categoryID,
			setup: func() {
				category := &entity.Category{
					CategoryID:   categoryID,
					CategoryName: "Empty Category",
					ParentID:     nil,
					CreatedAt:    now,
				}

				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(category, nil)

				mockRepo.EXPECT().
					GetAverageProductPrice(ctx, categoryID).
					Return(0.0, nil)
			},
			want: &schema.CategoryAverageProductPriceResponse{
				CategoryID:   categoryID,
				CategoryName: "Empty Category",
				AveragePrice: 0.0,
			},
			wantErr: false,
		},
		{
			name:       "category not found",
			categoryID: categoryID,
			setup: func() {
				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(nil, errors.New("category not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:       "average price calculation error",
			categoryID: categoryID,
			setup: func() {
				category := &entity.Category{
					CategoryID:   categoryID,
					CategoryName: "Electronics",
					ParentID:     nil,
					CreatedAt:    now,
				}

				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(category, nil)

				mockRepo.EXPECT().
					GetAverageProductPrice(ctx, categoryID).
					Return(0.0, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := service.GetCategoryAverageProductPrice(ctx, tt.categoryID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.CategoryID, result.CategoryID)
				assert.Equal(t, tt.want.CategoryName, result.CategoryName)
				assert.Equal(t, tt.want.AveragePrice, result.AveragePrice)
			}
		})
	}
}

func TestCategoryServiceImpl_DeleteCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCategoryRepository(ctrl)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		categoryID uuid.UUID
		setup      func()
		wantErr    bool
	}{
		{
			name:       "successful delete",
			categoryID: categoryID,
			setup: func() {
				category := &entity.Category{
					CategoryID:   categoryID,
					CategoryName: "Electronics",
					ParentID:     nil,
					CreatedAt:    now,
				}

				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(category, nil)

				mockRepo.EXPECT().
					SoftDelete(ctx, categoryID).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "category not found",
			categoryID: categoryID,
			setup: func() {
				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(nil, errors.New("category not found"))
			},
			wantErr: true,
		},
		{
			name:       "delete operation fails",
			categoryID: categoryID,
			setup: func() {
				category := &entity.Category{
					CategoryID:   categoryID,
					CategoryName: "Electronics",
					ParentID:     nil,
					CreatedAt:    now,
				}

				mockRepo.EXPECT().
					GetByID(ctx, categoryID).
					Return(category, nil)

				mockRepo.EXPECT().
					SoftDelete(ctx, categoryID).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := service.DeleteCategory(ctx, tt.categoryID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
