package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	createAuthor "github.com/JorgeGorrito/PT-News-API/internal/application/authors/commands/create_author"
	getAuthorSummary "github.com/JorgeGorrito/PT-News-API/internal/application/authors/queries/get_author_summary"
	getTopAuthors "github.com/JorgeGorrito/PT-News-API/internal/application/authors/queries/get_top_authors"
	"github.com/JorgeGorrito/PT-News-API/internal/web/dto"
)

type AuthorsHandler struct {
	createAuthorHandler     *createAuthor.Handler
	getAuthorSummaryHandler *getAuthorSummary.Handler
	getTopAuthorsHandler    *getTopAuthors.Handler
}

func NewAuthorsHandler(
	createAuthorHandler *createAuthor.Handler,
	getAuthorSummaryHandler *getAuthorSummary.Handler,
	getTopAuthorsHandler *getTopAuthors.Handler,
) *AuthorsHandler {
	return &AuthorsHandler{
		createAuthorHandler:     createAuthorHandler,
		getAuthorSummaryHandler: getAuthorSummaryHandler,
		getTopAuthorsHandler:    getTopAuthorsHandler,
	}
}

// Create godoc
// @Summary Crear un nuevo autor
// @Description Crea un nuevo autor con nombre, email y biografía opcional. El email debe ser único.
// @Tags Autores
// @Accept json
// @Produce json
// @Param request body dto.CreateAuthorRequest true "Solicitud de creación de autor"
// @Success 201 {object} dto.CreateAuthorResponse "Autor creado exitosamente"
// @Failure 400 {object} map[string]string "Entrada inválida o error de validación"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/authors [post]
func (h *AuthorsHandler) Create(c *gin.Context) {
	var req dto.CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := createAuthor.Command{
		Name:      req.Name,
		Email:     req.Email,
		Biography: req.Biography,
	}

	resp, err := h.createAuthorHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.CreateAuthorResponse{ID: resp.ID})
}

// GetSummary godoc
// @Summary Obtener resumen del autor
// @Description Obtiene resumen detallado de un autor incluyendo cantidad de artículos publicados, puntaje total e información del autor
// @Tags Autores
// @Accept json
// @Produce json
// @Param id path int true "ID del autor"
// @Success 200 {object} dto.AuthorSummaryResponse "Resumen del autor obtenido exitosamente"
// @Failure 400 {object} map[string]string "ID de autor inválido"
// @Failure 404 {object} map[string]string "Autor no encontrado"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/authors/{id}/summary [get]
func (h *AuthorsHandler) GetSummary(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de autor inválido"})
		return
	}

	query := getAuthorSummary.Query{AuthorID: id}

	resp, err := h.getAuthorSummaryHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.AuthorSummaryResponse{
		ID:                resp.ID,
		Name:              resp.Name,
		Email:             resp.Email,
		Biography:         resp.Biography,
		PublishedArticles: resp.PublishedArticles,
		TotalScore:        resp.TotalScore,
	})
}

// GetTop godoc
// @Summary Obtener top de autores por puntaje
// @Description Obtiene listado de mejores autores ordenados por puntaje total. El puntaje se calcula a partir del conteo de palabras de los artículos, cantidad de artículos publicados y bono de recencia.
// @Tags Autores
// @Accept json
// @Produce json
// @Param limit query int false "Cantidad de mejores autores a retornar (por defecto: 10)" default(10)
// @Success 200 {array} dto.TopAuthorResponse "Mejores autores obtenidos exitosamente"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /api/v1/authors/top [get]
func (h *AuthorsHandler) GetTop(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	query := getTopAuthors.Query{Limit: limit}

	resp, err := h.getTopAuthorsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	topAuthors := make([]dto.TopAuthorResponse, len(resp.Authors))
	for i, author := range resp.Authors {
		topAuthors[i] = dto.TopAuthorResponse{
			ID:                author.ID,
			Name:              author.Name,
			PublishedArticles: author.PublishedArticles,
			TotalScore:        author.TotalScore,
		}
	}

	c.JSON(http.StatusOK, topAuthors)
}
