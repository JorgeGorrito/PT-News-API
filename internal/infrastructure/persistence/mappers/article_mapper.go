package mappers

import (
	"database/sql"
	"time"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
)

func ScanArticle(rows *sql.Rows) (*entities.Article, error) {
	var id, authorID int64
	var title, body, statusStr string
	var wordCount uint
	var createdAt time.Time
	var publishedAt sql.NullTime

	if err := rows.Scan(&id, &authorID, &title, &body, &statusStr, &wordCount, &createdAt, &publishedAt); err != nil {
		return nil, err
	}

	article := entities.NewArticle(authorID, title)
	if err := article.SetBody(body); err != nil {
		return nil, err
	}

	article.SetID(id)

	if statusStr == "PUBLISHED" && publishedAt.Valid {
		if err := article.Publish(); err != nil {
			return nil, err
		}
	}

	return article, nil
}

func ScanArticleRow(row *sql.Row) (*entities.Article, error) {
	var id, authorID int64
	var title, body, statusStr string
	var wordCount uint
	var createdAt time.Time
	var publishedAt sql.NullTime

	if err := row.Scan(&id, &authorID, &title, &body, &statusStr, &wordCount, &createdAt, &publishedAt); err != nil {
		return nil, err
	}

	article := entities.NewArticle(authorID, title)
	if err := article.SetBody(body); err != nil {
		return nil, err
	}

	article.SetID(id)

	if statusStr == "PUBLISHED" && publishedAt.Valid {
		if err := article.Publish(); err != nil {
			return nil, err
		}
	}

	return article, nil
}

func ScanPublishedArticleWithScore(rows *sql.Rows) (*vo.PublishedArticleWithScore, error) {
	var articleID, authorID int64
	var authorName, title, body string
	var wordCount uint
	var publishedAt time.Time
	var authorPublishedCount int

	if err := rows.Scan(&articleID, &authorID, &authorName, &title, &body, &wordCount, &publishedAt, &authorPublishedCount); err != nil {
		return nil, err
	}

	return vo.NewPublishedArticleWithScore(
		articleID,
		authorID,
		authorName,
		title,
		body,
		wordCount,
		publishedAt,
		authorPublishedCount,
	), nil
}
