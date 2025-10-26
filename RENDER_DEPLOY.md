# 🚀 Guía de Despliegue en Render

## 📋 Configuración de Variables de Entorno

### Variables Requeridas en Render Dashboard

Configura estas variables en tu servicio de Render:

#### 1. Base de Datos
Render conecta automáticamente las bases de datos PostgreSQL. La variable `DATABASE_URL` se configura automáticamente si vinculaste una base de datos.

Si usas una base de datos externa, agrega:
```
DATABASE_URL=postgres://user:password@host:port/database?sslmode=require
```

#### 2. JWT (Requerido)
```
JWT_SECRET=tu-secret-jwt-muy-seguro-cambiar-en-produccion
```

#### 3. SMTP (Requerido para emails)
```
SMTP_HOST=smtp.tu-proveedor.com
SMTP_PORT=465
SMTP_USERNAME=tu-usuario-smtp
SMTP_PASSWORD=tu-password-smtp
SMTP_FROM=noreply@turivo.com
```

#### 4. Configuración del Servidor (Opcional)
```
HTTP_PORT=3000
LOG_LEVEL=info
CORS_ORIGINS=https://tu-frontend.com,https://otro-dominio.com
```

---

## 🔧 Configuración en Render Dashboard

1. **Servicio Web:**
   - Name: `turivo-backend`
   - Environment: `Docker`
   - Build Command: (automático con Dockerfile)
   - Start Command: (automático con Dockerfile)
   - Plan: `Free` o `Starter`

2. **Base de Datos PostgreSQL:**
   - Name: `turivo-db`
   - PostgreSQL Version: `15` o superior
   - Plan: `Free` (dev) o `Starter` (prod)

3. **Link Database:**
   - En la configuración del servicio web, haz clic en "Link Database"
   - Selecciona tu base de datos PostgreSQL

4. **Environment Variables:**
   - Agrega todas las variables listadas arriba
   - Guarda los cambios

---

## 🐳 Docker Configuration

El proyecto ya incluye un `Dockerfile` optimizado:
- Usa multi-stage build
- Imagen base: `golang:1.24-alpine`
- Crea un binario estático sin dependencias

---

## 🧪 Verificar el Despliegue

### 1. Health Check
```bash
curl https://tu-servicio.onrender.com/healthz
```

Deberías recibir:
```json
{
  "status": "ok",
  "message": "Turivo Backend API is running",
  "version": "1.0.0"
}
```

### 2. Swagger Documentation
Visita: `https://tu-servicio.onrender.com/swagger/index.html`

### 3. Probar Login
```bash
curl -X POST https://tu-servicio.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"tu-email@example.com","password":"tu-password"}'
```

---

## ⚠️ Solución de Problemas

### Error: "lookup dpg-xxx on 169.254.20.10:53: no such host"

**Problema:** Render está resolviendo el host interno en lugar del host completo.

**Solución:**
1. Verifica que la base de datos esté vinculada al servicio web
2. Si usas una base de datos externa, asegúrate de proporcionar el `DATABASE_URL` completo con el host correcto:
   ```
   DATABASE_URL=postgres://user:password@dpg-xxx.oregon-postgres.render.com:5432/database?sslmode=require
   ```

### Error: "Could not read config file"

**Problema:** El servicio está buscando un archivo `.env` que no existe en el contenedor.

**Solución:** Normal. El código está diseñado para funcionar solo con variables de entorno. Este warning es inofensivo.

### Error: "Failed to ping database"

**Problema:** No puede conectarse a la base de datos.

**Soluciones:**
1. Verifica que la base de datos esté "Linked" al servicio web
2. Revisa los logs de la base de datos
3. Asegúrate de que la variable `DATABASE_URL` esté configurada correctamente
4. Si usas SSL, verifica que `sslmode=require`

---

## 📊 Monitoreo

### Logs en Tiempo Real
Render proporciona logs en tiempo real en el dashboard:
```
Dashboard > Tu Servicio > Logs
```

### Logs Importantes
- `"Starting Turivo Backend API"` - ✅ Servicio inició correctamente
- `"Database connection established"` - ✅ Conectado a la BD
- `"Server starting"` - ✅ API está lista para recibir requests

---

## 🔐 Seguridad en Producción

### Checklist de Seguridad

- [ ] Cambiar `JWT_SECRET` por un valor fuerte y único
- [ ] Configurar `CORS_ORIGINS` con solo tus dominios permitidos
- [ ] Usar `sslmode=require` para la base de datos
- [ ] Configurar SMTP con credenciales seguras
- [ ] Revisar logs regularmente para detectar errores
- [ ] Hacer backup regular de la base de datos

---

## 🚀 Próximos Pasos

1. Configurar dominio personalizado (opcional)
2. Configurar SSL automático (ya incluido en Render)
3. Configurar webhooks para CI/CD
4. Configurar alertas y monitoreo

---

## 📞 Soporte

Si tienes problemas con el despliegue:
1. Revisa los logs en Render Dashboard
2. Verifica todas las variables de entorno
3. Consulta la documentación de Render: https://render.com/docs

---

**Nota:** La configuración ahora prioriza `DATABASE_URL` sobre los componentes individuales, lo cual es compatible con Render, Heroku y otros proveedores de PaaS.
