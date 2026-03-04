package valueobjects

type TopAuthor struct {
	authorID       int64
	authorName     string
	totalScore     float64
	publishedCount int
}

func NewTopAuthor(authorID int64, authorName string, totalScore float64, publishedCount int) *TopAuthor {
	return &TopAuthor{
		authorID:       authorID,
		authorName:     authorName,
		totalScore:     totalScore,
		publishedCount: publishedCount,
	}
}

func (t *TopAuthor) AuthorID() int64 {
	return t.authorID
}

func (t *TopAuthor) AuthorName() string {
	return t.authorName
}

func (t *TopAuthor) TotalScore() float64 {
	return t.totalScore
}

func (t *TopAuthor) PublishedCount() int {
	return t.publishedCount
}

func (t *TopAuthor) AverageScore() float64 {
	if t.publishedCount == 0 {
		return 0
	}
	return t.totalScore / float64(t.publishedCount)
}
