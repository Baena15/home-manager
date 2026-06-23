# Project Maturity Roadmap

> Flujo de madurez profesional para proyectos generados desde Plantilla Madre.
> Trigger: Cuando el usuario dice "roadmap de madurez", "fases del proyecto", "qué me falta para producción", o al finalizar un MVP.

---

## Fases de madurez recomendadas

Los proyectos de Gentleman Stack deben crecer en este orden. No saltar fases.

### Fase 0: Scaffold + MVP funcional
- Generar proyecto con `scripts/setup.ps1`
- Modelos, vistas, templates, URLs funcionales
- Auth básico (login/registro/perfil)
- Deploy a Railway/Render
- **Objetivo**: Producto usable por usuarios reales

### Fase 1: Tests automatizados
- `pytest-django` o `go test` según stack
- Tests de modelos, forms, vistas principales
- Cobertura mínima del 70% en lógica de negocio
- **Por qué**: En empresa, sin tests no hay merge

### Fase 2: Docker + Docker Compose
- `Dockerfile` para la app
- `docker-compose.yml` con DB + Redis
- Entorno de desarrollo idéntico al de producción
- **Por qué**: "En mi máquina sí funciona" no existe con Docker

### Fase 3: API REST (si aplica)
- Django REST Framework / FastAPI / chi + JSON
- Frontend separado (React/Vue) o app móvil
- Documentación de endpoints (OpenAPI/Swagger)
- **Por qué**: Escalabilidad, apps móviles, equipos separados frontend/backend

### Fase 4: CI/CD (GitHub Actions)
- Pipeline: Tests → Lint → Build → Deploy
- Deploy automático en push a `main`
- Block de deploy si tests fallan
- **Por qué**: Deploy seguro y repetible

### Fase 5: Performance (Redis, Celery, etc.)
- Caché con Redis en vistas de lectura frecuente
- Tareas en background con Celery / Go workers
- Optimización de queries (N+1)
- **Por qué**: El MVP aguanta 10 usuarios; esto aguanta 10.000

### Fase 6: Observabilidad
- Logging estructurado
- Métricas básicas (requests, errores)
- Alertas si el servicio cae
- **Por qué**: En producción sin logs no sabes qué ha roto

---

## Checklist de entrega mínima

Antes de considerar un proyecto "listo para entregar académica":
- [ ] MVP funcional con todas las features pedidas
- [ ] Deploy en producción funcionando
- [ ] README con URL, credenciales y cómo ejecutar en local
- [ ] `DEBUG=False` (o equivalente en otros stacks)
- [ ] Variables de entorno configuradas (no secrets hardcodeados)

Antes de considerar un proyecto "listo para trabajo profesional":
- [ ] Tests automatizados
- [ ] Docker funcionando
- [ ] CI/CD activo
- [ ] API REST documentada (si tiene frontend separado)
- [ ] Manejo de errores 500 con templates/pages personalizadas

---

## Cuándo NO aplicar una fase

| Fase | Omitir si... |
|------|-------------|
| API REST | El proyecto es solo server-side templates (proyecto académico simple) |
| Celery | No hay tareas lentas (emails, PDFs, reportes) |
| Redis | El proyecto tiene < 100 usuarios concurrentes |
| Docker | Es un MVP de 1 semana para demostración interna |

---

## Convención de nombre de ramas/fases

Si se trabaja por fases, usar prefijos en commits o ramas:
```
feat/mvp-reservas
feat/tests-reservas
feat/docker-setup
feat/api-rest
feat/ci-cd
feat/redis-cache
```
