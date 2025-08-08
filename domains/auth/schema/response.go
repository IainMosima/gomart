package schema

import (
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

type TokenResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	TokenType    string           `json:"token_type"`
	ExpiresIn    int              `json:"expires_in"`
	Customer     CustomerResponse `json:"customer"`
}

type CustomerResponse struct {
	CustomerID     uuid.UUID  `json:"customer_id"`
	Email          string     `json:"email"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Phone          *string    `json:"phone,omitempty"`
	Address        *string    `json:"address,omitempty"`
	City           *string    `json:"city,omitempty"`
	PostalCode     *string    `json:"postal_code,omitempty"`
	SetupCompleted bool       `json:"setup_completed"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type UserInfoResponse struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

type LogoutResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
