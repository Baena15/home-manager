# Tasks: MVP Shopping Foundation

## Phase 1: Foundation

- [x] 1.1 Create `migrations/001_initial_schema.sql` with users, products, product_prices, shopping_lists, shopping_list_items.
- [x] 1.2 Add `Makefile` target `db-migrate` to run migration files with psql.
- [x] 1.3 Implement `internal/config/config.go` to load DATABASE_URL, JWT_SECRET, ENV, and validate JWT length.
- [x] 1.4 Create `internal/store/db.go` with `NewDB()` helper and connection ping.
- [x] 1.5 Add bcrypt + JWT helpers in `pkg/auth/auth.go`.

## Phase 2: Auth

- [x] 2.1 Create `internal/store/user_store.go` with `GetByEmail()`.
- [x] 2.2 Create `internal/handlers/auth_handler.go` with `POST /api/v1/auth/login`.
- [x] 2.3 Create `internal/middleware/auth.go` with JWT validation middleware.
- [x] 2.4 Seed two users on first dev startup via `scripts/seed.sql` or startup function.
- [x] 2.5 Add auth handler tests for valid login, invalid password, and missing token.

## Phase 3: Products

- [x] 3.1 Create `internal/store/product_store.go` with CRUD and latest-price query.
- [x] 3.2 Create `internal/store/price_store.go` with `Create()` for product_prices.
- [x] 3.3 Create `internal/handlers/product_handler.go` with product and price endpoints.
- [x] 3.4 Add handler tests for create, duplicate name, list, and price endpoints.
- [x] 3.5 Wire product routes in `cmd/api/main.go`.

## Phase 4: Shopping Lists

- [x] 4.1 Create `internal/store/list_store.go` with list CRUD and item queries.
- [x] 4.2 Create `internal/store/list_item_store.go` with add, update, delete, and total calculation.
- [x] 4.3 Create `internal/handlers/list_handler.go` with list and item endpoints.
- [x] 4.4 Add handler tests for list creation, item addition with latest/custom price, and total.
- [x] 4.5 Wire list routes in `cmd/api/main.go`.

## Phase 5: Wiring & PWA Shell

- [x] 5.1 Wire all handlers and middleware in `cmd/api/main.go`.
- [x] 5.2 Serve static files from `web/` at `/`.
- [x] 5.3 Create `web/index.html`, `web/css/theme.css` (pistachio variables), `web/js/app.js`.
- [x] 5.4 Add `web/manifest.json` and `web/sw.js` for PWA installability.
- [x] 5.5 Implement login view and basic navigation between products and lists.

## Phase 6: Verification

- [x] 6.1 Run `go test ./...` and fix failures.
- [x] 6.2 Run `go vet ./...` (golangci-lint not installed locally).
- [x] 6.3 Build passes; local runtime verification blocked by unavailable Docker daemon/PostgreSQL.
- [x] 6.4 Created Dockerfile and railway.toml for production deploy.
