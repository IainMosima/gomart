package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/IainMosima/gomart/domains/auth/entity"
	"github.com/IainMosima/gomart/domains/auth/repository"
	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/IainMosima/gomart/domains/auth/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type MockCognitoService struct {
	GetAuthURLFunc            func(state string) (string, error)
	ExchangeCodeForTokensFunc func(ctx context.Context, code string) (*schema.TokenResponse, error)
	RefreshTokensFunc         func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error)
	ValidateAccessTokenFunc   func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error)
}

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
		GetUserByID(gomock.Any(), expectedUserInfo.UserID).
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
		GetUserByID(gomock.Any(), expectedUserInfo.UserID).
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
		State: "",
	}

	tokens, err := authSvc.HandleCallback(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, expectedTokens, tokens)
}

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

func TestAuthServiceImpl_HandleCallback_UserValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTokens := &schema.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
	}

	mockCognito := &MockCognitoService{
		ExchangeCodeForTokensFunc: func(ctx context.Context, code string) (*schema.TokenResponse, error) {
			return expectedTokens, nil
		},
		ValidateAccessTokenFunc: func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
			return nil, errors.New("token validation failed")
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
	assert.Contains(t, err.Error(), "failed to validate access token")
}

func TestAuthServiceImpl_HandleCallback_ExistingUser(t *testing.T) {
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
		GetUserByID(gomock.Any(), userID).
		Return(existingCustomer, nil).
		Times(1)

	// Should NOT call CreateUser since user already exists

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
		GetUserByID(gomock.Any(), userID).
		Return(nil, errors.New("user not found")).
		Times(1)

	expectedCustomer := &entity.Customer{
		UserID:      userID,
		PhoneNumber: expectedUserInfo.PhoneNumber,
		UserName:    expectedUserInfo.UserName,
		Email:       expectedUserInfo.Email,
	}

	mockAuthRepo.EXPECT().
		CreateUser(gomock.Any(), expectedCustomer).
		Return(nil, errors.New("failed to create user in database")).
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

	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestAuthServiceImpl_RefreshAccessToken_EmptyToken(t *testing.T) {
	mockCognito := &MockCognitoService{
		RefreshTokensFunc: func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
			assert.Empty(t, refreshToken)
			return nil, errors.New("refresh token is required")
		},
	}

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	req := &schema.RefreshTokenRequest{
		RefreshToken: "",
	}

	response, err := authSvc.RefreshAccessToken(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestAuthServiceImpl_ValidateToken_EmptyToken(t *testing.T) {
	mockCognito := &MockCognitoService{
		ValidateAccessTokenFunc: func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
			assert.Empty(t, accessToken)
			return nil, errors.New("access token is required")
		},
	}

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	userInfo, err := authSvc.ValidateToken(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
}

func TestAuthServiceImpl_ContextCancellation(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*AuthServiceImpl, context.Context) error
	}{
		{
			name: "HandleCallback context cancellation",
			testFunc: func(service *AuthServiceImpl, ctx context.Context) error {
				code := "test-code"
				req := &schema.HandleCallbackRequest{Code: &code}
				_, err := service.HandleCallback(ctx, req)
				return err
			},
		},
		{
			name: "RefreshAccessToken context cancellation",
			testFunc: func(service *AuthServiceImpl, ctx context.Context) error {
				req := &schema.RefreshTokenRequest{RefreshToken: "test-token"}
				_, err := service.RefreshAccessToken(ctx, req)
				return err
			},
		},
		{
			name: "ValidateToken context cancellation",
			testFunc: func(service *AuthServiceImpl, ctx context.Context) error {
				_, err := service.ValidateToken(ctx, "test-token")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCognito := &MockCognitoService{
				ExchangeCodeForTokensFunc: func(ctx context.Context, code string) (*schema.TokenResponse, error) {
					<-ctx.Done()
					return nil, ctx.Err()
				},
				RefreshTokensFunc: func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
					<-ctx.Done()
					return nil, ctx.Err()
				},
				ValidateAccessTokenFunc: func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
					<-ctx.Done()
					return nil, ctx.Err()
				},
			}

			authSvc := &AuthServiceImpl{
				cognitoService: mockCognito,
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel() // Cancel immediately

			err := tt.testFunc(authSvc, ctx)
			assert.Error(t, err)
		})
	}
}

func TestAuthServiceImpl_RefreshAccessToken_ResponseFormat(t *testing.T) {
	mockTokens := &schema.TokenResponse{
		AccessToken:  "new-access-token-123",
		RefreshToken: "new-refresh-token-456",
		IDToken:      "new-id-token-789",
	}

	mockCognito := &MockCognitoService{
		RefreshTokensFunc: func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
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
	assert.NotNil(t, response)
	assert.Equal(t, "new-access-token-123", response.AccessToken)
	assert.Equal(t, "Bearer", response.TokenType)
	assert.Equal(t, 3600, response.ExpiresIn)
}

func TestAuthServiceImpl_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() *AuthServiceImpl
		action func(*AuthServiceImpl) error
	}{
		{
			name: "HandleCallback with nil request code",
			setup: func() *AuthServiceImpl {
				return &AuthServiceImpl{
					cognitoService: &MockCognitoService{},
				}
			},
			action: func(svc *AuthServiceImpl) error {
				req := &schema.HandleCallbackRequest{Code: nil}
				_, err := svc.HandleCallback(context.Background(), req)
				return err
			},
		},
		{
			name: "RefreshAccessToken with empty token",
			setup: func() *AuthServiceImpl {
				return &AuthServiceImpl{
					cognitoService: &MockCognitoService{
						RefreshTokensFunc: func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
							if refreshToken == "" {
								return nil, errors.New("empty refresh token")
							}
							return nil, nil
						},
					},
				}
			},
			action: func(svc *AuthServiceImpl) error {
				req := &schema.RefreshTokenRequest{RefreshToken: ""}
				_, err := svc.RefreshAccessToken(context.Background(), req)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.setup()
			err := tt.action(svc)
			assert.Error(t, err)
		})
	}
}

func TestAuthServiceImpl_CompleteErrorCoverage(t *testing.T) {
	// Test various error scenarios not covered by existing tests
	mockCognito := &MockCognitoService{
		GetAuthURLFunc: func(state string) (string, error) {
			return "", errors.New("cognito service error")
		},
		ExchangeCodeForTokensFunc: func(ctx context.Context, code string) (*schema.TokenResponse, error) {
			return nil, errors.New("exchange error")
		},
		RefreshTokensFunc: func(ctx context.Context, refreshToken string) (*schema.TokenResponse, error) {
			return nil, errors.New("refresh error")
		},
		ValidateAccessTokenFunc: func(ctx context.Context, accessToken string) (*schema.UserInfoResponse, error) {
			return nil, errors.New("validation error")
		},
	}

	authSvc := &AuthServiceImpl{
		cognitoService: mockCognito,
	}

	// Test all error paths
	_, err := authSvc.GetAuthURL("test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get auth url")

	_, err = authSvc.ValidateToken(context.Background(), "token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation error")

	_, err = authSvc.RefreshAccessToken(context.Background(), &schema.RefreshTokenRequest{RefreshToken: "token"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refresh error")

	code := "test-code"
	_, err = authSvc.HandleCallback(context.Background(), &schema.HandleCallbackRequest{Code: &code})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to exchange code for tokens")
}
