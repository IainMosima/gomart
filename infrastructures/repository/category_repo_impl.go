package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/category/entity"
	domainRepo "github.com/IainMosima/gomart/domains/category/repository"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CategoryRepositoryImpl struct {
	store db.Store
}

func NewCategoryRepository(store db.Store) domainRepo.CategoryRepository {
	return &CategoryRepositoryImpl{
		store: store,
	}
}

func (r *CategoryRepositoryImpl) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	var parentID pgtype.UUID
	if category.ParentID != nil {
		parentID = pgtype.UUID{
			Bytes: *category.ParentID,
			Valid: true,
		}
	}

	params := db.CreateCategoryParams{
		CategoryName: category.CategoryName,
		ParentID:     parentID,
	}

	result, err := r.store.CreateCategory(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(result), nil
}

func (r *CategoryRepositoryImpl) GetByID(ctx context.Context, categoryID uuid.UUID) (*entity.Category, error) {
	result, err := r.store.GetCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(result), nil
}

func (r *CategoryRepositoryImpl) GetByParent(ctx context.Context, parentID *uuid.UUID) ([]entity.Category, error) {
	var pgParentID pgtype.UUID
	if parentID != nil {
		pgParentID = pgtype.UUID{
			Bytes: *parentID,
			Valid: true,
		}
	}

	results, err := r.store.GetCategoryChildren(ctx, pgParentID)
	if err != nil {
		return nil, err
	}

	categories := make([]entity.Category, len(results))
	for i, result := range results {
		categories[i] = *r.convertToEntity(result)
	}

	return categories, nil
}

func (r *CategoryRepositoryImpl) GetRootCategories(ctx context.Context) ([]entity.Category, error) {
	results, err := r.store.GetRootCategories(ctx)
	if err != nil {
		return nil, err
	}

	categories := make([]entity.Category, len(results))
	for i, result := range results {
		categories[i] = *r.convertToEntity(result)
	}

	return categories, nil
}

func (r *CategoryRepositoryImpl) GetAll(ctx context.Context) ([]entity.Category, error) {
	results, err := r.store.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	categories := make([]entity.Category, len(results))
	for i, result := range results {
		categories[i] = *r.convertToEntity(result)
	}

	return categories, nil
}

func (r *CategoryRepositoryImpl) Update(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	params := db.UpdateCategoryParams{
		CategoryID:   category.CategoryID,
		CategoryName: category.CategoryName,
	}

	result, err := r.store.UpdateCategory(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(result), nil
}

func (r *CategoryRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.store.SoftDeleteCategory(ctx, id)
}

func (r *CategoryRepositoryImpl) convertToEntity(dbCategory db.Category) *entity.Category {
	category := &entity.Category{
		CategoryID:   dbCategory.CategoryID,
		CategoryName: dbCategory.CategoryName,
		CreatedAt:    dbCategory.CreatedAt.Time,
		IsDeleted:    dbCategory.IsDeleted.Bool,
	}

	if dbCategory.ParentID.Valid {
		parentUUID := uuid.UUID(dbCategory.ParentID.Bytes)
		category.ParentID = &parentUUID
	}

	if dbCategory.UpdatedAt.Valid {
		updatedAt := dbCategory.UpdatedAt.Time
		category.UpdatedAt = &updatedAt
	}

	return category
}
