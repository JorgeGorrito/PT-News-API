package publish_article

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

func (h *Handler) Handle(ctx context.Context, cmd Command) (*Response, error) {
	article, err := h.articleRepo.FindByID(ctx, cmd.ArticleID)
	if err != nil {
		return nil, err
	}

	if err := article.Publish(); err != nil {
		return nil, err
	}

	if err := h.articleRepo.Update(ctx, article); err != nil {
		return nil, err
	}

	return &Response{ID: article.ID()}, nil
}
