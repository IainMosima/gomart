package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/IainMosima/gomart/domains/auth/entity"
	"github.com/IainMosima/gomart/domains/auth/repository"
	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/IainMosima/gomart/domains/auth/service"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockCognitoService is a mock implementation of CognitoServiceInterface for testing
type MockCognitoService struct {
	GetAuthURLFunc            func(state string) (string, error)
	ExchangeCodeForTokensFunc func(ctx context.Context, code string) (*schema.TokenResponse, error)
	RefreshTokensFunc         func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error)
	ValidateAccessTokenFunc   func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error)
}

// Ensure MockCognitoService implements CognitoServiceInterface
var _ service.CognitoServiceInterface = (*MockCognitoService)(nil)

func (m *MockCognitoService) GetAuthURL(state string) (string, error) {
	if m.GetAuthURLFunc != nil {
		return m.GetAuthURLFunc(state)
	}
	return "", nil
}

func (m *MockCognitoService) ExchangeCodeForTokens(ctx context.Context, code string) (*schema.TokenResponse, error) {
	if m.ExchangeCodeForTokensFunc != nil {
		return m.ExchangeCodeForTokensFunc(ctx, code)
	}
	return nil, nil
}

func (m *MockCognitoService) RefreshTokens(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
	if m.RefreshTokensFunc != nil {
		return m.RefreshTokensFunc(ctx, refreshToken)
	}
	return nil, nil
}

func (m *MockCognitoService) ValidateAccessToken(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
	if m.ValidateAccessTokenFunc != nil {
		return m.ValidateAccessTokenFunc(ctx, accessToken)
	}
	return nil, nil
}

func (m *MockCognitoService) ParseIDToken(ctx context.Context, idToken string) (*schema.CognitoUserInfoJWTClaims, error) {
	return nil, nil
}

func TestAuthServiceImpl_GetAuthURL_Success(t *testing.T) {
	mockCognito := &MockCognitoService{
		GetAuthURLFunc: func(state string) (string, error) {
			assert.Equal(t, "test-state", state)
			return "https://test-domain.com/login?state=test-state", nil
		},
	}

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	url, err := authSvc.GetAuthURL("test-state")

	require.NoError(t, err)
	assert.Equal(t, "https://test-domain.com/login?state=test-state", url)
}

func TestAuthServiceImpl_GetAuthURL_CognitoError(t *testing.T) {
	mockCognito := &MockCognitoService{
		GetAuthURLFunc: func(state string) (string, error) {
			return "", errors.New("cognito error")
		},
	}

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	url, err := authSvc.GetAuthURL("test-state")

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Contains(t, err.Error(), "failed to get auth url")
	assert.Contains(t, err.Error(), "cognito error")
}

func TestAuthServiceImpl_HandleCallback_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTokens := &schema.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
	}

	expectedUserInfo := &schema.UserInfoResponse{
		UserID:        uuid.New(),
		UserName:      "testuser",
		Email:         "test@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	mockCognito := &MockCognitoService{
		ExchangeCodeForTokensFunc: func(ctx context.Context, code string) (*schema.TokenResponse, error) {
			assert.Equal(t, "auth-code", code)
			return expectedTokens, nil
		},
		ValidateAccessTokenFunc: func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
			assert.Equal(t, "access-token", accessToken)
			return expectedUserInfo, nil
		},
	}

	mockAuthRepo := repository.NewMockAuthRepository(ctrl)
	mockAuthRepo.EXPECT().
		GetUserByID(gomock.Any(), expectedUserInfo.Email).
		Return(nil, errors.New("user not found")).
		Times(1)

	expectedCustomer := &entity.Customer{
		UserID:      expectedUserInfo.UserID,
		PhoneNumber: expectedUserInfo.PhoneNumber,
		UserName:    expectedUserInfo.UserName,
		Email:       expectedUserInfo.Email,
	}

	mockAuthRepo.EXPECT().
		CreateUser(gomock.Any(), expectedCustomer).
		Return(expectedCustomer, nil).
		Times(1)

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
	authSvc := &AuthServiceImpl{
		cognitoService: &MockCognitoService{},
	}

	req := &schema.HandleCallbackRequest{
		Code:  nil,
		State: "test-state",
	}

	tokens, err := authSvc.HandleCallback(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "code is required")
}

func TestAuthServiceImpl_HandleCallback_CognitoError(t *testing.T) {
	mockCognito := &MockCognitoService{
		ExchangeCodeForTokensFunc: func(ctx context.Context, code string) (*schema.TokenResponse, error) {
			return nil, errors.New("token exchange failed")
		},
	}

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
	assert.Contains(t, err.Error(), "token exchange failed")
}

func TestAuthServiceImpl_RefreshAccessToken_Success(t *testing.T) {
	mockTokens := &schema.TokenResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		IDToken:      "new-id-token",
	}

	mockCognito := &MockCognitoService{
		RefreshTokensFunc: func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
			assert.Equal(t, "valid-refresh-token", refreshToken)
			return mockTokens, nil
		},
	}

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
	assert.Equal(t, int(3600), response.ExpiresIn)
}

func TestAuthServiceImpl_RefreshAccessToken_CognitoError(t *testing.T) {
	mockCognito := &MockCognitoService{
		RefreshTokensFunc: func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
			return nil, errors.New("refresh failed")
		},
	}

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
	expectedUserInfo := &schema.UserInfoResponse{
		UserName:      "testuser",
		Email:         "test@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	mockCognito := &MockCognitoService{
		ValidateAccessTokenFunc: func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
			assert.Equal(t, "valid-token", accessToken)
			return expectedUserInfo, nil
		},
	}

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	userInfo, err := authSvc.ValidateToken(context.Background(), "valid-token")

	require.NoError(t, err)
	assert.Equal(t, expectedUserInfo, userInfo)
}

func TestAuthServiceImpl_ValidateToken_CognitoError(t *testing.T) {
	mockCognito := &MockCognitoService{
		ValidateAccessTokenFunc: func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
			return nil, errors.New("validation failed")
		},
	}

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	userInfo, err := authSvc.ValidateToken(context.Background(), "invalid-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestAuthServiceImpl_GetAuthURL_EmptyState(t *testing.T) {
	mockCognito := &MockCognitoService{
		GetAuthURLFunc: func(state string) (string, error) {
			assert.Empty(t, state)
			return "https://test-domain.com/login", nil
		},
	}

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	url, err := authSvc.GetAuthURL("")

	require.NoError(t, err)
	assert.Equal(t, "https://test-domain.com/login", url)
}

func TestAuthServiceImpl_HandleCallback_EmptyState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTokens := &schema.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
	}

	expectedUserInfo := &schema.UserInfoResponse{
		UserID:        uuid.New(),
		UserName:      "testuser",
		Email:         "test@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	mockCognito := &MockCognitoService{
		ExchangeCodeForTokensFunc: func(ctx context.Context, code string) (*schema.TokenResponse, error) {
			return expectedTokens, nil
		},
		ValidateAccessTokenFunc: func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
			return expectedUserInfo, nil
		},
	}

	mockAuthRepo := repository.NewMockAuthRepository(ctrl)
	mockAuthRepo.EXPECT().
		GetUserByID(gomock.Any(), expectedUserInfo.Email).
		Return(nil, errors.New("user not found")).
		Times(1)

	expectedCustomer := &entity.Customer{
		UserID:      expectedUserInfo.UserID,
		PhoneNumber: expectedUserInfo.PhoneNumber,
		UserName:    expectedUserInfo.UserName,
		Email:       expectedUserInfo.Email,
	}

	mockAuthRepo.EXPECT().
		CreateUser(gomock.Any(), expectedCustomer).
		Return(expectedCustomer, nil).
		Times(1)

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
		authRepo:       mockAuthRepo,
	}

	code := "auth-code"
	req := &schema.HandleCallbackRequest{
		Code:  &code,
		State: "", // Empty state should work
	}

	tokens, err := authSvc.HandleCallback(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, expectedTokens, tokens)
}

// Note: We focus on testing auth service business logic with mocked dependencies
// rather than testing actual AWS integration in unit tests

func TestAuthServiceImpl_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockCognitoService)
		testFunc    func(*AuthServiceImpl) error
		expectError string
	}{
		{
			name: "GetAuthURL with service error",
			setupMock: func(m *MockCognitoService) {
				m.GetAuthURLFunc = func(state string) (string, error) {
					return "", errors.New("cognito service unavailable")
				}
			},
			testFunc: func(service *AuthServiceImpl) error {
				_, err := service.GetAuthURL("test")
				return err
			},
			expectError: "failed to get auth url",
		},
		{
			name: "HandleCallback with exchange error",
			setupMock: func(m *MockCognitoService) {
				m.ExchangeCodeForTokensFunc = func(ctx context.Context, code string) (*schema.TokenResponse, error) {
					return nil, errors.New("token exchange failed")
				}
			},
			testFunc: func(service *AuthServiceImpl) error {
				code := "test-code"
				req := &schema.HandleCallbackRequest{Code: &code}
				_, err := service.HandleCallback(context.Background(), req)
				return err
			},
			expectError: "failed to exchange code for tokens",
		},
		{
			name: "RefreshAccessToken with refresh error",
			setupMock: func(m *MockCognitoService) {
				m.RefreshTokensFunc = func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
					return nil, errors.New("refresh failed")
				}
			},
			testFunc: func(service *AuthServiceImpl) error {
				req := &schema.RefreshTokenRequest{RefreshToken: "test-token"}
				_, err := service.RefreshAccessToken(context.Background(), req)
				return err
			},
			expectError: "refresh failed",
		},
		{
			name: "ValidateToken with validation error",
			setupMock: func(m *MockCognitoService) {
				m.ValidateAccessTokenFunc = func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
					return nil, errors.New("validation failed")
				}
			},
			testFunc: func(service *AuthServiceImpl) error {
				_, err := service.ValidateToken(context.Background(), "test-token")
				return err
			},
			expectError: "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCognito := &MockCognitoService{}
			tt.setupMock(mockCognito)

			authSvc := &AuthServiceImpl{
				cognitoService: mockCognito,
			}

			err := tt.testFunc(authSvc)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}
}
