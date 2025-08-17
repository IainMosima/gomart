package rest_server

import (
	"github.com/IainMosima/gomart/rest-server/handlers/auth"
	"github.com/IainMosima/gomart/rest-server/handlers/category"
	"github.com/IainMosima/gomart/rest-server/handlers/order"
	"github.com/IainMosima/gomart/rest-server/handlers/product"
	"github.com/IainMosima/gomart/rest-server/routes"
	"github.com/gin-gonic/gin"
)

type RestServer struct {
	router                *gin.Engine
	categoryHandler       category.CategoryHandlerInterface
	productHandler        product.ProductHandlerInterface
	authHandler           auth.AuthHandlerInterface
	orderHandlerInterface order.OrderHandlerInterface
}

func NewRestServer(categoryHandler category.CategoryHandlerInterface, productHandler product.ProductHandlerInterface, authHandler auth.AuthHandlerInterface, orderHandlerInterface order.OrderHandlerInterface) *RestServer {
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
	}

	server.setupRoutes()
	return server
}

func (s *RestServer) setupRoutes() {
	routes.SetupCategoryRoutes(s.router, s.categoryHandler)
	routes.SetupProductRoutes(s.router, s.productHandler)
	routes.SetupAuthRoutes(s.router, s.authHandler)
	routes.SetupOrderRoutes(s.router, s.orderHandlerInterface)
}

func (s *RestServer) Start(addr string) error {
	return s.router.Run(addr)
}
