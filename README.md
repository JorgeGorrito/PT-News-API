# PT News API

RESTful API para gestión de artículos y autores con sistema de puntuación dinámica y flujo de publicación con validaciones.

## 📋 Características

- **Gestión de Autores**: Crear autores, obtener resumen con estadísticas, ranking por score
- **Gestión de Artículos**: Crear artículos en borrador, publicar con validaciones
- **Sistema de Puntuación**: Score dinámico calculado en base a:
  - Conteo de palabras del artículo
  - Cantidad de artículos publicados del autor
  - Bonus por recencia de publicación
- **Validaciones de Publicación**:
  - Mínimo 120 palabras
  - Máximo 35% de palabras repetidas
- **Documentación Interactiva**: Swagger UI completo

## 🏗️ Arquitectura

Arquitectura Hexagonal (Ports & Adapters) inspirada en DDD con patrón CQRS:

```
├── cmd/api/                    # Entry point
├── internal/
│   ├── domain/                 # Capa de dominio (entities, value objects, interfaces)
│   ├── application/            # Casos de uso (CQRS: commands & queries)
│   ├── infrastructure/         # Implementaciones (MySQL, config)
│   └── web/                    # HTTP handlers, DTOs, middleware
└── docs/                       # Swagger documentation
```

### Capas

- **Domain**: Lógica de negocio pura, sin dependencias externas
- **Application**: Orquestación con CQRS (Commands retornan ID, Queries retornan data completa)
- **Infrastructure**: Implementaciones concretas (MySQL con transacciones, mappers)
- **Web**: HTTP layer con Gin (handlers, DTOs, middleware de errores)

## 🚀 Inicio Rápido con Docker

### Prerrequisitos

- Docker 20.10+
- Docker Compose 2.0+

### Ejecutar

```bash
# Construir y ejecutar todos los servicios
docker-compose up --build

# Ejecutar en segundo plano (detached)
docker-compose up -d

# Ver logs
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f api
docker-compose logs -f db

# Detener servicios
docker-compose down

# Detener y eliminar volúmenes (⚠️ elimina datos de la BD)
docker-compose down -v
```

La API estará disponible en:
- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **MySQL**: localhost:3306

### Credenciales de Base de Datos

```
Host: localhost (o 'db' desde dentro de Docker)
Port: 3306
Database: pt_news_api
User: app_user
Password: app_password
Root Password: root
```

### Acceder a los Contenedores

```bash
# Acceder al contenedor de la API
docker exec -it pt-news-api sh

# Acceder a MySQL
docker exec -it pt-news-db mysql -u root -proot pt_news_api

# Ver estado de salud de los servicios
docker-compose ps
```

### Documentación Swagger

La documentación Swagger se genera **automáticamente** durante la construcción de la imagen Docker. Los archivos `docs/docs.go`, `swagger.json` y `swagger.yaml` se regeneran en cada build.

**Para actualizar después de cambiar anotaciones Swagger:**
```bash
docker-compose up --build
```

## 🛠️ Desarrollo Local (sin Docker)

### Prerrequisitos

- Go 1.25.1+
- MySQL 8.0+

### Configuración

1. **Base de datos**:
```bash
mysql -u root -p
CREATE DATABASE pt_news_api;
```

**Las migraciones se ejecutan automáticamente** al iniciar la aplicación.

2. **Variables de entorno**:
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=root
export DB_NAME=pt_news_api
export PORT=8080
```

3. **Instalar dependencias**:
```bash
go mod download
```

4. **Generar documentación Swagger**:
```bash
# Instalar swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación
swag init -g cmd/api/main.go --output docs
```

5. **Ejecutar**:
```bash
go run cmd/api/main.go
```

## 🧪 Testing

El proyecto incluye pruebas unitarias e de integración siguiendo las mejores prácticas de Go.

### Ejecutar Tests Unitarios

Los tests unitarios no requieren base de datos (usan el flag `-short`):

```bash
# Ejecutar todos los tests unitarios
go test ./... -short

# Con detalles verbose
go test ./... -short -v

# Solo tests de un paquete específico
go test ./internal/domain/entities/... -v
go test ./internal/domain/value-objects/... -v
go test ./internal/application/authors/queries/get_top_authors/... -v

# Con cobertura
go test ./... -short -cover

# Reporte de cobertura en HTML
go test ./... -short -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Tests Implementados

#### ✅ **Pruebas Unitarias (3 requeridas)**

1. **Validaciones antes de publicar** ([`internal/domain/entities/article_test.go`](internal/domain/entities/article_test.go))
   - ✅ Validación mínimo 120 palabras
   - ✅ Validación máximo 35% palabras repetidas
   - ✅ No publicar artículo ya publicado

2. **Cálculo de score** ([`internal/domain/value-objects/score_test.go`](internal/domain/value-objects/score_test.go))
   - ✅ Fórmula: `(word_count * 0.1) + (author_published * 5) + bonus`
   - ✅ Bonus recencia < 24h: +50
   - ✅ Bonus recencia 24-72h: +20
   - ✅ Sin bonus > 72h: +0

3. **Endpoint Top Autores** ([`internal/application/authors/queries/get_top_authors/handler_test.go`](internal/application/authors/queries/get_top_authors/handler_test.go))
   - ✅ Retorna top N autores por score
   - ✅ Manejo de límites inválidos
   - ✅ Manejo de empates
   - ✅ Validación de datos mapeados

#### ✅ **Prueba de Integración (1 requerida)**

**Publicar artículo end-to-end** ([`internal/infrastructure/persistence/mysql/integration_test.go`](internal/infrastructure/persistence/mysql/integration_test.go))
- ✅ Crear autor → Crear artículo → Publicar → Verificar en BD
- ✅ Validar cambio de estado BORRADOR → PUBLICADO
- ✅ Validar asignación de `published_at`.

**Con Docker:**

```bash
# 1. Iniciar solo la base de datos
docker-compose up -d db

# 2. Esperar a que MySQL esté listo (~5 segundos)
sleep 5

# 3. Ejecutar tests de integración
go test ./internal/infrastructure/persistence/mysql/... -v

# 4. Ejecutar TODOS los tests (unitarios + integración)
go test ./... -v
```

**Dentro de Docker (todo en contenedor):**

```bash
# Tests unitarios (no requiere BD)
docker-compose run --rm api go test ./... -short -v

# Tests de integración (con BD en Docker)
docker-compose up -d db
sleep 5
docker-compose run --rm \
  -e DB_HOST=db \
  -e DB_USER=app_user \
  -e DB_PASSWORD=app_password \
  -e DB_NAME=pt_news_api \
  api go test ./... -v

# Limpiar
docker-compose down
```

**Sin Docker (MySQL local):**

```bash
# Configurar variables de entorno
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=root
export DB_NAME=pt_news_api

# Ejecutar tests

# Ejecutar tests de integración
go test ./internal/infrastructure/persistence/mysql/... -v

# Ejecutar TODOS los tests (unitarios + integración)
go test ./... -v
```

### Cobertura de Tests

```bash
# Generar reporte HTML de cobertura
go test ./... -short -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Ver porcentaje de cobertura por paquete
go test ./... -short -cover
```

## 📡 Endpoints

### Authors

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/api/v1/authors` | Crear autor |
| `GET` | `/api/v1/authors/{id}/summary` | Resumen del autor con estadísticas |
| `GET` | `/api/v1/authors/top?limit=10` | Top autores por score |

### Articles

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/api/v1/articles` | Crear artículo (DRAFT) |
| `PUT` | `/api/v1/articles/{id}/publish` | Publicar artículo (valida reglas) |
| `GET` | `/api/v1/articles?page=1&page_size=10&order_by=score` | Listar artículos publicados |
| `GET` | `/api/v1/articles/{id}` | Obtener artículo por ID |
| `GET` | `/api/v1/articles/author/{id}?status=PUBLISHED` | Listar artículos de un autor |

### Swagger

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/swagger/index.html` | Documentación interactiva |

## 📊 Fórmula de Score

```
Score = (word_count × 0.1) + (author_published_articles × 5) + recency_bonus

Donde recency_bonus:
- +50 si publicado hace menos de 24 horas
- +20 si publicado hace menos de 72 horas
- 0 en otros casos
```

El score se calcula **dinámicamente en la consulta SQL**, nunca se almacena.

## 📝 Ejemplos de Uso

### Crear un autor

```bash
curl -X POST http://localhost:8080/api/v1/authors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "biography": "Technical writer and blogger"
  }'
```

### Crear un artículo

```bash
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Clean Architecture",
    "content": "Clean architecture is a software design philosophy...",
    "author_id": 1
  }'
```

### Publicar un artículo

```bash
curl -X PUT http://localhost:8080/api/v1/articles/1/publish
```

### Listar artículos publicados

```bash
curl "http://localhost:8080/api/v1/articles?page=1&page_size=10&order_by=score"
```
BORRADOR', 'PUBLICADO
### Obtener top autores

```bash
curl "http://localhost:8080/api/v1/authors/top?limit=5"
```

## 🧪 Testing

Para probar los endpoints, usa:

1. **Swagger UI** (recomendado): http://localhost:8080/swagger/index.html
2. **cURL**: Ver ejemplos arriba
3. **Postman/Insomnia**: Importa la especificación desde `/swagger/doc.json`

## 🗄️ Base de Datos

### Esquema

**authors**
- `id` (INT, PK, AUTO_INCREMENT)
- `name` (VARCHAR(255))
- `email` (VARCHAR(255), UNIQUE)
- `biography` (TEXT)
- `created_at` (TIMESTAMP)

**articles**
- `id` (INT, PK, AUTO_INCREMENT)
- `author_id` (INT, FK → authors.id)
- `title` (VARCHAR(255))
- `body` (TEXT)
- `word_count` (INT)
- `status` (ENUM: 'BORRADOR', 'PUBLICADO')
- `created_at` (TIMESTAMP)
- `published_at` (TIMESTAMP, NULL)

**migrations**
- `id` (INT, PK, AUTO_INCREMENT)
- `name` (VARCHAR(255), UNIQUE)
- `executed_at` (TIMESTAMP)

### Migraciones Automáticas

La aplicación verifica y ejecuta migraciones automáticamente al iniciar:
1. Crea la tabla `migrations` si no existe
2. Verifica qué migraciones ya se ejecutaron
3. Ejecuta solo las migraciones pendientes
4. Cada migración se ejecuta en una transacción

**No necesitas ejecutar scripts SQL manualmente.**

### Índices

- `idx_email` en `authors.email`
- `idx_author_status` en `articles(author_id, status)`
- `idx_published_at` en `articles.published_at`
- `idx_status_published` en `articles(status, published_at)`

## 
- **Lenguaje**: Go 1.25.1
- **Framework Web**: Gin
- **Base de Datos**: MySQL 8.0
- **Documentación**: Swagger/OpenAPI 3.0
- **Contenerización**: Docker & Docker Compose
- **Patrón Arquitectónico**: Hexagonal + DDD + CQRS

## 📦 Dependencias

```go
require (
    github.com/gin-gonic/gin v1.12.0
    github.com/go-sql-driver/mysql v1.9.3
    github.com/swaggo/gin-swagger v1.6.1
    github.com/swaggo/files
    github.com/swaggo/swag v1.16.6
)
```

## 🎯 Decisiones Arquitectónicas

1. **CQRS**: Separación clara entre Commands (escritura) y Queries (lectura)
2. **Score dinámico**: Never stored, calculated in SQL for accuracy
3. **Transaction Manager**: Context-based para transacciones transparentes
4. **Dependency Inversion**: Application define interfaces, Infrastructure implementa
5. **Error Middleware**: Mapeo centralizado de errores de dominio a HTTP status
6. **DTOs en Web layer**: Separación entre modelo de dominio y contratos HTTP

## 📚 Documentación Adicional

- [docs/README.md](docs/README.md) - Documentación de Swagger
- Swagger UI: http://localhost:8080/swagger/index.html (cuando la app esté corriendo)

## 🐛 Troubleshooting

### La API no arranca

1. Verifica que MySQL esté corriendo y accesible
2. Revisa las variables de entorno de conexión
3. Las migraciones se ejecutan automáticamente - si falla alguna, la app no iniciará
4. Revisa los logs para ver el error de migración

### Error en migraciones

```bash
# Ver logs de migraciones
docker-compose logs api | grep -i migration

# Reiniciar migraciones (elimina tabla migrations)
docker exec -it pt-news-db mysql -u root -proot pt_news_api -e "DROP TABLE IF EXISTS migrations;"
docker-compose restart api
```

### Error de conexión a MySQL en Docker

El contenedor de la API espera a que MySQL esté healthy gracias al health check. Si persiste:

```bash
docker-compose logs db
docker-compose restart api
```

### Regenerar documentación Swagger

```bash
swag init -g cmd/api/main.go -o docs
```

## 👨‍💻 Autor

PT News API - Prueba Técnica

## 📄 Licencia

MIT
