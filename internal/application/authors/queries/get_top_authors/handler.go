package get_top_authors

import (
	"context"

	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
	services "github.com/JorgeGorrito/PT-News-API/internal/domain/services"
)

type Handler struct {
	articleRepo  interfaces.ArticleRepository
	scoreService services.IScoreService
}

func NewHandler(articleRepo interfaces.ArticleRepository, scoreService services.IScoreService) *Handler {
	return &Handler{articleRepo: articleRepo, scoreService: scoreService}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Response, error) {
	if query.Limit <= 0 {
		return nil, domerrs.NewDomainError(
			"limit must be greater than 0",
			domerrs.GeneralError,
		)
	}

	// El cálculo se realiza completamente en memoria tal como exigen los requisitos.
	// La BD solo provee los datos crudos (sin score); el Domain Service hace el resto.
	allArticles, err := h.articleRepo.FindAllPublished(ctx)
	if err != nil {
		return nil, err
	}

	topAuthors := h.scoreService.CalculateTopAuthors(allArticles, query.Limit)

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
