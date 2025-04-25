package commands

import (
	"context"
	"testing"
	"time"

	"books/core/storage/models"
	"books/core/storage/repositories"
	"books/core/storage/repositories/interfaces"
)

func TestUpdateBookCommandHandler_Handle(t *testing.T) {
	testBook := &models.Book{
		ISBN:        "123",
		Title:       "Test Book",
		Author:      "Test Author",
		PublishedAt: time.Now(),
	}

	tests := []struct {
		name           string
		setupRepo      func(*repositories.BookStorageInMemoryRepository)
		command        *UpdateBookCommand
		validateResult func(*testing.T, *repositories.BookStorageInMemoryRepository)
		wantErr        bool
		expectedErr    error
	}{
		{
			name: "Update all fields",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {
				repo.Save(context.Background(), testBook)
			},
			command: &UpdateBookCommand{
				ISBN:  testBook.ISBN,
				Title: "Updated Title",
				Author: "Updated Author",
			},
			wantErr: false,
			validateResult: func(t *testing.T, repo *repositories.BookStorageInMemoryRepository) {
				book, err := repo.FindByISBN(context.Background(), testBook.ISBN)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if book.Title != "Updated Title" {
					t.Errorf("expected title 'Updated Title', got '%s'", book.Title)
				}
				if book.Author != "Updated Author" {
					t.Errorf("expected author 'Updated Author', got '%s'", book.Author)
				}
			},
		},
		{
			name: "Partial update - title only",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {
				repo.Save(context.Background(), testBook)
			},
			command: &UpdateBookCommand{
				ISBN:  testBook.ISBN,
				Title: "New Title Only",
			},
			wantErr: false,
			validateResult: func(t *testing.T, repo *repositories.BookStorageInMemoryRepository) {
				book, _ := repo.FindByISBN(context.Background(), testBook.ISBN)
				if book.Title != "New Title Only" {
					t.Errorf("expected title 'New Title Only', got '%s'", book.Title)
				}
				if book.Author != testBook.Author {
					t.Errorf("expected author '%s', got '%s'", testBook.Author, book.Author)
				}
			},
		},
		{
			name: "book not found",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {
				// Empty repository
			},
			command: &UpdateBookCommand{
				ISBN:  "123",
				Title: "Updated Title",
				Author: "Updated Author",
			},
			wantErr:     true,
			expectedErr: interfaces.ErrBookNotFound,
		},
		{
			name: "Empty ID",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &UpdateBookCommand{
				ISBN:  "",
				Title: "New Title",
				Author: "New Author",
			},
			wantErr:   true,
			expectedErr: ErrInvalidCommandType,
		},
		{
			name: "All empty fields - nothing to update",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {
				repo.Save(context.Background(), testBook)
			},
			command: &UpdateBookCommand{
				ISBN:  testBook.ISBN,
			},
			wantErr: false,
			validateResult: func(t *testing.T, repo *repositories.BookStorageInMemoryRepository) {
				book, _ := repo.FindByISBN(context.Background(), testBook.ISBN)
				if book.Title != testBook.Title {
					t.Errorf("expected title '%s', got '%s'", testBook.Title, book.Title)
				}
				if book.Author != testBook.Author {
					t.Errorf("expected author '%s', got '%s'", testBook.Author, book.Author)
				}
			},
		},
		{
			name:      "Invalid command type",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command:   &UpdateBookCommand{},
			wantErr:   true,
			expectedErr: ErrInvalidCommandType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repositories.NewBookStorageInMemoryRepository()
			if tt.setupRepo != nil {
				tt.setupRepo(mockRepo)
			}

			handler := NewUpdateBookCommandHandler(mockRepo)
			err := handler.Handle(context.Background(), tt.command)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				if err != tt.expectedErr {
					t.Errorf("expected error %v but got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.validateResult != nil {
				tt.validateResult(t, mockRepo)
			}
		})
	}
}