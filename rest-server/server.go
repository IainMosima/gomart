package rest_server

import (
	"github.com/IainMosima/gomart/rest-server/handlers"
	"github.com/IainMosima/gomart/rest-server/routes"
	"github.com/gin-gonic/gin"
)

type RestServer struct {
	router          *gin.Engine
	categoryHandler handlers.CategoryHandlerInterface
	productHandler  handlers.ProductHandlerInterface
}

func NewRestServer(categoryHandler handlers.CategoryHandlerInterface, productHandler handlers.ProductHandlerInterface) *RestServer {
	router := gin.New()

	_ = router.SetTrustedProxies(nil)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := &RestServer{
		router:          router,
		categoryHandler: categoryHandler,
		productHandler:  productHandler,
	}

	server.setupRoutes()
	return server
}

func (s *RestServer) setupRoutes() {
	routes.SetupCategoryRoutes(s.router, s.categoryHandler)
	routes.SetupProductRoutes(s.router, s.productHandler)
}

func (s *RestServer) Start(addr string) error {
	return s.router.Run(addr)
}
