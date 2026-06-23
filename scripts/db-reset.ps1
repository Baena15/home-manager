# Script para resetear la base de datos (PowerShell)
# Template: Reemplaza placeholders [PROJECT_*] via setup.ps1

param(
    [switch]$Force = $false
)

$containerName = "home-manager-postgres"

if (-not $Force) {
    $response = Read-Host "âš ï¸  Esto eliminara TODOS los datos de la base de datos. Â¿Continuar? (y/N)"
    if ($response -ne 'y') {
        Write-Host "âŒ Cancelado" -ForegroundColor Red
        exit 0
    }
}

Write-Host "ðŸ—„ï¸  Reseteando base de datos..." -ForegroundColor Cyan

docker-compose down -v postgres 2>$null
docker-compose up -d postgres

# Esperar a que PostgreSQL estÃ© saludable
$retries = 30
while ($retries -gt 0) {
    $healthy = docker-compose ps | Select-String "$containerName.*healthy|Up.*healthy"
    if ($healthy) { break }
    Start-Sleep -Seconds 1
    $retries--
}

if ($retries -eq 0) {
    Write-Host "âŒ PostgreSQL no se volvio saludable a tiempo" -ForegroundColor Red
    exit 1
}

Write-Host "âœ… Base de datos reseteada y lista" -ForegroundColor Green
