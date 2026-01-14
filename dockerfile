# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/api

# Production stage
FROM alpine:3.18

WORKDIR /root/

# Устанавливаем зависимости времени выполнения
RUN apk --no-cache add ca-certificates

# Копируем бинарник из builder stage
COPY --from=builder /app/main .

# Создаем директорию для хранения данных
RUN mkdir -p /root/storage

# Экспонируем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]