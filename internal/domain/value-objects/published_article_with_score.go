package valueobjects

import "time"

type PublishedArticleWithScore struct {
	articleID   int64
	authorID    int64
	authorName  string
	title       string
	body        string
	wordCount   uint
	publishedAt time.Time
	score       float64
}

func NewPublishedArticleWithScore(
	articleID, authorID int64,
	authorName, title, body string,
	wordCount uint,
	publishedAt time.Time,
	score float64,
) *PublishedArticleWithScore {
	return &PublishedArticleWithScore{
		articleID:   articleID,
		authorID:    authorID,
		authorName:  authorName,
		title:       title,
		body:        body,
		wordCount:   wordCount,
		publishedAt: publishedAt,
		score:       score,
	}
}

func (p *PublishedArticleWithScore) ArticleID() int64 {
	return p.articleID
}

func (p *PublishedArticleWithScore) AuthorID() int64 {
	return p.authorID
}

func (p *PublishedArticleWithScore) AuthorName() string {
	return p.authorName
}

func (p *PublishedArticleWithScore) Title() string {
	return p.title
}

func (p *PublishedArticleWithScore) Body() string {
	return p.body
}

func (p *PublishedArticleWithScore) WordCount() uint {
	return p.wordCount
}

func (p *PublishedArticleWithScore) PublishedAt() time.Time {
	return p.publishedAt
}

func (p *PublishedArticleWithScore) Score() float64 {
	return p.score
}
