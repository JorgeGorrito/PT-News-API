package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	createArticle "github.com/JorgeGorrito/PT-News-API/internal/application/articles/commands/create_article"
	publishArticle "github.com/JorgeGorrito/PT-News-API/internal/application/articles/commands/publish_article"
	getArticleByID "github.com/JorgeGorrito/PT-News-API/internal/application/articles/queries/get_article_by_id"
	listArticlesByAuthor "github.com/JorgeGorrito/PT-News-API/internal/application/articles/queries/list_articles_by_author"
	listPublishedArticles "github.com/JorgeGorrito/PT-News-API/internal/application/articles/queries/list_published_articles"
	createAuthor "github.com/JorgeGorrito/PT-News-API/internal/application/authors/commands/create_author"
	getAuthorSummary "github.com/JorgeGorrito/PT-News-API/internal/application/authors/queries/get_author_summary"
	getTopAuthors "github.com/JorgeGorrito/PT-News-API/internal/application/authors/queries/get_top_authors"
	domainservices "github.com/JorgeGorrito/PT-News-API/internal/domain/services"
	"github.com/JorgeGorrito/PT-News-API/internal/infrastructure/config"
	"github.com/JorgeGorrito/PT-News-API/internal/infrastructure/persistence/mysql"
	"github.com/JorgeGorrito/PT-News-API/internal/infrastructure/persistence/mysql/migrations"
	"github.com/JorgeGorrito/PT-News-API/internal/web/handlers"
	"github.com/JorgeGorrito/PT-News-API/internal/web/routes"

	_ "github.com/JorgeGorrito/PT-News-API/docs" // Swagger docs
)

// @title PT News API
// @version 1.0
// @description API RESTful para gestión de artículos y autores con cálculo de puntaje y flujo de publicación
// @description
// @description ## Funcionalidades
// @description - Gestión de autores (crear, resumen, top por puntaje)
// @description - Gestión de artículos (crear como borrador, publicar con validaciones)
// @description - Cálculo dinámico de puntaje basado en conteo de palabras, actividad del autor y recencia
// @description - Validación de publicación: mínimo 120 palabras, máximo 35% de repetición de palabras
// @description
// @description ## Fórmula de Puntaje
// @description Puntaje = (word_count * 0.1) + (author_published_articles * 5) + recency_bonus
// @description - Bono de recencia: +50 si < 24h, +20 si < 72h, 0 de lo contrario
//
// @contact.name Soporte API
// @contact.email j0rg3.4b3ll4@gmail.com
//
//
// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	dbConfig := config.DefaultDatabaseConfig()

	if host := os.Getenv("DB_HOST"); host != "" {
		dbConfig.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &dbConfig.Port)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		dbConfig.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		dbConfig.Password = password
	}
	if database := os.Getenv("DB_NAME"); database != "" {
		dbConfig.Database = database
	}

	db, err := mysql.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	migrationRunner := migrations.NewRunner(db.DB)
	if err := migrationRunner.Run(context.Background()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	authorRepo := mysql.NewAuthorRepository(db)
	articleRepo := mysql.NewArticleRepository(db)

	scoreService := domainservices.NewScoreService()

	txManager := mysql.NewTxManager(db)

	createAuthorHandler := createAuthor.NewHandler(authorRepo)
	createArticleHandler := createArticle.NewHandler(articleRepo, authorRepo, txManager)
	publishArticleHandler := publishArticle.NewHandler(articleRepo, txManager)

	getAuthorSummaryHandler := getAuthorSummary.NewHandler(authorRepo, articleRepo, scoreService)
	getTopAuthorsHandler := getTopAuthors.NewHandler(articleRepo, scoreService)
	getArticleByIDHandler := getArticleByID.NewHandler(articleRepo, authorRepo, scoreService)
	listPublishedArticlesHandler := listPublishedArticles.NewHandler(articleRepo, scoreService)
	listArticlesByAuthorHandler := listArticlesByAuthor.NewHandler(articleRepo, authorRepo)

	authorsHandler := handlers.NewAuthorsHandler(
		createAuthorHandler,
		getAuthorSummaryHandler,
		getTopAuthorsHandler,
	)

	articlesHandler := handlers.NewArticlesHandler(
		createArticleHandler,
		publishArticleHandler,
		getArticleByIDHandler,
		listPublishedArticlesHandler,
		listArticlesByAuthorHandler,
	)

	router := routes.SetupRouter(authorsHandler, articlesHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
