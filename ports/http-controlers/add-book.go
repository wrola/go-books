package httpControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddBook(context *gin.Context) {
	var newBook Book

	if err := context.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)

	context.IndentedJSON(http.StatusCreated, newBook)
}