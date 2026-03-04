package get_article_by_id

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
	article, err := h.articleRepo.FindByID(ctx, query.ArticleID)
	if err != nil {
		return nil, err
	}

	response := &Response{
		ID:          article.ID(),
		AuthorID:    article.AuthorID(),
		Title:       article.Title(),
		Body:        article.Body(),
		Status:      article.Status().String(),
		WordCount:   article.WordCount(),
		PublishedAt: article.PublishedAt(),
	}

	if query.IncludeScore && article.IsPublished() {
		items, _, err := h.articleRepo.FindPublishedPaginated(ctx, 1, 1000, valueobjects.OrderByScore)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			if item.ArticleID() == article.ID() {
				score := item.Score()
				response.Score = &score
				break
			}
		}
	}

	return response, nil
}
