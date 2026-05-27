-- ============================================================
-- Setup inicial de PostgreSQL para MyInquisitor
-- Este script se ejecuta UNA SOLA VEZ cuando el volumen está vacío
-- ============================================================

-- Crear usuario limitado para la aplicación (runtime)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'myinquisitor_app') THEN
        CREATE USER myinquisitor_app WITH PASSWORD '${APP_DB_PASSWORD:-myinquisitor_dev}';
    END IF;
END
$$;

-- Crear base de datos si no existe
SELECT 'CREATE DATABASE myinquisitor'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'myinquisitor')\gexec

-- Conceder permisos al usuario de la app
GRANT CONNECT ON DATABASE myinquisitor TO myinquisitor_app;
GRANT USAGE ON SCHEMA public TO myinquisitor_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO myinquisitor_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO myinquisitor_app;

-- Nota: las migraciones (CREATE TABLE) se ejecutan con el usuario postgres
--       o myinquisitor_app solo si tiene permisos DDL adicionales.
--       En producción, las migraciones NO deben usar myinquisitor_app.
