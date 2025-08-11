package dtos

import "github.com/google/uuid"

type CreateCategoryRequestDTO struct {
	CategoryName string     `json:"category_name" validate:"required,min=1,max=255"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
}

type UpdateCategoryRequestDTO struct {
	CategoryName string `json:"category_name" validate:"required,min=1,max=255"`
}
