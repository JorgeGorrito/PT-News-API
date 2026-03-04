package entities

import (
	"fmt"
	"strings"
	"time"

	constants "github.com/JorgeGorrito/PT-News-API/internal/domain/constants"
	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
	"github.com/JorgeGorrito/PT-News-API/internal/domain/services"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
)

type Article struct {
	BaseEntity[int64]
	authorID    int64
	title       string
	body        string
	status      vo.ArticleStatus
	wordCount   uint
	createdAt   time.Time
	publishedAt *time.Time
}

func NewArticle(authorID int64, title string) *Article {
	return &Article{
		BaseEntity:  BaseEntity[int64]{id: 0},
		authorID:    authorID,
		title:       title,
		body:        "",
		status:      vo.Draft,
		wordCount:   0,
		createdAt:   time.Now().UTC(),
		publishedAt: nil,
	}
}

func (a *Article) checkWordRepetition() error {
	wordsFormatted := strings.Fields(strings.ToLower(a.body))
	totalWords := a.wordCount

	wordFrequency := make(map[string]int)
	maxFrecuency := 0

	for _, word := range wordsFormatted {
		wordFrequency[word]++

		if wordFrequency[word] > maxFrecuency {
			maxFrecuency = wordFrequency[word]
		}
	}

	percentaje := (float64(maxFrecuency) / float64(totalWords)) * 100

	if percentaje > constants.MaxAllowedRepetitionPercentage {
		return domerrs.NewDomainError(
			fmt.Sprintf("El artículo tiene una repetición de palabras superior al %.2f%%", constants.MaxAllowedRepetitionPercentage),
			domerrs.PercentageOfRepetitionError,
		)
	}

	return nil
}

func (a *Article) canBePublished() error {
	if a.status == vo.Published {
		return domerrs.NewDomainError("El artículo ya está publicado", domerrs.ArticleAlreadyPublishedError)
	}

	if a.wordCount < constants.MinWordsToPublish {
		return domerrs.NewDomainError(
			fmt.Sprintf("El artículo debe tener al menos %d palabras para ser publicado", constants.MinWordsToPublish),
			domerrs.MinWordsToPublishError,
		)
	}

	return a.checkWordRepetition()
}

func (a *Article) Publish() error {
	if err := a.canBePublished(); err != nil {
		return err
	}

	a.status = vo.Published

	now := time.Now().UTC()
	a.publishedAt = &now
	return nil
}

func (a *Article) SetTitle(title string) error {
	if title == "" {
		return domerrs.NewDomainError("El titulo no debe estar vacío", domerrs.EmptyArticleTitleError)
	}

	a.title = title
	return nil
}

func (a *Article) Title() string {
	return a.title
}

func (a *Article) updateWordCount() {
	words := strings.Fields(a.body)
	a.wordCount = uint(len(words))
}

func (a *Article) SetBody(body string) error {
	a.body = body
	a.updateWordCount()

	return nil
}

func (a *Article) Body() string {
	return a.body
}

func (a *Article) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Article) PublishedAt() *time.Time {
	return a.publishedAt
}

func (a *Article) Status() vo.ArticleStatus {
	return a.status
}

func (a *Article) WordCount() uint {
	return a.wordCount
}

func (a *Article) AuthorID() int64 {
	return a.authorID
}

func (a *Article) IsPublished() bool {
	return a.status == vo.Published
}

func (a *Article) IsDraft() bool {
	return a.status == vo.Draft
}

// CalculateScore calcula el score de relevancia del artículo delegando la lógica
// al servicio de dominio. Recibe como parámetro externo el número de artículos
// publicados del autor, ya que ese dato no pertenece a esta entidad.
func (a *Article) CalculateScore(svc services.IScoreService, authorPublishedCount int) float64 {
	if a.publishedAt == nil {
		return 0
	}
	return svc.CalculateArticleScore(a.wordCount, authorPublishedCount, *a.publishedAt)
}
