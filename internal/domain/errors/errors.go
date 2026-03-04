package domerrs

type ErrorType uint

const (
	GeneralError ErrorType = 0

	AuthorError                 = 100
	EmptyAuthorNameError        = 101
	InvalidAuthorNameError      = 102
	EmptyAuthorEmailError       = 103
	InvalidAuthorEmailError     = 104
	InvalidAuthorBiographyError = 105

	ArticleError                 = 200
	EmptyArticleTitleError       = 201
	MinWordsToPublishError       = 202
	PercentageOfRepetitionError  = 203
	ArticleAlreadyPublishedError = 204
)
