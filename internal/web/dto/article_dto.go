package dto

import "time"

// CreateArticleRequest represents the request to create an article
type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=255"`
	Content  string `json:"content" binding:"required,min=1"`
	AuthorID int64  `json:"author_id" binding:"required,gt=0"`
}

// CreateArticleResponse represents the response after creating an article
type CreateArticleResponse struct {
	ID int64 `json:"id"`
}

// PublishArticleResponse represents the response after publishing an article
type PublishArticleResponse struct {
	ID int64 `json:"id"`
}

// ArticleResponse represents a complete article
type ArticleResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	WordCount   int        `json:"word_count"`
	AuthorID    int64      `json:"author_id"`
	AuthorName  string     `json:"author_name"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// PublishedArticleResponse represents a published article with score
type PublishedArticleResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	WordCount   int       `json:"word_count"`
	AuthorID    int64     `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	PublishedAt time.Time `json:"published_at"`
	Score       float64   `json:"score"`
}

// ListArticlesResponse represents a paginated list of articles
type ListArticlesResponse struct {
	Articles   []PublishedArticleResponse `json:"articles"`
	TotalCount int                        `json:"total_count"`
	Page       int                        `json:"page"`
	PageSize   int                        `json:"page_size"`
}
