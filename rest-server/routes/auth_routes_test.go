package routes

import (
	"testing"

	"github.com/IainMosima/gomart/domains/auth/service"
	"github.com/IainMosima/gomart/rest-server/handlers/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetupAuthRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockAuthService(ctrl)
	handler := auth.NewAuthHandlerImpl(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	SetupAuthRoutes(router, handler)

	routes := router.Routes()

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/cognito/callback"},
		{"POST", "/auth/validate"},
		{"POST", "/auth/refresh"},
		{"GET", "/auth/login"},
	}

	assert.Equal(t, len(expectedRoutes), len(routes), "Number of registered routes should match expected")
	assert.NotEmpty(t, routes, "Routes should not be empty after setup")
	assert.NotNil(t, router, "Router should not be nil after route setup")

	// Verify specific routes exist
	routePaths := make(map[string]bool)
	for _, route := range routes {
		routePaths[route.Method+"::"+route.Path] = true
	}

	for _, expected := range expectedRoutes {
		key := expected.method + "::" + expected.path
		assert.True(t, routePaths[key], "Route %s %s should be registered", expected.method, expected.path)
	}
}

func TestSetupAuthRoutes_HandlerIntegration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockAuthService(ctrl)
	handler := auth.NewAuthHandlerImpl(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	assert.NotPanics(t, func() {
		SetupAuthRoutes(router, handler)
	}, "SetupAuthRoutes should not panic with valid inputs")

	routes := router.Routes()
	assert.NotEmpty(t, routes, "Router should have routes after setup")
	assert.Len(t, routes, 4, "Should have exactly 4 auth routes")
}

func TestSetupAuthRoutes_NilHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	assert.Panics(t, func() {
		SetupAuthRoutes(router, nil)
	}, "SetupAuthRoutes should panic with nil handler")
}
