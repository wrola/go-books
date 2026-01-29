package httpControllers

import (
	"database/sql"

	"books/core"
)

// Module provides a simplified API for starting the HTTP server
type Module struct {
	server *Server
}

// NewModule creates a new HTTP module
func NewModule(core *core.Core) *Module {
	return &Module{
		server: NewServer(core),
	}
}

// NewModuleWithDB creates a new HTTP module with database health check support
func NewModuleWithDB(core *core.Core, db *sql.DB) *Module {
	return &Module{
		server: NewServerWithDB(core, db),
	}
}

// Start starts the HTTP server on the specified address
func (m *Module) Start(address string) error {
	return m.server.Start(address)
}