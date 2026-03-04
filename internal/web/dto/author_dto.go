package dto

// CreateAuthorRequest represents the request to create an author
type CreateAuthorRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=255"`
	Email     string `json:"email" binding:"required,email"`
	Biography string `json:"biography" binding:"max=1000"`
}

// CreateAuthorResponse represents the response after creating an author
type CreateAuthorResponse struct {
	ID int64 `json:"id"`
}

// AuthorSummaryResponse represents the author summary with statistics
type AuthorSummaryResponse struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Email             string  `json:"email"`
	Biography         string  `json:"biography"`
	PublishedArticles int     `json:"published_articles"`
	TotalScore        float64 `json:"total_score"`
}

// TopAuthorResponse represents a top author by score
type TopAuthorResponse struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	PublishedArticles int     `json:"published_articles"`
	TotalScore        float64 `json:"total_score"`
}
