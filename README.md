# Product API

Go REST service for creating and partially updating products. PostgreSQL, Clean Architecture, constructor-based DI, Swagger UI at `/api-docs`.

## Requirements covered

- `POST /product` and `PATCH /product/{id}` (partial update, nullable `description` / `sale_price`)
- Unit tests
  - Domain / function tests: `internal/domain`
  - Service tests: `internal/service` (in-memory fake repository)
  - Repository integration tests: `internal/repository/postgres` (GORM + PostgreSQL)
  - Component / in-service E2E: `internal/component`
- Swagger: http://localhost:8080/api-docs
- PostgreSQL

## Architecture

```
HTTP handler  ->  service  ->  domain rules
                         |
                         v
                 PostgreSQL repository (GORM)
```

`cmd/api` loads config, opens GORM (`gorm.io/driver/postgres`), AutoMigrate, then injects:

`ProductRepository` → `ProductService` → `ProductHandler` → Chi router

## How to start

### 1. Prerequisites

- Go 1.22+
- Docker (for PostgreSQL and integration tests)

### 2. Start PostgreSQL

```bash
docker compose up -d
```

### 3. Run the API

```bash
go run ./cmd/api
```

Default env:

| Variable | Default |
|---|---|
| `HTTP_ADDR` | `:8080` |
| `DATABASE_URL` | `postgres://product:product@localhost:5433/product?sslmode=disable` |

Windows PowerShell:

```powershell
$env:HTTP_ADDR=":8080"
$env:DATABASE_URL="postgres://product:product@localhost:5433/product?sslmode=disable"
go run ./cmd/api
```

### 4. Open docs

- Health: http://localhost:8080/health
- Swagger UI: http://localhost:8080/api-docs
- OpenAPI JSON: http://localhost:8080/api-docs/doc.json

## API

### Create product

`POST /product`

```json
{
  "name": "Latte",
  "description": "hot milk coffee",
  "sale_price": 80.0,
  "price": 100.0
}
```

`description` and `sale_price` may be `null`.

Response:

```json
{
  "successful": true,
  "error_code": "SUCCESS",
  "data": {
    "id": "product-uuid",
    "name": "Latte",
    "description": "hot milk coffee",
    "sale_price": 80.0,
    "price": 100.0,
    "created_at": "2026-08-20T10:00:00Z",
    "updated_at": "2026-08-20T10:00:00Z"
  }
}
```

### Patch product

`PATCH /product/{id}`

Send **only fields to change**. Omitted fields stay as-is. `description` and `sale_price` can be set to `null` to clear them. `name` and `price` cannot be `null`.

```json
{
  "description": null,
  "price": 110
}
```

Response:

```json
{
  "successful": true,
  "error_code": "SUCCESS"
}
```

Error codes: `SUCCESS`, `VALIDATION_ERROR`, `NOT_FOUND`, `INTERNAL_ERROR`.

### curl examples

```bash
curl -s -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Latte\",\"description\":\"hot\",\"sale_price\":80,\"price\":100}"

curl -s -X PATCH http://localhost:8080/product/<id> \
  -H "Content-Type: application/json" \
  -d "{\"sale_price\":null,\"price\":95}"
```

## Tests

Fast unit tests (no Docker):

```bash
go test ./internal/domain ./internal/service ./internal/handler/http -count=1
```

All tests, including PostgreSQL integration and component tests:

```bash
docker compose up -d
go test ./... -count=1 -timeout 5m
```

If Docker is not running, integration/component tests **skip**. Fast unit tests still pass.

To skip integration even when Postgres is up:

```bash
SKIP_INTEGRATION=1 go test ./... -count=1
```

## Validation rules

- `name` required, 1–255 characters
- `price` >= 0
- `sale_price` if set must be >= 0 and <= `price`
