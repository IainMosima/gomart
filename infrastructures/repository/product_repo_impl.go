package repository

import (
	"context"
	"fmt"

	"github.com/IainMosima/gomart/domains/product/entity"
	domainRepo "github.com/IainMosima/gomart/domains/product/repository"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductRepositoryImpl struct {
	store db.Store
}

func NewProductRepository(store db.Store) domainRepo.ProductRepository {
	return &ProductRepositoryImpl{
		store: store,
	}
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	var description pgtype.Text
	if product.Description != nil {
		description = pgtype.Text{
			String: *product.Description,
			Valid:  true,
		}
	}

	price := pgtype.Numeric{}
	priceStr := fmt.Sprintf("%.2f", product.Price)
	if err := price.Scan(priceStr); err != nil {
		return nil, err
	}

	var isActive pgtype.Bool
	if err := isActive.Scan(product.IsActive); err != nil {
		return nil, err
	}

	params := db.CreateProductParams{
		ProductName:   product.ProductName,
		Description:   description,
		Price:         price,
		Sku:           product.SKU,
		StockQuantity: product.StockQuantity,
		CategoryID:    product.CategoryID,
		IsActive:      isActive,
	}

	result, err := r.store.CreateProduct(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(result), nil
}

func (r *ProductRepositoryImpl) GetByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {
	result, err := r.store.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(result), nil
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	var description pgtype.Text
	if product.Description != nil {
		description = pgtype.Text{
			String: *product.Description,
			Valid:  true,
		}
	}

	price := pgtype.Numeric{}
	priceStr := fmt.Sprintf("%.2f", product.Price)
	if err := price.Scan(priceStr); err != nil {
		return nil, err
	}

	var isActive pgtype.Bool
	if err := isActive.Scan(product.IsActive); err != nil {
		return nil, err
	}

	params := db.UpdateProductParams{
		ProductID:     product.ProductID,
		ProductName:   product.ProductName,
		Description:   description,
		Price:         price,
		StockQuantity: product.StockQuantity,
		CategoryID:    product.CategoryID,
		IsActive:      isActive,
	}

	result, err := r.store.UpdateProduct(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(result), nil
}

func (r *ProductRepositoryImpl) Delete(ctx context.Context, productID uuid.UUID) error {
	return r.store.DeleteProduct(ctx, productID)
}

func (r *ProductRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Product, error) {
	results, err := r.store.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*entity.Product, len(results))
	for i, result := range results {
		products[i] = r.convertToEntity(result)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.Product, error) {
	results, err := r.store.GetProductsByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	products := make([]*entity.Product, len(results))
	for i, result := range results {
		products[i] = r.convertToEntity(result)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) convertToEntity(dbProduct db.Product) *entity.Product {
	product := &entity.Product{
		ProductID:     dbProduct.ProductID,
		ProductName:   dbProduct.ProductName,
		SKU:           dbProduct.Sku,
		StockQuantity: dbProduct.StockQuantity,
		CategoryID:    dbProduct.CategoryID,
		CreatedAt:     dbProduct.CreatedAt.Time,
		IsDeleted:     dbProduct.IsDeleted.Bool,
	}

	if dbProduct.Description.Valid {
		product.Description = &dbProduct.Description.String
	}

	if dbProduct.Price.Valid {
		price, err := dbProduct.Price.Float64Value()
		if err == nil {
			product.Price = price.Float64
		}
	}

	if dbProduct.IsActive.Valid {
		product.IsActive = dbProduct.IsActive.Bool
	}

	if dbProduct.UpdatedAt.Valid {
		updatedAt := dbProduct.UpdatedAt.Time
		product.UpdatedAt = &updatedAt
	}

	return product
}
