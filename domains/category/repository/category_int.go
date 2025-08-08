package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/category/entity"
	"github.com/google/uuid"
)

type CategoryRepository interface {
	// Basic CRUD operations
	CreateCategory(ctx context.Context, categoryName string, parentID *uuid.UUID) (*entity.Category, error)
	CreateRootCategory(ctx context.Context, categoryName string) (*entity.Category, error)
	GetCategory(ctx context.Context, categoryID uuid.UUID) (*entity.Category, error)
	GetCategoryByName(ctx context.Context, categoryName string) (*entity.Category, error)
	UpdateCategory(ctx context.Context, categoryID uuid.UUID, categoryName string) (*entity.Category, error)
	SoftDeleteCategory(ctx context.Context, categoryID uuid.UUID) error

	// Hierarchy operations
	GetRootCategories(ctx context.Context) ([]*entity.Category, error)
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.Category, error)
	GetCategoryDescendants(ctx context.Context, categoryID uuid.UUID) ([]*entity.CategoryWithLevel, error)
	GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]*entity.CategoryWithLevel, error)
	MoveCategoryToParent(ctx context.Context, categoryID, newParentID uuid.UUID) (*entity.Category, error)

	// Listing and counting
	ListCategories(ctx context.Context) ([]*entity.Category, error)
	CountCategories(ctx context.Context) (int64, error)
	CountRootCategories(ctx context.Context) (int64, error)

	// Validation helpers
	CategoryExists(ctx context.Context, categoryID uuid.UUID) (bool, error)
	IsValidParent(ctx context.Context, categoryID, parentID uuid.UUID) (bool, error)
}
