package services

import (
	"sort"
	"time"

	constants "github.com/JorgeGorrito/PT-News-API/internal/domain/constants"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
)

// IScoreService define el contrato del Domain Service de puntuación.
// Permite inyectar la dependencia en la capa de aplicación facilitando testing
// y cumpliendo con el principio de inversión de dependencias (DIP).
// Solo expone las operaciones que los consumidores externos realmente necesitan.
type IScoreService interface {
	CalculateArticleScore(wordCount uint, authorPublishedCount int, publishedAt time.Time) float64
	CalculateTopAuthors(articles []*vo.PublishedArticleWithScore, limit int) []*vo.TopAuthor
}

// ScoreService es un Domain Service responsable de calcular el score de relevancia
// de los artículos. Encapsula la regla de negocio del cálculo del score que depende
// de datos de múltiples fuentes (artículo y estadísticas del autor).
type ScoreService struct{}

// NewScoreService crea una nueva instancia del servicio de score.
func NewScoreService() *ScoreService {
	return &ScoreService{}
}

// CalculateArticleScore calcula el score de relevancia de un artículo publicado.
//
// Fórmula: (cantidad_palabras * 0.1) + (articulos_publicados_autor * 5) + bonus_reciente
//
// bonus_reciente:
//   - 50 puntos si fue publicado hace menos de 24 horas
//   - 20 puntos si fue publicado hace menos de 72 horas
//   - 0 en caso contrario
//
// Parámetros:
//   - wordCount: cantidad de palabras del artículo
//   - authorPublishedCount: cantidad de artículos publicados por el autor
//   - publishedAt: fecha de publicación del artículo
//
// Retorna el score calculado dinámicamente.
func (s *ScoreService) CalculateArticleScore(wordCount uint, authorPublishedCount int, publishedAt time.Time) float64 {
	baseScore := s.calculateBaseScore(wordCount, authorPublishedCount)
	recencyBonus := s.calculateRecencyBonus(publishedAt)

	return baseScore + recencyBonus
}

// calculateBaseScore calcula el score base sin bonus de recencia.
func (s *ScoreService) calculateBaseScore(wordCount uint, authorPublishedCount int) float64 {
	wordScore := float64(wordCount) * constants.ScorePerWord
	authorScore := float64(authorPublishedCount) * constants.ScorePerAuthorArticle

	return wordScore + authorScore
}

// calculateRecencyBonus calcula el bonus por recencia según la fecha de publicación.
func (s *ScoreService) calculateRecencyBonus(publishedAt time.Time) float64 {
	hoursSincePublished := time.Since(publishedAt).Hours()

	switch {
	case hoursSincePublished < constants.Hours24:
		return constants.BonusRecent24Hours
	case hoursSincePublished < constants.Hours72:
		return constants.BonusRecent72Hours
	default:
		return 0
	}
}

// CalculateTopAuthors calcula los N autores con mayor suma acumulada de score.
// El cálculo se realiza completamente en memoria.
// Maneja empates ordenando por AuthorID ascendente.
func (s *ScoreService) CalculateTopAuthors(articles []*vo.PublishedArticleWithScore, limit int) []*vo.TopAuthor {
	// Agrupar artículos por autor y calcular scores
	authorScores := make(map[int64]*vo.TopAuthor)

	for _, article := range articles {
		score := s.CalculateArticleScore(
			article.WordCount(),
			article.AuthorPublishedCount(),
			article.PublishedAt(),
		)

		if existing, ok := authorScores[article.AuthorID()]; ok {
			authorScores[article.AuthorID()] = vo.NewTopAuthor(
				existing.AuthorID(),
				existing.AuthorName(),
				existing.TotalScore()+score,
				existing.PublishedCount()+1,
			)
		} else {
			authorScores[article.AuthorID()] = vo.NewTopAuthor(
				article.AuthorID(),
				article.AuthorName(),
				score,
				1,
			)
		}
	}

	// Convertir map a slice
	topAuthors := make([]*vo.TopAuthor, 0, len(authorScores))
	for _, author := range authorScores {
		topAuthors = append(topAuthors, author)
	}

	// Ordenar por total_score descendente, id ascendente en caso de empate
	sort.Slice(topAuthors, func(i, j int) bool {
		if topAuthors[i].TotalScore() != topAuthors[j].TotalScore() {
			return topAuthors[i].TotalScore() > topAuthors[j].TotalScore()
		}
		return topAuthors[i].AuthorID() < topAuthors[j].AuthorID()
	})

	// Aplicar límite (maneja cuando limit > len)
	if limit > len(topAuthors) {
		limit = len(topAuthors)
	}

	return topAuthors[:limit]
}
