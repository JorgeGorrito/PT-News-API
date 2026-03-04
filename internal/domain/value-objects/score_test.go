package valueobjects

import (
	"math"
	"testing"
	"time"

	constants "github.com/JorgeGorrito/PT-News-API/internal/domain/constants"
)

// CalculateArticleScore calcula el score de un artículo según la fórmula especificada
// Esta función replica la lógica de cálculo que se implementa en SQL
func CalculateArticleScore(wordCount uint, authorPublishedCount int, publishedAt time.Time) float64 {
	// Base score: (word_count * 0.1) + (author_published_articles * 5)
	baseScore := (float64(wordCount) * constants.ScorePerWord) + (float64(authorPublishedCount) * constants.ScorePerAuthorArticle)

	// Calcular bonus por recencia
	hoursSincePublished := time.Since(publishedAt).Hours()
	var recencyBonus float64

	if hoursSincePublished < constants.Hours24 {
		recencyBonus = constants.BonusRecent24Hours
	} else if hoursSincePublished < constants.Hours72 {
		recencyBonus = constants.BonusRecent72Hours
	} else {
		recencyBonus = 0
	}

	return baseScore + recencyBonus
}

// Test: Cálculo de score con diferentes conteos de palabras
func TestCalculateArticleScore_WordCount(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name                 string
		wordCount            uint
		authorPublishedCount int
		publishedAt          time.Time
		expectedScore        float64
	}{
		{
			name:                 "Score con 200 palabras, 5 artículos publicados, sin bonus",
			wordCount:            200,
			authorPublishedCount: 5,
			publishedAt:          now.Add(-100 * time.Hour), // Más de 72 horas
			expectedScore:        (200 * 0.1) + (5 * 5) + 0, // 20 + 25 + 0 = 45
		},
		{
			name:                 "Score con 1000 palabras, 10 artículos, sin bonus",
			wordCount:            1000,
			authorPublishedCount: 10,
			publishedAt:          now.Add(-80 * time.Hour),
			expectedScore:        (1000 * 0.1) + (10 * 5) + 0, // 100 + 50 + 0 = 150
		},
		{
			name:                 "Score con 120 palabras (mínimo), 1 artículo, sin bonus",
			wordCount:            120,
			authorPublishedCount: 1,
			publishedAt:          now.Add(-100 * time.Hour),
			expectedScore:        (120 * 0.1) + (1 * 5) + 0, // 12 + 5 + 0 = 17
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateArticleScore(tt.wordCount, tt.authorPublishedCount, tt.publishedAt)

			if math.Abs(score-tt.expectedScore) > 0.01 {
				t.Errorf("Expected score %.2f, got %.2f", tt.expectedScore, score)
			}
		})
	}
}

// Test: Bonus de recencia menor a 24 horas
func TestCalculateArticleScore_RecencyBonus24Hours(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		publishedAt   time.Time
		expectedBonus float64
	}{
		{
			name:          "Publicado hace 1 hora - bonus 50",
			publishedAt:   now.Add(-1 * time.Hour),
			expectedBonus: 50.0,
		},
		{
			name:          "Publicado hace 12 horas - bonus 50",
			publishedAt:   now.Add(-12 * time.Hour),
			expectedBonus: 50.0,
		},
		{
			name:          "Publicado hace 23 horas - bonus 50",
			publishedAt:   now.Add(-23 * time.Hour),
			expectedBonus: 50.0,
		},
		{
			name:          "Publicado hace 23.5 horas - bonus 50",
			publishedAt:   now.Add(-23*time.Hour - 30*time.Minute),
			expectedBonus: 50.0,
		},
	}

	wordCount := uint(200)
	authorPublishedCount := 5
	baseScore := (float64(wordCount) * 0.1) + (float64(authorPublishedCount) * 5.0) // 45

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateArticleScore(wordCount, authorPublishedCount, tt.publishedAt)
			expectedScore := baseScore + tt.expectedBonus

			if math.Abs(score-expectedScore) > 0.01 {
				t.Errorf("Expected score %.2f (base %.2f + bonus %.2f), got %.2f",
					expectedScore, baseScore, tt.expectedBonus, score)
			}
		})
	}
}

// Test: Bonus de recencia entre 24 y 72 horas
func TestCalculateArticleScore_RecencyBonus72Hours(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		publishedAt   time.Time
		expectedBonus float64
	}{
		{
			name:          "Publicado hace 25 horas - bonus 20",
			publishedAt:   now.Add(-25 * time.Hour),
			expectedBonus: 20.0,
		},
		{
			name:          "Publicado hace 48 horas - bonus 20",
			publishedAt:   now.Add(-48 * time.Hour),
			expectedBonus: 20.0,
		},
		{
			name:          "Publicado hace 71 horas - bonus 20",
			publishedAt:   now.Add(-71 * time.Hour),
			expectedBonus: 20.0,
		},
		{
			name:          "Publicado hace 71.5 horas - bonus 20",
			publishedAt:   now.Add(-71*time.Hour - 30*time.Minute),
			expectedBonus: 20.0,
		},
	}

	wordCount := uint(200)
	authorPublishedCount := 5
	baseScore := (float64(wordCount) * 0.1) + (float64(authorPublishedCount) * 5.0) // 45

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateArticleScore(wordCount, authorPublishedCount, tt.publishedAt)
			expectedScore := baseScore + tt.expectedBonus

			if math.Abs(score-expectedScore) > 0.01 {
				t.Errorf("Expected score %.2f (base %.2f + bonus %.2f), got %.2f",
					expectedScore, baseScore, tt.expectedBonus, score)
			}
		})
	}
}

// Test: Sin bonus de recencia después de 72 horas
func TestCalculateArticleScore_NoRecencyBonus(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		publishedAt time.Time
	}{
		{
			name:        "Publicado hace 73 horas - sin bonus",
			publishedAt: now.Add(-73 * time.Hour),
		},
		{
			name:        "Publicado hace 100 horas - sin bonus",
			publishedAt: now.Add(-100 * time.Hour),
		},
		{
			name:        "Publicado hace 1 semana - sin bonus",
			publishedAt: now.Add(-168 * time.Hour),
		},
		{
			name:        "Publicado hace 1 mes - sin bonus",
			publishedAt: now.Add(-720 * time.Hour),
		},
	}

	wordCount := uint(200)
	authorPublishedCount := 5
	expectedScore := (float64(wordCount) * 0.1) + (float64(authorPublishedCount) * 5.0) // 45

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateArticleScore(wordCount, authorPublishedCount, tt.publishedAt)

			if math.Abs(score-expectedScore) > 0.01 {
				t.Errorf("Expected score %.2f (no bonus), got %.2f", expectedScore, score)
			}
		})
	}
}

// Test: Fórmula completa con diferentes combinaciones
func TestCalculateArticleScore_CompleteFormula(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name                 string
		wordCount            uint
		authorPublishedCount int
		publishedAt          time.Time
		description          string
		expectedScore        float64
	}{
		{
			name:                 "Artículo reciente de autor prolífico",
			wordCount:            500,
			authorPublishedCount: 20,
			publishedAt:          now.Add(-5 * time.Hour),
			description:          "500 palabras, 20 artículos, <24h",
			expectedScore:        (500 * 0.1) + (20 * 5) + 50, // 50 + 100 + 50 = 200
		},
		{
			name:                 "Artículo medianamente reciente",
			wordCount:            300,
			authorPublishedCount: 10,
			publishedAt:          now.Add(-30 * time.Hour),
			description:          "300 palabras, 10 artículos, 24-72h",
			expectedScore:        (300 * 0.1) + (10 * 5) + 20, // 30 + 50 + 20 = 100
		},
		{
			name:                 "Artículo antiguo de autor nuevo",
			wordCount:            120,
			authorPublishedCount: 1,
			publishedAt:          now.Add(-200 * time.Hour),
			description:          "120 palabras, 1 artículo, >72h",
			expectedScore:        (120 * 0.1) + (1 * 5) + 0, // 12 + 5 + 0 = 17
		},
		{
			name:                 "Artículo largo y reciente",
			wordCount:            2000,
			authorPublishedCount: 5,
			publishedAt:          now.Add(-10 * time.Hour),
			description:          "2000 palabras, 5 artículos, <24h",
			expectedScore:        (2000 * 0.1) + (5 * 5) + 50, // 200 + 25 + 50 = 275
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateArticleScore(tt.wordCount, tt.authorPublishedCount, tt.publishedAt)

			if math.Abs(score-tt.expectedScore) > 0.01 {
				t.Errorf("%s: Expected score %.2f, got %.2f", tt.description, tt.expectedScore, score)
			}
		})
	}
}

// Test: Verificar que las constantes están correctamente definidas
func TestScoreConstants(t *testing.T) {
	if constants.ScorePerWord != 0.1 {
		t.Errorf("Expected ScorePerWord = 0.1, got %f", constants.ScorePerWord)
	}

	if constants.ScorePerAuthorArticle != 5.0 {
		t.Errorf("Expected ScorePerAuthorArticle = 5.0, got %f", constants.ScorePerAuthorArticle)
	}

	if constants.BonusRecent24Hours != 50.0 {
		t.Errorf("Expected BonusRecent24Hours = 50.0, got %f", constants.BonusRecent24Hours)
	}

	if constants.BonusRecent72Hours != 20.0 {
		t.Errorf("Expected BonusRecent72Hours = 20.0, got %f", constants.BonusRecent72Hours)
	}

	if constants.Hours24 != 24.0 {
		t.Errorf("Expected Hours24 = 24.0, got %f", constants.Hours24)
	}

	if constants.Hours72 != 72.0 {
		t.Errorf("Expected Hours72 = 72.0, got %f", constants.Hours72)
	}
}
