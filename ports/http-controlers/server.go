package httpControllers

import (
	"database/sql"

	"books/core"
	"books/ports/http-controlers/controllers"
	"books/ports/http-controlers/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router      *gin.Engine
	controllers *controllers.Controllers
	core        *core.Core
}

// NewServer creates a new HTTP server
func NewServer(core *core.Core) *Server {
	return NewServerWithDB(core, nil)
}

// NewServerWithDB creates a new HTTP server with database health check support
func NewServerWithDB(core *core.Core, db *sql.DB) *Server {
	router := gin.Default()

	// Apply middleware in order
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.SkipAuthPaths("/health"))

	// Create controllers
	var ctlrs *controllers.Controllers
	if db != nil {
		ctlrs = controllers.NewControllersWithDB(core, db)
	} else {
		ctlrs = controllers.NewControllers(core)
	}

	return &Server{
		router:      router,
		controllers: ctlrs,
		core:        core,
	}
}

// SetupRoutes configures all routes for the server
func (s *Server) SetupRoutes() {
	// Register all routes from controllers
	s.controllers.RegisterRoutes(s.router)
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	s.SetupRoutes()

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