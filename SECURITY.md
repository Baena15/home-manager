# Security Checklist - Gentleman Stack

## Pre-Deploy Checklist

### 🔐 Secrets & Credentials
- [ ] Ningún secret hardcodeado en código
- [ ] `.env` añadido a `.gitignore`
- [ ] `.env.example` actualizado con todas las variables necesarias
- [ ] API keys rotadas si se expusieron
- [ ] JWT_SECRET tiene al menos 32 caracteres

### 🗄️ Base de Datos
- [ ] PostgreSQL: SSL/TLS habilitado en producción
- [ ] Contraseñas hasheadas (bcrypt/Argon2)
- [ ] Índices en campos de búsqueda
- [ ] Conexiones con usuario de solo lectura donde aplique

### 🌐 API & CORS
- [ ] CORS configurado solo para orígenes permitidos
- [ ] Rate limiting en endpoints públicos
- [ ] Validación de input en todos los endpoints
- [ ] SQL injection: solo queries parametrizadas
- [ ] XSS: escapar output HTML

### 🔑 Autenticación
- [ ] JWT con expiración corta (24h max)
- [ ] Refresh tokens con rotación
- [ ] Secure cookies (httpOnly, secure, sameSite)
- [ ] CORS no permite credenciales con wildcards

### 📦 Dependencias
- [ ] `npm audit` / `go mod tidy` sin vulnerabilidades
- [ ] Dependencias actualizadas
- [ ] No dependencias abandonadas

### 🚀 Deploy
- [ ] HTTPS forzado
- [ ] Headers de seguridad (HSTS, CSP, X-Frame-Options)
- [ ] Logs no exponen PII (Personal Identifiable Information)
- [ ] Backups configurados

## Emergencias

Si se expone un secret:
1. Rotar inmediatamente
2. Revocar tokens antiguos
3. Revisar logs de acceso
4. Notificar a usuarios si aplica

## Reportar Vulnerabilidades

Email: security@[tu-dominio].com
