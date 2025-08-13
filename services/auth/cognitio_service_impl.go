package auth

import (
	"context"
	"fmt"

	"github.com/IainMosima/gomart/configs"
	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type CognitoService struct {
	client       *cognitoidentityprovider.Client
	config       *configs.Config
	oidcProvider *oidc.Provider
	oauth2Config oauth2.Config
}

func NewCognitoService(cfg *configs.Config) (*CognitoService, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWSRegion),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := cognitoidentityprovider.NewFromConfig(awsConfig)

	issuerURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", cfg.AWSRegion, cfg.CognitoUserPoolID)
	oidcProvider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.CognitoClientID,
		ClientSecret: cfg.CognitoClientSecret,
		RedirectURL:  cfg.CognitoRedirectURI,
		Endpoint:     oidcProvider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "phone", "profile", "offline_access"},
	}

	return &CognitoService{
		client:       client,
		config:       cfg,
		oidcProvider: oidcProvider,
		oauth2Config: oauth2Config,
	}, nil
}

func (c *CognitoService) ExchangeCodeForTokens(ctx context.Context, code string) (*schema.TokenResponse, error) {
	rawToken, err := c.oauth2Config.Exchange(ctx, code)

	if err != nil {
		return nil, err
	}

	tokenResp := &schema.TokenResponse{
		AccessToken:  rawToken.AccessToken,
		RefreshToken: rawToken.RefreshToken,
		IDToken:      rawToken.Extra("id_token").(string),
	}

	return tokenResp, nil

}

func (c *CognitoService) ValidateAccessToken(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
	verifier := c.oidcProvider.Verifier(&oidc.Config{ClientID: c.config.CognitoClientID})

	idToken, err := verifier.Verify(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	var claims struct {
		Sub           string `json:"sub"`
		Username      string `json:"cognito:username"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		PhoneNumber   string `json:"phone_number"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}

	userInfo := &schema.UserInfoResponse{
		UserName:      claims.Username,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		PhoneNumber:   claims.PhoneNumber,
	}

	return userInfo, nil
}

func (c *CognitoService) RefreshTokens(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := c.oauth2Config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	tokenResp := &schema.TokenResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		IDToken:      newToken.Extra("id_token").(string),
	}

	return tokenResp, nil
}
