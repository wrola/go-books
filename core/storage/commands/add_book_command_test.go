package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"books/core/storage/models"
	"books/core/storage/repositories"
)

func TestAddBookCommandHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*repositories.BookStorageInMemoryRepository)
		command        interface{}
		validateResult func(*testing.T, *repositories.BookStorageInMemoryRepository)
		wantErr        bool
		expectedErr    error
	}{
		{
			name:      "successful add",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &AddBookCommand{
				ISBN:  "1234567890",
				Title: "Test Book",
				Author: "Test Author",
			},
			wantErr: false,
			validateResult: func(t *testing.T, repo *repositories.BookStorageInMemoryRepository) {
				book, err := repo.FindByISBN(context.Background(), "1234567890")
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if book.Title != "Test Book" {
					t.Errorf("expected title 'Test Book', got '%s'", book.Title)
				}
				if book.Author != "Test Author" {
					t.Errorf("expected author 'Test Author', got '%s'", book.Author)
				}
			},
		},
		{
			name:      "duplicate ISBN",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {
				book, _ := models.NewBook("1234567890", "Existing Book", "Existing Author", time.Now())
				repo.Save(context.Background(), book)
			},
			command: &AddBookCommand{
				ISBN:  "1234567890",
				Title: "Test Book",
				Author: "Test Author",
			},
			wantErr: true,
			expectedErr: errors.New("book with ISBN 1234567890 already exists"),
		},
		{
			name:      "empty ISBN",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &AddBookCommand{
				ISBN:  "",
				Title: "Test Book",
				Author: "Test Author",
			},
			wantErr:   true,
			expectedErr: errors.New("ISBN cannot be empty"),
		},
		{
			name:      "empty title",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &AddBookCommand{
				ISBN:  "1234567890",
				Title: "",
				Author: "Test Author",
			},
			wantErr:   true,
			expectedErr: errors.New("title cannot be empty"),
		},
		{
			name:      "empty author",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &AddBookCommand{
				ISBN:  "1234567890",
				Title: "Test Book",
				Author: "",
			},
			wantErr:   true,
			expectedErr: errors.New("author cannot be empty"),
		},
		{
			name:      "invalid command type",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command:   &DeleteBookCommand{ISBN: "1234567890"},
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

			handler := NewAddBookCommandHandler(mockRepo)
			err := handler.Handle(context.Background(), tt.command)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				if err.Error() != tt.expectedErr.Error() {
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