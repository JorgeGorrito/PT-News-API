package get_author_summary

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
	services "github.com/JorgeGorrito/PT-News-API/internal/domain/services"
)

type Handler struct {
	authorRepo   interfaces.AuthorRepository
	articleRepo  interfaces.ArticleRepository
	scoreService services.IScoreService
}

func NewHandler(
	authorRepo interfaces.AuthorRepository,
	articleRepo interfaces.ArticleRepository,
	scoreService services.IScoreService,
) *Handler {
	return &Handler{
		authorRepo:   authorRepo,
		articleRepo:  articleRepo,
		scoreService: scoreService,
	}
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

	// Calcular el score total del autor en memoria usando el Domain Service.
	// Cada artículo delega el cálculo a través de Article.CalculateScore(svc, authorPublishedCount).
	publishedArticles, err := h.articleRepo.FindPublishedByAuthorID(ctx, query.AuthorID)
	if err != nil {
		return nil, err
	}

	var totalScore float64
	for _, article := range publishedArticles {
		totalScore += article.CalculateScore(h.scoreService, summary.TotalPublished())
	}

	draftArticles := summary.TotalArticles() - summary.TotalPublished()

	return &Response{
		ID:                author.ID(),
		Name:              author.Name(),
		Email:             author.Email(),
		Biography:         author.Biography(),
		PublishedArticles: summary.TotalPublished(),
		DraftArticles:     draftArticles,
		TotalScore:        totalScore,
	}, nil
}
