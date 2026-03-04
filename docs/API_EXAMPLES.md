# API Examples

Quick reference for manually testing every endpoint. All examples assume the API is running on `http://localhost:8080`.

---

## Authors

### POST /api/v1/authors — Create author

```bash
curl -s -X POST http://localhost:8080/api/v1/authors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laura Martínez",
    "email": "laura.martinez@example.com",
    "biography": "Periodista de tecnología con 10 años de experiencia en medios digitales."
  }'
```

**Response `201`**
```json
{ "id": 1 }
```

---

### GET /api/v1/authors/{id}/summary — Author summary

```bash
curl -s http://localhost:8080/api/v1/authors/1/summary
```

**Response `200`**
```json
{
  "id": 1,
  "name": "Laura Martínez",
  "email": "laura.martinez@example.com",
  "biography": "Periodista de tecnología con 10 años de experiencia en medios digitales.",
  "published_articles": 2,
  "total_score": 182.5
}
```

---

### GET /api/v1/authors/top — Top N authors by score

```bash
# Top 3 (default)
curl -s "http://localhost:8080/api/v1/authors/top?limit=3"

# Top 5
curl -s "http://localhost:8080/api/v1/authors/top?limit=5"
```

**Response `200`**
```json
{
  "authors": [
    { "id": 1, "name": "Laura Martínez", "published_articles": 2, "total_score": 182.5 },
    { "id": 2, "name": "Carlos Ruiz",    "published_articles": 1, "total_score": 87.0  }
  ]
}
```

---

## Articles

### POST /api/v1/articles — Create article (BORRADOR)

The body must reference an existing `author_id`.

```bash
curl -s -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -d '{
    "author_id": 1,
    "title": "El futuro de la inteligencia artificial en los medios",
    "content": "La inteligencia artificial está transformando radicalmente la manera en que los medios de comunicación producen y distribuyen contenido periodístico. Desde la automatización de noticias financieras hasta la personalización del contenido para cada lector, las redacciones de todo el mundo están adoptando herramientas basadas en modelos de lenguaje de gran escala. Esta tendencia plantea preguntas fundamentales sobre el rol del periodista humano, la veracidad de la información generada por máquinas y los sesgos algorítmicos que pueden filtrarse en el flujo noticioso. En este artículo exploramos los casos de uso más relevantes, los desafíos éticos que surgen y las oportunidades que la IA abre para el periodismo de datos, la verificación de hechos y la cobertura en tiempo real de eventos de alta complejidad informativa."
  }'
```

**Response `201`**
```json
{ "id": 1 }
```

> **Nota**: el `content` debe tener al menos **120 palabras** para poder publicarse después.

---

### PUT /api/v1/articles/{id}/publish — Publish article

Cambia el estado de `BORRADOR` a `PUBLICADO`. Valida mínimo 120 palabras y máximo 35 % de repetición.

```bash
curl -s -X PUT http://localhost:8080/api/v1/articles/1/publish
```

**Response `200`**
```json
{ "id": 1 }
```

**Error — menos de 120 palabras `400`**
```json
{ "error": "El artículo debe tener al menos 120 palabras para ser publicado" }
```

**Error — ya publicado `409`**
```json
{ "error": "El artículo ya está publicado" }
```

---

### GET /api/v1/articles/{id} — Get article by ID

```bash
curl -s http://localhost:8080/api/v1/articles/1
```

**Response `200`**
```json
{
  "id": 1,
  "title": "El futuro de la inteligencia artificial en los medios",
  "content": "La inteligencia artificial está transformando...",
  "word_count": 142,
  "author_id": 1,
  "author_name": "Laura Martínez",
  "status": "PUBLICADO",
  "created_at": "2026-03-04T10:00:00Z",
  "published_at": "2026-03-04T10:05:00Z"
}
```

---

### GET /api/v1/articles — List published articles (paginated)

```bash
# Ordenado por score (default), página 1, 10 por página
curl -s "http://localhost:8080/api/v1/articles"

# Ordenado por fecha de publicación descendente
curl -s "http://localhost:8080/api/v1/articles?order_by=published_at"

# Página 2, 5 por página, ordenado por score
curl -s "http://localhost:8080/api/v1/articles?page=2&page_size=5&order_by=score"
```

| Query param | Valores | Default |
|---|---|---|
| `page` | entero ≥ 1 | `1` |
| `page_size` | entero ≥ 1 | `10` |
| `order_by` | `score` \| `published_at` | `score` |

**Response `200`**
```json
{
  "articles": [
    {
      "id": 1,
      "title": "El futuro de la inteligencia artificial en los medios",
      "content": "La inteligencia artificial está transformando...",
      "word_count": 142,
      "author_id": 1,
      "author_name": "Laura Martínez",
      "published_at": "2026-03-04T10:05:00Z",
      "score": 82.2
    }
  ],
  "total_count": 1,
  "page": 1,
  "page_size": 10
}
```

---

### GET /api/v1/articles/author/{id} — List articles by author

```bash
# Todos los artículos del autor 1
curl -s "http://localhost:8080/api/v1/articles/author/1"

# Solo borradores
curl -s "http://localhost:8080/api/v1/articles/author/1?status=BORRADOR"

# Solo publicados, página 1
curl -s "http://localhost:8080/api/v1/articles/author/1?status=PUBLICADO&page=1&page_size=5"
```

| Query param | Valores | Default |
|---|---|---|
| `status` | `BORRADOR` \| `PUBLICADO` | todos |
| `page` | entero ≥ 1 | `1` |
| `page_size` | entero ≥ 1 | `10` |

**Response `200`**
```json
{
  "articles": [
    {
      "id": 1,
      "title": "El futuro de la inteligencia artificial en los medios",
      "content": "La inteligencia artificial está transformando...",
      "word_count": 142,
      "author_id": 1,
      "status": "PUBLICADO",
      "published_at": "2026-03-04T10:05:00Z"
    }
  ],
  "total_count": 1,
  "page": 1,
  "page_size": 10
}
```

---

## Flujo completo de prueba

Secuencia mínima para tener datos reales en la API:

```bash
BASE=http://localhost:8080/api/v1

# 1. Crear dos autores
curl -s -X POST $BASE/authors -H "Content-Type: application/json" -d '{
  "name": "Laura Martínez",
  "email": "laura@example.com",
  "biography": "Periodista de tecnología."
}'

curl -s -X POST $BASE/authors -H "Content-Type: application/json" -d '{
  "name": "Carlos Ruiz",
  "email": "carlos@example.com",
  "biography": "Editor de política y economía."
}'

# 2. Crear artículos (mínimo 120 palabras en content)
curl -s -X POST $BASE/articles -H "Content-Type: application/json" -d '{
  "author_id": 1,
  "title": "IA en medios digitales",
  "content": "La inteligencia artificial está transformando radicalmente la manera en que los medios de comunicación producen y distribuyen contenido periodístico. Desde la automatización de noticias financieras hasta la personalización del contenido para cada lector, las redacciones de todo el mundo están adoptando herramientas basadas en modelos de lenguaje de gran escala. Esta tendencia plantea preguntas fundamentales sobre el rol del periodista humano, la veracidad de la información generada por máquinas y los sesgos algorítmicos que pueden filtrarse en el flujo noticioso."
}'

curl -s -X POST $BASE/articles -H "Content-Type: application/json" -d '{
  "author_id": 2,
  "title": "Economía colombiana en 2026",
  "content": "La economía colombiana enfrenta en 2026 un escenario de recuperación moderada, impulsada por la inversión extranjera directa y el dinamismo del sector servicios. Los analistas proyectan un crecimiento del PIB cercano al tres por ciento, aunque advierten sobre los riesgos asociados a la volatilidad del peso frente al dólar y la presión inflacionaria en alimentos y energía. El gobierno nacional ha anunciado paquetes de estímulo fiscal orientados a la reactivación de la construcción y la exportación de productos agroindustriales, sectores que históricamente han jalado el empleo formal en las regiones."
}'

# 3. Publicar los artículos
curl -s -X PUT $BASE/articles/1/publish
curl -s -X PUT $BASE/articles/2/publish

# 4. Consultar resultados
curl -s "$BASE/articles?order_by=score"
curl -s "$BASE/authors/top?limit=3"
curl -s "$BASE/authors/1/summary"
```
