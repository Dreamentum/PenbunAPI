# PenbunAPI v3.1.0

RESTful API for managing distribution and supply of books and stationery — built with Go + Fiber + MSSQL.

---

## Features

- **Authentication** – JWT login/logout with bcrypt password hashing + token blacklist
- **16 Master Data Modules** – 8 standard CRUD functions per module (128 endpoints)
- **Global Error Handler** – Centralized `middleware.GlobalErrorHandler` for consistent JSON error responses
- **Transaction Safety** – `utils.ExecuteTransaction()` with panic recovery + step logging
- **Partial Updates** – `COALESCE(NULLIF(...))` pattern for PATCH-like PUT
- **Soft Delete** – `is_delete = 1` with `update_by` tracking
- **Pagination** – `?page=&limit=` on all Select Page endpoints
- **LIKE Search** – `?name=` with `%LIKE%` pattern for Select By Name
- **Consistent Response** – `{status, message, data}` across all endpoints
- **Audit Logging** – Transaction steps logged to `logs/transaction.log`
- **Graceful Shutdown** – Safe server stop on SIGINT/SIGTERM
- **Testing** – 109 unit tests with race detection

---

## Quick Start

```bash
# Prerequisites: Go 1.21+, MSSQL Server

git clone <repo> && cd PenbunAPI
go mod tidy

# Configure .env
cp .env.example .env
# Edit DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET

go run main.go
# Server starts on :8089 (configurable via FIBER_PORT)
```

---

## Project Structure

```
PenbunAPI/
├── main.go                    # Entry point + Fiber config + graceful shutdown
├── config/
│   ├── env.go                 # Environment variable loading
│   ├── database.go            # MSSQL connection pool
│   ├── blacklist.go           # Thread-safe token blacklist
│   └── logger.go              # Transaction log file
├── middleware/
│   ├── jwt.go                 # JWT Bearer validation
│   └── error.go               # Global error handler
├── controllers/
│   ├── auth.go                # Login / Logout
│   ├── publisher.go           # 8 CRUD functions
│   ├── publisherType.go
│   ├── customer.go
│   ├── customerType.go
│   ├── vendor.go
│   ├── vendorType.go
│   ├── book.go
│   ├── bookType.go
│   ├── discount.go
│   ├── discountType.go
│   ├── productFormatType.go
│   ├── productCategory.go
│   ├── productGroup.go
│   ├── unitType.go
│   └── warehouse.go
├── models/
│   ├── api.go                 # ApiResponse struct
│   ├── user.go                # User + LoginRequest
│   └── ...                    # 14 entity structs
├── routes/
│   ├── public.go              # /api/v1/public/*
│   ├── v1.go                  # /api/v1/protected/*
│   └── v2.go                  # Placeholder for future
├── utils/
│   ├── transaction.go         # ExecuteTransaction with rollback
│   └── response.go            # JSON response helpers
├── logs/                      # Transaction audit logs
├── docs/                      # Documentation
├── .env                       # Environment variables
├── *_test.go                  # 109 tests
└── go.mod
```

---

## API Reference

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/public/login` | Login, returns JWT |
| POST | `/api/v1/public/logout` | Blacklist token (requires auth) |

### Master Data (all protected)

Each module has 8 endpoints under `/api/v1/protected/{module}`:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/all` | Select all (is_delete = 0) |
| GET | `/page?page=1&limit=10` | Paginated results |
| GET | `/select/id/:id` | Select by code/ID |
| GET | `/select/name/:name` | LIKE search by name |
| POST | `/insert` | Create new record |
| PUT | `/update/:id` | Partial update |
| PUT | `/delete/:id?user=USER` | Soft delete |
| DELETE | `/remove/:id` | Hard delete |

#### Available Modules

| Module | Endpoint | Table |
|--------|----------|-------|
| Publisher | `/publisher/*` | `tb_publisher` |
| Publisher Type | `/publisher-type/*` | `tb_publisher_type` |
| Customer | `/customer/*` | `tb_customer` |
| Customer Type | `/customer-type/*` | `tb_customer_type` |
| Vendor | `/vendor/*` | `tb_vendor` |
| Vendor Type | `/vendor-type/*` | `tb_vendor_type` |
| Book | `/book/*` | `tb_book` |
| Book Type | `/book-type/*` | `tb_book_type` |
| Discount | `/discount/*` | `tb_discount` |
| Discount Type | `/discount-type/*` | `tb_discount_type` |
| Unit Type | `/unit-type/*` | `tb_unit_type` |
| Product Format Type | `/product-format-type/*` | `tb_product_format_type` |
| Product Category | `/product-category/*` | `tb_product_category` |
| Product Group | `/product-group/*` | `tb_product_group` |
| Warehouse | `/warehouse/*` | `tb_warehouse` |

---

## Response Format

```json
{
  "status": "success | fail | error",
  "message": "Human-readable message",
  "data": { ... }
}
```

**Login response** additionally includes `"token": "jwt..."`.

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | MSSQL host |
| `DB_PORT` | `1433` | MSSQL port |
| `DB_USER` | `sa` | Database user |
| `DB_PASSWORD` | — | Database password |
| `DB_NAME` | `PENBUN` | Database name |
| `FIBER_PORT` | `8089` | Server port |
| `JWT_SECRET` | `default-secret` | JWT signing key |
| `LOG_FILE` | `logs/transaction.log` | Transaction log path |

---

## Database Conventions

- `is_delete BIT` — soft delete (0=active, 1=deleted)
- `id_status BIT` — status (1=active, 0=inactive, default 1)
- `update_date` — auto-set by DB trigger (SE Asia Standard Time)
- Business codes (e.g. `PTG000001`) — auto-generated by DB trigger
- All queries use parameterized placeholders (`@p1`, `@p2`) — no SQL injection

---

## Testing

```bash
# Unit tests with race detection
go test -race -short ./...

# Coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Integration tests (requires live DB)
go test -tags=integration ./...
```

**Current:** 109 tests, 0 failed, 0 races across 6 packages.

---

## Libraries

| Library | Purpose |
|---------|---------|
| [Fiber v2](https://gofiber.io/) | Web framework |
| [go-mssqldb](https://github.com/denisenkom/go-mssqldb) | MSSQL driver |
| [Recover](https://docs.gofiber.io/api/middleware/recover) | Panic recovery middleware |
| [CORS](https://docs.gofiber.io/api/middleware/cors) | Cross-origin resource sharing |
| [golang-jwt v5](https://github.com/golang-jwt/jwt) | JWT auth |
| [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) | Password hashing |
| [godotenv](https://github.com/joho/godotenv) | .env loader |
| [testify](https://github.com/stretchr/testify) | Test assertions |

---

## Changelog

See [docs/CHANGELOG.md](docs/CHANGELOG.md).

## License

PENBUN License.
