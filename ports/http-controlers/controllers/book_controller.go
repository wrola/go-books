package controllers

import (
	"net/http"

	"books/core"

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
	ISBN   string `json:"isbn"`
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"book": gin.H{
			"id":     book.ID,
			"title":  book.Title,
			"author": book.Author,
			"isbn":   book.ISBN,
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
			"id":     book.ID,
			"title":  book.Title,
			"author": book.Author,
			"isbn":   book.ISBN,
		},
	})
}

// GetBook handles the request to get a book by ID
func (c *BookController) GetBook(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	book, err := c.core.GetBookByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"book": gin.H{
			"id":     book.ID,
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, book := range books {
		result = append(result, gin.H{
			"id":     book.ID,
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
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	var request UpdateBookRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := c.core.UpdateBook(ctx, id, request.Title, request.Author, request.ISBN)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Book updated successfully",
		"book": gin.H{
			"id":     book.ID,
			"title":  book.Title,
			"author": book.Author,
			"isbn":   book.ISBN,
		},
	})
}

// DeleteBook handles the request to delete a book
func (c *BookController) DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	err := c.core.DeleteBook(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

