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
	ListCategories(ctx context.Context) (*schema.CategoryListResponse, error)
	GetCategoryChildren(ctx context.Context, parentID *uuid.UUID) (*schema.CategoryListResponse, error)
}
