# Script de prueba rapida de la API
# Template: Ajusta la configuracion de endpoints a tu API
# Ejecutar despues de: docker-compose up

# â”€â”€â”€ CONFIGURACION â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
$baseUrl     = "http://localhost:8080"
$projectName = "Home Manager"

# Endpoints publicos (no requieren autenticacion)
$publicEndpoints = @(
    @{ Method = "GET"; Path = "/api/v1/items?limit=5"; Name = "Listar recursos" },
    @{ Method = "GET"; Path = "/api/v1/items/{id}"; Name = "Obtener recurso por ID"; DynamicId = $true }
)

# Endpoints de autenticacion
$authEndpoints = @(
    @{ Method = "POST"; Path = "/api/v1/auth/register"; Name = "Registro"; Body = @{ email = "apitest@example.com"; password = "TestPassword123!"; first_name = "API"; last_name = "Test" } },
    @{ Method = "POST"; Path = "/api/v1/auth/login";    Name = "Login";    Body = @{ email = "apitest@example.com"; password = "TestPassword123!" } }
)

# Endpoints protegidos (requieren token)
$protectedEndpoints = @(
    @{ Method = "GET"; Path = "/api/v1/me"; Name = "Perfil de usuario" },
    @{ Method = "GET"; Path = "/api/v1/dashboard"; Name = "Recurso protegido" }
)

# Endpoints de admin
$adminEndpoints = @(
    @{ Method = "GET"; Path = "/api/v1/admin/stats"; Name = "Estadisticas admin" }
)

# Usuario admin de seed (para probar endpoints admin)
$adminCredentials = @{
    email    = "admin@home-manager.com"
    password = "password"
}
# â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

$green  = "`e[32m"
$red    = "`e[31m"
$yellow = "`e[33m"
$reset  = "`e[0m"

function Test-Endpoint {
    param($Method, $Path, $Body = $null, $Token = $null)
    
    $url = "$baseUrl$Path"
    $headers = @{}
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }
    if ($Body) {
        $headers["Content-Type"] = "application/json"
    }

    try {
        if ($Method -eq "GET") {
            $response = Invoke-RestMethod -Uri $url -Method $Method -Headers $headers -ErrorAction Stop
        } else {
            $response = Invoke-RestMethod -Uri $url -Method $Method -Headers $headers -Body ($Body | ConvertTo-Json) -ErrorAction Stop
        }
        Write-Host "${green}âœ… $Method $Path - OK${reset}"
        return $response
    } catch {
        Write-Host "${red}âŒ $Method $Path - FAILED${reset}"
        Write-Host "   Error: $($_.Exception.Message)"
        return $null
    }
}

Write-Host "`nðŸ§ª Testing $projectName API...`n" -ForegroundColor Cyan
Write-Host "Base URL: $baseUrl`n"

# 1. Health check (endpoint fijo)
Write-Host "${yellow}--- Health Check ---${reset}"
Test-Endpoint -Method "GET" -Path "/health"

# 2. Endpoints publicos
Write-Host "`n${yellow}--- Public Endpoints ---${reset}"
$lastResponse = $null
foreach ($ep in $publicEndpoints) {
    $path = $ep.Path
    # Si el endpoint necesita un ID dinamico del anterior, intentamos extraerlo
    if ($ep.DynamicId -and $lastResponse -and $lastResponse.PSObject.Properties["items"]) {
        $items = $lastResponse.items
        if ($items -and $items.Count -gt 0 -and $items[0].id) {
            $path = $path -replace "\{id\}",$items[0].id
        } elseif ($items -and $items.Count -gt 0 -and $items[0].slug) {
            $path = $path -replace "\{id\}",$items[0].slug
        }
    }
    $lastResponse = Test-Endpoint -Method $ep.Method -Path $path
}

# 3. Autenticacion
Write-Host "`n${yellow}--- Auth Endpoints ---${reset}"
$token = $null
foreach ($ep in $authEndpoints) {
    $response = Test-Endpoint -Method $ep.Method -Path $ep.Path -Body $ep.Body
    if ($response -and $response.access_token) {
        $token = $response.access_token
        Write-Host "   Access token received: $($token.Substring(0,20))..."
    }
}

# 4. Endpoints protegidos
if ($token) {
    Write-Host "`n${yellow}--- Protected Endpoints ---${reset}"
    foreach ($ep in $protectedEndpoints) {
        Test-Endpoint -Method $ep.Method -Path $ep.Path -Token $token
    }
} else {
    Write-Host "${red}âš ï¸  No token available, skipping protected endpoints${reset}"
}

# 5. Admin endpoints (login como admin primero)
Write-Host "`n${yellow}--- Admin Endpoints ---${reset}"
$adminToken = $null
$adminLogin = Test-Endpoint -Method "POST" -Path "/api/v1/auth/login" -Body $adminCredentials
if ($adminLogin -and $adminLogin.access_token) {
    $adminToken = $adminLogin.access_token
    foreach ($ep in $adminEndpoints) {
        Test-Endpoint -Method $ep.Method -Path $ep.Path -Token $adminToken
    }
}

Write-Host "`n${green}âœ… Test completed!${reset}`n"
