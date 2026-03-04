package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
	"github.com/JorgeGorrito/PT-News-API/internal/infrastructure/persistence/mappers"
)

type ArticleRepository struct {
	db *DB
}

func NewArticleRepository(db *DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Save(ctx context.Context, article *entities.Article) error {
	query := `
		INSERT INTO articles (author_id, title, body, status, word_count, created_at)
		VALUES (?, ?, ?, ?, ?, NOW())
	`

	exec := getExecutor(ctx, r.db)
	result, err := exec.ExecContext(
		ctx,
		query,
		article.AuthorID(),
		article.Title(),
		article.Body(),
		article.Status().String(),
		article.WordCount(),
	)
	if err != nil {
		return fmt.Errorf("failed to save article: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	article.SetID(id)
	return nil
}

func (r *ArticleRepository) FindByID(ctx context.Context, id int64) (*entities.Article, error) {
	query := `
		SELECT id, author_id, title, body, status, word_count, created_at, published_at
		FROM articles
		WHERE id = ?
	`

	exec := getExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, id)
	article, err := mappers.ScanArticleRow(row)
	if err == sql.ErrNoRows {
		return nil, domerrs.NewDomainError(fmt.Sprintf("artículo no encontrado con id %d", id), domerrs.NotFoundError)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find article: %w", err)
	}

	return article, nil
}

func (r *ArticleRepository) Update(ctx context.Context, article *entities.Article) error {
	query := `
		UPDATE articles
		SET title = ?, body = ?, status = ?, word_count = ?, published_at = ?
		WHERE id = ?
	`

	exec := getExecutor(ctx, r.db)
	_, err := exec.ExecContext(
		ctx,
		query,
		article.Title(),
		article.Body(),
		article.Status().String(),
		article.WordCount(),
		article.PublishedAt(),
		article.ID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	return nil
}

func (r *ArticleRepository) FindByAuthorIDPaginated(
	ctx context.Context,
	authorID int64,
	status *vo.ArticleStatus,
	page, perPage int,
) ([]*entities.Article, int, error) {
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "author_id = ?")
	args = append(args, authorID)

	if status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, status.String())
	}

	whereClause := strings.Join(conditions, " AND ")

	exec := getExecutor(ctx, r.db)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM articles WHERE %s", whereClause)
	var total int
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, author_id, title, body, status, word_count, created_at, published_at
		FROM articles
		WHERE %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	offset := (page - 1) * perPage
	args = append(args, perPage, offset)

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find articles: %w", err)
	}
	defer rows.Close()

	var articles []*entities.Article
	for rows.Next() {
		article, err := mappers.ScanArticle(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating articles: %w", err)
	}

	return articles, total, nil
}

func (r *ArticleRepository) FindPublishedByAuthorID(ctx context.Context, authorID int64) ([]*entities.Article, error) {
	query := `
		SELECT id, author_id, title, body, status, word_count, created_at, published_at
		FROM articles
		WHERE author_id = ? AND status = 'PUBLICADO'
		ORDER BY published_at DESC
	`

	exec := getExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find published articles: %w", err)
	}
	defer rows.Close()

	var articles []*entities.Article
	for rows.Next() {
		article, err := mappers.ScanArticle(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating articles: %w", err)
	}

	return articles, nil
}

func (r *ArticleRepository) CountByAuthorID(ctx context.Context, authorID int64) (int, error) {
	query := "SELECT COUNT(*) FROM articles WHERE author_id = ?"

	exec := getExecutor(ctx, r.db)
	var count int
	if err := exec.QueryRowContext(ctx, query, authorID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count articles: %w", err)
	}

	return count, nil
}

func (r *ArticleRepository) CountPublishedByAuthorID(ctx context.Context, authorID int64) (int, error) {
	query := "SELECT COUNT(*) FROM articles WHERE author_id = ? AND status = 'PUBLICADO'"

	exec := getExecutor(ctx, r.db)
	var count int
	if err := exec.QueryRowContext(ctx, query, authorID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count published articles: %w", err)
	}

	return count, nil
}

func (r *ArticleRepository) FindPublishedPaginated(
	ctx context.Context,
	page, perPage int,
	orderBy vo.ArticleOrderBy,
) ([]*vo.PublishedArticleWithScore, int, error) {
	exec := getExecutor(ctx, r.db)
	countQuery := "SELECT COUNT(*) FROM articles WHERE status = 'PUBLICADO'"
	var total int
	if err := exec.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count published articles: %w", err)
	}

	// El score NO se calcula en SQL. Se retorna author_published_count para que
	// el Domain Service calcule el score dinámicamente en la capa de aplicación.
	// La paginación SQL se usa para orden por fecha; el orden por score se delega
	// a FindAllPublished + ScoreService en el handler.
	_ = orderBy // el ordenamiento por score ocurre en memoria en la capa de aplicación
	query := `
		SELECT 
			a.id,
			a.author_id,
			au.name,
			a.title,
			a.body,
			a.word_count,
			a.published_at,
			(SELECT COUNT(*) FROM articles WHERE author_id = a.author_id AND status = 'PUBLICADO') as author_published_count
		FROM articles a
		INNER JOIN authors au ON a.author_id = au.id
		WHERE a.status = 'PUBLICADO'
		ORDER BY a.published_at DESC
		LIMIT ? OFFSET ?
	`

	offset := (page - 1) * perPage

	rows, err := exec.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find published articles: %w", err)
	}
	defer rows.Close()

	var articles []*vo.PublishedArticleWithScore
	for rows.Next() {
		article, err := mappers.ScanPublishedArticleWithScore(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan published article: %w", err)
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating articles: %w", err)
	}

	return articles, total, nil
}

// FindAllPublished retorna todos los artículos publicados junto con la cantidad de
// publicaciones del autor. No aplica paginación ni calcula el score en SQL;
// el Domain Service calcula los scores y el orden en memoria.
func (r *ArticleRepository) FindAllPublished(ctx context.Context) ([]*vo.PublishedArticleWithScore, error) {
	query := `
		SELECT 
			a.id,
			a.author_id,
			au.name,
			a.title,
			a.body,
			a.word_count,
			a.published_at,
			(SELECT COUNT(*) FROM articles WHERE author_id = a.author_id AND status = 'PUBLICADO') as author_published_count
		FROM articles a
		INNER JOIN authors au ON a.author_id = au.id
		WHERE a.status = 'PUBLICADO'
		ORDER BY a.published_at DESC
	`

	exec := getExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all published articles: %w", err)
	}
	defer rows.Close()

	var articles []*vo.PublishedArticleWithScore
	for rows.Next() {
		article, err := mappers.ScanPublishedArticleWithScore(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan published article: %w", err)
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating published articles: %w", err)
	}

	return articles, nil
}
