package publish_article

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/application/transaction"
	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
)

type Handler struct {
	articleRepo interfaces.ArticleRepository
	txManager   transaction.Manager
}

func NewHandler(articleRepo interfaces.ArticleRepository, txManager transaction.Manager) *Handler {
	return &Handler{
		articleRepo: articleRepo,
		txManager:   txManager,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (*Response, error) {
	var articleID int64

	err := h.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		article, err := h.articleRepo.FindByID(txCtx, cmd.ArticleID)
		if err != nil {
			return err
		}

		if err := article.Publish(); err != nil {
			return err
		}

		if err := h.articleRepo.Update(txCtx, article); err != nil {
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
