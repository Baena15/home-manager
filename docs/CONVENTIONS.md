# Convenciones de Código - Gentleman Stack

## 🌍 Idiomas

| Contexto | Idioma | Ejemplo |
|----------|--------|---------|
| UI / User-facing | **Español** | `<h1>Bienvenido</h1>` |
| Código / Variables | **Inglés** | `const userName = "..."` |
| Comentarios | **Inglés** | `// Initialize user session` |
| Commits | **Inglés** | `feat: add user authentication` |
| Docs técnicas | **Español** | Este archivo |

---

## 📁 Estructura de Archivos

```
proyecto/
├── cmd/                     # Entry points (Go)
├── internal/                # Código privado
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Middleware
│   ├── models/              # Modelos de datos
│   ├── store/               # Acceso a datos
│   └── config/              # Configuración
├── pkg/                     # Código público/reusable
├── web/                     # Frontend (si aplica)
├── scripts/                 # Scripts de utilidad
├── docs/                    # Documentación
└── AGENTS.md                # Contexto para agentes
```

---

## 🎨 Estilo de Código

### Go
```go
// ─── Package Description ─────────────────
package handlers

import (
    "net/http"
)

// UserHandler maneja operaciones de usuario
type UserHandler struct {
    store  *store.Store
    config *config.Config
}

// NewUserHandler crea un nuevo handler
func NewUserHandler(s *store.Store, cfg *config.Config) *UserHandler {
    return &UserHandler{
        store:  s,
        config: cfg,
    }
}

// GetUser obtiene un usuario por ID
// Returns: User, error
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    user, err := h.store.GetUser(id)
    if err != nil {
        // Error messages: lowercase, no period
        respondWithError(w, http.StatusNotFound, "user not found")
        return
    }
    
    respondWithJSON(w, http.StatusOK, user)
}
```

### JavaScript
```javascript
// ─── Configuración API ────────────────────
const API_CONFIG = {
    baseUrl: 'https://api.ejemplo.com',
    endpoints: {
        users: '/api/v1/users',
        auth: {
            login: '/api/v1/auth/login',
            register: '/api/v1/auth/register'
        }
    }
};

/**
 * Obtiene usuario por ID
 * @param {string} id - User ID
 * @returns {Promise<User>}
 */
async function getUser(id) {
    const response = await fetch(`${API_CONFIG.baseUrl}/users/${id}`);
    if (!response.ok) {
        throw new Error('user not found');
    }
    return response.json();
}
```

### CSS
```css
/* ─── Variables Globales ───────────────── */
:root {
    --color-primary: #6366f1;
    --color-secondary: #a855f7;
    --spacing-unit: 1rem;
}

/* ─── Componente Card ───────────────── */
.card {
    padding: var(--spacing-unit);
    border-radius: 8px;
}

/* Mobile-first responsive */
@media (min-width: 768px) {
    .card {
        padding: calc(var(--spacing-unit) * 2);
    }
}
```

---

## 📝 Comentarios

### Secciones (obligatorio)
```go
// ─── Nombre de Sección ─────────────────
```

### Funciones
```go
// FunctionName hace X
// Parameters: ...
// Returns: ...
```

### Decisiones importantes
```go
// NOTE: Decidimos usar X en lugar de Y porque Z
// Ver: AGENTS.md#decisiones
```

---

## 🏷️ Naming Conventions

| Tipo | Convención | Ejemplo |
|------|------------|---------|
| Variables | camelCase | `userName`, `totalCount` |
| Funciones | PascalCase (export) / camelCase (private) | `GetUser()`, `validateInput()` |
| Structs/Types | PascalCase | `UserHandler`, `Config` |
| Constants | UPPER_SNAKE_CASE | `MAX_RETRY_COUNT` |
| Archivos | snake_case | `user_handler.go` |
| Paquetes | lowercase | `handlers`, `middleware` |

---

## 🧪 Testing

```go
// TestGetUser prueba obtener usuario
func TestGetUser(t *testing.T) {
    tests := []struct {
        name       string
        id         string
        wantStatus int
    }{
        {
            name:       "user exists",
            id:         "123",
            wantStatus: http.StatusOK,
        },
        {
            name:       "user not found",
            id:         "999",
            wantStatus: http.StatusNotFound,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

---

## 🗄️ SQL / Base de Datos

### Naming
- Tablas: plural, snake_case — `users`, `book_loans`
- Columnas: snake_case — `created_at`, `password_hash`
- Claves foráneas: `{tabla}_id` — `user_id`, `book_id`
- Índices: `idx_{tabla}_{columna}` — `idx_users_email`
- Migraciones: `{timestamp}_{descripcion}.sql` o usar herramienta (golang-migrate, goose)

### Queries
- **SIEMPRE** parametrizadas. Nunca concatenar strings de usuario.
- Preferir `RETURNING id` en INSERTs para evitar round-trips.
- Usar `COALESCE` para valores por defecto en consultas.
- Documentar índices en `indexes.sql` o migrations.

```go
// ✅ Bien
rows, err := db.QueryContext(ctx, "SELECT * FROM users WHERE email = $1", email)

// ❌ Mal — SQL Injection
rows, err := db.QueryContext(ctx, "SELECT * FROM users WHERE email = '"+email+"'")
```

### Migrations
- Una migración = un cambio atómico (una tabla o un índice).
- Incluir `DOWN` / `DROP` para rollback.
- Probar migraciones en una BD limpia antes de commit.

---

## 🌐 API Responses

### Formato estándar de error
```json
{
  "error": "user not found",
  "code": "NOT_FOUND",
  "status": 404
}
```

### Formato estándar de éxito (listas)
```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 150
  }
}
```

### Formato estándar de éxito (single resource)
```json
{
  "data": { "id": "123", "name": "..." }
}
```

### Códigos HTTP
| Escenario | Status |
|-----------|--------|
| Éxito | 200 OK / 201 Created |
| Validación fallida | 400 Bad Request |
| Autenticación fallida | 401 Unauthorized |
| Sin permisos | 403 Forbidden |
| Recurso no existe | 404 Not Found |
| Conflicto (duplicado) | 409 Conflict |
| Error servidor | 500 Internal Server Error |

### Headers obligatorios
- `Content-Type: application/json`
- `X-Request-ID` para tracing (si aplica)

---

## 📝 Conventional Commits

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Tipos (`type`)
| Tipo | Uso |
|------|-----|
| `feat` | Nueva funcionalidad |
| `fix` | Corrección de bug |
| `docs` | Cambios en documentación |
| `style` | Formato, punto y coma, etc. (no cambia lógica) |
| `refactor` | Refactorización de código |
| `test` | Agregar o corregir tests |
| `chore` | Tareas de mantenimiento, deps, build |
| `perf` | Mejora de performance |
| `security` | Fix de seguridad |

### Scopes (`scope`) comunes
- `api`, `auth`, `db`, `ui`, `config`, `ci`, `test`

### Ejemplos
```
feat(auth): add password reset endpoint

fix(db): correct foreign key constraint on loans table

refactor(store): extract user queries to prepared statements

docs(readme): update setup instructions for Windows
```

---

## 🔒 Seguridad

- Nunca commitear secrets (usar .env)
- Sanitizar input de usuarios
- Usar HTTPS en producción
- Rate limiting en endpoints públicos
- JWT con expiración corta (24h default)
- SQL injection: solo queries parametrizadas
- XSS: escapar output HTML

---

## 📚 Recursos

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Clean Code](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)

---

*Estas convenciones son nuestra "fuente de verdad" para mantener consistencia entre proyectos.*
