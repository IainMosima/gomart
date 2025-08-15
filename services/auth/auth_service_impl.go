package auth

import (
	"context"
	"fmt"

	"github.com/IainMosima/gomart/configs"
	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/IainMosima/gomart/domains/auth/service"
)

type AuthServiceImpl struct {
	cognitoService *CognitoService
}

func NewAuthServiceImpl(cfg *configs.Config) (service.AuthService, error) {
	cognitoService, err := NewCognitoService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create cognito service: %w", err)
	}

	return &AuthServiceImpl{
		cognitoService: cognitoService,
	}, nil
}

func (a *AuthServiceImpl) GetAuthURL(state string) (string, error) {
	url, err := a.cognitoService.GetAuthURL(state)
	if err != nil {
		return "", fmt.Errorf("failed to get auth url: %w", err)
	}
	return url, nil
}

func (a *AuthServiceImpl) HandleCallback(ctx context.Context, req *schema.HandleCallbackRequest) (*schema.TokenResponse, error) {
	if req.Code == nil {
		return nil, fmt.Errorf("code is required")
	}

	tokenResp, err := a.cognitoService.ExchangeCodeForTokens(ctx, *req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}

	return tokenResp, nil

}

func (a *AuthServiceImpl) RefreshAccessToken(ctx context.Context, req *schema.RefreshTokenRequest) (*schema.RefreshTokenResponse, error) {
	tokenResp, err := a.cognitoService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &schema.RefreshTokenResponse{
		AccessToken: tokenResp.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}, nil
}

func (a *AuthServiceImpl) ValidateToken(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
	return a.cognitoService.ValidateAccessToken(ctx, accessToken)
}
