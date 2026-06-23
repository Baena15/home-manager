# AGENTS.md â€” Home Manager

> This file contains essential information for AI coding agents working on this project.
> All documentation and comments are primarily in **Spanish** for user-facing content, **English** for code internals.
>
> âš ï¸ **NOTE**: This project was generated from Plantilla Madre (Gentleman Stack).

---

## Project Overview

[Breve descripciÃ³n del proyecto]

| Aspecto | Detalle |
|---------|---------|
| **Stack** | Go 1.23 + chi + PostgreSQL 16 + Redis 7 |
| **PropÃ³sito** | [QuÃ© hace este proyecto] |
| **Estado** | [Active/Development/Planning] |
| **Idioma UI** | ðŸ‡ªðŸ‡¸ EspaÃ±ol |
| **Idioma CÃ³digo** | ðŸ‡¬ðŸ‡§ InglÃ©s |

---

## Technology Stack

| Capa | TecnologÃ­a |
|------|------------|
| **Backend** | Go 1.23+ |
| **Router** | [chi](https://github.com/go-chi/chi) |
| **Database** | PostgreSQL 16 |
| **Cache** | Redis 7 |
| **Auth** | JWT con bcrypt |
| **Frontend** | Vanilla JavaScript (zero frameworks) |
| **Styling** | CSS custom properties, Catppuccin theme |
| **DevOps** | Docker Compose, GitHub Actions |

---

## Directory Structure

```
home-manager/
â”œâ”€â”€ cmd/api/                 # Entry point (main.go)
â”œâ”€â”€ internal/
â”‚   â”œâ”€â”€ handlers/            # HTTP handlers
â”‚   â”œâ”€â”€ middleware/          # Middleware
â”‚   â”œâ”€â”€ store/               # Data access layer
â”‚   â””â”€â”€ config/              # Configuration
â”œâ”€â”€ pkg/                     # Public/reusable code
â”œâ”€â”€ web/                     # Frontend assets (if applicable)
â”œâ”€â”€ scripts/
â”‚   â”œâ”€â”€ seed.ps1             # Database seeding
â”‚   â”œâ”€â”€ test-api.ps1         # API endpoint testing
â”‚   â””â”€â”€ db-reset.ps1         # Database reset
â”œâ”€â”€ docs/
â”‚   â””â”€â”€ CONVENTIONS.md       # Coding conventions
â”œâ”€â”€ .claude/skills/          # 17+ SDD skills
â”œâ”€â”€ .github/
â”‚   â”œâ”€â”€ workflows/ci.yml     # CI/CD pipeline
â”‚   â”œâ”€â”€ pull_request_template.md
â”‚   â””â”€â”€ ISSUE_TEMPLATE/      # Issue forms
â”œâ”€â”€ .vscode/                 # VS Code settings (optional)
â”œâ”€â”€ Dockerfile.dev           # Dev environment
â”œâ”€â”€ docker-compose.yml       # PostgreSQL + Redis + API
â”œâ”€â”€ Makefile                 # Standard commands
â”œâ”€â”€ .env.example             # Environment variables
â”œâ”€â”€ .gitignore
â”œâ”€â”€ SECURITY.md              # Pre-deploy security checklist
â”œâ”€â”€ .cursorrules             # Cursor IDE rules
â””â”€â”€ AGENTS.md                # This file
```

---

## Convenciones de CÃ³digo

### Idiomas
- **UI/User-facing**: EspaÃ±ol
- **CÃ³digo/Variables/Comentarios**: InglÃ©s
- **Git Commits**: InglÃ©s (Conventional Commits)
- **DocumentaciÃ³n tÃ©cnica**: EspaÃ±ol

### Estilo RÃ¡pido
```go
// â”€â”€â”€ Section Name â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
func FunctionName() error {
    // errors en lowercase sin punto final
    return fmt.Errorf("operation failed: %w", err)
}
```

Ver `docs/CONVENTIONS.md` para convenciones completas:
- Go, JavaScript, CSS
- SQL y base de datos
- API responses (formato JSON estÃ¡ndar)
- Conventional Commits detallado
- Naming conventions
- Testing patterns

---

## Comandos RÃ¡pidos

```bash
make dev          # Iniciar servidor de desarrollo
make test         # Ejecutar tests
make test-watch   # Tests en watch mode
make coverage     # Reporte de cobertura
make build        # Compilar binario
make lint         # Ejecutar linter
make fmt          # Formatear cÃ³digo
make db-seed      # Poblar base de datos
make db-reset     # Resetear base de datos
make setup        # Setup inicial (copia .env.example a .env)
make clean        # Limpiar binarios/cache
make security-check # Checklist pre-deploy
```

### Docker Compose
```bash
docker-compose up -d          # Iniciar PostgreSQL + Redis + API
docker-compose logs -f api    # Ver logs
docker-compose down           # Detener todo
```

---

## Variables de Entorno CrÃ­ticas

| Variable | DescripciÃ³n | Requerida |
|----------|-------------|-----------|
| `DATABASE_URL` | PostgreSQL connection string | âœ… SÃ­ |
| `JWT_SECRET` | JWT signing secret (>=32 chars) | âœ… SÃ­ |
| `JWT_EXPIRATION_HOURS` | Token expiration (default: 24) | âš ï¸ Opcional |
| `API_URL` | URL base de la API | âš ï¸ Opcional |
| `PORT` | Puerto del servidor (default: 8080) | âš ï¸ Opcional |
| `ENV` | Entorno (development/production) | âš ï¸ Opcional |
| `EMAIL_PROVIDER` | Proveedor de email | âš ï¸ Opcional |
| `EMAIL_API_KEY` | API key de email | âš ï¸ Opcional |
| `EMAIL_FROM` | Email remitente | âš ï¸ Opcional |

---

## SDD Workflow

Usamos Spec-Driven Development con 17+ skills:

| Fase | Skill | PropÃ³sito |
|------|-------|-----------|
| 1 | **`/sdd-init`** | Inicializar proyecto |
| 2 | **`/sdd-explore`** | Investigar/explorar antes de cambiar |
| 3 | **`/sdd-spec`** | Especificar cambios |
| 4 | **`/sdd-tasks`** | Dividir en tareas |
| 5 | **`/sdd-apply`** | Implementar |
| 6 | **`/sdd-verify`** | Verificar |
| 7 | **`/sdd-archive`** | Archivar |

**Workflows adicionales:**
- `/sdd-propose` â€” Crear propuestas de cambio
- `/sdd-design` â€” Documentos de diseÃ±o tÃ©cnico
- `/sdd-debug` â€” Debugging sistemÃ¡tico
- `/sdd-refactor` â€” RefactorizaciÃ³n segura
- `/issue-creation` â€” Crear issues en GitHub
- `/branch-pr` â€” Crear pull requests
- `/judgment-day` â€” Code review adversarial

### Persistence Modes
- **`engram`** â€” Persistent memory (cross-session)
- **`openspec`** â€” Filesystem-based (versionable)
- **`hybrid`** â€” Ambos
- **`none`** â€” EfÃ­mero

---

## Engram

Este proyecto usa **engram** para memoria persistente entre sesiones:

```bash
engram save "TÃ­tulo" "DescripciÃ³n"
engram search "tÃ©rmino"
engram context home-manager
```

---

## Seguridad

Revisar `SECURITY.md` antes de cada deploy:

- [ ] Secrets no commiteados, `.env` en `.gitignore`
- [ ] JWT_SECRET >= 32 caracteres, expiraciÃ³n <= 24h
- [ ] Passwords hasheadas (bcrypt/Argon2)
- [ ] CORS configurado, rate limiting activo
- [ ] SQL solo parametrizado
- [ ] HTTPS en producciÃ³n
- [ ] Headers de seguridad (HSTS, CSP, X-Frame-Options)

---

## Checklist de Inicio

- [ ] Personalizar este archivo (`AGENTS.md`)
- [ ] Configurar `.env` con variables necesarias
- [ ] Ejecutar `make setup` para inicializar
- [ ] Probar `make dev` (servidor en :8080)
- [ ] Probar `make test`
- [ ] Configurar engram: `engram setup`
- [ ] Revisar `docs/CONVENTIONS.md` con el equipo
- [ ] Primer commit con cambios especÃ­ficos del proyecto

---

## Contacto / Recursos

- DocumentaciÃ³n: [links]
- Repositorio: [link]
- Issues: [link]

---

*Generado desde Plantilla Madre â€” Gentleman Programming* ðŸ¤
