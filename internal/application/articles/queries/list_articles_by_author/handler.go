package list_articles_by_author

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
)

type Handler struct {
	articleRepo interfaces.ArticleRepository
	authorRepo  interfaces.AuthorRepository
}

func NewHandler(articleRepo interfaces.ArticleRepository, authorRepo interfaces.AuthorRepository) *Handler {
	return &Handler{
		articleRepo: articleRepo,
		authorRepo:  authorRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Response, error) {
	// Verify that author exists
	_, err := h.authorRepo.FindByID(ctx, query.AuthorID)
	if err != nil {
		return nil, err
	}

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
