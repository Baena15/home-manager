# Design: MVP Shopping Foundation

## Database Schema

### users
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| email | VARCHAR(255) UNIQUE | |
| password_hash | VARCHAR(255) | bcrypt |
| role | VARCHAR(50) | `owner` or `partner` |
| created_at | TIMESTAMPTZ | default now() |

### products
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| name | VARCHAR(255) UNIQUE | |
| unit | VARCHAR(50) | e.g. `g`, `kg`, `unit` |
| category | VARCHAR(100) | optional |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

### product_prices
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| product_id | UUID FK → products | cascade delete |
| store | VARCHAR(100) | e.g. `Mercadona` |
| amount | DECIMAL(10,2) | |
| recorded_at | TIMESTAMPTZ | default now() |

### shopping_lists
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| name | VARCHAR(255) | |
| status | VARCHAR(50) | `active` or `completed` |
| created_by | UUID FK → users | |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

### shopping_list_items
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| list_id | UUID FK → shopping_lists | cascade delete |
| product_id | UUID FK → products | |
| quantity | DECIMAL(10,3) | |
| unit_price | DECIMAL(10,2) | snapshot at add time |
| total | DECIMAL(10,2) | quantity * unit_price |
| purchased | BOOLEAN | default false |
| created_at | TIMESTAMPTZ | |

## API Endpoints

### Auth
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/login` | Returns JWT |

### Products
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/products` | Create product |
| GET | `/api/v1/products` | List products |
| GET | `/api/v1/products/{id}` | Get product |
| PUT | `/api/v1/products/{id}` | Update product |
| DELETE | `/api/v1/products/{id}` | Delete product |
| POST | `/api/v1/products/{id}/prices` | Add price record |

### Shopping Lists
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/lists` | Create list |
| GET | `/api/v1/lists` | List lists |
| GET | `/api/v1/lists/{id}` | Get list with items |
| POST | `/api/v1/lists/{id}/items` | Add item |
| PATCH | `/api/v1/lists/{id}/items/{item_id}` | Update item (purchased) |
| DELETE | `/api/v1/lists/{id}/items/{item_id}` | Remove item |

## Architecture Decisions

1. **No ORM**: Use `database/sql` + small typed helpers to keep dependencies minimal.
2. **Price snapshot**: `shopping_list_items.unit_price` stores the price at add time to keep list totals stable even if product prices change later.
3. **Latest price helper**: A SQL CTE or subquery returns the most recent `product_prices.amount` per product.
4. **Seed users**: `cmd/api/main.go` calls a seed function on startup if `ENV=development` and users table is empty.
5. **PWA minimal**: Single `index.html` with dynamic views; no framework to reduce complexity.

## Frontend Views

- `/` → Login
- `/products` → Product catalog + add price
- `/lists` → Shopping lists overview
- `/lists/{id}` → List detail + add items + total

## Security

- Passwords hashed with bcrypt.
- JWT signed with HS256; secret validated >= 32 chars on startup.
- All product/list endpoints require valid JWT.
- CORS restricted to the frontend origin in production.
