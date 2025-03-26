package repositories

import (
	"context"
	"errors"
	"sync"

	"books/core/storage/models"
)

// BookRepository defines the repository interface for book storage operations
type BookRepository interface {
	Save(ctx context.Context, book *models.Book) error
	FindAll(ctx context.Context) ([]*models.Book, error)
	FindByISBN(ctx context.Context, isbn string) (*models.Book, error)
	Delete(ctx context.Context, isbn string) error
}

// InMemoryBookRepository is a simple in-memory implementation of BookRepository
type InMemoryBookRepository struct {
	books []*models.Book
	mutex sync.RWMutex
}

// NewInMemoryBookRepository creates a new in-memory book repository
func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		books: make([]*models.Book, 0),
	}
}

// Save adds a book to the repository
func (r *InMemoryBookRepository) Save(ctx context.Context, book *models.Book) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if book with this ISBN already exists (for update case)
	for i, existingBook := range r.books {
		if existingBook.ISBN == book.ISBN {
			// Replace the existing book
			r.books[i] = book
			return nil
		}
	}

	// Book doesn't exist, add it
	r.books = append(r.books, book)
	return nil
}

// FindAll returns all books in the repository
func (r *InMemoryBookRepository) FindAll(ctx context.Context) ([]*models.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Return a copy to prevent external modifications
	result := make([]*models.Book, len(r.books))
	copy(result, r.books)
	return result, nil
}

// FindByISBN returns a book by its ISBN
func (r *InMemoryBookRepository) FindByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, book := range r.books {
		if book.ISBN == isbn {
			// Return a copy to prevent external modifications
			return &models.Book{
				ISBN:        book.ISBN,
				Title:       book.Title,
				Author:      book.Author,
				PublishedAt: book.PublishedAt,
			}, nil
		}
	}

	return nil, ErrBookNotFound
}

// Delete removes a book from the repository
func (r *InMemoryBookRepository) Delete(ctx context.Context, isbn string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, book := range r.books {
		if book.ISBN == isbn {
			// Remove the book by replacing it with the last element and truncating slice
			r.books[i] = r.books[len(r.books)-1]
			r.books = r.books[:len(r.books)-1]
			return nil
		}
	}

	return ErrBookNotFound
}

// ErrBookNotFound is returned when a book is not found in the repository
var ErrBookNotFound = errors.New("book not found")

// Ensure InMemoryBookRepository implements BookRepository interface
var _ BookRepository = (*InMemoryBookRepository)(nil)