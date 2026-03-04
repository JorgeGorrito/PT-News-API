package list_published_articles

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
	valueobjects "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
)

type Handler struct {
	articleRepo interfaces.ArticleRepository
}

func NewHandler(articleRepo interfaces.ArticleRepository) *Handler {
	return &Handler{articleRepo: articleRepo}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Response, error) {
	orderBy, err := valueobjects.NewArticleOrderBy(query.OrderBy)
	if err != nil {
		return nil, err
	}

	items, total, err := h.articleRepo.FindPublishedPaginated(ctx, query.Page, query.PerPage, orderBy)
	if err != nil {
		return nil, err
	}

	articles := make([]ArticleItem, 0, len(items))
	for _, item := range items {
		articles = append(articles, ArticleItem{
			ID:          item.ArticleID(),
			AuthorID:    item.AuthorID(),
			AuthorName:  item.AuthorName(),
			Title:       item.Title(),
			Body:        item.Body(),
			WordCount:   item.WordCount(),
			PublishedAt: item.PublishedAt(),
			Score:       item.Score(),
		})
	}

	return &Response{
		Articles: articles,
		Total:    total,
		Page:     query.Page,
		PerPage:  query.PerPage,
	}, nil
}
