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
	dbConfig := infrastructure.NewConfig(
		os.Getenv("DB_HOST"),
		5432,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := infrastructure.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := infrastructure.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	bookRepo := repositories.NewBookStoragePostgresRepository(db)

	appCore := core.NewCore(bookRepo)

	httpModule := httpControllers.NewModuleWithDB(appCore, db)
	if err := httpModule.Start(":8080"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

}