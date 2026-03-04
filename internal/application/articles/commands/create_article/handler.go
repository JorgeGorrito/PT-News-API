package create_article

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
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

func (h *Handler) Handle(ctx context.Context, cmd Command) (*Response, error) {
	_, err := h.authorRepo.FindByID(ctx, cmd.AuthorID)
	if err != nil {
		return nil, err
	}

	article := entities.NewArticle(cmd.AuthorID, cmd.Title)

	if err := article.SetBody(cmd.Body); err != nil {
		return nil, err
	}

	if err := h.articleRepo.Save(ctx, article); err != nil {
		return nil, err
	}

	return &Response{ID: article.ID()}, nil
}
