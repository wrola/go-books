package main

import (
	"books/core"
	"books/core/storage/repositories"
	httpControllers "books/ports/http-controlers"
	"log"
)

func main() {
	// Create repository
	bookRepo := repositories.NewInMemoryBookRepository()

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