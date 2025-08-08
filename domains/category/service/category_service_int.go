package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/category/schema"
	"github.com/google/uuid"
)

type CategoryService interface {
	// Basic operations
	CreateCategory(ctx context.Context, req *schema.CreateCategoryRequest) (*schema.CategoryResponse, error)
	GetCategory(ctx context.Context, categoryID uuid.UUID) (*schema.CategoryResponse, error)
	UpdateCategory(ctx context.Context, categoryID uuid.UUID, req *schema.UpdateCategoryRequest) (*schema.CategoryResponse, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error

	// Hierarchy operations
	GetCategoryTree(ctx context.Context, req *schema.GetCategoryTreeRequest) ([]*schema.CategoryTreeResponse, error)
	GetCategoryPath(ctx context.Context, categoryID uuid.UUID) (*schema.CategoryPathResponse, error)
	MoveCategory(ctx context.Context, categoryID uuid.UUID, req *schema.MoveCategoryRequest) (*schema.CategoryResponse, error)

	// Listing operations
	ListCategories(ctx context.Context, req *schema.ListCategoriesRequest) (*schema.CategoryListResponse, error)
	GetRootCategories(ctx context.Context) (*schema.CategoryListResponse, error)
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID) (*schema.CategoryListResponse, error)

	// Statistics and validation
	GetCategoryStats(ctx context.Context) (*schema.CategoryStatsResponse, error)
	ValidateCategory(ctx context.Context, categoryID uuid.UUID) error
	ValidateParentChild(ctx context.Context, categoryID, parentID uuid.UUID) error
}
