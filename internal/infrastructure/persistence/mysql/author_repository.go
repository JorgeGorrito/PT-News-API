package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
	"github.com/JorgeGorrito/PT-News-API/internal/infrastructure/persistence/mappers"
)

type AuthorRepository struct {
	db *DB
}

func NewAuthorRepository(db *DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) Save(ctx context.Context, author *entities.Author) error {
	query := `
		INSERT INTO authors (name, email, biography, created_at)
		VALUES (?, ?, ?, NOW())
	`

	exec := getExecutor(ctx, r.db)
	result, err := exec.ExecContext(ctx, query, author.Name(), author.Email(), author.Biography())
	if err != nil {
		return fmt.Errorf("failed to save author: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	author.SetID(id)
	return nil
}

func (r *AuthorRepository) FindByID(ctx context.Context, id int64) (*entities.Author, error) {
	query := `
		SELECT id, name, email, biography, created_at
		FROM authors
		WHERE id = ?
	`

	exec := getExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, id)
	author, err := mappers.ScanAuthorRow(row)
	if err == sql.ErrNoRows {
		return nil, domerrs.NewDomainError(fmt.Sprintf("autor no encontrado con id %d", id), domerrs.NotFoundError)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find author: %w", err)
	}

	return author, nil
}

func (r *AuthorRepository) FindByEmail(ctx context.Context, email string) (*entities.Author, error) {
	query := `
		SELECT id, name, email, biography, created_at
		FROM authors
		WHERE email = ?
	`

	exec := getExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, email)
	author, err := mappers.ScanAuthorRow(row)
	if err == sql.ErrNoRows {
		return nil, domerrs.NewDomainError(fmt.Sprintf("autor no encontrado con email %s", email), domerrs.NotFoundError)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find author: %w", err)
	}

	return author, nil
}

// GetSummary retorna las estadísticas básicas de publicación del autor.
// El total_score NO se calcula aquí; es responsabilidad de la capa de aplicación
// usando el Domain Service de puntuación sobre los artículos publicados.
func (r *AuthorRepository) GetSummary(ctx context.Context, authorID int64) (*vo.AuthorSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_articles,
			COUNT(CASE WHEN status = 'PUBLICADO' THEN 1 END) as total_published,
			MAX(published_at) as last_publication_at
		FROM articles
		WHERE author_id = ?
	`

	var totalArticles, totalPublished int
	var lastPublicationAt sql.NullTime

	exec := getExecutor(ctx, r.db)
	err := exec.QueryRowContext(ctx, query, authorID).Scan(&totalArticles, &totalPublished, &lastPublicationAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get author summary: %w", err)
	}

	var lastPubPtr *time.Time
	if lastPublicationAt.Valid {
		lastPubPtr = &lastPublicationAt.Time
	}

	return vo.NewAuthorSummary(authorID, totalArticles, totalPublished, lastPubPtr), nil
}

func (r *AuthorRepository) FindAll(ctx context.Context) ([]*entities.Author, error) {
	query := `
		SELECT id, name, email, biography, created_at
		FROM authors
		ORDER BY created_at DESC
	`

	exec := getExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all authors: %w", err)
	}
	defer rows.Close()

	var authors []*entities.Author
	for rows.Next() {
		author, err := mappers.ScanAuthor(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan author: %w", err)
		}
		authors = append(authors, author)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating authors: %w", err)
	}

	return authors, nil
}
