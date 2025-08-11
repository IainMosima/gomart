package category

import (
	"context"

	"github.com/IainMosima/gomart/domains/category/entity"
	"github.com/IainMosima/gomart/domains/category/repository"
	"github.com/IainMosima/gomart/domains/category/schema"
	"github.com/IainMosima/gomart/domains/category/service"
	"github.com/google/uuid"
)

type CategoryServiceImpl struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) service.CategoryService {
	return &CategoryServiceImpl{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryServiceImpl) CreateCategory(ctx context.Context, req *schema.CreateCategoryRequest) (*schema.CategoryResponse, error) {
	category := &entity.Category{
		CategoryID:   uuid.New(),
		CategoryName: req.CategoryName,
		ParentID:     req.ParentID,
	}

	createdCategory, err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(createdCategory), nil
}

func (s *CategoryServiceImpl) GetCategory(ctx context.Context, categoryID uuid.UUID) (*schema.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(category), nil
}

func (s *CategoryServiceImpl) UpdateCategory(ctx context.Context, categoryID uuid.UUID, req *schema.UpdateCategoryRequest) (*schema.CategoryResponse, error) {
	existingCategory, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	existingCategory.CategoryName = req.CategoryName

	updatedCategory, err := s.categoryRepo.Update(ctx, existingCategory)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(updatedCategory), nil
}

func (s *CategoryServiceImpl) ListCategories(ctx context.Context) (*schema.CategoryListResponse, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*schema.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = s.entityToResponse(&category)
	}

	return &schema.CategoryListResponse{
		Categories: responses,
		Total:      int64(len(responses)),
	}, nil
}

func (s *CategoryServiceImpl) GetCategoryChildren(ctx context.Context, parentID *uuid.UUID) (*schema.CategoryListResponse, error) {
	categories, err := s.categoryRepo.GetByParent(ctx, parentID)
	if err != nil {
		return nil, err
	}

	responses := make([]*schema.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = s.entityToResponse(&category)
	}

	return &schema.CategoryListResponse{
		Categories: responses,
		Total:      int64(len(responses)),
	}, nil
}

func (s *CategoryServiceImpl) entityToResponse(category *entity.Category) *schema.CategoryResponse {
	return &schema.CategoryResponse{
		CategoryID:   category.CategoryID,
		CategoryName: category.CategoryName,
		ParentID:     category.ParentID,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}
}
