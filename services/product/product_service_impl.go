package product

import (
	"context"
	"fmt"

	categoryRepo "github.com/IainMosima/gomart/domains/category/repository"
	"github.com/IainMosima/gomart/domains/product/entity"
	"github.com/IainMosima/gomart/domains/product/repository"
	"github.com/IainMosima/gomart/domains/product/schema"
	"github.com/IainMosima/gomart/domains/product/service"
	"github.com/google/uuid"
)

type ProductServiceImpl struct {
	productRepo  repository.ProductRepository
	categoryRepo categoryRepo.CategoryRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo categoryRepo.CategoryRepository) service.ProductService {
	return &ProductServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *schema.CreateProductRequest) (*schema.ProductResponse, error) {
	// Validate that category exists
	_, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	product := &entity.Product{
		ProductID:     uuid.New(),
		ProductName:   req.ProductName,
		Description:   req.Description,
		Price:         req.Price,
		SKU:           req.SKU,
		StockQuantity: req.StockQuantity,
		CategoryID:    req.CategoryID,
		IsActive:      isActive,
	}

	createdProduct, err := s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(createdProduct), nil
}

func (s *ProductServiceImpl) GetProduct(ctx context.Context, productID uuid.UUID) (*schema.ProductResponse, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(product), nil
}

func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, productID uuid.UUID, req *schema.UpdateProductRequest) (*schema.ProductResponse, error) {
	existingProduct, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Validate category if it's being updated
	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		existingProduct.CategoryID = *req.CategoryID
	}

	// Update fields if provided
	if req.ProductName != nil {
		existingProduct.ProductName = *req.ProductName
	}

	if req.Description != nil {
		existingProduct.Description = req.Description
	}

	if req.Price != nil {
		existingProduct.Price = *req.Price
	}

	if req.StockQuantity != nil {
		existingProduct.StockQuantity = *req.StockQuantity
	}

	if req.IsActive != nil {
		existingProduct.IsActive = *req.IsActive
	}

	updatedProduct, err := s.productRepo.Update(ctx, existingProduct)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(updatedProduct), nil
}

func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	// Check if product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	return s.productRepo.SoftDelete(ctx, productID)
}

func (s *ProductServiceImpl) ListProducts(ctx context.Context, req *schema.ProductSearchRequest) (*schema.ProductListResponse, error) {
	var products []*entity.Product
	var err error

	if req != nil && req.CategoryID != nil {
		products, err = s.productRepo.GetByCategory(ctx, *req.CategoryID)
	} else if req != nil && req.IsActive != nil && *req.IsActive {
		products, err = s.productRepo.GetActiveProducts(ctx)
	} else if req != nil && req.InStock != nil && *req.InStock {
		products, err = s.productRepo.GetInStock(ctx)
	} else {
		products, err = s.productRepo.GetAll(ctx)
	}

	if err != nil {
		return nil, err
	}

	responses := make([]*schema.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = s.entityToResponse(product)
	}

	// Apply pagination if specified
	total := int64(len(responses))
	if req != nil && req.Page > 0 && req.Limit > 0 {
		start := (req.Page - 1) * req.Limit
		if start >= len(responses) {
			responses = []*schema.ProductResponse{}
		} else {
			end := start + req.Limit
			if end > len(responses) {
				end = len(responses)
			}
			responses = responses[start:end]
		}
	}

	hasNext := false
	if req != nil && req.Page > 0 && req.Limit > 0 {
		hasNext = int64((req.Page * req.Limit)) < total
	}

	page := 1
	limit := len(responses)
	if req != nil && req.Page > 0 {
		page = req.Page
	}
	if req != nil && req.Limit > 0 {
		limit = req.Limit
	}

	return &schema.ProductListResponse{
		Products: responses,
		Total:    total,
		Page:     page,
		Limit:    limit,
		HasNext:  hasNext,
	}, nil
}

func (s *ProductServiceImpl) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) (*schema.ProductListResponse, error) {
	// Validate that category exists
	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	products, err := s.productRepo.GetByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	responses := make([]*schema.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = s.entityToResponse(product)
	}

	return &schema.ProductListResponse{
		Products: responses,
		Total:    int64(len(responses)),
		Page:     1,
		Limit:    len(responses),
		HasNext:  false,
	}, nil
}

func (s *ProductServiceImpl) entityToResponse(product *entity.Product) *schema.ProductResponse {
	return &schema.ProductResponse{
		ProductID:     product.ProductID,
		ProductName:   product.ProductName,
		Description:   product.Description,
		Price:         product.Price,
		SKU:           product.SKU,
		StockQuantity: product.StockQuantity,
		CategoryID:    product.CategoryID,
		IsActive:      product.IsActive,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}
}
