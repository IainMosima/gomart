package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/auth/schema"
)

type AuthService interface {
	GenerateAuthURL(req *schema.LoginRequest) (*schema.LoginResponse, error)
	ExchangeCodeForToken(ctx context.Context, req *schema.TokenExchangeRequest) (*schema.TokenResponse, error)
	RefreshAccessToken(ctx context.Context, req *schema.RefreshTokenRequest) (*schema.RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error)
}
