package schema

import "github.com/google/uuid"

type CreateCustomerRequest struct {
	Email          string  `json:"email" validate:"required,email,max=255"`
	FirstName      string  `json:"first_name" validate:"required,min=1,max=100"`
	LastName       string  `json:"last_name" validate:"required,min=1,max=100"`
	Phone          *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Address        *string `json:"address,omitempty"`
	City           *string `json:"city,omitempty" validate:"omitempty,max=100"`
	PostalCode     *string `json:"postal_code,omitempty" validate:"omitempty,max=20"`
	OpenIDSub      string  `json:"openid_sub" validate:"required,max=255"`
	SetupCompleted *bool   `json:"setup_completed,omitempty"`
}

type UpdateCustomerRequest struct {
	Email      *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	FirstName  *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName   *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	Phone      *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty" validate:"omitempty,max=100"`
	PostalCode *string `json:"postal_code,omitempty" validate:"omitempty,max=20"`
}

type UpdateSetupStatusRequest struct {
	SetupCompleted bool `json:"setup_completed"`
}

type CustomerSearchRequest struct {
	Query          *string `json:"query,omitempty"`
	Email          *string `json:"email,omitempty" validate:"omitempty,email"`
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	SetupCompleted *bool   `json:"setup_completed,omitempty"`
	Page           int     `json:"page" validate:"min=1"`
	Limit          int     `json:"limit" validate:"min=1,max=100"`
}

type CustomerListRequest struct {
	SetupCompleted *bool `json:"setup_completed,omitempty"`
	Page           int   `json:"page" validate:"min=1"`
	Limit          int   `json:"limit" validate:"min=1,max=100"`
}
