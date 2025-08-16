package routes

import (
	"testing"

	"github.com/IainMosima/gomart/domains/category/service"
	"github.com/IainMosima/gomart/rest-server/handlers/category"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSetupCategoryRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCategoryService(ctrl)
	handler := category.NewCategoryHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	SetupCategoryRoutes(router, handler)

	routes := router.Routes()

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/categories"},
		{"GET", "/categories"},
		{"GET", "/categories/:id"},
		{"PUT", "/categories/:id"},
		{"DELETE", "/categories/:id"},
		{"GET", "/categories/:id/children"},
		{"GET", "/categories/:id/average-price"},
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

func TestSetupCategoryRoutes_HandlerIntegration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCategoryService(ctrl)
	handler := category.NewCategoryHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	assert.NotPanics(t, func() {
		SetupCategoryRoutes(router, handler)
	}, "SetupCategoryRoutes should not panic with valid inputs")

	routes := router.Routes()
	assert.NotEmpty(t, routes, "Router should have routes after setup")
	assert.Len(t, routes, 7, "Should have exactly 7 category routes")
}

func TestSetupCategoryRoutes_NilHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	assert.Panics(t, func() {
		SetupCategoryRoutes(router, nil)
	}, "SetupCategoryRoutes should panic with nil handler")
}
