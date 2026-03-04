package valueobjects

import "time"

type AuthorSummary struct {
	authorID          int64
	totalArticles     int
	totalPublished    int
	totalScore        float64
	lastPublicationAt *time.Time
}

func NewAuthorSummary(authorID int64, totalArticles, totalPublished int, totalScore float64, lastPublicationAt *time.Time) *AuthorSummary {
	return &AuthorSummary{
		authorID:          authorID,
		totalArticles:     totalArticles,
		totalPublished:    totalPublished,
		totalScore:        totalScore,
		lastPublicationAt: lastPublicationAt,
	}
}

func (s *AuthorSummary) AuthorID() int64 {
	return s.authorID
}

func (s *AuthorSummary) TotalArticles() int {
	return s.totalArticles
}

func (s *AuthorSummary) TotalPublished() int {
	return s.totalPublished
}

func (s *AuthorSummary) TotalScore() float64 {
	return s.totalScore
}

func (s *AuthorSummary) LastPublicationAt() *time.Time {
	return s.lastPublicationAt
}

func (s *AuthorSummary) HasPublications() bool {
	return s.totalPublished > 0
}
