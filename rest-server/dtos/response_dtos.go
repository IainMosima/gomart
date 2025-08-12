package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResponseDTO struct {
	CategoryID   uuid.UUID  `json:"category_id"`
	CategoryName string     `json:"category_name"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type CategoryListResponseDTO struct {
	Categories []*CategoryResponseDTO `json:"categories"`
	Total      int64                  `json:"total"`
}

type CategoryAverageProductPriceResponseDTO struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	AveragePrice float64   `json:"average_price"`
}
