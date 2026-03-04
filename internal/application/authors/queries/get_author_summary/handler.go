package get_author_summary

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
)

type Handler struct {
	authorRepo interfaces.AuthorRepository
}

func NewHandler(authorRepo interfaces.AuthorRepository) *Handler {
	return &Handler{authorRepo: authorRepo}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Response, error) {
	author, err := h.authorRepo.FindByID(ctx, query.AuthorID)
	if err != nil {
		return nil, err
	}

	summary, err := h.authorRepo.GetSummary(ctx, query.AuthorID)
	if err != nil {
		return nil, err
	}

	draftArticles := summary.TotalArticles() - summary.TotalPublished()

	return &Response{
		ID:                author.ID(),
		Name:              author.Name(),
		Email:             author.Email(),
		Biography:         author.Biography(),
		PublishedArticles: summary.TotalPublished(),
		DraftArticles:     draftArticles,
	}, nil
}
