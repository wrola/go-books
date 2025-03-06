package main

import (
	"github.com/gin-gonic/gin"
	"books/controllers"
)

func main() {
	router := gin.Default()
		router.GET("/books", controllers.GetBooks)
	
	router.POST("/books", controllers.AddBook)

	router.Run("localhost:8080")
}
