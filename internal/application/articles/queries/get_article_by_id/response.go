package get_article_by_id

import "time"

type Response struct {
	ID          int64
	AuthorID    int64
	Title       string
	Body        string
	Status      string
	WordCount   uint
	PublishedAt *time.Time
	Score       *float64
}
