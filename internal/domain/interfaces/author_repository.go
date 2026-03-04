package interfaces

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
)

type AuthorRepository interface {
	Save(ctx context.Context, author *entities.Author) error
	FindByID(ctx context.Context, id int64) (*entities.Author, error)
	FindByEmail(ctx context.Context, email string) (*entities.Author, error)
	GetSummary(ctx context.Context, authorID int64) (*vo.AuthorSummary, error)
	FindAll(ctx context.Context) ([]*entities.Author, error)
}
