# AGENTS.md — Home Manager

> Este archivo contiene información esencial para agentes de código AI que trabajan en Home Manager.
> Documentación y UI en **español**. Código, variables y commits en **inglés**.
>
> ⚠️ **NOTA**: Proyecto generado desde Plantilla Madre (Gentleman Stack) y personalizado para uso doméstico.

---

## Project Overview

**Home Manager** es una aplicación web/PWA privada para gestionar el hogar de una pareja: listas de la compra inteligentes con histórico de precios, escaneo de productos (foto + OCR), registro de facturas y seguimiento de préstamos/gastos fijos.

| Aspecto | Detalle |
|---------|---------|
| **Nombre** | Home Manager |
| **Slug** | home-manager |
| **Stack** | Go 1.23 + chi + PostgreSQL 16 + Redis 7 |
| **Frontend** | Vanilla JavaScript PWA (instalable en iPhone) |
| **Tema UI** | Verde pistacho (#B4D89E) — color favorito de la usuaria |
| **Idioma UI** | 🇪🇸 Español |
| **Idioma código** | 🇬🇧 Inglés |
| **Deploy target** | Railway |
| **Estado** | En desarrollo / SDD activo |

---

## Technology Stack

| Capa | Tecnología |
|------|------------|
| **Backend** | Go 1.23+ |
| **Router** | [chi](https://github.com/go-chi/chi) v5 |
| **Base de datos** | PostgreSQL 16 |
| **Caché / sesiones** | Redis 7 |
| **Auth** | JWT + bcrypt (login simple para 2 usuarios) |
| **Frontend** | Vanilla JS, PWA, cámara + OCR (Tesseract.js) |
| **Estilos** | CSS custom properties, tema verde pistacho |
| **DevOps** | Docker Compose (dev), Railway (prod) |

---

## Funcionalidades principales

1. **Listas de la compra**
   - Crear listas de compra.
   - Añadir productos reales con nombre, cantidad, precio y tienda.
   - Histórico de precios por producto.
   - Precio estimado/cerrado de la lista basado en el último precio conocido.

2. **Escaneo de productos**
   - Tomar foto al producto + precio.
   - OCR con Tesseract.js para reconocer nombre y precio.
   - Guardar producto en la base de datos.
   - Añadir directamente a la lista activa.

3. **Gastos personales y compartidos**
   - Cada usuario registra sus gastos con descripción, importe, categoría y fecha.
   - Visibilidad `privada` (solo el creador) o `compartida` (ambos usuarios).
   - Gastos compartidos con porcentaje de división configurable.
   - Soporte para gastos recurrentes mensuales.
   - Listado filtrado por mes y visibilidad.

4. **Ingresos personales y compartidos**
   - Cada usuario registra sus ingresos con descripción, importe, categoría y fecha.
   - Visibilidad `privada` o `compartida` configurable por ingreso.
   - Soporte para ingresos recurrentes mensuales.
   - Listado filtrado por mes y visibilidad.

5. **Dashboard económico**
   - Resumen mensual: ingresos, gastos y balance.
   - Gráfico de barras: ingresos vs gastos por mes.
   - Gráfico de líneas: evolución del balance mensual.
   - Filtro por año.

6. **Facturas del hogar (pendiente)**
   - Registrar facturas recurrentes (luz, agua, gas, internet, etc.).
   - Desglose de gastos semanal y mensual.

7. **Préstamos y gastos fijos (pendiente)**
   - Registrar préstamos con cuota mensual.
   - Marcar gastos fijos recurrentes.
   - Impacto en el balance mensual.

---

## Directory Structure

```
home-manager/
├── cmd/api/                 # Entry point (main.go)
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Middleware (auth, logging)
│   ├── store/               # Data access layer (PostgreSQL)
│   └── config/              # Configuration
├── pkg/                     # Public/reusable code
├── web/                     # Frontend PWA (HTML, CSS, JS)
│   ├── index.html
│   ├── manifest.json
│   ├── sw.js
│   └── css/
├── scripts/
│   ├── seed.ps1             # Database seeding
│   ├── test-api.ps1         # API smoke tests
│   └── db-reset.ps1         # Database reset
├── docs/
│   └── CONVENTIONS.md       # Coding conventions
├── .claude/skills/          # SDD skills
├── .github/                 # GitHub templates y CI
├── Dockerfile.dev
├── docker-compose.yml
├── Makefile
├── .env.example
├── .gitignore
├── SECURITY.md
└── AGENTS.md                # This file
```

---

## Convenciones de Código

Ver `docs/CONVENTIONS.md` para el detalle completo. Puntos clave:

- **UI / user-facing**: español.
- **Código / variables / funciones / comentarios**: inglés.
- **Commits**: inglés, Conventional Commits.
- **Errores**: lowercase, sin punto final.
- **SQL**: tablas en plural snake_case, columnas snake_case, queries parametrizadas.
- **API responses**:
  - Error: `{ "error": "...", "code": "...", "status": N }`
  - Lista: `{ "data": [...], "meta": {...} }`
  - Recurso: `{ "data": {...} }`

---

## Comandos Rápidos

```bash
make setup        # Copia .env.example a .env + go mod tidy
make dev          # Iniciar servidor de desarrollo
make test         # Ejecutar tests
make build        # Compilar binario
make lint         # Linter
make fmt          # Formatear código
make db-seed      # Poblar base de datos
make db-reset     # Resetear base de datos
make security-check  # Checklist pre-deploy
```

### Docker Compose

```bash
docker-compose up -d          # Levantar PostgreSQL + Redis + API
docker-compose logs -f api    # Ver logs
docker-compose down           # Detener todo
```

---

## Variables de Entorno Críticas

| Variable | Descripción | Requerida |
|----------|-------------|-----------|
| `DATABASE_URL` | PostgreSQL connection string | ✅ Sí |
| `JWT_SECRET` | JWT signing secret (>=32 chars) | ✅ Sí |
| `JWT_EXPIRATION_HOURS` | Expiración del token (default: 24) | ⚠️ Opcional |
| `API_URL` | URL base de la API | ⚠️ Opcional |
| `PORT` | Puerto del servidor (default: 8080) | ⚠️ Opcional |
| `ENV` | Entorno (development/production) | ⚠️ Opcional |
| `REDIS_URL` | Redis connection string | ⚠️ Opcional |

---

## SDD Workflow

Usamos Spec-Driven Development:

| Fase | Skill | Propósito |
|------|-------|-----------|
| 1 | `/sdd-init` | Inicializar contexto del proyecto |
| 2 | `/sdd-explore` | Investigar/explorar antes de cambiar |
| 3 | `/sdd-spec` | Especificar cambios |
| 4 | `/sdd-tasks` | Dividir en tareas |
| 5 | `/sdd-apply` | Implementar |
| 6 | `/sdd-verify` | Verificar |
| 7 | `/sdd-archive` | Archivar |

**Persistencia**: `hybrid` — usa `openspec/` para especificaciones versionables y `engram` para contexto entre sesiones.

---

## Engram

Este proyecto usa **engram** para memoria persistente entre sesiones:

```bash
engram save "Título" "Descripción"
engram search "término"
engram context home-manager
```

---

## Seguridad

Revisar `SECURITY.md` antes de cada deploy:

- [ ] Secrets no commiteados, `.env` en `.gitignore`
- [ ] `JWT_SECRET` >= 32 caracteres, expiración <= 24h
- [ ] Passwords hasheadas con bcrypt/Argon2
- [ ] CORS configurado, rate limiting activo
- [ ] SQL solo parametrizado
- [ ] HTTPS forzado en producción (Railway)
- [ ] Headers de seguridad (HSTS, CSP, X-Frame-Options)

---

## Checklist de Inicio

- [x] Generar proyecto con `scripts/setup.ps1`
- [x] Personalizar `AGENTS.md`
- [x] Configurar `.env`
- [x] Diseñar esquema de base de datos
- [x] Implementar autenticación simple (2 usuarios)
- [x] Implementar CRUD de productos
- [x] Implementar listas de la compra
- [x] Crear PWA con tema verde pistacho
- [x] Preparar Dockerfile y Railway deploy
- [x] Registro de usuarios
- [x] Gastos privados/compartidos con división configurable
- [x] Ingresos privados/compartidos
- [x] Dashboard económico con gráficos
- [ ] Implementar escaneo OCR en PWA
- [ ] Implementar facturas y préstamos
- [ ] Deploy en Railway

---

*Home Manager — Gentleman Programming* 🤝
