package core

import (
	"books/core/storage/commands"
	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"
	"context"
	"time"
)

// Core represents the application core with all available commands and queries
type Core struct {
	commandBus  commands.CommandBus
	repository  interfaces.BookRepository // private field - not exported
}

// NewCore creates a new instance of the application core
func NewCore(bookRepository interfaces.BookRepository) *Core {
	// Create command bus
	commandBus := commands.NewCommandBus()

	// Register command handlers
	addBookHandler := commands.NewAddBookCommandHandler(bookRepository)
	updateBookHandler := commands.NewUpdateBookCommandHandler(bookRepository)
	deleteBookHandler := commands.NewDeleteBookCommandHandler(bookRepository)

	commandBus.RegisterHandler("*commands.AddBookCommand", addBookHandler)
	commandBus.RegisterHandler("*commands.UpdateBookCommand", updateBookHandler)
	commandBus.RegisterHandler("*commands.DeleteBookCommand", deleteBookHandler)

	// Return the configured core
	return &Core{
		commandBus: commandBus,
		repository: bookRepository,
	}
}

// Commands

// AddBook handles the add book command
func (c *Core) AddBook(ctx context.Context, title, author, isbn string) (*models.Book, error) {
	cmd := &commands.AddBookCommand{
		Title:  title,
		Author: author,
		ISBN:   isbn,
	}

	err := c.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Create a book response
	// Note: In a real implementation, you might want to return the actual created book
	// from the repository or use a query to fetch it
	book, err := models.NewBook(isbn, title, author, time.Now())
	if err != nil {
		return nil, err
	}
	return book, nil
}

// UpdateBook handles updating a book
func (c *Core) UpdateBook(ctx context.Context, isbn, title, author string) (*models.Book, error) {
	cmd := &commands.UpdateBookCommand{
		ISBN:  isbn,
		Title: title,
		Author: author,
	}

	err := c.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// After updating, get the book
	return c.GetBookByISBN(ctx, isbn)
}

// DeleteBook handles deleting a book
func (c *Core) DeleteBook(ctx context.Context, isbn string) error {
	cmd := &commands.DeleteBookCommand{
		ISBN: isbn,
	}

	return c.commandBus.Dispatch(ctx, cmd)
}

// Queries

// GetAllBooks returns all books
func (c *Core) GetAllBooks(ctx context.Context) ([]*models.Book, error) {
	// In a true CQRS implementation, this would use a query bus
	// For simplicity, we're directly using the repository here
	// In production, consider implementing a proper query bus
	return c.repository.FindAll(ctx)
}

// GetBookByID returns a book by its ID
func (c *Core) GetBookByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	// In a true CQRS implementation, this would use a query bus
	// For simplicity, we're directly using the repository here
	return c.repository.FindByISBN(ctx, isbn)
}
