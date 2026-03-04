package valueobjects

import "time"

// PublishedArticleWithScore es un Value Object de lectura (read model) que
// proyecta los datos de un artículo publicado junto con el contexto del autor
// necesario para calcular el score dinámicamente en la capa de dominio.
type PublishedArticleWithScore struct {
	articleID            int64
	authorID             int64
	authorName           string
	title                string
	body                 string
	wordCount            uint
	publishedAt          time.Time
	authorPublishedCount int
}

func NewPublishedArticleWithScore(
	articleID, authorID int64,
	authorName, title, body string,
	wordCount uint,
	publishedAt time.Time,
	authorPublishedCount int,
) *PublishedArticleWithScore {
	return &PublishedArticleWithScore{
		articleID:            articleID,
		authorID:             authorID,
		authorName:           authorName,
		title:                title,
		body:                 body,
		wordCount:            wordCount,
		publishedAt:          publishedAt,
		authorPublishedCount: authorPublishedCount,
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

// AuthorPublishedCount retorna la cantidad de artículos publicados del autor,
// dato necesario para calcular el score de relevancia dinámicamente.
func (p *PublishedArticleWithScore) AuthorPublishedCount() int {
	return p.authorPublishedCount
}
