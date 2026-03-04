package get_top_authors

type AuthorItem struct {
	ID                int64
	Name              string
	TotalScore        float64
	PublishedArticles int
}

type Response struct {
	Authors []AuthorItem
}
