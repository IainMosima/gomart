package routes

import (
	"testing"

	"github.com/IainMosima/gomart/domains/order/service"
	"github.com/IainMosima/gomart/rest-server/handlers/order"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetupOrderRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockOrderService(ctrl)
	handler := order.NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	SetupOrderRoutes(router, handler)

	routes := router.Routes()

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/orders"},
		{"GET", "/orders/:id/status"},
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

func TestSetupOrderRoutes_HandlerIntegration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockOrderService(ctrl)
	handler := order.NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	assert.NotPanics(t, func() {
		SetupOrderRoutes(router, handler)
	}, "SetupOrderRoutes should not panic with valid inputs")

	routes := router.Routes()
	assert.NotEmpty(t, routes, "Router should have routes after setup")
	assert.Len(t, routes, 2, "Should have exactly 2 order routes")
}

func TestSetupOrderRoutes_NilHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This should panic when trying to access nil handler methods
	assert.Panics(t, func() {
		SetupOrderRoutes(router, nil)
	}, "SetupOrderRoutes should panic with nil handler")
}

func TestSetupOrderRoutes_RouteGrouping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockOrderService(ctrl)
	handler := order.NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	SetupOrderRoutes(router, handler)

	routes := router.Routes()

	// Verify all routes are under the /orders group
	for _, route := range routes {
		assert.Contains(t, route.Path, "/orders", "All order routes should be under /orders group")
	}
}

func TestSetupOrderRoutes_MethodValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockOrderService(ctrl)
	handler := order.NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	SetupOrderRoutes(router, handler)

	routes := router.Routes()

	methodCount := make(map[string]int)
	for _, route := range routes {
		methodCount[route.Method]++
	}

	// Verify we have the expected HTTP methods
	assert.Equal(t, 1, methodCount["POST"], "Should have 1 POST route")
	assert.Equal(t, 1, methodCount["GET"], "Should have 1 GET route")
	assert.Equal(t, 0, methodCount["PUT"], "Should have 0 PUT routes")
	assert.Equal(t, 0, methodCount["DELETE"], "Should have 0 DELETE routes")
}
