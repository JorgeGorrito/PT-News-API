package list_published_articles

import (
	"context"
	"sort"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/interfaces"
	services "github.com/JorgeGorrito/PT-News-API/internal/domain/services"
	valueobjects "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
)

type Handler struct {
	articleRepo  interfaces.ArticleRepository
	scoreService services.IScoreService
}

func NewHandler(articleRepo interfaces.ArticleRepository, scoreService services.IScoreService) *Handler {
	return &Handler{articleRepo: articleRepo, scoreService: scoreService}
}

func (h *Handler) Handle(ctx context.Context, query Query) (*Response, error) {
	orderBy, err := valueobjects.NewArticleOrderBy(query.OrderBy)
	if err != nil {
		return nil, err
	}

	var pagedItems []*valueobjects.PublishedArticleWithScore
	var total int

	if orderBy.String() == "score" {
		// Para orden por score: se trae todo en memoria, se calcula, se ordena y pagina
		// en la capa de dominio. El score nunca toca la base de datos.
		allItems, err := h.articleRepo.FindAllPublished(ctx)
		if err != nil {
			return nil, err
		}
		total = len(allItems)

		// Ordenar en memoria usando el Domain Service (cálculo dinámico de score)
		sorted := sortByScore(allItems, h.scoreService)
		pagedItems = paginate(sorted, query.Page, query.PerPage)
	} else {
		// Para orden por fecha: la paginación ocurre en SQL de forma eficiente.
		items, t, err := h.articleRepo.FindPublishedPaginated(ctx, query.Page, query.PerPage, orderBy)
		if err != nil {
			return nil, err
		}
		total = t
		pagedItems = items
	}

	// Calcular score de cada artículo para la respuesta usando el Domain Service.
	articles := make([]ArticleItem, 0, len(pagedItems))
	for _, item := range pagedItems {
		score := h.scoreService.CalculateArticleScore(
			item.WordCount(),
			item.AuthorPublishedCount(),
			item.PublishedAt(),
		)
		articles = append(articles, ArticleItem{
			ID:          item.ArticleID(),
			AuthorID:    item.AuthorID(),
			AuthorName:  item.AuthorName(),
			Title:       item.Title(),
			Body:        item.Body(),
			WordCount:   item.WordCount(),
			PublishedAt: item.PublishedAt(),
			Score:       score,
		})
	}

	return &Response{
		Articles: articles,
		Total:    total,
		Page:     query.Page,
		PerPage:  query.PerPage,
	}, nil
}

// sortByScore ordena artículos por score descendente calculado dinámicamente
// y por fecha de publicación descendente en caso de empate.
func sortByScore(
	items []*valueobjects.PublishedArticleWithScore,
	svc services.IScoreService,
) []*valueobjects.PublishedArticleWithScore {
	result := make([]*valueobjects.PublishedArticleWithScore, len(items))
	copy(result, items)

	sort.Slice(result, func(i, j int) bool {
		si := svc.CalculateArticleScore(result[i].WordCount(), result[i].AuthorPublishedCount(), result[i].PublishedAt())
		sj := svc.CalculateArticleScore(result[j].WordCount(), result[j].AuthorPublishedCount(), result[j].PublishedAt())
		if si != sj {
			return si > sj
		}
		return result[i].PublishedAt().After(result[j].PublishedAt())
	})

	return result
}

// paginate aplica paginación a un slice de artículos.
func paginate(
	items []*valueobjects.PublishedArticleWithScore,
	page, perPage int,
) []*valueobjects.PublishedArticleWithScore {
	offset := (page - 1) * perPage
	if offset >= len(items) {
		return []*valueobjects.PublishedArticleWithScore{}
	}
	end := offset + perPage
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end]
}
