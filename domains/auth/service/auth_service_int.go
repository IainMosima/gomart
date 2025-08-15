package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/auth/schema"
)

type AuthService interface {
	RefreshAccessToken(ctx context.Context, req *schema.RefreshTokenRequest) (*schema.RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error)
	HandleCallback(ctx context.Context, req *schema.HandleCallbackRequest) (*schema.TokenResponse, error)
	GetAuthURL(state string) (string, error)
}
