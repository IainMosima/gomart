package product

import (
	"net/http"

	"github.com/IainMosima/gomart/domains/product/schema"
	"github.com/IainMosima/gomart/domains/product/service"
	"github.com/IainMosima/gomart/rest-server/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandlerImpl struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandlerInterface {
	return &ProductHandlerImpl{
		productService: productService,
	}
}

func (h *ProductHandlerImpl) CreateProduct(c *gin.Context) {
	var req dtos.CreateProductRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schemaReq := &schema.CreateProductRequest{
		ProductName:   req.ProductName,
		Description:   req.Description,
		Price:         req.Price,
		SKU:           req.SKU,
		StockQuantity: req.StockQuantity,
		CategoryID:    req.CategoryID,
		IsActive:      req.IsActive,
	}

	result, err := h.productService.CreateProduct(c.Request.Context(), schemaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &dtos.ProductResponseDTO{
		ProductID:     result.ProductID,
		ProductName:   result.ProductName,
		Description:   result.Description,
		Price:         result.Price,
		SKU:           result.SKU,
		StockQuantity: result.StockQuantity,
		CategoryID:    result.CategoryID,
		IsActive:      result.IsActive,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ProductHandlerImpl) GetProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	result, err := h.productService.GetProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := &dtos.ProductResponseDTO{
		ProductID:     result.ProductID,
		ProductName:   result.ProductName,
		Description:   result.Description,
		Price:         result.Price,
		SKU:           result.SKU,
		StockQuantity: result.StockQuantity,
		CategoryID:    result.CategoryID,
		IsActive:      result.IsActive,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandlerImpl) UpdateProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req dtos.UpdateProductRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schemaReq := &schema.UpdateProductRequest{
		ProductName:   req.ProductName,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		CategoryID:    req.CategoryID,
		IsActive:      req.IsActive,
	}

	result, err := h.productService.UpdateProduct(c.Request.Context(), productID, schemaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &dtos.ProductResponseDTO{
		ProductID:     result.ProductID,
		ProductName:   result.ProductName,
		Description:   result.Description,
		Price:         result.Price,
		SKU:           result.SKU,
		StockQuantity: result.StockQuantity,
		CategoryID:    result.CategoryID,
		IsActive:      result.IsActive,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandlerImpl) DeleteProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *ProductHandlerImpl) ListProducts(c *gin.Context) {
	var req dtos.ProductSearchRequestDTO
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schemaReq := &schema.ProductSearchRequest{
		Query:      req.Query,
		CategoryID: req.CategoryID,
		IsActive:   req.IsActive,
		InStock:    req.InStock,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	result, err := h.productService.ListProducts(c.Request.Context(), schemaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	products := make([]*dtos.ProductResponseDTO, len(result.Products))
	for i, product := range result.Products {
		products[i] = &dtos.ProductResponseDTO{
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

	response := &dtos.ProductListResponseDTO{
		Products: products,
		Total:    result.Total,
		Page:     result.Page,
		Limit:    result.Limit,
		HasNext:  result.HasNext,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandlerImpl) GetProductsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	result, err := h.productService.GetProductsByCategory(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	products := make([]*dtos.ProductResponseDTO, len(result.Products))
	for i, product := range result.Products {
		products[i] = &dtos.ProductResponseDTO{
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

	response := &dtos.ProductListResponseDTO{
		Products: products,
		Total:    result.Total,
		Page:     result.Page,
		Limit:    result.Limit,
		HasNext:  result.HasNext,
	}

	c.JSON(http.StatusOK, response)
}
