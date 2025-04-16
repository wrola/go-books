package commands

import (
	"books/core/storage/models"
	"context"
	"errors"
	"testing"
	"time"
)

func TestUpdateBookCommandHandler_Handle(t *testing.T) {
	// Create a test book
	testBook, _ := models.NewBook("1234567890", "Original Title", "Original Author", time.Now())

	// Test cases
	tests := []struct {
		name           string
		setupRepo      func(*MockBookRepository)
		command        interface{}
		wantErr        bool
		expectedErr    error
		validateResult func(*testing.T, *MockBookRepository)
	}{
		{
			name: "Update all fields",
			setupRepo: func(repo *MockBookRepository) {
				repo.books = append(repo.books, testBook)
			},
			command: UpdateBookCommand{
				ISBN:  testBook.ISBN,
				Title: "New Title",
				Author: "New Author",
			},
			wantErr: false,
			validateResult: func(t *testing.T, repo *MockBookRepository) {
				// Find the book
				book, err := repo.FindByISBN(context.Background(), testBook.ISBN)
				if err != nil {
					t.Errorf("Couldn't find updated book: %v", err)
					return
				}

				// Check fields were updated
				if book.Title != "New Title" {
					t.Errorf("Title not updated. Expected 'New Title', got '%s'", book.Title)
				}
				if book.Author != "New Author" {
					t.Errorf("Author not updated. Expected 'New Author', got '%s'", book.Author)
				}
				if book.ISBN != testBook.ISBN {
					t.Errorf("ISBN changed unexpectedly. Expected '%s', got '%s'", testBook.ISBN, book.ISBN)
				}
			},
		},
		{
			name: "Partial update - title only",
			setupRepo: func(repo *MockBookRepository) {
				repo.books = append(repo.books, testBook)
			},
			command: UpdateBookCommand{
				ISBN:  testBook.ISBN,
				Title: "New Title Only",
			},
			wantErr: false,
			validateResult: func(t *testing.T, repo *MockBookRepository) {
				book, _ := repo.FindByISBN(context.Background(), testBook.ISBN)
				if book.Title != "New Title Only" {
					t.Errorf("Title not updated. Expected 'New Title Only', got '%s'", book.Title)
				}
				if book.Author != testBook.Author {
					t.Errorf("Author changed unexpectedly. Expected '%s', got '%s'", testBook.Author, book.Author)
				}
				if book.ISBN != testBook.ISBN {
					t.Errorf("ISBN changed unexpectedly. Expected '%s', got '%s'", testBook.ISBN, book.ISBN)
				}
			},
		},
		{
			name: "Non-existent book",
			setupRepo: func(repo *MockBookRepository) {
				// Empty repository
			},
			command: UpdateBookCommand{
				ISBN:  "non-existent-isbn",
				Title: "New Title",
				Author: "New Author",
			},
			wantErr:     true,
			expectedErr: ErrBookNotFound,
		},
		{
			name: "Empty ID",
			setupRepo: func(repo *MockBookRepository) {},
			command: UpdateBookCommand{
				ISBN:  "",
				Title: "New Title",
				Author: "New Author",
			},
			wantErr:     true,
			expectedErr: errors.New("book ISBN cannot be empty"),
		},
		{
			name: "All empty fields - nothing to update",
			setupRepo: func(repo *MockBookRepository) {
				repo.books = append(repo.books, testBook)
			},
				command: UpdateBookCommand{
				ISBN: testBook.ISBN,
			},
			wantErr:     true,
			expectedErr: errors.New("at least one field must be provided for update"),
		},
		{
			name:      "Invalid command type",
			setupRepo: func(repo *MockBookRepository) {},
			command:   "not a command",
			wantErr:   true,
			expectedErr: ErrInvalidCommandType,
		},
		{
			name: "Repository error during lookup",
			setupRepo: func(repo *MockBookRepository) {
				repo.findByIDError = errors.New("database lookup error")
			},
			command: UpdateBookCommand{
				ISBN:  testBook.ISBN,
				Title: "New Title",
				Author: "New Author",
			},
			wantErr:     true,
			expectedErr: errors.New("database lookup error"),
		},
		{
			name: "Repository error during save",
			setupRepo: func(repo *MockBookRepository) {
				repo.books = append(repo.books, testBook)
				repo.saveError = errors.New("database save error")
			},
			command: UpdateBookCommand{
				ISBN:  testBook.ISBN,
				Title: "New Title",
				Author: "New Author",
			},
			wantErr:     true,
			expectedErr: errors.New("database save error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockBookRepository()
			tt.setupRepo(mockRepo)

			handler := NewUpdateBookCommandHandler(mockRepo)

			// Execute
			err := handler.Handle(context.Background(), tt.command)

			// Verify
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check specific error type if expected
			if tt.expectedErr != nil && err != nil {
				if tt.expectedErr.Error() != err.Error() {
					t.Errorf("Expected error '%v', got '%v'", tt.expectedErr, err)
				}
			}

			// Validate result for successful updates
			if !tt.wantErr && tt.validateResult != nil {
				tt.validateResult(t, mockRepo)
			}
		})
	}
}