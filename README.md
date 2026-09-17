# Go Gin Microservices Blueprint

Blueprint microservices Go yang rapi dan mudah dikembangkan menggunakan Gin, validator, Logrus, godotenv, PostgreSQL, GORM, dan OpenAPI.

## Struktur

```text
.
├── cmd/users                 # entrypoint users service
├── internal/config           # environment configuration
├── internal/database         # PostgreSQL connection
├── internal/httpmiddleware   # request ID dan access logging
├── internal/logger           # Logrus setup
├── internal/server           # HTTP server, health routes, dan lifecycle
├── internal/transport        # shared HTTP response formatter
├── internal/users            # model, payload, response, handler, router, service, repository, tests
├── migrations                # SQL migration yang dijalankan Compose
├── docs/openapi.yaml         # kontrak OpenAPI 3
├── docker-compose.yml        # local PostgreSQL
├── Makefile
└── .env.example
```

Pola ini mengikuti SOLID: handler hanya mengurus HTTP, service mengurus business logic, dan repository diakses melalui interface. GORM berada di adapter repository sehingga service tetap mudah di-unit-test tanpa database. Payload request dan response DTO dipisahkan dari domain model. Formatter response bersama menerapkan envelope `data`, `meta`, dan `error` secara DRY.

Setiap response memiliki `meta.request_id`, `meta.timestamp`, dan header `X-Request-ID`. Client boleh mengirim `X-Request-ID`; jika tidak, service akan membuat UUID baru. ID yang sama ditulis pada access log untuk tracing.

## Prasyarat

- Go 1.23+
- Docker dan Docker Compose
- `curl` untuk smoke test API

## Development flow

1. Siapkan environment:

   ```bash
   cp .env.example .env
   go mod download
   ```

2. Jalankan PostgreSQL dan migration awal:

   ```bash
   make infra-up
   ```

3. Jalankan service:

   ```bash
   make dev
   ```

4. Jalankan quality checks sebelum membuat pull request:

   ```bash
   make lint
   make test
   ```

5. Coba endpoint dasar:

   ```bash
   curl http://localhost:8080/health
   curl http://localhost:8080/ready
   curl -X POST http://localhost:8080/api/v1/users \
     -H 'Content-Type: application/json' \
     -d '{"name":"Ada Lovelace","email":"ada@example.com"}'
   curl 'http://localhost:8080/api/v1/users?page=1&limit=20'
   ```

6. Setelah selesai bekerja:

   ```bash
   make infra-down
   ```

## Testing

Unit test memakai stub repository sehingga tidak membutuhkan PostgreSQL. Test mencakup normalisasi input pada service, pagination, validasi payload handler, duplicate email conflict, response envelope, request ID, health route, dan graceful shutdown.

```bash
go test ./...
go test -race ./...
```

Integration test yang membutuhkan database dapat ditambahkan terpisah dengan build tag atau test package khusus agar unit test tetap cepat.

## OpenAPI

Kontrak tersedia di `docs/openapi.yaml`. File ini dapat dibuka dengan Swagger UI, Redoc, atau tool OpenAPI lain. Setiap perubahan endpoint harus mengubah handler, unit test, dan kontrak OpenAPI dalam pull request yang sama.

## Production notes

- Ganti `DATABASE_URL` dengan secret manager, jangan commit `.env`.
- Tambahkan migration tool versioned ketika jumlah migration bertambah.
- Tambahkan authentication, authorization, request ID, metrics, tracing, dan rate limiting sesuai kebutuhan domain.
- Pisahkan database/schema per service jika batas ownership datanya sudah jelas.
