package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/JorgeGorrito/PT-News-API/internal/web/handlers"
	"github.com/JorgeGorrito/PT-News-API/internal/web/middleware"
)

// SetupRouter configures all routes and middleware
func SetupRouter(
	authorsHandler *handlers.AuthorsHandler,
	articlesHandler *handlers.ArticlesHandler,
) *gin.Engine {
	router := gin.Default()

	// Global middleware
	router.Use(middleware.ErrorHandler())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authors routes
		authors := v1.Group("/authors")
		{
			authors.POST("", authorsHandler.Create)
			authors.GET("/:id/summary", authorsHandler.GetSummary)
			authors.GET("/top", authorsHandler.GetTop)
		}

		// Articles routes
		articles := v1.Group("/articles")
		{
			articles.POST("", articlesHandler.Create)
			articles.PUT("/:id/publish", articlesHandler.Publish)
			articles.GET("", articlesHandler.List)
			articles.GET("/:id", articlesHandler.GetByID)
			articles.GET("/author/:id", articlesHandler.ListByAuthor)
		}
	}

	return router
}
