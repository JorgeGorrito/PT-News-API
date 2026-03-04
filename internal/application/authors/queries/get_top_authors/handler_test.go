package get_top_authors

import (
	"context"
	"errors"
	"testing"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
)

// Mock del repositorio de artículos
type mockArticleRepository struct {
	topAuthors []*vo.TopAuthor
	err        error
}

func (m *mockArticleRepository) GetTopAuthorsByScore(ctx context.Context, limit int) ([]*vo.TopAuthor, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.topAuthors, nil
}

// Implementación de otros métodos requeridos por la interfaz (no usados en este test)
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

// Test: Top N autores ordenados por score
func TestHandler_GetTopAuthors_Success(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		mockTopAuthors []*vo.TopAuthor
		expectedCount  int
	}{
		{
			name:  "Top 3 autores",
			limit: 3,
			mockTopAuthors: []*vo.TopAuthor{
				vo.NewTopAuthor(1, "Autor A", 500.5, 10),
				vo.NewTopAuthor(2, "Autor B", 450.0, 9),
				vo.NewTopAuthor(3, "Autor C", 400.0, 8),
			},
			expectedCount: 3,
		},
		{
			name:  "Top 5 autores con solo 3 disponibles",
			limit: 5,
			mockTopAuthors: []*vo.TopAuthor{
				vo.NewTopAuthor(1, "Autor A", 500.5, 10),
				vo.NewTopAuthor(2, "Autor B", 450.0, 9),
				vo.NewTopAuthor(3, "Autor C", 400.0, 8),
			},
			expectedCount: 3,
		},
		{
			name:           "Top 10 autores sin datos",
			limit:          10,
			mockTopAuthors: []*vo.TopAuthor{},
			expectedCount:  0,
		},
		{
			name:  "Top 1 autor",
			limit: 1,
			mockTopAuthors: []*vo.TopAuthor{
				vo.NewTopAuthor(1, "Autor A", 800.0, 15),
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockArticleRepository{
				topAuthors: tt.mockTopAuthors,
			}

			handler := NewHandler(mockRepo)
			query := Query{Limit: tt.limit}

			response, err := handler.Handle(context.Background(), query)

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if len(response.Authors) != tt.expectedCount {
				t.Errorf("Expected %d authors, got %d", tt.expectedCount, len(response.Authors))
			}

			// Verificar que los autores están en el orden correcto (score descendente)
			for i := 0; i < len(response.Authors)-1; i++ {
				if response.Authors[i].TotalScore < response.Authors[i+1].TotalScore {
					t.Errorf("Authors not sorted correctly by score: %f should be >= %f",
						response.Authors[i].TotalScore, response.Authors[i+1].TotalScore)
				}
			}
		})
	}
}

// Test: Validation de límite <= 0
func TestHandler_GetTopAuthors_InvalidLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
	}{
		{
			name:  "Límite 0 debe fallar",
			limit: 0,
		},
		{
			name:  "Límite negativo debe fallar",
			limit: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockArticleRepository{}
			handler := NewHandler(mockRepo)
			query := Query{Limit: tt.limit}

			response, err := handler.Handle(context.Background(), query)

			if err == nil {
				t.Error("Expected error for invalid limit, got none")
			}

			if response != nil {
				t.Error("Expected nil response for invalid limit")
			}

			// Verificar que es un error de dominio con el tipo correcto
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
	mockRepo := &mockArticleRepository{
		err: expectedError,
	}

	handler := NewHandler(mockRepo)
	query := Query{Limit: 5}

	response, err := handler.Handle(context.Background(), query)

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
func TestHandler_GetTopAuthors_ScoreCalculation(t *testing.T) {
	// El repositorio debería retornar solo los top 3 (simulamos que el LIMIT se aplica en BD)
	mockTopAuthors := []*vo.TopAuthor{
		vo.NewTopAuthor(1, "Autor A", 750.5, 15), // Más score, más publicaciones
		vo.NewTopAuthor(2, "Autor B", 600.0, 12),
		vo.NewTopAuthor(3, "Autor C", 550.0, 11),
	}

	mockRepo := &mockArticleRepository{
		topAuthors: mockTopAuthors,
	}

	handler := NewHandler(mockRepo)
	query := Query{Limit: 3}

	response, err := handler.Handle(context.Background(), query)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verificar que retorna exactamente 3
	if len(response.Authors) != 3 {
		t.Errorf("Expected 3 authors, got %d", len(response.Authors))
	}

	// Verificar que el primer autor tiene el score más alto
	if response.Authors[0].TotalScore != 750.5 {
		t.Errorf("Expected top author score 750.5, got %f", response.Authors[0].TotalScore)
	}

	// Verificar que todos los datos se mapean correctamente
	for i, author := range response.Authors {
		expectedAuthor := mockTopAuthors[i]
		if author.ID != expectedAuthor.AuthorID() {
			t.Errorf("Author %d: expected ID %d, got %d", i, expectedAuthor.AuthorID(), author.ID)
		}
		if author.Name != expectedAuthor.AuthorName() {
			t.Errorf("Author %d: expected name %s, got %s", i, expectedAuthor.AuthorName(), author.Name)
		}
		if author.TotalScore != expectedAuthor.TotalScore() {
			t.Errorf("Author %d: expected score %f, got %f", i, expectedAuthor.TotalScore(), author.TotalScore)
		}
		if author.PublishedArticles != expectedAuthor.PublishedCount() {
			t.Errorf("Author %d: expected %d published articles, got %d", i, expectedAuthor.PublishedCount(), author.PublishedArticles)
		}
	}
}

// Test: Manejo de empates en scores
func TestHandler_GetTopAuthors_TieScores(t *testing.T) {
	// Autores con scores empatados
	mockTopAuthors := []*vo.TopAuthor{
		vo.NewTopAuthor(1, "Autor A", 500.0, 10),
		vo.NewTopAuthor(2, "Autor B", 500.0, 10), // Mismo score
		vo.NewTopAuthor(3, "Autor C", 400.0, 8),
	}

	mockRepo := &mockArticleRepository{
		topAuthors: mockTopAuthors,
	}

	handler := NewHandler(mockRepo)
	query := Query{Limit: 3}

	response, err := handler.Handle(context.Background(), query)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verificar que retorna todos, incluyendo los empatados
	if len(response.Authors) != 3 {
		t.Errorf("Expected 3 authors (including ties), got %d", len(response.Authors))
	}

	// Verificar que los dos primeros tienen el mismo score
	if response.Authors[0].TotalScore != response.Authors[1].TotalScore {
		t.Errorf("Expected tied scores: %f and %f should be equal",
			response.Authors[0].TotalScore, response.Authors[1].TotalScore)
	}
}
