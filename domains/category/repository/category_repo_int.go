package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/category/entity"
	"github.com/google/uuid"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) (*entity.Category, error)
	GetByID(ctx context.Context, categoryID uuid.UUID) (*entity.Category, error)
	GetByParent(ctx context.Context, parentID *uuid.UUID) ([]entity.Category, error)
	GetRootCategories(ctx context.Context) ([]entity.Category, error)
	GetAll(ctx context.Context) ([]entity.Category, error)
	Update(ctx context.Context, category *entity.Category) (*entity.Category, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
