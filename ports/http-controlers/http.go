package httpControllers

import (
	"books/core"
	"books/ports/http-controlers/controllers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPServer represents the HTTP server
type HTTPServer struct {
	router      *gin.Engine
	controllers *controllers.Controllers
	core        *core.Core
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(core *core.Core) *HTTPServer {
	router := gin.Default()

	// Create controllers
	ctlrs := controllers.NewControllers(core)

	return &HTTPServer{
		router:      router,
		controllers: ctlrs,
		core:        core,
	}
}

// SetupRoutes configures all routes for the server
func (s *HTTPServer) SetupRoutes() {
	// Register all routes from controllers
	s.controllers.RegisterRoutes(s.router)
}

// Start starts the HTTP server
func (s *HTTPServer) Start(addr string) error {
	s.SetupRoutes()

	log.Printf("Starting HTTP server on %s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// Start server in a goroutine so we can handle graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
		return err
	}

	log.Println("Server exiting")
	return nil
}