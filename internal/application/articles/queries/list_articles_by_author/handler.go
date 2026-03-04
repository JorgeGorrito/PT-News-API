package list_articles_by_author

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
)

type Handler struct {
	articleRepo interfaces.ArticleRepository
}

func NewHandler(articleRepo interfaces.ArticleRepository) *Handler {
	return &Handler{articleRepo: articleRepo}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Response, error) {
	articles, total, err := h.articleRepo.FindByAuthorIDPaginated(ctx, query.AuthorID, query.Status, query.Page, query.PerPage)
	if err != nil {
		return nil, err
	}

	items := make([]ArticleItem, 0, len(articles))
	for _, article := range articles {
		items = append(items, ArticleItem{
			ID:          article.ID(),
			AuthorID:    article.AuthorID(),
			Title:       article.Title(),
			Body:        article.Body(),
			Status:      article.Status().String(),
			WordCount:   article.WordCount(),
			PublishedAt: article.PublishedAt(),
		})
	}

	return &Response{
		Articles: items,
		Total:    total,
		Page:     query.Page,
		PerPage:  query.PerPage,
	}, nil
}
