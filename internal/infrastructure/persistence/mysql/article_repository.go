package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
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
		return nil, fmt.Errorf("article not found with id %d", id)
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
		WHERE author_id = ? AND status = 'PUBLISHED'
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
	query := "SELECT COUNT(*) FROM articles WHERE author_id = ? AND status = 'PUBLISHED'"

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
	countQuery := "SELECT COUNT(*) FROM articles WHERE status = 'PUBLISHED'"
	var total int
	if err := exec.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count published articles: %w", err)
	}

	scoreFormula := `
		(a.word_count * 0.1) + 
		(SELECT COUNT(*) FROM articles WHERE author_id = a.author_id AND status = 'PUBLISHED') * 5 +
		CASE
			WHEN TIMESTAMPDIFF(HOUR, a.published_at, NOW()) < 24 THEN 50
			WHEN TIMESTAMPDIFF(HOUR, a.published_at, NOW()) < 72 THEN 20
			ELSE 0
		END
	`

	orderClause := "a.published_at DESC"
	if orderBy.String() == "score" {
		orderClause = "score DESC, a.published_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT 
			a.id,
			a.author_id,
			au.name,
			a.title,
			a.body,
			a.word_count,
			a.published_at,
			(%s) as score
		FROM articles a
		INNER JOIN authors au ON a.author_id = au.id
		WHERE a.status = 'PUBLISHED'
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, scoreFormula, orderClause)

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

func (r *ArticleRepository) GetTopAuthorsByScore(ctx context.Context, limit int) ([]*vo.TopAuthor, error) {
	exec := getExecutor(ctx, r.db)
	scoreFormula := `
		(a.word_count * 0.1) + 
		(SELECT COUNT(*) FROM articles WHERE author_id = a.author_id AND status = 'PUBLISHED') * 5 +
		CASE
			WHEN TIMESTAMPDIFF(HOUR, a.published_at, NOW()) < 24 THEN 50
			WHEN TIMESTAMPDIFF(HOUR, a.published_at, NOW()) < 72 THEN 20
			ELSE 0
		END
	`

	query := fmt.Sprintf(`
		SELECT 
			au.id,
			au.name,
			SUM(%s) as total_score,
			COUNT(*) as published_count
		FROM articles a
		INNER JOIN authors au ON a.author_id = au.id
		WHERE a.status = 'PUBLISHED'
		GROUP BY au.id, au.name
		ORDER BY total_score DESC, au.id ASC
		LIMIT ?
	`, scoreFormula)

	rows, err := exec.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top authors: %w", err)
	}
	defer rows.Close()

	var topAuthors []*vo.TopAuthor
	for rows.Next() {
		topAuthor, err := mappers.ScanTopAuthor(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan top author: %w", err)
		}
		topAuthors = append(topAuthors, topAuthor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating top authors: %w", err)
	}

	return topAuthors, nil
}
