package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/google/uuid"
)

type AuthService interface {
	// OpenID Connect flow
	GenerateAuthURL(ctx context.Context, req *schema.LoginRequest) (*schema.LoginResponse, error)
	ExchangeCodeForToken(ctx context.Context, req *schema.TokenExchangeRequest) (*schema.TokenResponse, error)
	RefreshAccessToken(ctx context.Context, req *schema.RefreshTokenRequest) (*schema.RefreshTokenResponse, error)

	// Token validation and user info
	ValidateToken(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error)
	GetUserInfo(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error)

	// Customer management
	GetOrCreateCustomer(ctx context.Context, userInfo *schema.UserInfoResponse) (*schema.CustomerResponse, error)
	CompleteCustomerSetup(ctx context.Context, customerID uuid.UUID, req *schema.CompleteSetupRequest) (*schema.CustomerResponse, error)
	GetCustomerProfile(ctx context.Context, customerID uuid.UUID) (*schema.CustomerResponse, error)

	// Logout
	RevokeToken(ctx context.Context, req *schema.LogoutRequest) (*schema.LogoutResponse, error)
}
