package repositories

import (
	"context"
	"sync"

	"books/core/commands"
	"books/core/models"
)

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

	// Check if book with this ID already exists (for update case)
	for i, existingBook := range r.books {
		if existingBook.ID == book.ID {
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

// FindByID returns a book by its ID
func (r *InMemoryBookRepository) FindByID(ctx context.Context, id string) (*models.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, book := range r.books {
		if book.ID == id {
			// Return a copy to prevent external modifications
			return &models.Book{
				ID:     book.ID,
				Title:  book.Title,
				Author: book.Author,
				ISBN:   book.ISBN,
			}, nil
		}
	}

	return nil, commands.ErrBookNotFound
}

// Delete removes a book from the repository
func (r *InMemoryBookRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, book := range r.books {
		if book.ID == id {
			// Remove the book by replacing it with the last element and truncating slice
			r.books[i] = r.books[len(r.books)-1]
			r.books = r.books[:len(r.books)-1]
			return nil
		}
	}

	return commands.ErrBookNotFound
}

// Ensure InMemoryBookRepository implements BookRepository interface
var _ commands.BookRepository = (*InMemoryBookRepository)(nil)