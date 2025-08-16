package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/IainMosima/gomart/domains/auth/service"
	"github.com/IainMosima/gomart/rest-server/dtos"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandlerImpl_LoginHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	expectedURL := "https://test-domain.com/login?client_id=test&state=test-state"
	mockAuthService.EXPECT().
		GetAuthURL("test-state").
		Return(expectedURL, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/auth/login", handler.LoginHandler)

	req := httptest.NewRequest("GET", "/auth/login?state=test-state", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

	var response dtos.LoginResponseDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedURL, response.LoginUrl)
	assert.Equal(t, "test-state", response.State)
	assert.Contains(t, response.Message, "Use the login url")
}

func TestAuthHandlerImpl_LoginHandler_NoState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	expectedURL := "https://test-domain.com/login?client_id=test"
	mockAuthService.EXPECT().
		GetAuthURL("").
		Return(expectedURL, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/auth/login", handler.LoginHandler)

	req := httptest.NewRequest("GET", "/auth/login", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

	var response dtos.LoginResponseDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedURL, response.LoginUrl)
	assert.Empty(t, response.State)
}

func TestAuthHandlerImpl_LoginHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	mockAuthService.EXPECT().
		GetAuthURL("test-state").
		Return("", errors.New("service error"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/auth/login", handler.LoginHandler)

	req := httptest.NewRequest("GET", "/auth/login?state=test-state", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "service error")
}

func TestAuthHandlerImpl_HandleCognitoCallback_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	expectedTokens := &schema.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
	}

	code := "auth-code"
	expectedRequest := &schema.HandleCallbackRequest{
		Code:  &code,
		State: "test-state",
	}

	mockAuthService.EXPECT().
		HandleCallback(gomock.Any(), expectedRequest).
		Return(expectedTokens, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/cognito/callback", handler.HandleCognitoCallback)

	req := httptest.NewRequest("GET", "/cognito/callback?code=auth-code&state=test-state", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schema.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedTokens.AccessToken, response.AccessToken)
	assert.Equal(t, expectedTokens.RefreshToken, response.RefreshToken)
	assert.Equal(t, expectedTokens.IDToken, response.IDToken)
}

func TestAuthHandlerImpl_HandleCognitoCallback_MissingCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/cognito/callback", handler.HandleCognitoCallback)

	req := httptest.NewRequest("GET", "/cognito/callback?state=test-state", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "code is required")
}

func TestAuthHandlerImpl_HandleCognitoCallback_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	code := "auth-code"
	expectedRequest := &schema.HandleCallbackRequest{
		Code:  &code,
		State: "test-state",
	}

	mockAuthService.EXPECT().
		HandleCallback(gomock.Any(), expectedRequest).
		Return(nil, errors.New("token exchange failed"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/cognito/callback", handler.HandleCognitoCallback)

	req := httptest.NewRequest("GET", "/cognito/callback?code=auth-code&state=test-state", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "token exchange failed")
}

func TestAuthHandlerImpl_ValidateToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	expectedUserInfo := &schema.UserInfoResponse{
		UserName:      "testuser",
		Email:         "test@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "valid-token").
		Return(expectedUserInfo, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/validate", handler.ValidateToken)

	requestBody := dtos.ValidateTokenRequestDTO{
		AccessToken: "valid-token",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/auth/validate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schema.UserInfoResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedUserInfo.UserName, response.UserName)
	assert.Equal(t, expectedUserInfo.Email, response.Email)
	assert.Equal(t, expectedUserInfo.EmailVerified, response.EmailVerified)
	assert.Equal(t, expectedUserInfo.PhoneNumber, response.PhoneNumber)
}

func TestAuthHandlerImpl_ValidateToken_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/validate", handler.ValidateToken)

	req := httptest.NewRequest("POST", "/auth/validate", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandlerImpl_ValidateToken_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "invalid-token").
		Return(nil, errors.New("token validation failed"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/validate", handler.ValidateToken)

	requestBody := dtos.ValidateTokenRequestDTO{
		AccessToken: "invalid-token",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/auth/validate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "token validation failed")
}

func TestAuthHandlerImpl_RefreshAccessToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	expectedResponse := &schema.RefreshTokenResponse{
		AccessToken: "new-access-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	expectedRequest := &schema.RefreshTokenRequest{
		RefreshToken: "valid-refresh-token",
	}

	mockAuthService.EXPECT().
		RefreshAccessToken(gomock.Any(), expectedRequest).
		Return(expectedResponse, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/refresh", handler.RefreshAccessToken)

	requestBody := dtos.RefreshTokenRequestDTO{
		RefreshToken: "valid-refresh-token",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schema.RefreshTokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedResponse.AccessToken, response.AccessToken)
	assert.Equal(t, expectedResponse.TokenType, response.TokenType)
	assert.Equal(t, expectedResponse.ExpiresIn, response.ExpiresIn)
}

func TestAuthHandlerImpl_RefreshAccessToken_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/refresh", handler.RefreshAccessToken)

	req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandlerImpl_RefreshAccessToken_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	expectedRequest := &schema.RefreshTokenRequest{
		RefreshToken: "invalid-refresh-token",
	}

	mockAuthService.EXPECT().
		RefreshAccessToken(gomock.Any(), expectedRequest).
		Return(nil, errors.New("refresh failed"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/refresh", handler.RefreshAccessToken)

	requestBody := dtos.RefreshTokenRequestDTO{
		RefreshToken: "invalid-refresh-token",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "refresh failed")
}

func TestNewAuthHandlerImpl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	handler := NewAuthHandlerImpl(mockAuthService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockAuthService, handler.authService)
}
