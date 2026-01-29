package controllers

import (
	"database/sql"
	"net/http"

	"books/core"

	"github.com/gin-gonic/gin"
)

// DBPinger interface for database ping
type DBPinger interface {
	Ping() error
}

// Controllers contains all HTTP controllers
type Controllers struct {
	BookController *BookController
	db             DBPinger
	// Add other controllers here as needed
}

// NewControllers creates and initializes all HTTP controllers
func NewControllers(core *core.Core) *Controllers {
	return &Controllers{
		BookController: NewBookController(core),
		db:             nil, // No DB for simple setup
		// Initialize other controllers here
	}
}

// NewControllersWithDB creates and initializes all HTTP controllers with DB health check
func NewControllersWithDB(core *core.Core, db *sql.DB) *Controllers {
	return &Controllers{
		BookController: NewBookController(core),
		db:             db,
		// Initialize other controllers here
	}
}

// RegisterRoutes registers all controller routes to the router
func (c *Controllers) RegisterRoutes(router *gin.Engine) {
	// Register book routes
	booksGroup := router.Group("/books")
	{
		// Create
		booksGroup.POST("", c.BookController.AddBook)

		// Read
		booksGroup.GET("", c.BookController.GetAllBooks)
		booksGroup.GET("/isbn/:isbn", c.BookController.GetBookByISBN)
		booksGroup.GET("/:isbn", c.BookController.GetBook)

		// Update
		booksGroup.PUT("/:isbn", c.BookController.UpdateBook)

		// Delete
		booksGroup.DELETE("/:isbn", c.BookController.DeleteBook)
	}

	// Register health check with optional DB ping
	router.GET("/health", c.healthCheck)

	// Add other controller route registrations here
}

// healthCheck returns the health status of the service
func (c *Controllers) healthCheck(ctx *gin.Context) {
	status := "ok"
	statusCode := http.StatusOK
	var dbStatus string

	if c.db != nil {
		if err := c.db.Ping(); err != nil {
			status = "degraded"
			dbStatus = "unhealthy"
			statusCode = http.StatusServiceUnavailable
		} else {
			dbStatus = "healthy"
		}
	} else {
		dbStatus = "not configured"
	}

	ctx.JSON(statusCode, gin.H{
		"status":   status,
		"database": dbStatus,
	})
}