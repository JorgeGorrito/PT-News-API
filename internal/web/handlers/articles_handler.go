package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	createArticle "github.com/JorgeGorrito/PT-News-API/internal/application/articles/commands/create_article"
	publishArticle "github.com/JorgeGorrito/PT-News-API/internal/application/articles/commands/publish_article"
	getArticleByID "github.com/JorgeGorrito/PT-News-API/internal/application/articles/queries/get_article_by_id"
	listArticlesByAuthor "github.com/JorgeGorrito/PT-News-API/internal/application/articles/queries/list_articles_by_author"
	listPublishedArticles "github.com/JorgeGorrito/PT-News-API/internal/application/articles/queries/list_published_articles"
	valueobjects "github.com/JorgeGorrito/PT-News-API/internal/domain/value-objects"
	"github.com/JorgeGorrito/PT-News-API/internal/web/dto"
)

type ArticlesHandler struct {
	createArticleHandler         *createArticle.Handler
	publishArticleHandler        *publishArticle.Handler
	getArticleByIDHandler        *getArticleByID.Handler
	listPublishedArticlesHandler *listPublishedArticles.Handler
	listArticlesByAuthorHandler  *listArticlesByAuthor.Handler
}

func NewArticlesHandler(
	createArticleHandler *createArticle.Handler,
	publishArticleHandler *publishArticle.Handler,
	getArticleByIDHandler *getArticleByID.Handler,
	listPublishedArticlesHandler *listPublishedArticles.Handler,
	listArticlesByAuthorHandler *listArticlesByAuthor.Handler,
) *ArticlesHandler {
	return &ArticlesHandler{
		createArticleHandler:         createArticleHandler,
		publishArticleHandler:        publishArticleHandler,
		getArticleByIDHandler:        getArticleByIDHandler,
		listPublishedArticlesHandler: listPublishedArticlesHandler,
		listArticlesByAuthorHandler:  listArticlesByAuthorHandler,
	}
}

// Create godoc
// @Summary Crear un nuevo artículo
// @Description Crea un nuevo artículo en estado BORRADOR. El artículo debe tener título, contenido e ID de autor válido. El conteo de palabras se calcula automáticamente.
// @Tags Artículos
// @Accept json
// @Produce json
// @Param request body dto.CreateArticleRequest true "Solicitud de creación de artículo"
// @Success 201 {object} dto.CreateArticleResponse "Artículo creado exitosamente en estado BORRADOR"
// @Failure 400 {object} map[string]string "Entrada inválida o error de validación"
// @Failure 404 {object} map[string]string "Autor no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/articles [post]
func (h *ArticlesHandler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := createArticle.Command{
		Title:    req.Title,
		Body:     req.Content,
		AuthorID: req.AuthorID,
	}

	resp, err := h.createArticleHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.CreateArticleResponse{ID: resp.ID})
}

// Publish godoc
// @Summary Publicar un artículo
// @Description Publica un artículo en BORRADOR. Valida que el artículo tenga mínimo 120 palabras y máximo 35% de palabras repetidas. Cambia el estado de BORRADOR a PUBLICADO.
// @Tags Artículos
// @Accept json
// @Produce json
// @Param id path int true "ID del artículo"
// @Success 200 {object} dto.PublishArticleResponse "Artículo publicado exitosamente"
// @Failure 400 {object} map[string]string "ID de artículo inválido o validación fallida (conteo de palabras o repetición)"
// @Failure 404 {object} map[string]string "Artículo no encontrado"
// @Failure 409 {object} map[string]string "Artículo ya fue publicado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/articles/{id}/publish [put]
func (h *ArticlesHandler) Publish(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de artículo inválido"})
		return
	}

	cmd := publishArticle.Command{ArticleID: id}

	resp, err := h.publishArticleHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.PublishArticleResponse{ID: resp.ID})
}

// GetByID godoc
// @Summary Obtener artículo por ID
// @Description Obtiene información detallada de un artículo incluyendo título, contenido, conteo de palabras, información del autor, estado y fecha de creación.
// @Tags Artículos
// @Accept json
// @Produce json
// @Param id path int true "ID del artículo"
// @Success 200 {object} dto.ArticleResponse "Artículo obtenido exitosamente"
// @Failure 400 {object} map[string]string "ID de artículo inválido"
// @Failure 404 {object} map[string]string "Artículo no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/articles/{id} [get]
func (h *ArticlesHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de artículo inválido"})
		return
	}

	query := getArticleByID.Query{
		ArticleID:    id,
		IncludeScore: true,
	}

	resp, err := h.getArticleByIDHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	article := dto.ArticleResponse{
		ID:         resp.ID,
		Title:      resp.Title,
		Content:    resp.Body,
		WordCount:  int(resp.WordCount),
		AuthorID:   resp.AuthorID,
		AuthorName: resp.AuthorName,
		Status:     resp.Status,
		CreatedAt:  resp.CreatedAt,
	}

	if resp.PublishedAt != nil {
		article.PublishedAt = resp.PublishedAt
	}

	c.JSON(http.StatusOK, article)
}

// List godoc
// @Summary Listar artículos publicados
// @Description Obtiene listado paginado de artículos publicados. Los artículos se ordenan por puntaje (por defecto) o fecha de publicación. El puntaje se calcula a partir del conteo de palabras, cantidad de artículos publicados del autor y bono de recencia.
// @Tags Artículos
// @Accept json
// @Produce json
// @Param page query int false "Número de página (por defecto: 1)" default(1)
// @Param page_size query int false "Cantidad de elementos por página (por defecto: 10)" default(10)
// @Param order_by query string false "Campo de ordenamiento: 'score' o 'published_at' (por defecto: score)" default(score)
// @Success 200 {object} dto.ListArticlesResponse "Artículos publicados obtenidos exitosamente"
// @Failure 400 {object} map[string]string "Parámetros de consulta inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/articles [get]
func (h *ArticlesHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	orderBy := c.DefaultQuery("order_by", "score")

	query := listPublishedArticles.Query{
		Page:    page,
		PerPage: pageSize,
		OrderBy: orderBy,
	}

	resp, err := h.listPublishedArticlesHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	articles := make([]dto.PublishedArticleResponse, len(resp.Articles))
	for i, article := range resp.Articles {
		articles[i] = dto.PublishedArticleResponse{
			ID:          article.ID,
			Title:       article.Title,
			Content:     article.Body,
			WordCount:   int(article.WordCount),
			AuthorID:    article.AuthorID,
			AuthorName:  article.AuthorName,
			PublishedAt: article.PublishedAt,
			Score:       article.Score,
		}
	}

	c.JSON(http.StatusOK, dto.ListArticlesResponse{
		Articles:   articles,
		TotalCount: resp.Total,
		Page:       resp.Page,
		PageSize:   resp.PerPage,
	})
}

// ListByAuthor godoc
// @Summary Listar artículos por autor
// @Description Obtiene listado paginado de artículos de un autor específico. Puede filtrar por estado (DRAFT o PUBLISHED).
// @Tags Artículos
// @Accept json
// @Produce json
// @Param id path int true "ID del autor"
// @Param page query int false "Número de página (por defecto: 1)" default(1)
// @Param page_size query int false "Cantidad de elementos por página (por defecto: 10)" default(10)
// @Param status query string false "Filtrar por estado: 'DRAFT' o 'PUBLISHED'"
// @Success 200 {object} object "Artículos obtenidos exitosamente con información de paginación"
// @Failure 400 {object} map[string]string "Parámetros inválidos o valor de estado inválido"
// @Failure 404 {object} map[string]string "Autor no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/articles/author/{id} [get]
func (h *ArticlesHandler) ListByAuthor(c *gin.Context) {
	idParam := c.Param("id")
	authorID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de autor inválido"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	status := c.DefaultQuery("status", "")

	query := listArticlesByAuthor.Query{
		AuthorID: authorID,
		Page:     page,
		PerPage:  pageSize,
	}

	// Validate and set Status if provided
	if status != "" {
		articleStatus, err := valueobjects.NewArticleStatus(status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido: debe ser DRAFT o PUBLISHED"})
			return
		}
		query.Status = &articleStatus
	}

	resp, err := h.listArticlesByAuthorHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	type ArticleByAuthorResponse struct {
		ID          int64      `json:"id"`
		Title       string     `json:"title"`
		Content     string     `json:"content"`
		WordCount   int        `json:"word_count"`
		AuthorID    int64      `json:"author_id"`
		Status      string     `json:"status"`
		PublishedAt *time.Time `json:"published_at,omitempty"`
	}

	type ListByAuthorResponse struct {
		Articles   []ArticleByAuthorResponse `json:"articles"`
		TotalCount int                       `json:"total_count"`
		Page       int                       `json:"page"`
		PageSize   int                       `json:"page_size"`
	}

	articles := make([]ArticleByAuthorResponse, len(resp.Articles))
	for i, article := range resp.Articles {
		articles[i] = ArticleByAuthorResponse{
			ID:          article.ID,
			Title:       article.Title,
			Content:     article.Body,
			WordCount:   int(article.WordCount),
			AuthorID:    article.AuthorID,
			Status:      article.Status,
			PublishedAt: article.PublishedAt,
		}
	}

	c.JSON(http.StatusOK, ListByAuthorResponse{
		Articles:   articles,
		TotalCount: resp.Total,
		Page:       resp.Page,
		PageSize:   resp.PerPage,
	})
}
