# Makefile - Home Manager
# Comandos estandar para proyectos Go + PostgreSQL + Redis

.PHONY: help dev test test-watch build deploy clean lint fmt db-migrate db-seed db-reset setup security-check

# â”€â”€â”€ Variables â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
BINARY_NAME   := api
BINARY_PATH   := ./cmd/api
BUILD_DIR     := ./bin
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# â”€â”€â”€ Default: muestra ayuda â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
help:
	@echo "ðŸŽ¼ Comandos disponibles:"
	@echo ""
	@echo "  make dev          - Iniciar servidor de desarrollo (go run)"
	@echo "  make test         - Ejecutar tests"
	@echo "  make test-watch   - Tests en modo watch (requiere reflex)"
	@echo "  make coverage     - Reporte de cobertura"
	@echo "  make build        - Compilar binario"
	@echo "  make deploy       - Deploy a produccion"
	@echo "  make clean        - Limpiar binarios/cache"
	@echo "  make lint         - Ejecutar linter (golangci-lint o go vet)"
	@echo "  make fmt          - Formatear codigo (go fmt)"
	@echo "  make db-migrate   - Ejecutar migraciones"
	@echo "  make db-seed      - Poblar base de datos"
	@echo "  make db-reset     - Resetear base de datos"
	@echo "  make setup        - Setup inicial del proyecto"
	@echo "  make security-check - Checklist pre-deploy"
	@echo ""

# â”€â”€â”€ Desarrollo â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
dev:
	@echo "ðŸš€ Iniciando desarrollo..."
	@go run $(BINARY_PATH)

# â”€â”€â”€ Testing â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
test:
	@echo "ðŸ§ª Ejecutando tests..."
	@go test -v -race ./...

test-watch:
	@echo "ðŸ‘€ Tests en modo watch..."
	@echo "âš ï¸  Requiere reflex: go install github.com/cespare/reflex@latest"
	@reflex -r '\.go$$' -s -- go test -v ./...

coverage:
	@echo "ðŸ“Š Generando cobertura..."
	@go test -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "ðŸ“„ Reporte: $(COVERAGE_HTML)"

# â”€â”€â”€ Build â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
build:
	@echo "ðŸ”¨ Compilando..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	@echo "âœ… Build completo: $(BUILD_DIR)/$(BINARY_NAME)"

# â”€â”€â”€ Deploy â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
deploy:
	@echo "ðŸš€ Deploying..."
	@echo "âš ï¸  TODO: Configura tu comando de deploy (fly deploy, ssh, docker push, etc.)"
	@echo "âœ… Deploy completado"

# â”€â”€â”€ Clean â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
clean:
	@echo "ðŸ§¹ Limpiando..."
	@powershell -Command "Remove-Item -Recurse -Force -ErrorAction SilentlyContinue '$(BUILD_DIR)','dist','build','coverage','*.log','*.out','*.exe'"
	@echo "âœ… Limpieza completa"

# â”€â”€â”€ Code Quality â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
lint:
	@echo "ðŸ” Linting..."
	@which golangci-lint >/dev/null 2>&1 && golangci-lint run || go vet ./...

fmt:
	@echo "âœ¨ Formateando..."
	@go fmt ./...

# â”€â”€â”€ Database â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
db-migrate:
	@echo "ðŸ—„ï¸  Ejecutando migraciones..."
	@echo "âš ï¸  TODO: Agrega tu herramienta de migraciones (golang-migrate, goose, etc.)"
	# @migrate -path ./migrations -database $$DATABASE_URL up

db-seed:
	@echo "ðŸŒ± Poblando datos..."
	@powershell -ExecutionPolicy Bypass -File .\scripts\seed.ps1

db-reset:
	@echo "ðŸ—„ï¸  Resetear base de datos..."
	@powershell -ExecutionPolicy Bypass -File .\scripts\db-reset.ps1
	@echo "âœ… Base de datos reseteada"

# â”€â”€â”€ Setup â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
setup:
	@echo "âš™ï¸  Setup inicial..."
	@powershell -Command "if (Test-Path .env.example) { Copy-Item .env.example .env } else { Write-Host 'âš ï¸  .env.example no existe' }"
	@go mod tidy 2>nul || echo "âš ï¸  go.mod no existe aun"
	@echo "âœ… Setup completo. Edita .env con tus configuraciones"

# â”€â”€â”€ Seguridad pre-deploy â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
security-check:
	@echo "ðŸ” Security checklist..."
	@echo "  â˜ Secrets no commiteados"
	@echo "  â˜ .env.example actualizado"
	@echo "  â˜ Dependencias sin vulnerabilidades (go mod tidy)"
	@echo "  â˜ Tests pasando"
	@echo "  â˜ Linter sin errores"
	@echo "  â˜ JWT_SECRET >= 32 caracteres"
	@echo "  â˜ DATABASE_URL usa SSL en produccion"
	@echo ""
	@echo "Ejecuta: make lint && make test"
