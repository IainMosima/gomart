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
	"github.com/IainMosima/gomart/rest-server/handlers/product"
	authService "github.com/IainMosima/gomart/services/auth"
	categoryService "github.com/IainMosima/gomart/services/category"
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

	// Initialize services
	categoryServiceImpl := categoryService.NewCategoryService(categoryRepository)
	productServiceImpl := productService.NewProductService(productRepository, categoryRepository)
	authServiceImpl, err := authService.NewAuthServiceImpl(&config, authRepository)
	if err != nil {
		log.Fatal("cannot create auth service:", err)
	}

	// Initialize handlers
	categoryHandler := category.NewCategoryHandler(categoryServiceImpl)
	productHandler := product.NewProductHandler(productServiceImpl)
	authHandler := auth.NewAuthHandlerImpl(authServiceImpl)

	// Initialize REST server
	server := rest_server.NewRestServer(categoryHandler, productHandler, authHandler)

	// Start HTTP server
	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal("cannot start HTTP server:", err)
	}
}
