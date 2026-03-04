package get_top_authors

import (
	"context"

	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
)

type Handler struct {
	articleRepo interfaces.ArticleRepository
}

func NewHandler(articleRepo interfaces.ArticleRepository) *Handler {
	return &Handler{articleRepo: articleRepo}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Response, error) {
	if query.Limit <= 0 {
		return nil, domerrs.NewDomainError(
			"limit must be greater than 0",
			domerrs.GeneralError,
		)
	}

	topAuthors, err := h.articleRepo.GetTopAuthorsByScore(ctx, query.Limit)
	if err != nil {
		return nil, err
	}

	authors := make([]AuthorItem, 0, len(topAuthors))
	for _, ta := range topAuthors {
		authors = append(authors, AuthorItem{
			ID:                ta.AuthorID(),
			Name:              ta.AuthorName(),
			TotalScore:        ta.TotalScore(),
			PublishedArticles: ta.PublishedCount(),
		})
	}

	return &Response{Authors: authors}, nil
}
