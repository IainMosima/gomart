package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

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
		Scopes:       []string{oidc.ScopeOpenID, "email", "phone"},
	}

	return &CognitoService{
		client:       client,
		config:       cfg,
		oidcProvider: oidcProvider,
		oauth2Config: oauth2Config,
	}, nil
}

func (c *CognitoService) GetAuthURL(state string) (string, error) {
	baseURL := fmt.Sprintf("https://%s/login", c.config.CognitoDomain)

	authURL := fmt.Sprintf("%s?client_id=%s&response_type=code&scope=email+openid+phone&redirect_uri=%s",
		baseURL,
		c.config.CognitoClientID,
		url.QueryEscape(c.config.CognitoRedirectURI))

	if state != "" {
		authURL += "&state=" + url.QueryEscape(state)
	}

	log.Printf("Generated Cognito auth URL: %s", authURL)

	return authURL, nil
}

func (c *CognitoService) ParseIDToken(ctx context.Context, idToken string) (*schema.CognitoUserInfoJWTClaims, error) {
	verifier := c.oidcProvider.Verifier(&oidc.Config{ClientID: c.config.CognitoClientID})

	token, err := verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	var claims schema.CognitoUserInfoJWTClaims
	if err := token.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}

	return &claims, nil
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
	userInfoURL := fmt.Sprintf("https://%s/oauth2/userInfo", c.config.CognitoDomain)

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call userInfo endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userInfo request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfoClaims struct {
		Sub           string `json:"sub"`
		Username      string `json:"username"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		PhoneNumber   string `json:"phone_number"`
	}

	if err := json.Unmarshal(body, &userInfoClaims); err != nil {
		return nil, fmt.Errorf("failed to parse userInfo response: %w", err)
	}

	userInfo := &schema.UserInfoResponse{
		UserName:      userInfoClaims.Username,
		Email:         userInfoClaims.Email,
		EmailVerified: userInfoClaims.EmailVerified == "true",
		PhoneNumber:   userInfoClaims.PhoneNumber,
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
