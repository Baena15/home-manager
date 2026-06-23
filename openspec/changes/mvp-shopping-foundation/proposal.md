# Proposal: MVP Shopping Foundation

## Intent

Enable the two household users to log in, manage a product catalog with price history, and create shopping lists that calculate a closed total based on the latest known price.

## Scope

### In Scope
- Simple JWT auth for 2 fixed users (no registration).
- Product catalog: create, list, update, delete products.
- Price history per product (store, amount, recorded at).
- Shopping lists: create, list, add/remove items, mark as purchased.
- Closed price estimate for a list using the latest price of each product.
- Minimal PWA shell with pistachio theme and offline placeholder.

### Out of Scope
- OCR camera scanning.
- Bills and loans modules.
- Dashboard analytics.
- Multi-tenancy or invite flows.

## Approach

Use standard Go layout: `internal/handlers` for HTTP, `internal/store` for PostgreSQL with `database/sql`, `internal/config` for env vars, `web/` for PWA. Schema managed via migration files in `migrations/`. Frontend uses vanilla JS, custom CSS properties, and `fetch`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/config` | New | Load DATABASE_URL, JWT_SECRET, ENV |
| `internal/store` | New | ProductStore, PriceStore, ListStore, UserStore |
| `internal/handlers` | New | Auth, product, list HTTP handlers |
| `internal/middleware` | New | JWT validation middleware |
| `cmd/api/main.go` | Modified | Wire routes and dependencies |
| `web/` | New | PWA shell, login, products, lists pages |
| `migrations/` | New | PostgreSQL schema files |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Price fallback when no history exists | Med | Default to zero and show "sin precio" in UI |
| JWT secret too short in dev | Low | Validate length on startup, fail fast |
| Schema changes after first deploy | Low | Version migrations with sequential IDs |

## Rollback Plan

1. Stop the running service.
2. Revert to the previous git commit.
3. Run `make db-reset` only in dev; in prod, run down-migrations manually.

## Dependencies

- PostgreSQL 16 (local via Docker Compose or Railway).
- Redis 7 for session/cache (optional for MVP, wired for future use).

## Success Criteria

- [ ] Both users can log in and receive a JWT.
- [ ] Users can create products and see them listed.
- [ ] Adding a price to a product updates its latest price.
- [ ] Creating a list and adding products shows a closed total.
- [ ] `make test` passes with handler and store tests.
- [ ] PWA loads on iPhone and shows the login screen.
