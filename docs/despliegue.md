# Guía de Despliegue — MyInquisitor

## Requisitos

- Servidor Linux (Debian/Ubuntu recomendado)
- Docker + Docker Compose
- Git
- Acceso al repositorio en GitHub

## Paso 1: Clonar el repositorio

```bash
git clone https://github.com/Eazir/MyInquisitor.git
cd MyInquisitor
```

## Paso 2: Crear archivo `.env`

Copiar desde el servidor actual o crear nuevo en `docker/.env`:

```bash
cd docker
nano .env
```

Variables requeridas:

| Variable | Descripción |
|----------|-------------|
| `DB_PASSWORD` | Contraseña de PostgreSQL |
| `JWT_SECRET` | Clave secreta para firmar JWT (cambiar = invalidar sesiones) |
| `ENCRYPTION_KEY` | Clave AES-256-GCM de 32 bytes (hex) para cifrar PII |
| `ENVIRONMENT` | `development` o `production` |

> **Importante**: `JWT_SECRET` y `ENCRYPTION_KEY` deben ser las mismas del servidor original
> si querés mantener sesiones activas y datos PII descifrables.

Ejemplo mínimo:

```env
DB_PASSWORD=myinquisitor_dev
JWT_SECRET=clave_secreta_segura_aqui
ENCRYPTION_KEY=clave_hex_32_bytes_aqui
ENVIRONMENT=development
```

## Paso 3: Restaurar backup de la base de datos

### Opción A: Servidor nuevo (sin datos)

```bash
# Copiar el backup al proyecto y levantarlo
docker compose up -d db
cat ../backup_2026-06-15.sql | docker exec -i myinquisitor-db psql -U myinquisitor_app myinquisitor
```

### Opción B: Migrar desde servidor actual

En el servidor actual hacer:

```bash
docker exec myinquisitor-db pg_dump -U myinquisitor_app myinquisitor > backup_$(date +%F).sql
```

Copiar ese `.sql` al nuevo servidor y restaurar:

```bash
docker compose up -d db
cat backup_2026-06-15.sql | docker exec -i myinquisitor-db psql -U myinquisitor_app myinquisitor
```

## Paso 4: Construir y levantar todo

```bash
docker compose up -d --build
```

Esto construye backend + frontend y aplica migraciones automáticamente.

## Paso 5: Verificar

```bash
docker compose ps
docker compose logs backend | tail -20
```

Acceder en `http://servidor:3000` (frontend) o `http://servidor:8080` (API).

## Comandos útiles

| Comando | Descripción |
|---------|-------------|
| `docker compose logs backend -f` | Ver logs del backend en tiempo real |
| `docker compose logs db -f` | Ver logs de PostgreSQL |
| `docker compose restart backend` | Reiniciar solo el backend |
| `docker compose down` | Bajar servicios (no borra datos) |
| `docker compose down -v` | **Baja servicios y BORRA el volumen de DB** |
| `docker exec -it myinquisitor-db psql -U myinquisitor_app myinquisitor` | Conectarse a la DB |
| `docker exec myinquisitor-db pg_dump -U myinquisitor_app myinquisitor > respaldo.sql` | Backup manual |

## Notas importantes

- El volumen de PostgreSQL persiste en `infoDbMyInquisitor`. Mientras no uses `-v` en `down`, los datos sobreviven.
- Las migraciones se aplican automáticamente al iniciar el backend. No hacerlas manualmente.
- Si cambiás `JWT_SECRET`, todos los usuarios deben volver a iniciar sesión.
- Si cambiás `ENCRYPTION_KEY`, los PII existentes (email, nombre) se vuelven ilegibles.
