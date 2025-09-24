# Script para ejecutar Turivo Backend

Write-Host "🚀 Iniciando Turivo Backend..." -ForegroundColor Green

# Verificar que existe .env
if (-not (Test-Path ".env")) {
    Write-Host "❌ Archivo .env no encontrado. Ejecuta setup.ps1 primero." -ForegroundColor Red
    exit 1
}

# Verificar que PostgreSQL está corriendo
Write-Host "🗄️ Verificando PostgreSQL..." -ForegroundColor Yellow
try {
    $result = docker compose ps --services --filter "status=running"
    if ($result -contains "postgres") {
        Write-Host "✅ PostgreSQL está corriendo" -ForegroundColor Green
    } else {
        Write-Host "⚠️ PostgreSQL no está corriendo. Iniciando..." -ForegroundColor Yellow
        docker compose up -d
        Start-Sleep -Seconds 5
    }
} catch {
    Write-Host "❌ Error verificando PostgreSQL" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "🌐 Iniciando servidor en http://localhost:8080" -ForegroundColor Cyan
Write-Host "📚 Swagger UI: http://localhost:8080/swagger/index.html" -ForegroundColor Cyan
Write-Host ""
Write-Host "Presiona Ctrl+C para detener el servidor" -ForegroundColor Yellow
Write-Host ""

# Ejecutar el servidor
go run cmd/api/main.go

