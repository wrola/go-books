package httpControllers

import (
	"books/core"
)

// Module provides a simplified API for starting the HTTP server
type Module struct {
	httpServer *HTTPServer
}

// NewModule creates a new HTTP module
func NewModule(core *core.Core) *Module {
	return &Module{
		httpServer: NewHTTPServer(core),
	}
}

// Start starts the HTTP server on the specified address
func (m *Module) Start(address string) error {
	return m.httpServer.Start(address)
}