package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IainMosima/gomart/domains/auth/schema"
	"github.com/IainMosima/gomart/domains/auth/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_RequireAuth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	authMiddleware := NewAuthMiddleware(mockAuthService)

	expectedUser := &schema.UserInfoResponse{
		UserName:      "testuser",
		Email:         "test@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "valid-token").
		Return(expectedUser, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(authMiddleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		user, exists := GetUserFromContext(c)
		assert.True(t, exists)
		assert.Equal(t, expectedUser, user)
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_RequireAuth_MissingAuthHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	authMiddleware := NewAuthMiddleware(mockAuthService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(authMiddleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestAuthMiddleware_RequireAuth_InvalidBearerFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	authMiddleware := NewAuthMiddleware(mockAuthService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(authMiddleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	testCases := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "invalid-token"},
		{"Wrong prefix", "Basic valid-token"},
		{"Empty token", "Bearer "},
		{"Only Bearer", "Bearer"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tc.header)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Contains(t, w.Body.String(), "Unauthorized")
		})
	}
}

func TestAuthMiddleware_RequireAuth_TokenValidationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	authMiddleware := NewAuthMiddleware(mockAuthService)

	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "invalid-token").
		Return(nil, errors.New("token validation failed"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(authMiddleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestGetUserFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedUser := &schema.UserInfoResponse{
		UserName:      "testuser",
		Email:         "test@example.com",
		EmailVerified: true,
		PhoneNumber:   "+1234567890",
	}

	c.Set("user", expectedUser)

	user, exists := GetUserFromContext(c)

	assert.True(t, exists)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserFromContext_UserNotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	user, exists := GetUserFromContext(c)

	assert.False(t, exists)
	assert.Nil(t, user)
}

func TestGetUserFromContext_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set("user", "not-a-user-struct")

	user, exists := GetUserFromContext(c)

	assert.False(t, exists)
	assert.Nil(t, user)
}

func TestNewAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service.NewMockAuthService(ctrl)
	authMiddleware := NewAuthMiddleware(mockAuthService)

	assert.NotNil(t, authMiddleware)
	assert.Equal(t, mockAuthService, authMiddleware.authService)
}
