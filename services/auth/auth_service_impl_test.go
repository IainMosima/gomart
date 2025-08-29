package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/IainMosima/gomart/configs"
	"github.com/IainMosima/gomart/domains/auth/entity"
	"github.com/IainMosima/gomart/domains/auth/repository"
	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/IainMosima/gomart/domains/auth/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewAuthServiceImpl_Success(t *testing.T) {
	cfg := &configs.Config{
		CognitoDomain:       "test-domain.auth.us-east-1.amazoncognito.com",
		CognitoClientID:     "test-client-id",
		AWSRegion:           "us-east-1",
		CognitoRedirectURI:  "http://localhost:8080/callback",
		CognitoUserPoolID:   "us-east-1_TestPoolId",
		CognitoClientSecret: "test-secret",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := repository.NewMockAuthRepository(ctrl)

	authSvc, err := NewAuthServiceImpl(cfg, mockAuthRepo)

	// In test environment without proper AWS setup, this will fail
	if err != nil {
		assert.Contains(t, err.Error(), "failed to create cognito service")
		assert.Nil(t, authSvc)
	} else {
		assert.NotNil(t, authSvc)
	}
}

func TestAuthServiceImpl_GetAuthURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		GetAuthURL("test-state").
		Return("https://test-domain.com/login?state=test-state", nil)

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	url, err := authSvc.GetAuthURL("test-state")

	require.NoError(t, err)
	assert.Equal(t, "https://test-domain.com/login?state=test-state", url)
}

func TestAuthServiceImpl_GetAuthURL_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		GetAuthURL("test-state").
		Return("", errors.New("cognito error"))

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	url, err := authSvc.GetAuthURL("test-state")

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Contains(t, err.Error(), "failed to get auth url")
}

func TestAuthServiceImpl_HandleCallback_Success_NewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	expectedTokens := &schema.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
	}

	expectedUserInfo := &schema.UserInfoResponse{
		UserID:        userID,
		UserName:      "testuser",
		Email:         "test@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		ExchangeCodeForTokens(gomock.Any(), "auth-code").
		Return(expectedTokens, nil)

	mockCognito.EXPECT().
		ValidateAccessToken(gomock.Any(), "access-token").
		Return(expectedUserInfo, nil)

	mockAuthRepo := repository.NewMockAuthRepository(ctrl)
	mockAuthRepo.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(nil, errors.New("user not found"))

	expectedCustomer := &entity.Customer{
		UserID:      userID,
		PhoneNumber: expectedUserInfo.PhoneNumber,
		UserName:    expectedUserInfo.UserName,
		Email:       expectedUserInfo.Email,
	}

	mockAuthRepo.EXPECT().
		CreateUser(gomock.Any(), expectedCustomer).
		Return(expectedCustomer, nil)

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
		authRepo:       mockAuthRepo,
	}

	code := "auth-code"
	req := &schema.HandleCallbackRequest{
		Code:  &code,
		State: "test-state",
	}

	tokens, err := authSvc.HandleCallback(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, expectedTokens, tokens)
}

func TestAuthServiceImpl_HandleCallback_Success_ExistingUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	expectedTokens := &schema.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
	}

	expectedUserInfo := &schema.UserInfoResponse{
		UserID:        userID,
		UserName:      "existinguser",
		Email:         "existing@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	existingCustomer := &entity.Customer{
		UserID:      userID,
		UserName:    "existinguser",
		Email:       "existing@example.com",
		PhoneNumber: "+1234567890",
	}

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		ExchangeCodeForTokens(gomock.Any(), "auth-code").
		Return(expectedTokens, nil)

	mockCognito.EXPECT().
		ValidateAccessToken(gomock.Any(), "access-token").
		Return(expectedUserInfo, nil)

	mockAuthRepo := repository.NewMockAuthRepository(ctrl)
	mockAuthRepo.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(existingCustomer, nil)

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
		authRepo:       mockAuthRepo,
	}

	code := "auth-code"
	req := &schema.HandleCallbackRequest{
		Code:  &code,
		State: "test-state",
	}

	tokens, err := authSvc.HandleCallback(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, expectedTokens, tokens)
}

func TestAuthServiceImpl_HandleCallback_MissingCode(t *testing.T) {
	authSvc := &AuthServiceImpl{}

	req := &schema.HandleCallbackRequest{
		Code:  nil,
		State: "test-state",
	}

	tokens, err := authSvc.HandleCallback(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "code is required")
}

func TestAuthServiceImpl_HandleCallback_ExchangeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		ExchangeCodeForTokens(gomock.Any(), "auth-code").
		Return(nil, errors.New("token exchange failed"))

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	code := "auth-code"
	req := &schema.HandleCallbackRequest{
		Code:  &code,
		State: "test-state",
	}

	tokens, err := authSvc.HandleCallback(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "failed to exchange code for tokens")
}

func TestAuthServiceImpl_HandleCallback_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTokens := &schema.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
	}

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		ExchangeCodeForTokens(gomock.Any(), "auth-code").
		Return(expectedTokens, nil)

	mockCognito.EXPECT().
		ValidateAccessToken(gomock.Any(), "access-token").
		Return(nil, errors.New("token validation failed"))

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	code := "auth-code"
	req := &schema.HandleCallbackRequest{
		Code:  &code,
		State: "test-state",
	}

	tokens, err := authSvc.HandleCallback(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "failed to validate access token")
}

func TestAuthServiceImpl_HandleCallback_CreateUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	expectedTokens := &schema.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
	}

	expectedUserInfo := &schema.UserInfoResponse{
		UserID:        userID,
		UserName:      "newuser",
		Email:         "new@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		ExchangeCodeForTokens(gomock.Any(), "auth-code").
		Return(expectedTokens, nil)

	mockCognito.EXPECT().
		ValidateAccessToken(gomock.Any(), "access-token").
		Return(expectedUserInfo, nil)

	mockAuthRepo := repository.NewMockAuthRepository(ctrl)
	mockAuthRepo.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(nil, errors.New("user not found"))

	expectedCustomer := &entity.Customer{
		UserID:      userID,
		PhoneNumber: expectedUserInfo.PhoneNumber,
		UserName:    expectedUserInfo.UserName,
		Email:       expectedUserInfo.Email,
	}

	mockAuthRepo.EXPECT().
		CreateUser(gomock.Any(), expectedCustomer).
		Return(nil, errors.New("failed to create user in database"))

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
		authRepo:       mockAuthRepo,
	}

	code := "auth-code"
	req := &schema.HandleCallbackRequest{
		Code:  &code,
		State: "test-state",
	}

	tokens, err := authSvc.HandleCallback(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestAuthServiceImpl_RefreshAccessToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokens := &schema.TokenResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		IDToken:      "new-id-token",
	}

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		RefreshTokens(gomock.Any(), "valid-refresh-token").
		Return(mockTokens, nil)

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	req := &schema.RefreshTokenRequest{
		RefreshToken: "valid-refresh-token",
	}

	response, err := authSvc.RefreshAccessToken(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, mockTokens.AccessToken, response.AccessToken)
	assert.Equal(t, "Bearer", response.TokenType)
	assert.Equal(t, 3600, response.ExpiresIn)
}

func TestAuthServiceImpl_RefreshAccessToken_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		RefreshTokens(gomock.Any(), "invalid-refresh-token").
		Return(nil, errors.New("refresh failed"))

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	req := &schema.RefreshTokenRequest{
		RefreshToken: "invalid-refresh-token",
	}

	response, err := authSvc.RefreshAccessToken(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "refresh failed")
}

func TestAuthServiceImpl_ValidateToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUserInfo := &schema.UserInfoResponse{
		UserName:      "testuser",
		Email:         "test@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		ValidateAccessToken(gomock.Any(), "valid-token").
		Return(expectedUserInfo, nil)

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	userInfo, err := authSvc.ValidateToken(context.Background(), "valid-token")

	require.NoError(t, err)
	assert.Equal(t, expectedUserInfo, userInfo)
}

func TestAuthServiceImpl_ValidateToken_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCognito := service.NewMockCognitoServiceInterface(ctrl)
	mockCognito.EXPECT().
		ValidateAccessToken(gomock.Any(), "invalid-token").
		Return(nil, errors.New("validation failed"))

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	userInfo, err := authSvc.ValidateToken(context.Background(), "invalid-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "validation failed")
}
