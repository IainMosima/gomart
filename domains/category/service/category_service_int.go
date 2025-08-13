//go:generate mockgen -source=category_service_int.go -destination=category_service_mock.go -package=service

package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/category/schema"
	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *schema.CreateCategoryRequest) (*schema.CategoryResponse, error)
	GetCategory(ctx context.Context, categoryID uuid.UUID) (*schema.CategoryResponse, error)
	UpdateCategory(ctx context.Context, categoryID uuid.UUID, req *schema.UpdateCategoryRequest) (*schema.CategoryResponse, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	ListCategories(ctx context.Context, req *schema.ListCategoriesRequest) (*schema.CategoryListResponse, error)
	GetCategoryChildren(ctx context.Context, parentID *uuid.UUID) (*schema.CategoryListResponse, error)
	GetCategoryAverageProductPrice(ctx context.Context, categoryID uuid.UUID) (*schema.CategoryAverageProductPriceResponse, error)
}
