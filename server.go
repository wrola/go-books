package main

import (
    "net/http"
	"github.com/gin-gonic/gin"
)

type book struct {
	isbn string `json:"isbn"`
	title string `json:"title"`
	author string `json:"author"`
}

var books = []book{
	{isbn: "1", title: "Bible", author: "Pope"},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func addBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)

	c.IndentedJSON(http.StatusCreated, newBook)
}
func main() {
    router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", addBook)

	router.Run("localhost:8080")
}