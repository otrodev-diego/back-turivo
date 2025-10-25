# Turivo Backend Setup Script para Windows

Write-Host "🚀 Configurando Turivo Backend..." -ForegroundColor Green

# 1. Crear archivo .env
Write-Host "📝 Creando archivo .env..." -ForegroundColor Yellow
$envContent = @"
HTTP_PORT=8080
ENV=local
LOG_LEVEL=info
CORS_ORIGINS=*

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=turivo
DB_SSLMODE=disable
DB_TIMEZONE=UTC
DB_DSN=postgres://postgres:postgres@localhost:5432/turivo?sslmode=disable

JWT_SECRET=change-me-in-prod
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
"@

$envContent | Out-File -FilePath ".env" -Encoding UTF8
Write-Host "✅ Archivo .env creado" -ForegroundColor Green

# 2. Verificar Docker
Write-Host "🐳 Verificando Docker..." -ForegroundColor Yellow
try {
    docker --version | Out-Null
    Write-Host "✅ Docker encontrado" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker no encontrado. Instala Docker Desktop primero." -ForegroundColor Red
    exit 1
}

# 3. Levantar PostgreSQL
Write-Host "🗄️ Levantando PostgreSQL..." -ForegroundColor Yellow
docker compose up -d
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ PostgreSQL iniciado" -ForegroundColor Green
} else {
    Write-Host "❌ Error al iniciar PostgreSQL" -ForegroundColor Red
    exit 1
}

# 4. Esperar a que PostgreSQL esté listo
Write-Host "⏳ Esperando a que PostgreSQL esté listo..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# 5. Verificar si migrate está instalado
Write-Host "🔄 Verificando migrate..." -ForegroundColor Yellow
try {
    migrate -version | Out-Null
    Write-Host "✅ migrate encontrado" -ForegroundColor Green
} catch {
    Write-Host "📦 Instalando migrate..." -ForegroundColor Yellow
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    Write-Host "✅ migrate instalado" -ForegroundColor Green
}

# 6. Ejecutar migraciones
Write-Host "📊 Ejecutando migraciones..." -ForegroundColor Yellow
migrate -database "postgres://postgres:postgres@localhost:5432/turivo?sslmode=disable" -path migrations up
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Migraciones ejecutadas" -ForegroundColor Green
} else {
    Write-Host "❌ Error en migraciones" -ForegroundColor Red
    exit 1
}

# 7. Generar datos demo
Write-Host "🌱 Generando datos demo..." -ForegroundColor Yellow
go run cmd/seed/main.go -demo
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Datos demo generados" -ForegroundColor Green
} else {
    Write-Host "❌ Error generando datos demo" -ForegroundColor Red
}

Write-Host ""
Write-Host "🎉 ¡Setup completado!" -ForegroundColor Green
Write-Host ""
Write-Host "Para iniciar el servidor ejecuta:" -ForegroundColor Cyan
Write-Host "  go run cmd/api/main.go" -ForegroundColor White
Write-Host ""
Write-Host "URLs importantes:" -ForegroundColor Cyan
Write-Host "  Health Check: http://localhost:8080/healthz" -ForegroundColor White
Write-Host "  Swagger UI:   http://localhost:8080/swagger/index.html" -ForegroundColor White
Write-Host ""
Write-Host "Usuarios demo:" -ForegroundColor Cyan
Write-Host "  Las credenciales por defecto han sido removidas por seguridad" -ForegroundColor Yellow
Write-Host "  Crea tu primer usuario administrador usando el endpoint de registro" -ForegroundColor Yellow
