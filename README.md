# Turivo Backend API

API backend para el sistema de gestión de transporte Turivo, construido con Go siguiendo Clean Architecture.

## 🚀 Características

- **Clean Architecture** con separación clara de responsabilidades
- **Autenticación JWT** con refresh tokens
- **RBAC** (Role-Based Access Control) con roles: ADMIN, USER, DRIVER, COMPANY
- **Base de datos PostgreSQL** con migraciones automáticas
- **Documentación OpenAPI/Swagger** completa
- **Mock de Webpay** para simulación de pagos
- **Logging estructurado** con Zap
- **Middleware completo** (CORS, Recovery, Request ID, Auth)
- **Seeds de datos demo** alineados con el frontend
- **Tests unitarios** con mocks

## 🏗️ Arquitectura

```
back/
├── cmd/
│   ├── api/           # Punto de entrada de la aplicación
│   └── seed/          # Comando para generar datos demo
├── internal/
│   ├── domain/        # Entidades y reglas de negocio
│   ├── usecase/       # Casos de uso (interactors)
│   ├── interface/     # Handlers HTTP y middlewares
│   └── infrastructure/ # Implementaciones (DB, Auth, Config)
├── migrations/        # Migraciones SQL
├── docs/             # Documentación Swagger generada
└── sqlc/             # Queries SQL y configuración
```

## 📋 Requisitos

- Go 1.22+
- PostgreSQL 14+
- Docker & Docker Compose (opcional)

## ⚡ Inicio Rápido

### **Para Windows (PowerShell) - RECOMENDADO:**

```powershell
# 1. Setup automático completo
.\setup.ps1

# 2. Ejecutar servidor
.\run.ps1
```

### **Para Linux/Mac:**

```bash
# 1. Copiar variables de entorno
cp .env.example .env

# 2. Levantar PostgreSQL
make up

# 3. Ejecutar migraciones
make migrate-up

# 4. Generar datos demo
go run cmd/seed/main.go -demo

# 5. Ejecutar servidor
make run
```

### **Comandos manuales (cualquier SO):**

```bash
# 1. Levantar PostgreSQL
docker compose up -d

# 2. Instalar migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 3. Ejecutar migraciones
migrate -database "postgres://postgres:postgres@localhost:5432/turivo?sslmode=disable" -path migrations up

# 4. Generar datos demo
go run cmd/seed/main.go -demo

# 5. Ejecutar servidor
go run cmd/api/main.go
```

### 4. Acceder a la documentación

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/healthz

## 🔧 Variables de Entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `HTTP_PORT` | Puerto del servidor | `8080` |
| `LOG_LEVEL` | Nivel de logging | `info` |
| `DB_DSN` | Conexión a PostgreSQL | `postgres://postgres:postgres@localhost:5432/turivo?sslmode=disable` |
| `JWT_SECRET` | Secreto para JWT | `change-me-in-prod` |
| `JWT_ACCESS_TTL` | Duración del access token | `15m` |
| `JWT_REFRESH_TTL` | Duración del refresh token | `168h` |
| `CORS_ORIGINS` | Orígenes permitidos para CORS | `*` |

## 🛠️ Comandos Disponibles

```bash
# Docker
make up          # Levantar servicios
make down        # Bajar servicios
make logs        # Ver logs

# Base de datos
make migrate-up    # Ejecutar migraciones
make migrate-down  # Revertir migraciones
make psql         # Conectar a PostgreSQL

# Desarrollo
make run          # Ejecutar aplicación
make test         # Ejecutar tests
make build        # Compilar aplicación
make lint         # Ejecutar linter

# Documentación
make swagger      # Generar documentación Swagger
```

## 🔐 Autenticación

### Login
```bash
POST /api/v1/auth/login
{
  "email": "admin@turivo.com",
  "password": "password"
}
```

### Uso del token
```bash
Authorization: Bearer <access_token>
```

### Refresh token
```bash
POST /api/v1/auth/refresh
{
  "refresh_token": "<refresh_token>"
}
```

## 📊 Endpoints Principales

### Usuarios
- `GET /api/v1/users` - Listar usuarios (Admin)
- `POST /api/v1/users` - Crear usuario (Admin)
- `GET /api/v1/users/:id` - Obtener usuario (Admin)
- `PATCH /api/v1/users/:id` - Actualizar usuario (Admin)
- `DELETE /api/v1/users/:id` - Eliminar usuario (Admin)

### Conductores
- `GET /api/v1/drivers` - Listar conductores
- `POST /api/v1/drivers` - Crear conductor
- `GET /api/v1/drivers/:id` - Obtener conductor
- `PATCH /api/v1/drivers/:id` - Actualizar conductor
- `GET /api/v1/drivers/:id/kpis` - KPIs del conductor

### Reservas
- `GET /api/v1/reservations` - Listar reservas
- `POST /api/v1/reservations` - Crear reserva
- `GET /api/v1/reservations/:id` - Obtener reserva
- `PATCH /api/v1/reservations/:id/status` - Cambiar estado
- `GET /api/v1/reservations/:id/timeline` - Timeline de eventos

### Pagos
- `POST /api/v1/payments` - Crear pago
- `GET /api/v1/payments/:id` - Obtener pago
- `POST /api/v1/payments/:id/simulate` - Simular resultado (testing)

## 🎭 Roles y Permisos

| Rol | Descripción | Permisos |
|-----|-------------|----------|
| `ADMIN` | Administrador del sistema | Acceso completo |
| `COMPANY` | Empresa de transporte | CRUD conductores, reservas propias |
| `HOTEL` | Hotel cliente | Crear solicitudes, ver reservas propias |
| `DRIVER` | Conductor | Ver trips asignados |
| `USER` | Usuario final | Crear reservas, realizar pagos |

## 💰 Cálculo de Precios

Los precios se calculan automáticamente basado en:

- **Precio base por vehículo**:
  - Bus: $150,000 CLP
  - Van: $120,000 CLP
  - Sedan: $80,000 CLP
  - SUV: $100,000 CLP

- **Costos adicionales**:
  - Pasajero adicional (>1): +$5,000 CLP
  - Parada adicional: +$15,000 CLP
  - Idioma especial: +$25,000 CLP

## 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./...

# Tests con cobertura
go test -cover ./...

# Tests verbosos
go test -v ./internal/usecase
```

## 🔍 Datos Demo

El comando seed genera datos alineados con el frontend:

```bash
go run cmd/seed/main.go -demo
```

**Crear primer usuario administrador**:

```bash
# Crear usuario administrador
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin Sistema",
    "email": "admin@tudominio.com",
    "password": "tu_password_seguro",
    "role": "ADMIN"
  }'
```

**Nota**: Las credenciales por defecto han sido removidas por seguridad.

**Conductores**: CON-001, CON-002, CON-003
**Reservas**: RSV-1001, RSV-1002, RSV-1003

## 🐳 Docker

```bash
# Build imagen
docker build -t turivo-backend .

# Ejecutar contenedor
docker run -p 8080:8080 --env-file .env turivo-backend
```

## 📈 Monitoreo y Logs

- **Logs estructurados** en formato JSON
- **Request ID** para trazabilidad
- **Métricas de performance** en logs
- **Health check** en `/healthz`

## 🔧 Desarrollo

### Agregar nueva entidad

1. Crear migración en `migrations/`
2. Definir queries en `sqlc/queries/`
3. Regenerar código: `sqlc generate`
4. Crear entidad en `internal/domain/`
5. Implementar caso de uso en `internal/usecase/`
6. Crear handler en `internal/interface/http/handler/`
7. Agregar rutas en router
8. Actualizar documentación Swagger

### Regenerar Swagger

```bash
swag init -g cmd/api/main.go -o docs
```

## 🚀 Producción

### Variables críticas a cambiar:
- `JWT_SECRET`: Usar secreto seguro
- `DB_DSN`: Conexión a base de datos de producción
- `CORS_ORIGINS`: Especificar dominios permitidos
- `LOG_LEVEL`: Usar `warn` o `error`

### Consideraciones:
- Usar HTTPS
- Configurar rate limiting
- Monitoreo con Prometheus/Grafana
- Backup automático de base de datos
- Rotación de logs

## 📝 Contribución

1. Fork del repositorio
2. Crear feature branch
3. Implementar cambios con tests
4. Ejecutar linter: `make lint`
5. Verificar cobertura: `make test`
6. Crear Pull Request

## 📄 Licencia

MIT License - ver archivo LICENSE para detalles.

---

**Desarrollado con ❤️ para Turivo**
# back-turivo
