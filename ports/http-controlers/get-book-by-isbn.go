package httpControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBookByISBN(c *gin.Context) {
	isbn := c.Param("isbn")

	for _, aBook := range books {
		if aBook.ISBN == isbn {
			c.IndentedJSON(http.StatusOK, aBook)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}
