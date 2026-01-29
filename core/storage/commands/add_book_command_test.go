package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"books/core/storage/models"
	"books/core/storage/repositories"
)

type testCase struct {
	name           string
	setupRepo      func(*repositories.BookStorageInMemoryRepository)
	command        interface{}
	validateResult func(*testing.T, *repositories.BookStorageInMemoryRepository)
	wantErr        bool
	expectedErr    error
}

func getTestCases() []testCase {
	// Valid ISBN-13: 978-3-16-148410-0 (checksum valid)
	validISBN := "9783161484100"

	return []testCase{
		{
			name:      "successful add",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &AddBookCommand{
				ISBN:   validISBN,
				Title:  "Test Book",
				Author: "Test Author",
			},
			wantErr: false,
			validateResult: func(t *testing.T, repo *repositories.BookStorageInMemoryRepository) {
				book, err := repo.FindByISBN(context.Background(), validISBN)
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
			name: "duplicate ISBN",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {
				book, _ := models.NewBook(validISBN, "Existing Book", "Existing Author", time.Now())
				repo.Save(context.Background(), book)
			},
			command: &AddBookCommand{
				ISBN:   validISBN,
				Title:  "Test Book",
				Author: "Test Author",
			},
			wantErr:     true,
			expectedErr: errors.New("failed to save book: book with ISBN " + validISBN + " already exists"),
		},
		{
			name:      "empty ISBN",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &AddBookCommand{
				ISBN:   "",
				Title:  "Test Book",
				Author: "Test Author",
			},
			wantErr:     true,
			expectedErr: errors.New("ISBN cannot be empty"),
		},
		{
			name:      "empty title",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &AddBookCommand{
				ISBN:   validISBN,
				Title:  "",
				Author: "Test Author",
			},
			wantErr:     true,
			expectedErr: errors.New("title cannot be empty"),
		},
		{
			name:      "empty author",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &AddBookCommand{
				ISBN:   validISBN,
				Title:  "Test Book",
				Author: "",
			},
			wantErr:     true,
			expectedErr: errors.New("author cannot be empty"),
		},
		{
			name:      "invalid command type",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command:   &DeleteBookCommand{ISBN: validISBN},
			wantErr:   true,
			expectedErr: ErrInvalidCommandType,
		},
		{
			name:      "invalid ISBN format",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &AddBookCommand{
				ISBN:   "invalid-isbn",
				Title:  "Test Book",
				Author: "Test Author",
			},
			wantErr:     true,
			expectedErr: errors.New("ISBN must be 10 or 13 characters (excluding hyphens)"),
		},
	}
}

func TestAddBookCommandHandler_Handle(t *testing.T) {
	tests := getTestCases()

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