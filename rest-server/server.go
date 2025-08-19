package rest_server

import (
	"net/http"

	"github.com/IainMosima/gomart/domains/auth/service"
	"github.com/IainMosima/gomart/rest-server/handlers/auth"
	"github.com/IainMosima/gomart/rest-server/handlers/category"
	"github.com/IainMosima/gomart/rest-server/handlers/order"
	"github.com/IainMosima/gomart/rest-server/handlers/product"
	"github.com/IainMosima/gomart/rest-server/middleware"
	"github.com/IainMosima/gomart/rest-server/routes"
	"github.com/gin-gonic/gin"
)

type RestServer struct {
	router                *gin.Engine
	categoryHandler       category.CategoryHandlerInterface
	productHandler        product.ProductHandlerInterface
	authHandler           auth.AuthHandlerInterface
	orderHandlerInterface order.OrderHandlerInterface
	authMiddleware        *middleware.AuthMiddleware
}

func NewRestServer(categoryHandler category.CategoryHandlerInterface, productHandler product.ProductHandlerInterface, authHandler auth.AuthHandlerInterface, orderHandlerInterface order.OrderHandlerInterface, authService service.AuthService) *RestServer {
	router := gin.New()

	_ = router.SetTrustedProxies(nil)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := &RestServer{
		router:                router,
		categoryHandler:       categoryHandler,
		productHandler:        productHandler,
		authHandler:           authHandler,
		orderHandlerInterface: orderHandlerInterface,
		authMiddleware:        middleware.NewAuthMiddleware(authService),
	}

	server.setupRoutes()
	return server
}

func (s *RestServer) setupRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	routes.SetupCategoryRoutes(s.router, s.categoryHandler)
	routes.SetupProductRoutes(s.router, s.productHandler)
	routes.SetupAuthRoutes(s.router, s.authHandler)
	routes.SetupOrderRoutes(s.router, s.orderHandlerInterface, s.authMiddleware)
}

func (s *RestServer) Start(addr string) error {
	return s.router.Run(addr)
}
