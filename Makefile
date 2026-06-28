# Makefile - Home Manager
# Comandos estandar para proyectos Go + PostgreSQL + Redis

.PHONY: help dev test test-watch build deploy clean lint fmt db-migrate db-seed db-reset setup security-check

# ─── Variables ──────────────────────────────────────────────────────
BINARY_NAME   := api
BINARY_PATH   := ./cmd/api
BUILD_DIR     := ./bin
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html
MIGRATIONS_DIR := ./internal/store/migrations

# ─── Default: muestra ayuda ─────────────────────────────────────────
help:
	@echo "🎵 Comandos disponibles:"
	@echo ""
	@echo "  make dev          - Iniciar servidor de desarrollo (go run)"
	@echo "  make test         - Ejecutar tests"
	@echo "  make test-watch   - Tests en modo watch (requiere reflex)"
	@echo "  make coverage     - Reporte de cobertura"
	@echo "  make build        - Compilar binario"
	@echo "  make deploy       - Deploy a produccion"
	@echo "  make clean        - Limpiar binarios/cache"
	@echo "  make lint         - Ejecutar linter (golangci-lint o go vet)"
	@echo "  make fmt          - Formatear codigo"
	@echo "  make db-migrate   - Ejecutar migraciones"
	@echo "  make db-seed      - Poblar base de datos"
	@echo "  make db-reset     - Resetear base de datos"
	@echo "  make setup        - Setup inicial del proyecto"
	@echo "  make security-check - Checklist pre-deploy"
	@echo ""

# ─── Desarrollo ─────────────────────────────────────────────────────
dev:
	@echo "🚀 Iniciando desarrollo..."
	@go run $(BINARY_PATH)

# ─── Testing ────────────────────────────────────────────────────────
test:
	@echo "🧪 Ejecutando tests..."
	@go test -v -race ./...

test-watch:
	@echo "👀 Tests en modo watch..."
	@echo "⚠️  Requiere reflex: go install github.com/cespare/reflex@latest"
	@reflex -r '\.go$$' -s -- go test -v ./...

coverage:
	@echo "📊 Generando cobertura..."
	@go test -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "📄 Reporte: $(COVERAGE_HTML)"

# ─── Build ──────────────────────────────────────────────────────────
build:
	@echo "🔨 Compilando..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	@echo "✅ Build completo: $(BUILD_DIR)/$(BINARY_NAME)"

# ─── Deploy ─────────────────────────────────────────────────────────
deploy:
	@echo "🚀 Deploying to Railway..."
	@which railway >/dev/null 2>&1 || (echo "⚠️  Railway CLI no encontrado. Instálalo con: npm install -g @railway/cli" && exit 1)
	@railway up
	@echo "✅ Deploy completado"

# ─── Clean ──────────────────────────────────────────────────────────
clean:
	@echo "🧹 Limpiando..."
	@powershell -Command "Remove-Item -Recurse -Force -ErrorAction SilentlyContinue '$(BUILD_DIR)','dist','build','coverage','*.log','*.out','*.exe'"
	@echo "✅ Limpieza completa"

# ─── Code Quality ───────────────────────────────────────────────────
lint:
	@echo "🔍 Linting..."
	@which golangci-lint >/dev/null 2>&1 && golangci-lint run || go vet ./...

fmt:
	@echo "✨ Formateando..."
	@go fmt ./...

# ─── Database ───────────────────────────────────────────────────────
db-migrate:
	@echo "🗄️  Ejecutando migraciones..."
	@if "$(DATABASE_URL)" == "" ( \
		echo "⚠️  DATABASE_URL no esta definida"; \
		exit 1; \
	)
	@for %%f in ($(MIGRATIONS_DIR)\*.sql) do ( \
		@echo "Aplicando %%f"; \
		psql "$(DATABASE_URL)" -f "%%f"; \
	)

db-seed:
	@echo "🌱 Poblando datos..."
	@powershell -ExecutionPolicy Bypass -File .\scripts\seed.ps1

db-reset:
	@echo "🗑️  Resetear base de datos..."
	@powershell -ExecutionPolicy Bypass -File .\scripts\db-reset.ps1
	@echo "✅ Base de datos reseteada"

# ─── Setup ──────────────────────────────────────────────────────────
setup:
	@echo "⚙️  Setup inicial..."
	@powershell -Command "if (Test-Path .env.example) { Copy-Item .env.example .env } else { Write-Host '⚠️  .env.example no existe' }"
	@go mod tidy 2>nul || echo "⚠️  go.mod no existe aun"
	@echo "✅ Setup completo. Edita .env con tus configuraciones"

# ─── Seguridad pre-deploy ───────────────────────────────────────────
security-check:
	@echo "🔐 Security checklist..."
	@echo "  ☐ Secrets no commiteados"
	@echo "  ☐ .env.example actualizado"
	@echo "  ☐ Dependencias sin vulnerabilidades (go mod tidy)"
	@echo "  ☐ Tests pasando"
	@echo "  ☐ Linter sin errores"
	@echo "  ☐ JWT_SECRET >= 32 caracteres"
	@echo "  ☐ DATABASE_URL usa SSL en produccion"
	@echo ""
	@echo "Ejecuta: make lint && make test"
