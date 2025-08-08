package entity

import (
	"time"

	"github.com/google/uuid"
)

type CustomerAuth struct {
	CustomerID     uuid.UUID  `json:"customer_id" db:"customer_id"`
	Email          string     `json:"email" db:"email"`
	FirstName      string     `json:"first_name" db:"first_name"`
	LastName       string     `json:"last_name" db:"last_name"`
	OpenIDSub      string     `json:"openid_sub" db:"openid_sub"`
	SetupCompleted bool       `json:"setup_completed" db:"setup_completed"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`
}

type TokenClaims struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Exp           int64  `json:"exp"`
	Iat           int64  `json:"iat"`
	Aud           string `json:"aud"`
	Iss           string `json:"iss"`
}
