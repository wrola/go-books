package controllers

import (
	"errors"
	"net/http"

	"books/core"
	"books/core/storage/repositories/interfaces"

	"github.com/gin-gonic/gin"
)

// BookController handles HTTP requests related to books
type BookController struct {
	core *core.Core
}

// NewBookController creates a new book controller
func NewBookController(core *core.Core) *BookController {
	return &BookController{core: core}
}

// AddBookRequest represents the request body for adding a book
type AddBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	ISBN   string `json:"isbn" binding:"required"`
}

type UpdateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// AddBook handles the request to add a new book
func (c *BookController) AddBook(ctx *gin.Context) {
	var request AddBookRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := c.core.AddBook(ctx, request.Title, request.Author, request.ISBN)
	if err != nil {
		status := mapErrorToStatus(err)
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"book": gin.H{
			"isbn":   book.ISBN,
			"title":  book.Title,
			"author": book.Author,
		},
	})
}

// GetBookByISBN handles the request to get a book by ISBN
func (c *BookController) GetBookByISBN(ctx *gin.Context) {
	isbn := ctx.Param("isbn")

	if isbn == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ISBN parameter is required"})
		return
	}

	book, err := c.core.GetBookByISBN(ctx, isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"book": gin.H{
			"title":  book.Title,
			"author": book.Author,
			"isbn":   book.ISBN,
		},
	})
}

// GetBook handles the request to get a book by ID
func (c *BookController) GetBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")

	if isbn == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ISBN parameter is required"})
		return
	}

	book, err := c.core.GetBookByISBN(ctx, isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"book": gin.H{
			"title":  book.Title,
			"author": book.Author,
			"isbn":   book.ISBN,
		},
	})
}

// GetAllBooks handles the request to get all books
func (c *BookController) GetAllBooks(ctx *gin.Context) {
	books, err := c.core.GetAllBooks(ctx)
	if err != nil {
		status := mapErrorToStatus(err)
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, book := range books {
		result = append(result, gin.H{
			"title":  book.Title,
			"author": book.Author,
			"isbn":   book.ISBN,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"books": result,
	})
}

// UpdateBook handles the request to update a book
func (c *BookController) UpdateBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")

	if isbn == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ISBN parameter is required"})
		return
	}

	var request UpdateBookRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := c.core.UpdateBook(ctx, isbn, request.Title, request.Author)

	if err != nil {
		status := mapErrorToStatus(err)
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Book updated successfully",
		"book": gin.H{
			"title":  book.Title,
			"author": book.Author,
			"isbn":   book.ISBN,
		},
	})
}

// DeleteBook handles the request to delete a book
func (c *BookController) DeleteBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")

	if isbn == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ISBN parameter is required"})
		return
	}

	err := c.core.DeleteBook(ctx, isbn)

	if err != nil {
		status := mapErrorToStatus(err)
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// mapErrorToStatus maps domain errors to HTTP status codes
func mapErrorToStatus(err error) int {
	if errors.Is(err, interfaces.ErrBookNotFound) {
		return http.StatusNotFound
	}
	// Check for validation errors (contains common validation error messages)
	errMsg := err.Error()
	if contains(errMsg, "cannot be empty", "invalid", "required", "already exists", "ISBN must be", "checksum") {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

// contains checks if the string contains any of the substrings
func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

