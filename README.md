# Yorkshire API

Backend RESTful API service built with **Go**, following a **Contract-First (OpenAPI 3.0)** and **Type-Safe SQL (`sqlc`)** architecture.

---

## Tech Stack

- **Language**: Go `1.27+`
- **HTTP Router**: [Chi v5](https://github.com/go-chi/chi)
- **API Contract & Codegen**: OpenAPI 3.0.3 + [Redocly CLI](https://redocly.com/docs/cli) + [oapi-codegen v2](https://github.com/oapi-codegen/oapi-codegen)
- **Database**: PostgreSQL 15 + [pgx/v5](https://github.com/jackc/pgx) + [sqlc](https://sqlc.dev/)
- **Database Migrations**: [Goose v3](https://github.com/pressly/goose)
- **Cache**: Redis 7.4 ([go-redis/v9](https://github.com/redis/go-redis))
- **Authentication**: JWT RS256 ([golang-jwt/v5](https://github.com/golang-jwt/jwt)) & Bcrypt password hashing
- **Validation**: [go-playground/validator v10](https://github.com/go-playground/validator)
- **Logging**: Structured JSON logging (`log/slog`)

---

## Project Structure

```text
.
├── cmd/
│   └── server/
│       └── main.go                 # Application entrypoint & graceful shutdown
├── contract/                       # OpenAPI 3.0.3 multi-file specification
│   ├── openapi.yaml                # Root OpenAPI spec
│   ├── paths/                      # Endpoint path definitions
│   └── components/                 # Reusable schemas, responses, parameters, securitySchemes
├── files/                          # RSA key pair for JWT RS256 (private_key.pem, public_key.pem)
├── internal/
│   ├── api/                        # Generated Chi server interface & DTO models (oapi-codegen)
│   ├── apperror/                   # Application error types & PostgreSQL error translator
│   ├── cache/                      # Redis cache wrapper
│   ├── config/                     # Environment configuration loader
│   ├── database/
│   │   ├── postgres.go             # PostgreSQL connection pool (pgxpool)
│   │   ├── queries/                # SQL queries for sqlc
│   │   └── sqlc/                   # Generated type-safe Go database code (sqlc)
│   ├── handler/                    # HTTP handlers implementing api.ServerInterface
│   ├── jwt/                        # RS256 JWT token generator & validator
│   ├── logging/                    # Multi-writer JSON logger (stdout + file)
│   ├── response/                   # Standardized HTTP JSON response helpers
│   ├── security/                   # Bcrypt password hashing utilities
│   ├── server/                     # Chi router & middleware configuration
│   ├── service/                    # Business logic layer (e.g., auth service)
│   └── validator/                  # Struct validation & user-friendly error formatting
├── migrations/                     # SQL schema migrations (Goose)
├── docker-compose.yaml             # Local PostgreSQL & Redis services
├── oapi-codegen.yaml               # oapi-codegen configuration
├── sqlc.yaml                       # sqlc configuration
└── makefile                        # Development automation tasks
```

---

## Prerequisites

- **Go** `1.27+`
- **Docker** & **Docker Compose** (for PostgreSQL & Redis)
- **Node.js / npx** (used by `@redocly/cli` to bundle OpenAPI specs and build docs)
- **OpenSSL** (for generating RSA private & public keys)

---

## Getting Started

### 1. Generate RSA Key Pair (Private & Public Key)

The authentication service uses **RS256 (RSA Signature with SHA-256)** to sign and verify JWT tokens. Because `*.pem` files are ignored by Git, you need to generate the RSA key pair first:

```bash
mkdir -p files

# 1. Generate 2048-bit RSA Private Key
openssl genpkey -algorithm RSA -out files/private_key.pem -pkeyopt rsa_keygen_bits:2048

# 2. Extract Public Key from the Private Key
openssl rsa -pubout -in files/private_key.pem -out files/public_key.pem
```

### 2. Configure Environment Variables

Copy `.env.example` to `.env` and adjust the values as needed:

```bash
cp .env.example .env
```

Example `.env` configuration aligned with `internal/config/config.go`:

```dotenv
HTTP_ADDR=:8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/yorkshires?sslmode=disable
REDIS_URL=localhost:6379
REDIS_PASSWORD=secretpassword
REDIS_DB=0
LOG_FILE=logs/app.log
PRIVKEY_PATH=files/private_key.pem
PUBKEY_PATH=files/public_key.pem
JWT_TOKEN_DURATION=60m
JWT_ISSUER=yorkshire-service
JWT_AUDIENCE=yorkshire-client
```

### 3. Start Infrastructure (PostgreSQL & Redis)

```bash
REDIS_PASSWORD=secretpassword docker compose up -d
```

### 4. Run Database Migrations

```bash
make goose-up
```

### 5. Generate Code (OpenAPI & SQLC)

```bash
make oapi-generate
make sqlc-generate
```

### 6. Run the Server

```bash
make run
```

The server will start on `http://localhost:8080` (base API path: `/api/v1`).

---

## Makefile Commands

| Command                 | Description                                                              |
| :---------------------- | :----------------------------------------------------------------------- |
| `make run`              | Run the HTTP server (`go run ./cmd/server/main.go`)                      |
| `make oapi-bundle`      | Bundle multi-file OpenAPI spec into `tmp/openapi.bundled.yaml`           |
| `make oapi-generate`    | Bundle OpenAPI spec and generate Go code into `internal/api/api.gen.go`  |
| `make oapi-docs`        | Build and open interactive Redoc HTML documentation (`tmp/docs.html`)    |
| `make sqlc-generate`    | Generate type-safe Go database queries from `internal/database/queries/` |
| `make migration <name>` | Create a new SQL migration file in `migrations/`                         |
| `make goose-up`         | Apply all pending database migrations                                    |

---

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

| Method | Endpoint                | Description                                             | Auth   |
| :----- | :---------------------- | :------------------------------------------------------ | :----- |
| `POST` | `/api/v1/auth/register` | Register a new user account                             | Public |
| `POST` | `/api/v1/auth/login`    | Authenticate user and receive an RS256 JWT access token | Public |

### Standard Response Format

**Success Response (`200 OK` / `201 Created`)**

```json
{
  "status": true,
  "data": {
    "user": {
      "id": "01972d80-7c7e-7bb3-94dd-1f873e0b43f6",
      "username": "alexplaju",
      "email": "alexplaju@gmail.com",
      "avatar_url": null,
      "created_at": "2026-11-01T10:00:00Z",
      "updated_at": "2026-11-01T10:00:00Z"
    },
    "token": {
      "access_token": "eyJhbGciOiJSUzI1NiIs...",
      "token_type": "Bearer",
      "expires_in": 3600
    }
  }
}
```

**Validation Error Response (`422 Unprocessable Entity`)**

```json
{
  "status": false,
  "message": "validation failed",
  "errors": {
    "email": "Invalid email format",
    "password": "Minimum length is 8"
  }
}
```

**General Error Response (`400`, `401`, `404`, `409`, `500`)**

```json
{
  "status": false,
  "message": "invalid email or password"
}
```
