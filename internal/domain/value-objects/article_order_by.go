package valueobjects

import domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"

type ArticleOrderBy string

const (
	OrderByPublishedAt ArticleOrderBy = "published_at"
	OrderByScore       ArticleOrderBy = "score"
)

func (o ArticleOrderBy) IsValid() bool {
	switch o {
	case OrderByPublishedAt, OrderByScore:
		return true
	default:
		return false
	}
}

func (o ArticleOrderBy) String() string {
	return string(o)
}

func NewArticleOrderBy(orderBy string) (ArticleOrderBy, error) {
	o := ArticleOrderBy(orderBy)
	if !o.IsValid() {
		return "", domerrs.NewDomainError("valor de orderBy inválido: debe ser 'published_at' o 'score'", domerrs.InvalidArticleOrderByError)
	}
	return o, nil
}
