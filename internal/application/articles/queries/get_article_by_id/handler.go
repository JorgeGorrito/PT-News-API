package get_article_by_id

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
	services "github.com/JorgeGorrito/PT-News-API/internal/domain/services"
)

type Handler struct {
	articleRepo  interfaces.ArticleRepository
	authorRepo   interfaces.AuthorRepository
	scoreService services.IScoreService
}

func NewHandler(
	articleRepo interfaces.ArticleRepository,
	authorRepo interfaces.AuthorRepository,
	scoreService services.IScoreService,
) *Handler {
	return &Handler{
		articleRepo:  articleRepo,
		authorRepo:   authorRepo,
		scoreService: scoreService,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Response, error) {
	article, err := h.articleRepo.FindByID(ctx, query.ArticleID)
	if err != nil {
		return nil, err
	}

	author, err := h.authorRepo.FindByID(ctx, article.AuthorID())
	if err != nil {
		return nil, err
	}

	response := &Response{
		ID:          article.ID(),
		AuthorID:    article.AuthorID(),
		AuthorName:  author.Name(),
		Title:       article.Title(),
		Body:        article.Body(),
		Status:      article.Status().String(),
		WordCount:   article.WordCount(),
		CreatedAt:   article.CreatedAt(),
		PublishedAt: article.PublishedAt(),
	}

	if query.IncludeScore && article.IsPublished() {
		// Calcular score dinámicamente desde el dominio sin consultar la BD por score.
		publishedCount, err := h.articleRepo.CountPublishedByAuthorID(ctx, article.AuthorID())
		if err != nil {
			return nil, err
		}
		score := article.CalculateScore(h.scoreService, publishedCount)
		response.Score = &score
	}

	return response, nil
}
