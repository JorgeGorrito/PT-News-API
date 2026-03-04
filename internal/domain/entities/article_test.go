package entities

import (
	"strings"
	"testing"

	constants "github.com/JorgeGorrito/PT-News-API/internal/domain/constants"
	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
)

// Test de validación antes de publicar: Mínimo 120 palabras
func TestArticle_Publish_MinimumWords(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expectError bool
		errorType   domerrs.ErrorType
	}{
		{
			name:        "Artículo con menos de 120 palabras debe fallar",
			body:        strings.Repeat("alpha beta gamma delta epsilon zeta eta theta iota kappa ", 10), // 100 palabras
			expectError: true,
			errorType:   domerrs.MinWordsToPublishError,
		},
		{
			name:        "Artículo con exactamente 120 palabras debe pasar",
			body:        strings.Repeat("alpha beta gamma delta epsilon zeta eta theta iota kappa ", 12), // 120 palabras
			expectError: false,
		},
		{
			name:        "Artículo con más de 120 palabras debe pasar",
			body:        strings.Repeat("alpha beta gamma delta epsilon zeta eta theta iota kappa ", 15), // 150 palabras
			expectError: false,
		},
		{
			name:        "Artículo con 50 palabras debe fallar",
			body:        strings.Repeat("word1 word2 word3 word4 word5 ", 10), // 50 palabras
			expectError: true,
			errorType:   domerrs.MinWordsToPublishError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := NewArticle(1, "Test Article")
			err := article.SetBody(tt.body)
			if err != nil {
				t.Fatalf("SetBody failed: %v", err)
			}

			err = article.Publish()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if domErr, ok := err.(*domerrs.DomainError); ok {
					if domErr.Type() != tt.errorType {
						t.Errorf("Expected error type %d, got %d", tt.errorType, domErr.Type())
					}
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error but got: %v", err)
				}
			}
		})
	}
}

// Test de validación antes de publicar: Máximo 35% de palabras repetidas
func TestArticle_Publish_WordRepetitionValidation(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expectError bool
		description string
	}{
		{
			name: "Artículo con más de 35% de repetición debe fallar",
			body: func() string {
				// 120 palabras: 50 veces "repetida" (41.6%) + 70 palabras únicas
				repeated := strings.Repeat("repetida ", 50)
				unique := strings.Repeat("unica1 unica2 unica3 unica4 unica5 unica6 unica7 ", 10)
				return repeated + unique
			}(),
			expectError: true,
			description: "41.6% de repetición",
		},
		{
			name: "Artículo con exactamente 35% de repetición debe pasar",
			body: func() string {
				// 120 palabras: 42 veces "palabra" (35%) + 78 palabras únicas
				repeated := strings.Repeat("palabra ", 42)
				unique := strings.Repeat("uno dos tres cuatro cinco seis ", 13)
				return repeated + unique
			}(),
			expectError: false,
			description: "35% de repetición exacta",
		},
		{
			name: "Artículo con 20% de repetición debe pasar",
			body: func() string {
				// 150 palabras: 30 veces "común" (20%) + 120 palabras únicas
				repeated := strings.Repeat("común ", 30)
				unique := strings.Repeat("alpha beta gamma delta epsilon zeta ", 20)
				return repeated + unique
			}(),
			expectError: false,
			description: "20% de repetición",
		},
		{
			name: "Artículo con 50% de repetición debe fallar",
			body: func() string {
				// 120 palabras: 60 veces "test" (50%) + 60 palabras únicas
				repeated := strings.Repeat("test ", 60)
				unique := strings.Repeat("word1 word2 word3 word4 word5 word6 ", 10)
				return repeated + unique
			}(),
			expectError: true,
			description: "50% de repetición",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := NewArticle(1, "Test Article")
			err := article.SetBody(tt.body)
			if err != nil {
				t.Fatalf("SetBody failed: %v", err)
			}

			err = article.Publish()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s but got none", tt.description)
				} else if domErr, ok := err.(*domerrs.DomainError); ok {
					if domErr.Type() != domerrs.PercentageOfRepetitionError {
						t.Errorf("Expected PercentageOfRepetitionError, got %d", domErr.Type())
					}
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error for %s but got: %v", tt.description, err)
				}
			}
		})
	}
}

// Test de validación: artículo ya publicado no puede publicarse de nuevo
func TestArticle_Publish_AlreadyPublished(t *testing.T) {
	article := NewArticle(1, "Test Article")
	body := strings.Repeat("alpha beta gamma delta epsilon zeta eta theta iota kappa ", 15) // 150 palabras variadas
	err := article.SetBody(body)
	if err != nil {
		t.Fatalf("SetBody failed: %v", err)
	}

	// Primera publicación debe funcionar
	err = article.Publish()
	if err != nil {
		t.Fatalf("First publish should succeed: %v", err)
	}

	// Segunda publicación debe fallar
	err = article.Publish()
	if err == nil {
		t.Error("Expected error when publishing already published article")
	} else if domErr, ok := err.(*domerrs.DomainError); ok {
		if domErr.Type() != domerrs.ArticleAlreadyPublishedError {
			t.Errorf("Expected ArticleAlreadyPublishedError, got %d", domErr.Type())
		}
	}
}

// Test de conteo de palabras
func TestArticle_WordCount(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		expectedCount uint
	}{
		{
			name:          "Artículo con 10 palabras",
			body:          "uno dos tres cuatro cinco seis siete ocho nueve diez",
			expectedCount: 10,
		},
		{
			name:          "Artículo con 120 palabras",
			body:          strings.Repeat("palabra ", 120),
			expectedCount: 120,
		},
		{
			name:          "Artículo vacío",
			body:          "",
			expectedCount: 0,
		},
		{
			name:          "Artículo con múltiples espacios",
			body:          "palabra1    palabra2     palabra3",
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := NewArticle(1, "Test Article")
			err := article.SetBody(tt.body)
			if err != nil {
				t.Fatalf("SetBody failed: %v", err)
			}

			if article.WordCount() != tt.expectedCount {
				t.Errorf("Expected word count %d, got %d", tt.expectedCount, article.WordCount())
			}
		})
	}
}

// Test que verifica la validación completa antes de publicar
func TestArticle_Publish_CompleteValidation(t *testing.T) {
	t.Run("Artículo válido se publica correctamente", func(t *testing.T) {
		article := NewArticle(1, "Valid Article")
		// 120 palabras con buena diversidad (sin exceder 35% repetición)
		body := strings.Repeat("alpha beta gamma delta epsilon zeta eta theta iota kappa ", 12)
		err := article.SetBody(body)
		if err != nil {
			t.Fatalf("SetBody failed: %v", err)
		}

		if article.WordCount() < constants.MinWordsToPublish {
			t.Errorf("Expected at least %d words, got %d", constants.MinWordsToPublish, article.WordCount())
		}

		err = article.Publish()
		if err != nil {
			t.Errorf("Valid article should publish successfully: %v", err)
		}

		if !article.IsPublished() {
			t.Error("Article should be marked as published")
		}

		if article.PublishedAt() == nil {
			t.Error("PublishedAt should be set after publishing")
		}
	})
}
