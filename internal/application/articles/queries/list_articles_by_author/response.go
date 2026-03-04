package list_articles_by_author

import "time"

type ArticleItem struct {
	ID          int64
	AuthorID    int64
	Title       string
	Body        string
	Status      string
	WordCount   uint
	PublishedAt *time.Time
}

type Response struct {
	Articles []ArticleItem
	Total    int
	Page     int
	PerPage  int
}
