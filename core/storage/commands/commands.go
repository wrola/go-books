package commands

import (
	"context"
	"errors"
	"reflect"

	"books/core/storage/models"
)

// BookRepository defines the interface for book storage
type BookRepository interface {
	Save(ctx context.Context, book *models.Book) error
	FindAll(ctx context.Context) ([]*models.Book, error)
	FindByID(ctx context.Context, id string) (*models.Book, error)
	Delete(ctx context.Context, id string) error
}

// CommandHandler is the signature for all command handlers
type CommandHandler interface {
	Handle(ctx context.Context, command interface{}) error
}

// CommandBus is the interface for dispatching commands
type CommandBus interface {
	Dispatch(ctx context.Context, command interface{}) error
}

// DefaultCommandBus is a simple implementation of CommandBus
type DefaultCommandBus struct {
	handlers map[string]CommandHandler
}

// NewCommandBus creates a new command bus
func NewCommandBus() *DefaultCommandBus {
	return &DefaultCommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

// RegisterHandler registers a command handler for a specific command type
func (b *DefaultCommandBus) RegisterHandler(commandType string, handler CommandHandler) {
	b.handlers[commandType] = handler
}

// Dispatch sends a command to its appropriate handler
func (b *DefaultCommandBus) Dispatch(ctx context.Context, command interface{}) error {
	commandType := getCommandType(command)
	handler, exists := b.handlers[commandType]
	if !exists {
		return ErrHandlerNotFound
	}
	return handler.Handle(ctx, command)
}

// getCommandType returns the name of the command type
func getCommandType(command interface{}) string {
	return reflect.TypeOf(command).String()
}

// Error definitions
var (
	ErrHandlerNotFound   = errors.New("handler not found for command")
	ErrInvalidCommandType = errors.New("invalid command type")
	ErrBookNotFound       = errors.New("book not found")
)

// BookRepository is the interface for the book repository