package controllers

import (
	"errors"
	"log"
	"net/http"

	"books/core"
	"books/core/storage/repositories/interfaces"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	core *core.Core
}

func NewBookController(core *core.Core) *BookController {
	return &BookController{core: core}
}

type AddBookRequest struct {
	Title  string `json:"title" binding:"required,max=255"`
	Author string `json:"author" binding:"required,max=255"`
	ISBN   string `json:"isbn" binding:"required,max=20"`
}

type UpdateBookRequest struct {
	Title  string `json:"title" binding:"max=255"`
	Author string `json:"author" binding:"max=255"`
}

func (c *BookController) AddBook(ctx *gin.Context) {
	var request AddBookRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	book, err := c.core.AddBook(ctx, request.Title, request.Author, request.ISBN)
	if err != nil {
		status := mapErrorToStatus(err)
		ctx.JSON(status, gin.H{"error": sanitizeError(err, status)})
		log.Printf("AddBook error: %v", err)
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

func (c *BookController) GetAllBooks(ctx *gin.Context) {
	books, err := c.core.GetAllBooks(ctx)
	if err != nil {
		status := mapErrorToStatus(err)
		ctx.JSON(status, gin.H{"error": sanitizeError(err, status)})
		log.Printf("GetAllBooks error: %v", err)
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

func (c *BookController) UpdateBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")

	if isbn == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ISBN parameter is required"})
		return
	}

	var request UpdateBookRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	book, err := c.core.UpdateBook(ctx, isbn, request.Title, request.Author)

	if err != nil {
		status := mapErrorToStatus(err)
		ctx.JSON(status, gin.H{"error": sanitizeError(err, status)})
		log.Printf("UpdateBook error for ISBN %s: %v", isbn, err)
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

func (c *BookController) DeleteBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")

	if isbn == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ISBN parameter is required"})
		return
	}

	err := c.core.DeleteBook(ctx, isbn)

	if err != nil {
		status := mapErrorToStatus(err)
		ctx.JSON(status, gin.H{"error": sanitizeError(err, status)})
		log.Printf("DeleteBook error for ISBN %s: %v", isbn, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func mapErrorToStatus(err error) int {
	if errors.Is(err, interfaces.ErrBookNotFound) {
		return http.StatusNotFound
	}
	errMsg := err.Error()
	if contains(errMsg, "cannot be empty", "invalid", "required", "already exists", "ISBN must be", "checksum") {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

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

func sanitizeError(err error, status int) string {
	switch status {
	case http.StatusNotFound:
		return "resource not found"
	case http.StatusBadRequest:
		return "invalid request"
	default:
		return "internal server error"
	}
}

