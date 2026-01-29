package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/lib/pq"
)

var db *sql.DB
var repo *BookStoragePostgresRepository

func TestMain(m *testing.M) {
	// Uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Printf("Could not construct pool: %s - skipping integration tests", err)
		os.Exit(0)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Printf("Could not connect to Docker: %s - skipping integration tests", err)
		os.Exit(0)
	}

	// Pull postgres image and create container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15-alpine",
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=testdb",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://test:test@%s/testdb?sslmode=disable", hostAndPort)

	log.Printf("Connecting to database on url: %s", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Run migrations
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			isbn VARCHAR(13) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			published_at TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Could not run migrations: %s", err)
	}

	// Create repository
	repo = NewBookStoragePostgresRepository(db)

	// Run tests
	code := m.Run()

	// Cleanup
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func cleanupDB(t *testing.T) {
	_, err := db.Exec("DELETE FROM books")
	if err != nil {
		t.Fatalf("Failed to cleanup database: %v", err)
	}
}

func TestSave(t *testing.T) {
	cleanupDB(t)

	validISBN := "9783161484100"
	book, _ := models.NewBook(validISBN, "Test Book", "Test Author", time.Now())

	err := repo.Save(context.Background(), book)
	if err != nil {
		t.Fatalf("Failed to save book: %v", err)
	}

	// Verify it was saved
	savedBook, err := repo.FindByISBN(context.Background(), validISBN)
	if err != nil {
		t.Fatalf("Failed to find saved book: %v", err)
	}

	if savedBook.Title != book.Title {
		t.Errorf("expected title %s, got %s", book.Title, savedBook.Title)
	}
	if savedBook.Author != book.Author {
		t.Errorf("expected author %s, got %s", book.Author, savedBook.Author)
	}
}

func TestSaveUpsert(t *testing.T) {
	cleanupDB(t)

	validISBN := "9783161484100"
	book, _ := models.NewBook(validISBN, "Original Title", "Original Author", time.Now())

	// Save initial
	err := repo.Save(context.Background(), book)
	if err != nil {
		t.Fatalf("Failed to save book: %v", err)
	}

	// Update via upsert
	book.Title = "Updated Title"
	book.Author = "Updated Author"
	err = repo.Save(context.Background(), book)
	if err != nil {
		t.Fatalf("Failed to upsert book: %v", err)
	}

	// Verify update
	savedBook, _ := repo.FindByISBN(context.Background(), validISBN)
	if savedBook.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", savedBook.Title)
	}
	if savedBook.Author != "Updated Author" {
		t.Errorf("expected author 'Updated Author', got %s", savedBook.Author)
	}
}

func TestFindAll(t *testing.T) {
	cleanupDB(t)

	// Add multiple books
	isbns := []string{"9783161484100", "9780306406157", "9780596517748"}
	for i, isbn := range isbns {
		book, _ := models.NewBook(isbn, fmt.Sprintf("Book %d", i+1), fmt.Sprintf("Author %d", i+1), time.Now())
		repo.Save(context.Background(), book)
	}

	books, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("Failed to find all books: %v", err)
	}

	if len(books) != 3 {
		t.Errorf("expected 3 books, got %d", len(books))
	}
}

func TestFindByISBN(t *testing.T) {
	cleanupDB(t)

	validISBN := "9783161484100"
	book, _ := models.NewBook(validISBN, "Test Book", "Test Author", time.Now())
	repo.Save(context.Background(), book)

	// Test found
	foundBook, err := repo.FindByISBN(context.Background(), validISBN)
	if err != nil {
		t.Fatalf("Failed to find book: %v", err)
	}
	if foundBook.ISBN != validISBN {
		t.Errorf("expected ISBN %s, got %s", validISBN, foundBook.ISBN)
	}

	// Test not found
	_, err = repo.FindByISBN(context.Background(), "9780306406157")
	if err != interfaces.ErrBookNotFound {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	cleanupDB(t)

	validISBN := "9783161484100"
	book, _ := models.NewBook(validISBN, "Test Book", "Test Author", time.Now())
	repo.Save(context.Background(), book)

	// Delete
	err := repo.Delete(context.Background(), validISBN)
	if err != nil {
		t.Fatalf("Failed to delete book: %v", err)
	}

	// Verify deleted
	_, err = repo.FindByISBN(context.Background(), validISBN)
	if err != interfaces.ErrBookNotFound {
		t.Errorf("expected ErrBookNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	cleanupDB(t)

	err := repo.Delete(context.Background(), "nonexistent")
	if err != interfaces.ErrBookNotFound {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}
