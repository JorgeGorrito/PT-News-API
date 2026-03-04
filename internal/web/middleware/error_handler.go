package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
)

// ErrorHandler middleware handles errors set by handlers using c.Error()
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last().Err

		// Map error to HTTP status code and response
		statusCode, response := mapErrorToHTTPResponse(err)

		// Prevent writing headers twice
		if !c.Writer.Written() {
			c.JSON(statusCode, response)
		}
	}
}

func mapErrorToHTTPResponse(err error) (int, gin.H) {
	// Check if it's a domain error
	var domainErr *domerrs.DomainError
	if errors.As(err, &domainErr) {
		return mapDomainError(domainErr)
	}

	// Default to internal server error
	return http.StatusInternalServerError, gin.H{
		"error": "Algo salió mal al procesar la solicitud",
	}
}

func mapDomainError(err *domerrs.DomainError) (int, gin.H) {
	errorType := err.Type()

	// Not found errors (404)
	if errorType == domerrs.NotFoundError {
		return http.StatusNotFound, gin.H{
			"error": err.Error(),
		}
	}

	// Validation errors (400)
	validationErrors := []domerrs.ErrorType{
		domerrs.EmptyAuthorNameError,
		domerrs.InvalidAuthorNameError,
		domerrs.EmptyAuthorEmailError,
		domerrs.InvalidAuthorEmailError,
		domerrs.InvalidAuthorBiographyError,
		domerrs.EmptyArticleTitleError,
		domerrs.MinWordsToPublishError,
		domerrs.PercentageOfRepetitionError,
		domerrs.InvalidArticleStatusError,
		domerrs.InvalidArticleOrderByError,
	}

	for _, valErr := range validationErrors {
		if errorType == valErr {
			return http.StatusBadRequest, gin.H{
				"error": err.Error(),
			}
		}
	}

	// Conflict errors (409)
	if errorType == domerrs.ArticleAlreadyPublishedError {
		return http.StatusConflict, gin.H{
			"error": err.Error(),
		}
	}

	// Default to bad request for domain errors
	return http.StatusBadRequest, gin.H{
		"error": err.Error(),
	}
}
