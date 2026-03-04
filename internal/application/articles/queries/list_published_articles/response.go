package list_published_articles

import "time"

type ArticleItem struct {
	ID          int64
	AuthorID    int64
	AuthorName  string
	Title       string
	Body        string
	WordCount   uint
	PublishedAt time.Time
	Score       float64
}

type Response struct {
	Articles []ArticleItem
	Total    int
	Page     int
	PerPage  int
}
