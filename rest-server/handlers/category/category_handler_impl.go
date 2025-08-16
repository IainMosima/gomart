package category

import (
	"net/http"

	"github.com/IainMosima/gomart/domains/category/schema"
	"github.com/IainMosima/gomart/domains/category/service"
	"github.com/IainMosima/gomart/rest-server/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandlerImpl struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandlerInterface {
	return &CategoryHandlerImpl{
		categoryService: categoryService,
	}
}

func (h *CategoryHandlerImpl) CreateCategory(c *gin.Context) {
	var req dtos.CreateCategoryRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schemaReq := &schema.CreateCategoryRequest{
		CategoryName: req.CategoryName,
		ParentID:     req.ParentID,
	}

	result, err := h.categoryService.CreateCategory(c.Request.Context(), schemaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &dtos.CategoryResponseDTO{
		CategoryID:   result.CategoryID,
		CategoryName: result.CategoryName,
		ParentID:     result.ParentID,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *CategoryHandlerImpl) GetCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	result, err := h.categoryService.GetCategory(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := &dtos.CategoryResponseDTO{
		CategoryID:   result.CategoryID,
		CategoryName: result.CategoryName,
		ParentID:     result.ParentID,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandlerImpl) UpdateCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req dtos.UpdateCategoryRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schemaReq := &schema.UpdateCategoryRequest{
		CategoryName: req.CategoryName,
	}

	result, err := h.categoryService.UpdateCategory(c.Request.Context(), categoryID, schemaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &dtos.CategoryResponseDTO{
		CategoryID:   result.CategoryID,
		CategoryName: result.CategoryName,
		ParentID:     result.ParentID,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandlerImpl) ListCategories(c *gin.Context) {
	var req dtos.ListCategoriesRequestDTO

	if c.Request.Method == "GET" {
		req.RootOnly = c.Query("root_only") == "true"

		if parentIDStr := c.Query("parent_id"); parentIDStr != "" {
			parentID, err := uuid.Parse(parentIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent ID format"})
				return
			}
			req.ParentID = &parentID
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	schemaReq := &schema.ListCategoriesRequest{
		ParentID: req.ParentID,
		RootOnly: req.RootOnly,
	}

	result, err := h.categoryService.ListCategories(c.Request.Context(), schemaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categories := make([]*dtos.CategoryResponseDTO, len(result.Categories))
	for i, cat := range result.Categories {
		categories[i] = &dtos.CategoryResponseDTO{
			CategoryID:   cat.CategoryID,
			CategoryName: cat.CategoryName,
			ParentID:     cat.ParentID,
			CreatedAt:    cat.CreatedAt,
			UpdatedAt:    cat.UpdatedAt,
		}
	}

	response := &dtos.CategoryListResponseDTO{
		Categories: categories,
		Total:      result.Total,
	}

	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandlerImpl) DeleteCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *CategoryHandlerImpl) GetCategoryChildren(c *gin.Context) {
	categoryIDStr := c.Param("id")

	var parentID *uuid.UUID
	if categoryIDStr != "" {
		id, err := uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		parentID = &id
	}

	result, err := h.categoryService.GetCategoryChildren(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categories := make([]*dtos.CategoryResponseDTO, len(result.Categories))
	for i, cat := range result.Categories {
		categories[i] = &dtos.CategoryResponseDTO{
			CategoryID:   cat.CategoryID,
			CategoryName: cat.CategoryName,
			ParentID:     cat.ParentID,
			CreatedAt:    cat.CreatedAt,
			UpdatedAt:    cat.UpdatedAt,
		}
	}

	response := &dtos.CategoryListResponseDTO{
		Categories: categories,
		Total:      result.Total,
	}

	c.JSON(http.StatusOK, response)
}
func (h *CategoryHandlerImpl) GetCategoryAverageProductPrice(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	result, err := h.categoryService.GetCategoryAverageProductPrice(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := &dtos.CategoryAverageProductPriceResponseDTO{
		CategoryID:   result.CategoryID,
		CategoryName: result.CategoryName,
		AveragePrice: result.AveragePrice,
	}

	c.JSON(http.StatusOK, response)
}
