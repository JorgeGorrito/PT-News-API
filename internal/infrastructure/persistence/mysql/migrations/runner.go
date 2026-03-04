package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Migration struct {
	Name string
	SQL  string
}

var AllMigrations = []Migration{
	{
		Name: "001_create_authors_table",
		SQL:  CreateAuthorsTable,
	},
	{
		Name: "002_create_articles_table",
		SQL:  CreateArticlesTable,
	},
	{
		Name: "003_alter_articles_status",
		SQL:  AlterArticlesStatus,
	},
}

type Runner struct {
	db *sql.DB
}

func NewRunner(db *sql.DB) *Runner {
	return &Runner{db: db}
}

// Run executes all pending migrations
func (r *Runner) Run(ctx context.Context) error {
	log.Println("Starting database migrations...")

	// Create migrations table if it doesn't exist
	if err := r.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get executed migrations
	executedMigrations, err := r.getExecutedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	// Execute pending migrations
	for _, migration := range AllMigrations {
		if executedMigrations[migration.Name] {
			log.Printf("Migration '%s' already executed, skipping...", migration.Name)
			continue
		}

		log.Printf("Executing migration '%s'...", migration.Name)
		if err := r.executeMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to execute migration '%s': %w", migration.Name, err)
		}
		log.Printf("Migration '%s' executed successfully", migration.Name)
	}

	log.Println("All migrations completed successfully")
	return nil
}

func (r *Runner) createMigrationsTable(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, CreateMigrationsTable)
	return err
}

func (r *Runner) getExecutedMigrations(ctx context.Context) (map[string]bool, error) {
	query := "SELECT name FROM migrations"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		executed[name] = true
	}

	return executed, rows.Err()
}

func (r *Runner) executeMigration(ctx context.Context, migration Migration) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
		return err
	}

	// Record migration
	insertQuery := "INSERT INTO migrations (name, executed_at) VALUES (?, ?)"
	if _, err := tx.ExecContext(ctx, insertQuery, migration.Name, time.Now().UTC()); err != nil {
		return err
	}

	return tx.Commit()
}
