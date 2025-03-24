package core

import (
	"context"

	"books/core/commands"
	"books/core/models"
)

// Core represents the application core with all available commands and queries
type Core struct {
	commandBus  commands.CommandBus
	repository  commands.BookRepository // private field - not exported
}

// NewCore creates a new instance of the application core
func NewCore(bookRepository commands.BookRepository) *Core {
	// Create command bus
	commandBus := commands.NewCommandBus()

	// Register command handlers
	addBookHandler := commands.NewAddBookCommandHandler(bookRepository)
	updateBookHandler := commands.NewUpdateBookCommandHandler(bookRepository)
	deleteBookHandler := commands.NewDeleteBookCommandHandler(bookRepository)

	commandBus.RegisterHandler("commands.AddBookCommand", addBookHandler)
	commandBus.RegisterHandler("commands.UpdateBookCommand", updateBookHandler)
	commandBus.RegisterHandler("commands.DeleteBookCommand", deleteBookHandler)

	// Return the configured core
	return &Core{
		commandBus: commandBus,
		repository: bookRepository,
	}
}

// Commands

// AddBook handles the add book command
func (c *Core) AddBook(ctx context.Context, title, author, isbn string) (*models.Book, error) {
	cmd := commands.AddBookCommand{
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
	book, _ := models.NewBook(title, author, isbn)
	return book, nil
}

// UpdateBook handles updating a book
func (c *Core) UpdateBook(ctx context.Context, id, title, author, isbn string) (*models.Book, error) {
	cmd := commands.UpdateBookCommand{
		ID:     id,
		Title:  title,
		Author: author,
		ISBN:   isbn,
	}

	err := c.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// After updating, get the book
	return c.GetBookByID(ctx, id)
}

// DeleteBook handles deleting a book
func (c *Core) DeleteBook(ctx context.Context, id string) error {
	cmd := commands.DeleteBookCommand{
		ID: id,
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
func (c *Core) GetBookByID(ctx context.Context, id string) (*models.Book, error) {
	// In a true CQRS implementation, this would use a query bus
	// For simplicity, we're directly using the repository here
	return c.repository.FindByID(ctx, id)
}

// GetBookByISBN returns a book by its ISBN
func (c *Core) GetBookByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	// In a true CQRS implementation, this would use a query bus
	// For now, we'll implement a simple search through all books
	books, err := c.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, book := range books {
		if book.ISBN == isbn {
			return book, nil
		}
	}

	return nil, commands.ErrBookNotFound
}