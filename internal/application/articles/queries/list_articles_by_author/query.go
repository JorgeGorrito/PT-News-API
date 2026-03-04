package list_articles_by_author

import valueobjects "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"

type Query struct {
	AuthorID int64
	Status   *valueobjects.ArticleStatus
	Page     int
	PerPage  int
}
