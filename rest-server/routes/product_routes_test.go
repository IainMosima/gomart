package routes

import (
	"testing"

	"github.com/IainMosima/gomart/domains/product/service"
	"github.com/IainMosima/gomart/rest-server/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSetupProductRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockProductService(ctrl)
	handler := handlers.NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	SetupProductRoutes(router, handler)

	routes := router.Routes()

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/products"},
		{"GET", "/products"},
		{"GET", "/products/:id"},
		{"PUT", "/products/:id"},
		{"DELETE", "/products/:id"},
		{"GET", "/products/category/:categoryId"},
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

func TestSetupProductRoutes_HandlerIntegration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockProductService(ctrl)
	handler := handlers.NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	assert.NotPanics(t, func() {
		SetupProductRoutes(router, handler)
	}, "SetupProductRoutes should not panic with valid inputs")

	routes := router.Routes()
	assert.NotEmpty(t, routes, "Router should have routes after setup")
	assert.Len(t, routes, 6, "Should have exactly 6 product routes")
}

func TestSetupProductRoutes_NilHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This should panic when trying to access nil handler methods
	assert.Panics(t, func() {
		SetupProductRoutes(router, nil)
	}, "SetupProductRoutes should panic with nil handler")
}
