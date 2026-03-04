package create_article

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/application/transaction"
	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
)

type Handler struct {
	articleRepo interfaces.ArticleRepository
	authorRepo  interfaces.AuthorRepository
	txManager   transaction.Manager
}

func NewHandler(
	articleRepo interfaces.ArticleRepository,
	authorRepo interfaces.AuthorRepository,
	txManager transaction.Manager,
) *Handler {
	return &Handler{
		articleRepo: articleRepo,
		authorRepo:  authorRepo,
		txManager:   txManager,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (*Response, error) {
	var articleID int64

	err := h.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		_, err := h.authorRepo.FindByID(txCtx, cmd.AuthorID)
		if err != nil {
			return err
		}

		article := entities.NewArticle(cmd.AuthorID, cmd.Title)

		if err := article.SetBody(cmd.Body); err != nil {
			return err
		}

		if err := h.articleRepo.Save(txCtx, article); err != nil {
			return err
		}

		articleID = article.ID()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Response{ID: articleID}, nil
}
