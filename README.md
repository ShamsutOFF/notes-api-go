# Notes API and CLI Client

Простое CRUD API для управления заметками с CLI клиентом.

## Структура проекта
notes-api/
├── cmd/
│ ├── api/ # API сервер
│ └── client/ # CLI клиент
├── internal/ # Внутренние пакеты
│ ├── domain/ # Доменные модели
│ ├── repository/ # Репозитории (хранилище)
│ ├── service/ # Бизнес-логика
│ └── handler/ # HTTP обработчики
├── storage/ # Файловое хранилище
└── Makefile # Утилиты сборки

text

## Быстрый старт

### 1. Установка зависимостей
```bash
go mod tidy
```

### 2. Сборка
```bash
make build-all
```

### 3. Запуск
```bash
make run
```

Сервер будет доступен по адресу: http://localhost:8080

API Endpoints
POST /api/notes - Создать заметку

GET /api/notes - Получить все заметки

GET /api/notes/:id - Получить заметку по ID

PUT /api/notes/:id - Обновить заметку

DELETE /api/notes/:id - Удалить заметку