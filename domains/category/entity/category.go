package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	CategoryID   uuid.UUID  `json:"category_id" db:"category_id"`
	CategoryName string     `json:"category_name" db:"category_name"`
	ParentID     *uuid.UUID `json:"parent_id" db:"parent_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted    bool       `json:"is_deleted" db:"is_deleted"`
}

type CategoryWithLevel struct {
	Category
	Level int `json:"level" db:"level"`
}

type CategoryTree struct {
	Category
	Children []*CategoryTree `json:"children,omitempty"`
}

type CategoryPath struct {
	Categories []*Category `json:"path"`
	Depth      int         `json:"depth"`
}

type CategoryHierarchy struct {
	Category
	Level    int                  `json:"level" db:"level"`
	Children []*CategoryHierarchy `json:"children,omitempty"`
}
