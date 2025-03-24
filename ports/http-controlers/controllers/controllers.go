package controllers

import (
	"books/core"

	"github.com/gin-gonic/gin"
)

// Controllers contains all HTTP controllers
type Controllers struct {
	BookController *BookController
	// Add other controllers here as needed
}

// NewControllers creates and initializes all HTTP controllers
func NewControllers(core *core.Core) *Controllers {
	return &Controllers{
		BookController: NewBookController(core),
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
		booksGroup.GET("/:id", c.BookController.GetBook)

		// Update
		booksGroup.PUT("/:id", c.BookController.UpdateBook)

		// Delete
		booksGroup.DELETE("/:id", c.BookController.DeleteBook)
	}

	// Register health check
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	// Add other controller route registrations here
}