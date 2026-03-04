package get_author_summary

type Response struct {
	ID                int64
	Name              string
	Email             string
	Biography         string
	PublishedArticles int
	DraftArticles     int
}
