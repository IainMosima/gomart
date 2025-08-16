package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/auth/schema"
)

type CognitoServiceInterface interface {
	ExchangeCodeForTokens(ctx context.Context, code string) (*schema.TokenResponse, error)
	ValidateAccessToken(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*schema.TokenResponse, error)
	ParseIDToken(ctx context.Context, idToken string) (*schema.CognitoUserInfoJWTClaims, error)
	GetAuthURL(state string) (string, error)
}
