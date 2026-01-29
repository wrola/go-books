package interfaces

import (
	"context"
	"errors"

	"books/core/storage/models"
)

type BookRepository interface {
	Save(ctx context.Context, book *models.Book) error
	FindAll(ctx context.Context) ([]*models.Book, error)
	FindByISBN(ctx context.Context, isbn string) (*models.Book, error)
	Delete(ctx context.Context, isbn string) error
}

var ErrBookNotFound = errors.New("book not found")