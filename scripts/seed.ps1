# Script para poblar la base de datos con datos de ejemplo (PowerShell)
# Template: Ajusta los placeholders [PROJECT_*] al dominio de tu proyecto

Write-Host "ðŸŒ± Seeding database for Home Manager..." -ForegroundColor Green

# Verificar que PostgreSQL esta corriendo
$containerRunning = docker-compose ps | Select-String "postgres.*Up"

if (-not $containerRunning) {
    Write-Host "âš ï¸  PostgreSQL no esta corriendo. Iniciando..." -ForegroundColor Yellow
    docker-compose up -d postgres
    Start-Sleep -Seconds 5
}

# Configuracion leida de variables de entorno o defaults
$dbUser = if ($env:DB_USER) { $env:DB_USER } else { "home_manager" }
$dbName = if ($env:DB_NAME) { $env:DB_NAME } else { "home_manager" }

# â”€â”€â”€ SQL DE SEED â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
# NOTA: Edita este bloque para reflejar el dominio de tu proyecto.
# Los placeholders [PROJECT_*] fueron pre-llenados por setup.ps1.

$seedSQL = @"
-- Crear usuario admin
INSERT INTO users (email, password_hash, first_name, last_name, is_admin)
VALUES (
    'admin@home-manager.com',
    '\$2a\$10\$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Admin',
    'User',
    true
)
ON CONFLICT (email) DO NOTHING;

-- Crear usuario de prueba
INSERT INTO users (email, password_hash, first_name, last_name, is_admin)
VALUES (
    'test@test.com',
    '\$2a\$10\$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'Test',
    'User',
    false
)
ON CONFLICT (email) DO NOTHING;

-- Mostrar usuarios creados
SELECT 'Usuarios creados:' as info;
SELECT id, email, first_name, is_admin FROM users;

-- TODO: Agrega aqui seed de tu dominio
-- Ejemplo:
-- INSERT INTO items (name, description)
-- VALUES ('Ejemplo 1', 'Descripcion de prueba')
-- ON CONFLICT DO NOTHING;
--
-- SELECT 'Registros en items:' as info;
-- SELECT id, name, created_at, updated_at FROM items LIMIT 5;
"@

$seedSQL | docker-compose exec -T postgres psql -U $dbUser -d $dbName

Write-Host "âœ… Seed completado!" -ForegroundColor Green
Write-Host ""
Write-Host "Usuarios de prueba:"
Write-Host "  Admin: admin@home-manager.com / password"
Write-Host "  User:  test@test.com / password"
