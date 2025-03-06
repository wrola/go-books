package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{
		Title:  "s",
		Author: "",
	},
}

func AddBook(context *gin.Context) {
	var newBook Book

	if err := context.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)

	context.IndentedJSON(http.StatusCreated, newBook)
}

func GetBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}
