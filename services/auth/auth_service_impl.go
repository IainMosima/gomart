package auth

import (
	"context"
	"fmt"

	"github.com/IainMosima/gomart/configs"
	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/IainMosima/gomart/domains/auth/service"
	"github.com/google/uuid"
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

func (a *AuthServiceImpl) GenerateAuthURL(req *schema.LoginRequest) (*schema.LoginResponse, error) {
	state := req.State
	if state == "" {
		state = uuid.New().String()
	}

	authURL := a.cognitoService.oauth2Config.AuthCodeURL(state)

	return &schema.LoginResponse{
		AuthURL: authURL,
		State:   state,
	}, nil
}

func (a *AuthServiceImpl) ExchangeCodeForToken(ctx context.Context, req *schema.TokenExchangeRequest) (*schema.TokenResponse, error) {
	return a.cognitoService.ExchangeCodeForTokens(ctx, req.Code)
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
