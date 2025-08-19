package main

import (
	"context"
	"log"

	"github.com/IainMosima/gomart/configs"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	repos "github.com/IainMosima/gomart/infrastructures/repository"
	"github.com/IainMosima/gomart/rest-server"
	"github.com/IainMosima/gomart/rest-server/handlers/auth"
	"github.com/IainMosima/gomart/rest-server/handlers/category"
	"github.com/IainMosima/gomart/rest-server/handlers/order"
	"github.com/IainMosima/gomart/rest-server/handlers/product"
	authService "github.com/IainMosima/gomart/services/auth"
	categoryService "github.com/IainMosima/gomart/services/category"
	orderService "github.com/IainMosima/gomart/services/order"
	"github.com/IainMosima/gomart/services/order/notification"
	productService "github.com/IainMosima/gomart/services/product"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer connPool.Close()

	store := db.NewStore(connPool)

	// Initialize repositories
	categoryRepository := repos.NewCategoryRepository(store)
	productRepository := repos.NewProductRepository(store)
	authRepository := repos.NewAuthRepository(store)
	orderRepository := repos.NewOrderRepository(store)

	// Initialize services
	categoryServiceImpl := categoryService.NewCategoryService(categoryRepository)
	productServiceImpl := productService.NewProductService(productRepository, categoryRepository)
	authServiceImpl, err := authService.NewAuthServiceImpl(&config, authRepository)
	notificationServiceImpl := notification.NewNotificationServiceImpl(&config, authRepository)
	orderServiceImpl := orderService.NewOrderServiceImpl(orderRepository, productRepository, authRepository, notificationServiceImpl)
	if err != nil {
		log.Fatal("cannot create auth service:", err)
	}

	// Initialize handlers
	categoryHandler := category.NewCategoryHandler(categoryServiceImpl)
	productHandler := product.NewProductHandler(productServiceImpl)
	authHandler := auth.NewAuthHandlerImpl(authServiceImpl)
	orderHandler := order.NewOrderHandler(orderServiceImpl)

	// Initialize REST server
	server := rest_server.NewRestServer(categoryHandler, productHandler, authHandler, orderHandler, authServiceImpl)

	// Start HTTP server
	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal("cannot start HTTP server:", err)
	}
}
