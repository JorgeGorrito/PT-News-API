package create_author

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
)

type Handler struct {
	authorRepo interfaces.AuthorRepository
}

func NewHandler(authorRepo interfaces.AuthorRepository) *Handler {
	return &Handler{authorRepo: authorRepo}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (*Response, error) {
	author, err := entities.NewAuthor(cmd.Name, cmd.Email)
	if err != nil {
		return nil, err
	}

	if err := author.SetBiography(cmd.Biography); err != nil {
		return nil, err
	}

	if err := h.authorRepo.Save(ctx, author); err != nil {
		return nil, err
	}

	return &Response{ID: author.ID()}, nil
}
