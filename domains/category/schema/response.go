package schema

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	CategoryID   uuid.UUID  `json:"category_id"`
	CategoryName string     `json:"category_name"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type CategoryWithLevelResponse struct {
	CategoryResponse
	Level int `json:"level"`
}

type CategoryTreeResponse struct {
	CategoryResponse
	Children []*CategoryTreeResponse `json:"children,omitempty"`
}

type CategoryPathResponse struct {
	Path  []*CategoryWithLevelResponse `json:"path"`
	Depth int                          `json:"depth"`
}

type CategoryListResponse struct {
	Categories []*CategoryResponse `json:"categories"`
	Total      int64               `json:"total"`
}

type CategoryTreeListResponse struct {
	Trees []*CategoryTreeResponse `json:"trees"`
	Total int64                   `json:"total"`
}

type CategoryStatsResponse struct {
	TotalCategories int64 `json:"total_categories"`
	RootCategories  int64 `json:"root_categories"`
	MaxDepth        int   `json:"max_depth"`
}

type CategoryAverageProductPriceResponse struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	AveragePrice float64   `json:"average_price"`
}
