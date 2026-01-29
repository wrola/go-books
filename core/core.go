package core

import (
	"books/core/storage/commands"
	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"
	"context"
	"time"
)

type Core struct {
	commandBus  commands.CommandBus
	repository  interfaces.BookRepository
}

func NewCore(bookRepository interfaces.BookRepository) *Core {
	commandBus := commands.NewCommandBus()

	addBookHandler := commands.NewAddBookCommandHandler(bookRepository)
	updateBookHandler := commands.NewUpdateBookCommandHandler(bookRepository)
	deleteBookHandler := commands.NewDeleteBookCommandHandler(bookRepository)

	commandBus.RegisterHandler("*commands.AddBookCommand", addBookHandler)
	commandBus.RegisterHandler("*commands.UpdateBookCommand", updateBookHandler)
	commandBus.RegisterHandler("*commands.DeleteBookCommand", deleteBookHandler)

	return &Core{
		commandBus: commandBus,
		repository: bookRepository,
	}
}

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

	book, err := models.NewBook(isbn, title, author, time.Now())
	if err != nil {
		return nil, err
	}
	return book, nil
}

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

	return c.GetBookByISBN(ctx, isbn)
}

func (c *Core) DeleteBook(ctx context.Context, isbn string) error {
	cmd := &commands.DeleteBookCommand{
		ISBN: isbn,
	}

	return c.commandBus.Dispatch(ctx, cmd)
}

func (c *Core) GetAllBooks(ctx context.Context) ([]*models.Book, error) {
	return c.repository.FindAll(ctx)
}

func (c *Core) GetBookByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	return c.repository.FindByISBN(ctx, isbn)
}
