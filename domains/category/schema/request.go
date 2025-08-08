package schema

import "github.com/google/uuid"

type CreateCategoryRequest struct {
	CategoryName string     `json:"category_name" validate:"required,min=1,max=255"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
}

type UpdateCategoryRequest struct {
	CategoryName string `json:"category_name" validate:"required,min=1,max=255"`
}

type MoveCategoryRequest struct {
	NewParentID *uuid.UUID `json:"new_parent_id"`
}

type GetCategoryTreeRequest struct {
	RootCategoryID *uuid.UUID `json:"root_category_id,omitempty"`
	MaxDepth       *int       `json:"max_depth,omitempty"`
}

type ListCategoriesRequest struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	RootOnly    bool       `json:"root_only,omitempty"`
	IncludeTree bool       `json:"include_tree,omitempty"`
}
