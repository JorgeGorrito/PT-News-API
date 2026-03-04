package get_top_authors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
	"github.com/JorgeGorrito/PT-News-API/internal/domain/services"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
)

// mockArticleRepository implementa interfaces.ArticleRepository para testing.
type mockArticleRepository struct {
	allPublished []*vo.PublishedArticleWithScore
	err          error
}

func (m *mockArticleRepository) FindAllPublished(ctx context.Context) ([]*vo.PublishedArticleWithScore, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allPublished, nil
}

// Métodos restantes requeridos por la interfaz (no usados en estos tests)
func (m *mockArticleRepository) Save(ctx context.Context, article *entities.Article) error {
	return nil
}

func (m *mockArticleRepository) FindByID(ctx context.Context, id int64) (*entities.Article, error) {
	return nil, nil
}

func (m *mockArticleRepository) Update(ctx context.Context, article *entities.Article) error {
	return nil
}

func (m *mockArticleRepository) FindByAuthorIDPaginated(ctx context.Context, authorID int64, status *vo.ArticleStatus, page, perPage int) ([]*entities.Article, int, error) {
	return nil, 0, nil
}

func (m *mockArticleRepository) FindPublishedByAuthorID(ctx context.Context, authorID int64) ([]*entities.Article, error) {
	return nil, nil
}

func (m *mockArticleRepository) CountByAuthorID(ctx context.Context, authorID int64) (int, error) {
	return 0, nil
}

func (m *mockArticleRepository) CountPublishedByAuthorID(ctx context.Context, authorID int64) (int, error) {
	return 0, nil
}

func (m *mockArticleRepository) FindPublishedPaginated(ctx context.Context, page, perPage int, orderBy vo.ArticleOrderBy) ([]*vo.PublishedArticleWithScore, int, error) {
	return nil, 0, nil
}

// newTestArticle crea un PublishedArticleWithScore con fecha >72h para score predecible (sin bonus).
// score = (wordCount * 0.1) + (authorPublishedCount * 5)
func newTestArticle(articleID, authorID int64, authorName string, wordCount uint, authorPublishedCount int) *vo.PublishedArticleWithScore {
	return vo.NewPublishedArticleWithScore(
		articleID,
		authorID,
		authorName,
		"Título de prueba",
		"Cuerpo de prueba",
		wordCount,
		time.Now().Add(-100*time.Hour), // >72h → bonus_reciente = 0
		authorPublishedCount,
	)
}

// Test: Top N autores ordenados por score (cálculo en memoria con ScoreService real)
func TestHandler_GetTopAuthors_Success(t *testing.T) {
	// Scores predecibles (>72h, sin bonus):
	// Autor A: wordCount=1000, authorPublished=1 → 100+5 = 105
	// Autor B: wordCount=500,  authorPublished=1 → 50+5  = 55
	// Autor C: wordCount=200,  authorPublished=1 → 20+5  = 25
	articles := []*vo.PublishedArticleWithScore{
		newTestArticle(1, 1, "Autor A", 1000, 1),
		newTestArticle(2, 2, "Autor B", 500, 1),
		newTestArticle(3, 3, "Autor C", 200, 1),
	}

	tests := []struct {
		name          string
		limit         int
		expectedCount int
	}{
		{name: "Top 3 autores", limit: 3, expectedCount: 3},
		{name: "Top 5 autores con solo 3 disponibles (n > total)", limit: 5, expectedCount: 3},
		{name: "Top 1 autor", limit: 1, expectedCount: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockArticleRepository{allPublished: articles}
			handler := NewHandler(mockRepo, services.NewScoreService())

			response, err := handler.Handle(context.Background(), Query{Limit: tt.limit})
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}
			if len(response.Authors) != tt.expectedCount {
				t.Errorf("Expected %d authors, got %d", tt.expectedCount, len(response.Authors))
			}
			// Verificar orden descendente por score
			for i := 0; i < len(response.Authors)-1; i++ {
				if response.Authors[i].TotalScore < response.Authors[i+1].TotalScore {
					t.Errorf("Authors not sorted correctly: %.2f < %.2f",
						response.Authors[i].TotalScore, response.Authors[i+1].TotalScore)
				}
			}
		})
	}
}

// Test: Sin artículos publicados
func TestHandler_GetTopAuthors_NoArticles(t *testing.T) {
	mockRepo := &mockArticleRepository{allPublished: []*vo.PublishedArticleWithScore{}}
	handler := NewHandler(mockRepo, services.NewScoreService())

	response, err := handler.Handle(context.Background(), Query{Limit: 10})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(response.Authors) != 0 {
		t.Errorf("Expected 0 authors, got %d", len(response.Authors))
	}
}

// Test: Validación de límite <= 0
func TestHandler_GetTopAuthors_InvalidLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
	}{
		{name: "Límite 0 debe fallar", limit: 0},
		{name: "Límite negativo debe fallar", limit: -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockArticleRepository{}
			handler := NewHandler(mockRepo, services.NewScoreService())

			response, err := handler.Handle(context.Background(), Query{Limit: tt.limit})
			if err == nil {
				t.Error("Expected error for invalid limit, got none")
			}
			if response != nil {
				t.Error("Expected nil response for invalid limit")
			}
			if domErr, ok := err.(*domerrs.DomainError); ok {
				if domErr.Type() != domerrs.GeneralError {
					t.Errorf("Expected GeneralError type, got %d", domErr.Type())
				}
			}
		})
	}
}

// Test: Manejo de errores del repositorio
func TestHandler_GetTopAuthors_RepositoryError(t *testing.T) {
	expectedError := errors.New("database connection error")
	mockRepo := &mockArticleRepository{err: expectedError}
	handler := NewHandler(mockRepo, services.NewScoreService())

	response, err := handler.Handle(context.Background(), Query{Limit: 5})
	if err == nil {
		t.Fatal("Expected error from repository, got none")
	}
	if response != nil {
		t.Error("Expected nil response when repository fails")
	}
	if err != expectedError {
		t.Errorf("Expected repository error, got: %v", err)
	}
}

// Test: Verificar cálculo y agregación correcta de scores
// Scores esperados (>72h, sin bonus):
// Autor A (id=1): wordCount=1400, authorPublished=1 → 140+5 = 145.0
// Autor B (id=2): wordCount=1100, authorPublished=1 → 110+5 = 115.0
// Autor C (id=3): wordCount=900,  authorPublished=1 → 90+5  = 95.0
func TestHandler_GetTopAuthors_ScoreCalculation(t *testing.T) {
	articles := []*vo.PublishedArticleWithScore{
		newTestArticle(1, 1, "Autor A", 1400, 1),
		newTestArticle(2, 2, "Autor B", 1100, 1),
		newTestArticle(3, 3, "Autor C", 900, 1),
	}

	mockRepo := &mockArticleRepository{allPublished: articles}
	handler := NewHandler(mockRepo, services.NewScoreService())

	response, err := handler.Handle(context.Background(), Query{Limit: 3})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(response.Authors) != 3 {
		t.Errorf("Expected 3 authors, got %d", len(response.Authors))
	}

	expectedScores := []float64{145.0, 115.0, 95.0}
	expectedIDs := []int64{1, 2, 3}

	for i, author := range response.Authors {
		if author.TotalScore != expectedScores[i] {
			t.Errorf("Author %d: expected score %.1f, got %.2f", i, expectedScores[i], author.TotalScore)
		}
		if author.ID != expectedIDs[i] {
			t.Errorf("Author %d: expected ID %d, got %d", i, expectedIDs[i], author.ID)
		}
	}
}

// Test: Manejo de empates en scores — el de menor ID aparece primero
// Autor 1 (id=1): wordCount=500, authorPublished=2 → 50+10 = 60.0
// Autor 2 (id=2): wordCount=500, authorPublished=2 → 50+10 = 60.0 (empate)
// Autor 3 (id=3): wordCount=200, authorPublished=1 → 20+5  = 25.0
func TestHandler_GetTopAuthors_TieScores(t *testing.T) {
	articles := []*vo.PublishedArticleWithScore{
		newTestArticle(1, 1, "Autor A", 500, 2),
		newTestArticle(2, 2, "Autor B", 500, 2),
		newTestArticle(3, 3, "Autor C", 200, 1),
	}

	mockRepo := &mockArticleRepository{allPublished: articles}
	handler := NewHandler(mockRepo, services.NewScoreService())

	response, err := handler.Handle(context.Background(), Query{Limit: 3})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(response.Authors) != 3 {
		t.Errorf("Expected 3 authors (including ties), got %d", len(response.Authors))
	}
	// Los dos primeros deben tener el mismo score
	if response.Authors[0].TotalScore != response.Authors[1].TotalScore {
		t.Errorf("Expected tied scores but got %.2f and %.2f",
			response.Authors[0].TotalScore, response.Authors[1].TotalScore)
	}
	// En caso de empate, el de menor ID va primero
	if response.Authors[0].ID > response.Authors[1].ID {
		t.Errorf("Tie: author with lower ID should come first (got IDs %d, %d)",
			response.Authors[0].ID, response.Authors[1].ID)
	}
}
