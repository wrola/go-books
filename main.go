package main

import (
	"books/core"
	"books/core/storage/repositories"
	"books/infrastructure"
	httpControllers "books/ports/http-controlers"
	"log"
	"os"
)

func main() {
	// Database configuration
	dbConfig := infrastructure.NewConfig(
		os.Getenv("DB_HOST"),
		5432,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Connect to database
	db, err := infrastructure.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := infrastructure.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repository
	bookRepo := repositories.NewBookStoragePostgresRepository(db)

	// Create application core
	appCore := core.NewCore(bookRepo)

	// Start HTTP server
	httpModule := httpControllers.NewModule(appCore)
	if err := httpModule.Start(":8080"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	// In the future, you could start gRPC server here as well
	// grpcServer := grpcModule.NewModule(appCore)
	// if err := grpcServer.Start(":9090"); err != nil {
	//     log.Fatalf("Failed to start gRPC server: %v", err)
	// }
}