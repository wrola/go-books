package repositories

import (
	"context"
	"sync"

	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"
)

type BookStorageInMemoryRepository struct {
	books []*models.Book
	mutex sync.RWMutex
}

func NewBookStorageInMemoryRepository() *BookStorageInMemoryRepository {
	return &BookStorageInMemoryRepository{
		books: make([]*models.Book, 0),
	}
}

func (r *BookStorageInMemoryRepository) Save(ctx context.Context, book *models.Book) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, existingBook := range r.books {
		if existingBook.ISBN == book.ISBN {
			r.books[i] = book
			return nil
		}
	}

	r.books = append(r.books, book)
	return nil
}

func (r *BookStorageInMemoryRepository) FindAll(ctx context.Context) ([]*models.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*models.Book, len(r.books))
	copy(result, r.books)
	return result, nil
}

func (r *BookStorageInMemoryRepository) FindByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, book := range r.books {
		if book.ISBN == isbn {
			return &models.Book{
				ISBN:        book.ISBN,
				Title:       book.Title,
				Author:      book.Author,
				PublishedAt: book.PublishedAt,
			}, nil
		}
	}

	return nil, interfaces.ErrBookNotFound
}

func (r *BookStorageInMemoryRepository) Delete(ctx context.Context, isbn string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, book := range r.books {
		if book.ISBN == isbn {
			r.books[i] = r.books[len(r.books)-1]
			r.books = r.books[:len(r.books)-1]
			return nil
		}
	}

	return interfaces.ErrBookNotFound
}

var _ interfaces.BookRepository = (*BookStorageInMemoryRepository)(nil)