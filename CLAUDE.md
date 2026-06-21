# CLAUDE.md

Guidance for Claude Code when working in this repository.

## What this is

The **Go + Gin + GORM** REST API powering **Todo Manager** — a two-service portfolio project. It owns auth, persistence, and business logic. The client is a separate React SPA ([todo-app-react](https://github.com/YawDev/todo-app-react)) served from `https://todo-manager.app`. This API is served from `https://api.todo-manager.app`, base path `/api/v1`.

For local development the API runs at `http://localhost:8080` and the SPA at `http://localhost:5173`.

## Commands

```bash
go mod download
go run .            # start server (localhost:8080)
go test ./...       # run all tests
swag init           # regenerate Swagger docs after changing handler annotations
```

Default config uses SQLite (`todo.db`), so no DB setup is needed to run.

## Architecture

Layered, with interface-based storage selected at startup:

- **`main.go`** → builds config + logger → `server.NewService(...).Start(r)`.
- **`server/service.go`** wires CORS, Swagger, the DB connection, and registers routes. **All routes are defined here** in `RouteSetup` — public group vs. `AuthMiddleware`-protected group.
- **`controllers/`** parse/validate requests and write responses. They depend on the **storage interfaces** (`s.UserManager`, `s.ListManager`, `s.TaskManager`), never on a concrete driver.
- **`storage/database.go`** defines the interfaces (`IUserManager`, `IListManager`, `ITaskManager`, `IDatabase`) and `ConfigureDb(useSQLite)` picks the implementation set:
  - `storagelite/` = SQLite impls (default)
  - `storage/` = MySQL impls
- **`authentication/jwt.go`** issues/parses JWTs and holds the **in-memory active-token maps** (`ActiveTokens`, `refreshTokens`, mutex-guarded).
- **`models/models.go`** = GORM models (User, List, Task). Auto-migrated on connect.
- **`helpers/api_helpers.go`** = request/response DTOs (with `binding` tags and Swagger `example` tags).
- **`loggerutils/`** = logrus + context-aware helpers (`InfoLog`/`ErrorLog` pull the request id from context).
- **`messages/messages.go`** = centralized message/error strings — reuse these, don't inline literals.

## Conventions

- **Adding an endpoint:** add the handler in `controllers/`, register it in `server/service.go` (`RouteSetup`), and add Swagger annotations (the `// @Summary`, `// @Router` comment block) — then run `swag init`.
- **Adding a storage method:** add it to the interface in `storage/database.go` AND implement it in **both** `storage/` (MySQL) and `storagelite/` (SQLite), or the build breaks / a driver silently lacks it.
- Error responses use the typed structs in `helpers/` (`BadRequestResponse`, `ErrorResponse`, `NotFoundResponse`, …) with explicit status codes. Match this pattern.
- Log via `loggerutils.ErrorLog(ctx, status, err)` / `InfoLog(...)` so request-id correlation works; pass `c.Request.Context()`.
- Passwords are bcrypt-hashed (`controllers.Hash`). Never store or compare plaintext.

## Auth model (important)

- JWT (HS256) access token (30 min) + refresh token (1 hr), delivered as **HttpOnly, Secure, SameSite=None cookies**.
- A token is valid only while its username is present in the in-memory `ActiveTokens` map — Logout deletes it. This means **tokens reset on server restart** and won't work across multiple instances.
- `AuthMiddleware` accepts the token from `Authorization: Bearer …` header OR the `access_token` cookie, then injects `user_id`/`username` into the Gin context.

## Known shortcuts (don't treat as production-ready)

- `config.yaml` has DB credentials and `authentication/jwt.go` has a hardcoded `jwtKey` — both should move to env/secrets.
- In-memory session store (see above).
- Handlers trust path ids without verifying the caller owns the resource.

If asked to "fix" or "harden," these are the priorities. Otherwise leave behavior as-is.

## Database switching

`config.yaml` → `database.useSQLite: true|false`. No code change needed; `ConfigureDb` swaps the implementation set at startup.
