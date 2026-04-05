# Finance Data Processing API

Role-based finance backend using Go Fiber + PostgreSQL (pgx) with JWT auth, CRUD financial records, and dashboard analytics.

## Features

- JWT authentication
- RBAC with `viewer`, `analyst`, and `admin`
- Financial record CRUD with soft delete
- Filters by date/category/type + pagination
- Dashboard insights:
  - total income
  - total expenses
  - net balance
  - category totals
  - monthly trends
  - recent transactions

## Project Structure

- `cmd/server` app entrypoint
- `internal/handlers` HTTP layer
- `internal/services` business logic
- `internal/repositories` raw SQL with pgx
- `internal/middleware` auth and RBAC
- `internal/models` domain models
- `migrations` schema and indexes

## Setup

1. Copy environment template:
	- `cp .env.example .env`
2. Update `DATABASE_URL` and `JWT_SECRET`.
3. Apply SQL migrations in `migrations/` to your Postgres database.
4. Install dependencies:
	- `go mod tidy`
5. Run API:
	- `go run ./cmd/server`

## Environment Variables

- `APP_PORT` HTTP port
- `DATABASE_URL` postgres connection string
- `JWT_SECRET` JWT signing secret
- `JWT_ACCESS_TOKEN_TTL` token TTL (e.g. `1h`)
- `DEFAULT_PAGE_LIMIT` list default limit
- `MAX_PAGE_LIMIT` maximum list limit
- `DEFAULT_ADMIN_EMAIL` optional seeded admin email
- `DEFAULT_ADMIN_NAME` optional seeded admin name
- `DEFAULT_ADMIN_PASSWORD` optional seeded admin password

## API Endpoints

### Auth
- `POST /api/auth/login`

### Users (admin only)
- `POST /api/users`
- `GET /api/users`
- `PATCH /api/users/:id`
- `PATCH /api/users/:id/active`

### Records
- `GET /api/records` viewer/analyst/admin
- `GET /api/records/:id` viewer/analyst/admin (ownership enforced for non-admin)
- `POST /api/records` analyst/admin
- `PATCH /api/records/:id` analyst/admin (ownership enforced for non-admin)
- `DELETE /api/records/:id` analyst/admin (soft delete)

### Dashboard
- `GET /api/dashboard/summary`
- `GET /api/dashboard/categories`
- `GET /api/dashboard/trends`
- `GET /api/dashboard/recent`

Dashboard endpoints are available to viewer/analyst/admin; non-admin results are scoped to the caller.

## Assumptions

- Monetary values are stored as `NUMERIC(14,2)` and returned as strings to avoid floating point drift.
- Date filters expect RFC3339 format.
- `analyst` can create/update/delete records but cannot manage users.

## Test and Build

- `go build ./...`
- `go test ./...`