package interfaces

import (
	"context"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
)

type ArticleRepository interface {
	Save(ctx context.Context, article *entities.Article) error
	FindByID(ctx context.Context, id int64) (*entities.Article, error)
	FindByAuthorIDPaginated(ctx context.Context, authorID int64, status *vo.ArticleStatus, page, perPage int) (articles []*entities.Article, total int, err error)
	FindPublishedPaginated(ctx context.Context, page, perPage int, orderBy vo.ArticleOrderBy) (articles []*vo.PublishedArticleWithScore, total int, err error)
	FindPublishedByAuthorID(ctx context.Context, authorID int64) ([]*entities.Article, error)
	CountByAuthorID(ctx context.Context, authorID int64) (int, error)
	CountPublishedByAuthorID(ctx context.Context, authorID int64) (int, error)
	Update(ctx context.Context, article *entities.Article) error
	GetTopAuthorsByScore(ctx context.Context, limit int) ([]*vo.TopAuthor, error)
}
