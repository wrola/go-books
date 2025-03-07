package ports

import (
	core "books/core"
	httpControllers "books/ports/http-controlers"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	app *core.Core
}

func NewHttpServer(app *core.Core) *HttpServer {
	return &HttpServer{
		app: app,
	}
}

func (s *HttpServer) Start() {
	router := gin.Default()
	router.GET("/books", httpControllers.GetBooks)
	router.POST("/books", httpControllers.AddBook)
	err := router.Run("localhost:8080")
	if err != nil {
		panic("Failed to start HTTP server: " + err.Error())
	}
}
