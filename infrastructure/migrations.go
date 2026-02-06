package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Migration struct {
	ID          int
	Name        string
	SQL         string
	Description string
}

var Migrations = []Migration{
	{
		ID:          1,
		Name:        "create_books_table",
		Description: "Creates the initial books table",
		SQL: `
			CREATE TABLE IF NOT EXISTS books (
				isbn VARCHAR(13) PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				author VARCHAR(255) NOT NULL,
				published_at TIMESTAMP NOT NULL
			);
		`,
	},
}

func RunMigrations(db *sql.DB) error {
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			applied_at TIMESTAMP NOT NULL
		);
	`

	_, err := db.Exec(createMigrationsTable)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	appliedMigrations := make(map[int]bool)
	rows, err := db.Query("SELECT id FROM migrations")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan migration id: %w", err)
		}
		appliedMigrations[id] = true
	}

	for _, migration := range Migrations {
		if !appliedMigrations[migration.ID] {
			log.Printf("Applying migration: %s (%s)", migration.Name, migration.Description)

			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.ID, err)
			}

		_, err = tx.Exec(migration.SQL)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", migration.ID, err)
		}

		_, err = tx.Exec(
			"INSERT INTO migrations (id, name, description, applied_at) VALUES ($1, $2, $3, $4)",
			migration.ID, migration.Name, migration.Description, time.Now(),
		)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.ID, err)
		}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit transaction for migration %d: %w", migration.ID, err)
			}

			log.Printf("Successfully applied migration: %s", migration.Name)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}