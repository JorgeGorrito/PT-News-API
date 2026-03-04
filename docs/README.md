# API Documentation

Esta carpeta contiene la documentación Swagger/OpenAPI generada automáticamente para PT News API.

## Acceder a la Documentación

Una vez que el servidor esté corriendo, puedes acceder a la documentación interactiva de Swagger en:

```
http://localhost:8080/swagger/index.html
```

## Características de la Documentación

La documentación incluye:

- **8 Endpoints completamente documentados**:
  - `POST /api/v1/authors` - Crear autor
  - `GET /api/v1/authors/{id}/summary` - Resumen del autor
  - `GET /api/v1/authors/top` - Top autores por score
  - `POST /api/v1/articles` - Crear artículo (draft)
  - `PUT /api/v1/articles/{id}/publish` - Publicar artículo
  - `GET /api/v1/articles` - Listar artículos publicados
  - `GET /api/v1/articles/{id}` - Obtener artículo por ID
  - `GET /api/v1/articles/author/{id}` - Listar artículos por autor

- **Descripción detallada** de cada endpoint:
  - Parámetros requeridos y opcionales
  - Formatos de request body
  - Ejemplos de responses (success y error)
  - Códigos de estado HTTP
  - Reglas de validación

- **Fórmula de Score documentada**:
  - `Score = (word_count * 0.1) + (author_published_articles * 5) + recency_bonus`
  - Bonus de recencia: +50 si < 24h, +20 si < 72h

- **Validaciones de publicación**:
  - Mínimo 120 palabras
  - Máximo 35% de palabras repetidas

## Probar los Endpoints

Desde la interfaz de Swagger puedes:
1. Ver todos los endpoints disponibles
2. Leer la documentación detallada
3. Ejecutar requests directamente desde el navegador
4. Ver ejemplos de responses
5. Descargar la especificación OpenAPI (JSON/YAML)

## 🔄 Generación Automática

**Los archivos de documentación se generan automáticamente** - no los edites manualmente.

### Con Docker (recomendado)

```bash
# Regenera automáticamente al hacer build
docker-compose up --build
```

La documentación se genera automáticamente durante el build de la imagen Docker.

### Sin Docker (desarrollo local)

```bash
# Instalar swag CLI (solo una vez)
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación
swag init -g cmd/api/main.go --output docs
```

### Archivos Generados (en .gitignore)

- `docs.go` - Documentación embebida en Go (generado automáticamente)
- `swagger.json` - Especificación OpenAPI en JSON (generado automáticamente)
- `swagger.yaml` - Especificación OpenAPI en YAML (generado automáticamente)

**Nota**: Estos archivos NO están en control de versiones (`git`). Se generan automáticamente en cada build para garantizar que la documentación siempre esté actualizada con el código.
