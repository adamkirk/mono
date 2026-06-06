# Panoptes API — Claude Skills

## Meta — keeping this file up to date

Whenever the user points out something that was missed, wrong, or should be done differently, **immediately update this file** with the correct guidance so it won't happen again. Don't wait to be asked.

All project-specific guidance belongs here, not in external memory files outside the repo. Prefer updating this file over writing to `~/.claude/` so the knowledge is shared with all contributors.

## Stack

- **Language**: Go
- **HTTP framework**: [huma v2](https://github.com/danielgtaylor/huma) (OpenAPI-native)
- **Database**: PostgreSQL via `pgx/v5`
- **Query generation**: sqlc (`api/sqlc.yaml` → `internal/infra/repository/postgres/db/`)
- **UUID generation**: `github.com/google/uuid` — always UUID V7 (`uuid.NewV7()`)

---

## Directory layout

```
api/
  cmd/server/application.go       # Dependency wiring (lazy-init getters)
  etc/postgres/migrations/        # SQL migration files
  internal/
    app/                          # Application (business logic) layer
      main.go                     # Repository interfaces
      environments.go             # Environments handler + DTOs
      environment_components.go   # EnvironmentComponents handler + DTOs
    common/
      model.go                    # Shared domain models + IsValidSlug
      errors.go                   # Shared error types
      validation_error.go         # ValidationError + FieldError types
    infra/
      repository/postgres/
        db/                       # sqlc-generated code (DO NOT EDIT)
        queries.sql               # Source of truth for all SQL queries
        environments.go           # EnvironmentsRepository
        environment_components.go # EnvironmentComponentsRepository
      server/
        handlers.go               # Handler interfaces (used by controllers)
        error_handler.go          # ErrorHandler wrapper + buildValidationError
        v1beta.<resource>.go      # One file per resource
```

---

## Adding a new resource — checklist

### 1. Migration
Create `api/etc/postgres/migrations/NNNN_create_<table>.up.sql` and `.down.sql`.

Conventions:
- Table names are plural and snake_case: `environments`, `environment_components`
- Nested resources are prefixed with the parent: `environment_components` (not `components`)
- No `ON DELETE CASCADE` — handle cascades in code
- Primary key is always `UUID`
- Unique constraints are named `<table>_<columns>_unique`
- FK constraints are named `<table>_<column>_fk`

### 2. Model
Add to `internal/common/model.go`:
```go
type MyResource struct {
    ID   uuid.UUID
    Name string
    // ...
}
```

### 3. SQL queries
Add to `internal/infra/repository/postgres/queries.sql`, then run:
```
make api-sqlc-gen
```

Query conventions:
- Upsert by ID for Save: `ON CONFLICT (id) DO UPDATE SET ...` with `RETURNING *`
- Delete by ID (not name)
- Lookup by natural key for existence checks

### 4. Repository
Create `internal/infra/repository/postgres/<resource>.go` with a struct and methods. Always accept/return the common model type. Log errors before returning them. Follow the lazy-nil pattern for `pgx.ErrNoRows` → return `nil, nil`.

Add the repository interface to `internal/app/main.go`.

### 5. App handler
Create `internal/app/<resource>.go` with:
- DTOs (one per operation)
- `Validate` method on each DTO — see below
- Handler struct with injected repository interface(s)
- Handler methods that follow: **load → validate → mutate → save**

**DTO Validate conventions:**
- Every DTO **must** have a `Validate` method — never omit it
- Takes the repository as an argument when a DB lookup is needed, otherwise no args
- Call `dto.Validate(...)` at the top of the handler method, before any DB lookups
- Collects all field errors before returning (don't short-circuit)
- Skip DB uniqueness check if format is already invalid (use `else` branch)
- Returns `common.ValidationError{FieldErrors: ...}`
- Use `common.IsValidSlug` for name fields

**Handler conventions:**
- Load the entity first; return `common.ErrNotFound{}` if nil
- Call `dto.Validate(...)` after loading (so you have the ID if needed)
- Generate UUID V7 with `uuid.NewV7()` at creation time
- Save via the upsert-based `Save(model)` method

### 6. Server handler interface
Add the interface to `internal/infra/server/handlers.go`:
```go
type myResourceHandler interface {
    Create(dto app.CreateMyResourceDTO) (*common.MyResource, error)
}
```

### 7. Controller
Create `internal/infra/server/v1beta.<resource>.go`.

Conventions:
- One file per resource
- Route registration in `RegisterRoutes(v ApiVersion, api huma.API)`
- All handlers wrapped with `ErrorHandler(c.Method, http.MethodXxx)`
- Operation IDs follow dot-separated path hierarchy: `%s.environments.components.create`
- Request structs implement `MapErrorKey(string) string` for field-name mapping to JSON paths
- Responses use a `Body` field for the JSON payload and a `Status int` field
- 204 responses have no `Body` field

### 8. Wire up
In `cmd/server/application.go`:
- Add fields to `Application` struct
- Add lazy-init `Get<Resource>Repository()` and `Get<Resource>Handler()` methods
- Register the controller in `GetV1BetaControllers()`

---

## Error handling

Errors are mapped centrally in `internal/infra/server/error_handler.go`:

| Error type | HTTP status |
|---|---|
| `common.ValidationError` | 422 |
| `common.ErrNotFound` | 404 |
| `common.ErrUnauthorised` | 401 |
| `common.ErrConflict` (or anything wrapping it) | 409 |
| anything else | passed through as-is |

Add new error types to `internal/common/errors.go` and handle them in the switch.

---

## Pagination

List endpoints use page-based pagination (not offset/limit in the API):
- Query params: `page` (min 1, default 1) and `per_page` (min 1, max 100, default 20)
- Conversion to offset happens in the app handler: `offset = (page - 1) * perPage`
- Validate that `page` is not out of bounds after fetching the total count
- Response includes a `meta.pagination` object: `total`, `total_pages`, `page`, `per_page`

---

## Make targets

The `Makefile` is at the **repo root** (`/panoptes/Makefile`), not inside `api/`. Always run `make` from the repo root — running from `api/` produces "no rule to make target" errors.

| Target | Description |
|---|---|
| `make api-sqlc-gen` | Regenerate sqlc db package from `queries.sql` |

## Building / compilation

Do **not** run `go build` (or `go vet`, etc.) to verify compilation. A Docker watcher is always running and recompiles on file changes. The user will explicitly report any compilation errors — skip any build verification step after editing Go files.
