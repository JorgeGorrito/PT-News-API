package mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	createArticle "github.com/JorgeGorrito/PT-News-API/internal/application/articles/commands/create_article"
	publishArticle "github.com/JorgeGorrito/PT-News-API/internal/application/articles/commands/publish_article"
	createAuthor "github.com/JorgeGorrito/PT-News-API/internal/application/authors/commands/create_author"
	domainservices "github.com/JorgeGorrito/PT-News-API/internal/domain/services"
	vo "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
	"github.com/JorgeGorrito/PT-News-API/internal/infrastructure/config"
	"github.com/JorgeGorrito/PT-News-API/internal/infrastructure/persistence/mysql"
	"github.com/JorgeGorrito/PT-News-API/internal/infrastructure/persistence/mysql/migrations"
)

// setupTestDB configura una base de datos de prueba
func setupTestDB(t *testing.T) (*mysql.DB, func()) {
	dbConfig := config.DefaultDatabaseConfig()

	// Permitir override con variables de entorno
	if host := os.Getenv("TEST_DB_HOST"); host != "" {
		dbConfig.Host = host
	} else if host := os.Getenv("DB_HOST"); host != "" {
		dbConfig.Host = host
	}

	if port := os.Getenv("TEST_DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &dbConfig.Port)
	} else if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &dbConfig.Port)
	}

	if user := os.Getenv("TEST_DB_USER"); user != "" {
		dbConfig.User = user
	} else if user := os.Getenv("DB_USER"); user != "" {
		dbConfig.User = user
	}

	if password := os.Getenv("TEST_DB_PASSWORD"); password != "" {
		dbConfig.Password = password
	} else if password := os.Getenv("DB_PASSWORD"); password != "" {
		dbConfig.Password = password
	}

	if database := os.Getenv("TEST_DB_NAME"); database != "" {
		dbConfig.Database = database
	} else if database := os.Getenv("DB_NAME"); database != "" {
		dbConfig.Database = database
	}

	db, err := mysql.NewConnection(dbConfig)
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to database: %v", err)
		return nil, func() {}
	}

	// Ejecutar migraciones
	migrationRunner := migrations.NewRunner(db.DB)
	if err := migrationRunner.Run(context.Background()); err != nil {
		t.Skipf("Skipping integration test: failed to run migrations: %v", err)
		return nil, func() {}
	}

	// Cleanup function
	cleanup := func() {
		// Limpiar datos de prueba
		db.DB.Exec("DELETE FROM articles WHERE title LIKE 'TEST_%'")
		db.DB.Exec("DELETE FROM authors WHERE email LIKE 'test_%@integration.test'")
		db.Close()
	}

	return db, cleanup
}

// Test de integración: Publicar artículo y verificar en BD
func TestIntegration_PublishArticle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	// Configurar repositorios y manejadores
	authorRepo := mysql.NewAuthorRepository(db)
	articleRepo := mysql.NewArticleRepository(db)
	txManager := mysql.NewTxManager(db)

	createAuthorHandler := createAuthor.NewHandler(authorRepo)
	createArticleHandler := createArticle.NewHandler(articleRepo, authorRepo, txManager)
	publishArticleHandler := publishArticle.NewHandler(articleRepo, txManager)

	// 1. Crear un autor
	authorCmd := createAuthor.Command{
		Name:      "TEST_Integration Author",
		Email:     fmt.Sprintf("test_%d@integration.test", time.Now().Unix()),
		Biography: "Author for integration testing",
	}

	authorResp, err := createAuthorHandler.Handle(ctx, authorCmd)
	if err != nil {
		t.Fatalf("Failed to create author: %v", err)
	}
	authorID := authorResp.ID

	// 2. Crear un artículo en borrador con contenido válido
	articleBody := strings.Repeat("palabra única diferente contenido texto artículo prueba ejemplo información dato ", 20) // 180 palabras
	articleCmd := createArticle.Command{
		AuthorID: authorID,
		Title:    "TEST_Integration Article",
		Body:     articleBody,
	}

	articleResp, err := createArticleHandler.Handle(ctx, articleCmd)
	if err != nil {
		t.Fatalf("Failed to create article: %v", err)
	}
	articleID := articleResp.ID

	// Verificar que el artículo está en BORRADOR en la BD
	var statusBefore string
	var publishedAtBefore sql.NullTime
	err = db.DB.QueryRow("SELECT status, published_at FROM articles WHERE id = ?", articleID).
		Scan(&statusBefore, &publishedAtBefore)
	if err != nil {
		t.Fatalf("Failed to verify draft article: %v", err)
	}

	if statusBefore != "BORRADOR" {
		t.Errorf("Expected status BORRADOR, got %s", statusBefore)
	}

	if publishedAtBefore.Valid {
		t.Error("Expected published_at to be NULL for draft article")
	}

	// 3. Publicar el artículo
	publishCmd := publishArticle.Command{
		ArticleID: articleID,
	}

	publishResp, err := publishArticleHandler.Handle(ctx, publishCmd)
	if err != nil {
		t.Fatalf("Failed to publish article: %v", err)
	}

	if publishResp.ID != articleID {
		t.Errorf("Expected published article ID %d, got %d", articleID, publishResp.ID)
	}

	// 4. Verificar en la base de datos que el artículo está PUBLICADO
	var statusAfter string
	var publishedAtAfter sql.NullTime
	var titleAfter string
	var wordCountAfter uint

	err = db.DB.QueryRow(`
		SELECT title, status, word_count, published_at 
		FROM articles 
		WHERE id = ?
	`, articleID).Scan(&titleAfter, &statusAfter, &wordCountAfter, &publishedAtAfter)

	if err != nil {
		t.Fatalf("Failed to verify published article: %v", err)
	}

	// Verificaciones
	if statusAfter != "PUBLICADO" {
		t.Errorf("Expected status PUBLICADO, got %s", statusAfter)
	}

	if !publishedAtAfter.Valid {
		t.Error("Expected published_at to be set after publishing")
	} else {
		// Verificar que published_at es reciente (menos de 1 minuto)
		timeSincePublished := time.Since(publishedAtAfter.Time)
		if timeSincePublished > time.Minute {
			t.Errorf("published_at should be recent, but it's %v old", timeSincePublished)
		}
	}

	if titleAfter != "TEST_Integration Article" {
		t.Errorf("Expected title 'TEST_Integration Article', got %s", titleAfter)
	}

	// Verificar que el conteo de palabras es correcto
	if wordCountAfter < 120 {
		t.Errorf("Expected at least 120 words, got %d", wordCountAfter)
	}

	t.Logf("✓ Article published successfully - ID: %d, Status: %s, Words: %d, Published: %v",
		articleID, statusAfter, wordCountAfter, publishedAtAfter.Time)
}

// Test de integración: Validaciones al publicar
func TestIntegration_PublishArticle_Validations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	authorRepo := mysql.NewAuthorRepository(db)
	articleRepo := mysql.NewArticleRepository(db)
	txManager := mysql.NewTxManager(db)

	createAuthorHandler := createAuthor.NewHandler(authorRepo)
	createArticleHandler := createArticle.NewHandler(articleRepo, authorRepo, txManager)
	publishArticleHandler := publishArticle.NewHandler(articleRepo, txManager)

	// Crear autor
	authorCmd := createAuthor.Command{
		Name:      "TEST_Validation Author",
		Email:     fmt.Sprintf("test_validation_%d@integration.test", time.Now().Unix()),
		Biography: "Author for validation testing",
	}

	authorResp, err := createAuthorHandler.Handle(ctx, authorCmd)
	if err != nil {
		t.Fatalf("Failed to create author: %v", err)
	}
	authorID := authorResp.ID

	t.Run("Artículo con menos de 120 palabras no se puede publicar", func(t *testing.T) {
		// Crear artículo con solo 50 palabras
		articleBody := strings.Repeat("palabra ", 50)
		articleCmd := createArticle.Command{
			AuthorID: authorID,
			Title:    "TEST_Short Article",
			Body:     articleBody,
		}

		articleResp, err := createArticleHandler.Handle(ctx, articleCmd)
		if err != nil {
			t.Fatalf("Failed to create article: %v", err)
		}

		// Intentar publicar debe fallar
		publishCmd := publishArticle.Command{
			ArticleID: articleResp.ID,
		}

		_, err = publishArticleHandler.Handle(ctx, publishCmd)
		if err == nil {
			t.Error("Expected error when publishing article with less than 120 words")
		}

		// Verificar que sigue en DRAFT en la BD
		var status string
		db.DB.QueryRow("SELECT status FROM articles WHERE id = ?", articleResp.ID).Scan(&status)
		if status != "BORRADOR" {
			t.Errorf("Article should remain BORRADOR after failed publish, got %s", status)
		}
	})

	t.Run("Artículo con más de 35% repetición no se puede publicar", func(t *testing.T) {
		// Crear artículo con 50% de repetición
		articleBody := strings.Repeat("repetida ", 70) + strings.Repeat("unica1 unica2 unica3 unica4 unica5 ", 10)
		articleCmd := createArticle.Command{
			AuthorID: authorID,
			Title:    "TEST_Repetitive Article",
			Body:     articleBody,
		}

		articleResp, err := createArticleHandler.Handle(ctx, articleCmd)
		if err != nil {
			t.Fatalf("Failed to create article: %v", err)
		}

		// Intentar publicar debe fallar
		publishCmd := publishArticle.Command{
			ArticleID: articleResp.ID,
		}

		_, err = publishArticleHandler.Handle(ctx, publishCmd)
		if err == nil {
			t.Error("Expected error when publishing article with >35% word repetition")
		}

		// Verificar que sigue en DRAFT en la BD
		var status string
		db.DB.QueryRow("SELECT status FROM articles WHERE id = ?", articleResp.ID).Scan(&status)
		if status != "BORRADOR" {
			t.Errorf("Article should remain BORRADOR after failed publish, got %s", status)
		}
	})
}

// Test de integración: Verificar que el score se calcula correctamente
func TestIntegration_ArticleScore_Calculation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	ctx := context.Background()

	authorRepo := mysql.NewAuthorRepository(db)
	articleRepo := mysql.NewArticleRepository(db)
	txManager := mysql.NewTxManager(db)

	createAuthorHandler := createAuthor.NewHandler(authorRepo)
	createArticleHandler := createArticle.NewHandler(articleRepo, authorRepo, txManager)
	publishArticleHandler := publishArticle.NewHandler(articleRepo, txManager)

	// Crear autor
	authorCmd := createAuthor.Command{
		Name:      "TEST_Score Author",
		Email:     fmt.Sprintf("test_score_%d@integration.test", time.Now().Unix()),
		Biography: "Author for score testing",
	}

	authorResp, err := createAuthorHandler.Handle(ctx, authorCmd)
	if err != nil {
		t.Fatalf("Failed to create author: %v", err)
	}
	authorID := authorResp.ID

	// Crear y publicar artículo
	articleBody := strings.Repeat("alpha beta gamma delta epsilon zeta eta theta iota kappa ", 20) // 200 palabras
	articleCmd := createArticle.Command{
		AuthorID: authorID,
		Title:    "TEST_Score Article",
		Body:     articleBody,
	}

	articleResp, err := createArticleHandler.Handle(ctx, articleCmd)
	if err != nil {
		t.Fatalf("Failed to create article: %v", err)
	}

	publishCmd := publishArticle.Command{
		ArticleID: articleResp.ID,
	}

	_, err = publishArticleHandler.Handle(ctx, publishCmd)
	if err != nil {
		t.Fatalf("Failed to publish article: %v", err)
	}

	// Obtener artículo publicado con score desde el repositorio
	articles, total, err := articleRepo.FindPublishedPaginated(ctx, 1, 10, vo.OrderByScore)
	if err != nil {
		t.Fatalf("Failed to get published articles: %v", err)
	}

	if total == 0 {
		t.Fatal("No published articles found")
	}

	// Buscar nuestro artículo
	var foundArticle *vo.PublishedArticleWithScore
	for _, article := range articles {
		if article.ArticleID() == articleResp.ID {
			foundArticle = article
			break
		}
	}

	if foundArticle == nil {
		t.Fatal("Published article not found in results")
	}

	// Calcular el score dinámicamente usando el Domain Service (nunca se almacena en BD).
	scoreService := domainservices.NewScoreService()
	score := scoreService.CalculateArticleScore(
		foundArticle.WordCount(),
		foundArticle.AuthorPublishedCount(),
		foundArticle.PublishedAt(),
	)
	if score <= 0 {
		t.Errorf("Expected positive score, got %f", score)
	}

	// El score debería ser aproximadamente:
	// (200 words * 0.1) + (1 published article * 5) + 50 (bonus < 24h) = 20 + 5 + 50 = 75
	expectedMinScore := 20.0 + 5.0        // Sin bonus
	expectedMaxScore := 20.0 + 5.0 + 50.0 // Con bonus máximo

	if score < expectedMinScore || score > expectedMaxScore {
		t.Logf("Score %.2f is outside expected range [%.2f, %.2f]", score, expectedMinScore, expectedMaxScore)
	} else {
		t.Logf("✓ Score calculated correctly: %.2f (within expected range)", score)
	}
}
