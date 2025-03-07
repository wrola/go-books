package main

import (
	"books/application"
	"books/ports"
	"context"
)

func main() {
	ctx := context.Background()

	app := application.NewApplication(ctx)

	httpServer := ports.NewHttpServer(app)
	httpServer.Start()
}
