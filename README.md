# Palge

A production-oriented banking backend API built with Go's standard library and PostgreSQL.

Palge is an incremental engineering project focused on applying real-world backend practices — from structured error handling and token-based authentication to database transactions with row-level locking.

> **Version:** 0.1.0 · **Status:** In Development

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26 |
| HTTP | `net/http` (standard library, Go 1.22+ routing) |
| Database | PostgreSQL 16 (via `pgx` driver) |
| Auth | Bearer token authentication (SHA-256 hashed) |
| Email | SMTP via `go-mail` |
| Password | bcrypt (cost 12) |
| Infrastructure | Docker Compose |

---

## Getting Started

### Prerequisites

- Go 1.26+
- Docker & Docker Compose
- `make` (optional, for convenience commands)

### 1. Clone the repository

```bash
git clone git@github.com:nongpal/Palge-Backend.git
cd Palge-Backend
```

### 2. Start the database

```bash
docker compose up -d
```

This starts a PostgreSQL 16 instance on port `5433` with a persistent volume.

### 3. Configure environment variables

Create a `.env` file in the project root:

```env
PORT=4000
ENV=development
DB_DSN=postgres://postgres:postgres@localhost:5433/palge?sslmode=disable

SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password
SMTP_SENDER="Palge <no-reply@example.com>"
```

### 4. Run the application

```bash
make run
```

The server starts at `http://localhost:4000`.

---

## Configuration

| Variable | Default | Description |
|---|---|---|
| `PORT` | `4000` | HTTP server port |
| `ENV` | `development` | Environment name (`development`, `staging`, `production`) |
| `DB_DSN` | `postgres://postgres:postgres@localhost:5433/palge?sslmode=disable` | PostgreSQL connection string |
| `SMTP_HOST` | `sandbox.smtp.mailtrap.io` | SMTP server host |
| `SMTP_PORT` | `2525` | SMTP server port |
| `SMTP_USERNAME` | — | SMTP authentication username |
| `SMTP_PASSWORD` | — | SMTP authentication password |
| `SMTP_SENDER` | `Palge <no-reply@github.com/nongpal/Palge-Backend>` | Sender email address |

---

## API Endpoints

All endpoints are prefixed with `/v1`.

### Public

| Method | Path | Description |
|---|---|---|
| `GET` | `/v1/healthcheck` | Server health and version info |
| `POST` | `/v1/users` | Register a new user (sends activation email) |
| `PUT` | `/v1/users/activated` | Activate account via token |
| `POST` | `/v1/tokens/authentication` | Authenticate and receive a bearer token |

### Protected

All protected endpoints require a valid `Authorization: Bearer <token>` header and an activated account.

| Method | Path | Permission | Description |
|---|---|---|---|
| `POST` | `/v1/accounts` | `accounts:create` | Create a new bank account |
| `GET` | `/v1/accounts` | `accounts:read` | List all accounts |
| `GET` | `/v1/accounts/{id}` | `accounts:read` | Get account by ID |
| `POST` | `/v1/accounts/{id}/deposit` | `accounts:deposit` | Deposit funds |
| `POST` | `/v1/accounts/{id}/withdraw` | `accounts:withdraw` | Withdraw funds |
| `POST` | `/v1/transfers` | `accounts:transfer` | Transfer funds between accounts |

### Error Responses

All errors return JSON:

```json
{
  "error": "Description of what went wrong"
}
```

Common status codes: `400`, `401`, `403`, `404`, `409`, `422`, `429`, `500`.

---

## Project Structure

```
cmd/api/              Application entrypoint
internal/
  api/                HTTP layer (handlers, middleware, routing, server)
  data/               Data access layer (models, database queries)
  mailer/             Email sending and templates
  validator/          Request validation framework
migrations/           PostgreSQL migration SQL files
```

### Key Design Decisions

- **Standard library only** — no web framework. Routes use Go 1.22+ method-pattern syntax.
- **Token-based auth** — tokens are SHA-256 hashed before storage; plaintext returned only once at creation.
- **Database transactions** — fund transfers use `SELECT ... FOR UPDATE` row locking to prevent race conditions.
- **RBAC** — five granular permissions (`accounts:read`, `create`, `deposit`, `withdraw`, `transfer`) enforced per endpoint.
- **Audit logging** — every significant action is recorded with user, action, resource, IP, user agent, and result.
- **Rate limiting** — token bucket algorithm, IP-based, with proxy header support (`X-Forwarded-For`, `X-Real-IP`).

---

## Middleware Pipeline

Requests pass through the following middleware chain:

1. **Request ID** — extracts or generates a UUID (`X-Request-ID` header)
2. **Request Logger** — logs method, URI, status, duration, request ID, client IP
3. **Rate Limiter** — token bucket (1 token/sec refill, capacity 10)
4. **Authenticate** — validates bearer tokens, injects user into context

---

## Database

PostgreSQL 16 with six migrations:

| Table | Purpose |
|---|---|
| `accounts` | Bank accounts with owner and balance |
| `users` | User accounts with hashed passwords |
| `tokens` | Authentication and activation tokens (SHA-256 hashes) |
| `permissions` | Permission codes for RBAC |
| `users_permissions` | Many-to-many user-permission mapping |
| `audit_logs` | Action audit trail with IP, user agent, and result |

Connection pool: 25 max open, 25 max idle, 15 min idle timeout.

---

## Make Commands

| Command | Description |
|---|---|
| `make run` | Run the server |
| `make build` | Build binary to `./bin/palge` |
| `make test` | Run all tests |
| `make tidy` | Clean up Go module dependencies |

---

## Roadmap

- [x] v0.1.0 — Foundation
- [x] v0.2.0 — Core Banking API
- [x] v0.3.0 — PostgreSQL Persistence
- [x] v0.4.0 — Authentication
- [ ] v0.5.0 — Production Readiness
- [ ] v0.6.0 — Performance
- [ ] v0.7.0 — Observability
- [ ] v0.8.0 — Distributed Architecture

---

## License

This project is for learning purposes.
