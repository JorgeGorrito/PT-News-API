package get_article_by_id

type Query struct {
	ArticleID    int64
	IncludeScore bool
}
