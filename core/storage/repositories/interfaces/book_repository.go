package interfaces

import (
	"context"
	"errors"

	"books/core/storage/models"
)

// BookStoragePostgresRepository defines the repository interface for book storage operations
type BookStoragePostgresRepository interface {
	Save(ctx context.Context, book *models.Book) error
	FindAll(ctx context.Context) ([]*models.Book, error)
	FindByISBN(ctx context.Context, isbn string) (*models.Book, error)
	Delete(ctx context.Context, isbn string) error
}

// ErrBookNotFound is returned when a book is not found in the repository
var ErrBookNotFound = errors.New("book not found")