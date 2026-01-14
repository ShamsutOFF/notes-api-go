-- Этот скрипт выполнится автоматически при запуске контейнера PostgreSQL
-- Таблицы создаст GORM через AutoMigrate, но можно добавить дополнительные настройки

-- Создаем расширение для UUID если нужно
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаем схему если не существует
CREATE SCHEMA IF NOT EXISTS public;

-- Комментарий к базе данных
COMMENT ON DATABASE notesdb IS 'Notes API Database';