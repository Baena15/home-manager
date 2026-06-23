# Home Manager

Gestor del hogar para parejas. PWA para organizar listas de la compra con histórico de precios, registrar facturas y controlar préstamos/gastos fijos.

> Proyecto generado desde [Plantilla Madre — Gentleman Stack](https://github.com/gentleman-programming/plantilla-madre).

## 🚀 Stack

- **Backend**: Go 1.23 + chi
- **Base de datos**: PostgreSQL 16
- **Caché**: Redis 7 (preparado para futuro uso)
- **Frontend**: Vanilla JavaScript PWA
- **Tema**: Verde pistacho (#B4D89E)
- **Deploy**: Railway

## ⚡ Quick Start

### Requisitos

- Go 1.23+
- Docker + Docker Compose (para BD local)
- Make (opcional)

### 1. Clonar y entrar

```bash
cd "Home Manager"
```

### 2. Configurar variables de entorno

```bash
cp .env.example .env
# Edita .env con tus valores (JWT_SECRET, OWNER_PASSWORD, PARTNER_PASSWORD)
```

### 3. Levantar PostgreSQL y Redis

```bash
docker-compose up -d postgres redis
```

### 4. Ejecutar migraciones

```bash
make db-migrate
```

O arranca el servidor directamente en desarrollo (las migraciones se ejecutan automáticamente):

```bash
make dev
```

### 5. Acceder

- API: http://localhost:8080
- PWA: http://localhost:8080
- Health: http://localhost:8080/health

### Usuarios por defecto (desarrollo)

| Email | Rol |
|-------|-----|
| `owner@home.local` | owner |
| `partner@home.local` | partner |

Las contraseñas se configuran en `.env` (`OWNER_PASSWORD`, `PARTNER_PASSWORD`).

## 🧪 Testing

```bash
make test
make lint
```

## 🚂 Deploy en Railway

1. Crea un nuevo proyecto en [Railway](https://railway.app).
2. Conecta tu repositorio de GitHub.
3. Añade un servicio PostgreSQL desde el dashboard de Railway.
4. Configura las variables de entorno:
   - `DATABASE_URL` (Railway la genera automáticamente)
   - `JWT_SECRET` (mínimo 32 caracteres)
   - `OWNER_EMAIL`, `OWNER_PASSWORD`
   - `PARTNER_EMAIL`, `PARTNER_PASSWORD`
   - `AUTO_MIGRATE=true`
   - `ENV=production`
5. Railway usará el `Dockerfile` y `railway.toml` incluidos.

## 📁 Estructura

```
home-manager/
├── cmd/api/              # Punto de entrada
├── internal/
│   ├── config/           # Configuración
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # Middleware
│   └── store/            # Capa de acceso a datos
├── pkg/auth/             # Bcrypt + JWT
├── web/                  # PWA (HTML, CSS, JS)
├── migrations/           # Esquema PostgreSQL
├── Dockerfile            # Imagen de producción
├── railway.toml          # Configuración de Railway
└── docker-compose.yml    # Entorno de desarrollo
```

## 📱 Funcionalidades MVP

- [x] Login JWT para 2 usuarios
- [x] Catálogo de productos con histórico de precios
- [x] Listas de la compra con total cerrado
- [x] PWA instalable en iPhone
- [ ] Escaneo OCR de productos (próximamente)
- [ ] Facturas del hogar (próximamente)
- [ ] Préstamos y gastos fijos (próximamente)
- [ ] Dashboard de gastos (próximamente)

## 🤝 Convenciones

- Documentación y UI en español.
- Código, variables y commits en inglés.
- Conventional Commits.
