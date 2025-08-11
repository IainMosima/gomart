package main

import (
	"context"
	"log"

	"github.com/IainMosima/gomart/configs"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	categoryRepo "github.com/IainMosima/gomart/infrastructures/repository"
	"github.com/IainMosima/gomart/rest-server"
	"github.com/IainMosima/gomart/rest-server/handlers"
	categoryService "github.com/IainMosima/gomart/services/category"
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
	categoryRepository := categoryRepo.NewCategoryRepository(store)

	// Initialize services
	categoryServiceImpl := categoryService.NewCategoryService(categoryRepository)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryServiceImpl)

	// Initialize REST server
	server := rest_server.NewRestServer(categoryHandler)

	// Start HTTP server
	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal("cannot start HTTP server:", err)
	}
}
