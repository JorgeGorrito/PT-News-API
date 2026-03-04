package get_article_by_id

import "time"

type Response struct {
	ID          int64
	AuthorID    int64
	AuthorName  string
	Title       string
	Body        string
	Status      string
	WordCount   uint
	CreatedAt   time.Time
	PublishedAt *time.Time
	Score       *float64
}
