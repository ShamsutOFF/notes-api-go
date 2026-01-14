.PHONY: build-api build-client run test clean

# Сборка API сервера
build-api:
	@echo "Building API server..."
	@go build -o bin/api cmd/api/main.go

# Сборка CLI клиента (версия с Cobra)
build-client-cobra:
	@echo "Building CLI client with Cobra..."
	@go build -o bin/client-cobra cmd/client/main_cobra.go

# Запуск API сервера
run:
	@echo "Starting API server..."
	@go run cmd/api/main.go

# Запуск тестов
test:
	@echo "Running tests..."
	@go test ./...

# Очистка бинарных файлов
clean:
	@echo "Cleaning up..."
	@rm -rf bin/

# Сборка всего
build-all: build-api build-client-cobra

# Запуск примера использования
example:
	@echo "=== Example Usage ==="
	@echo "1. Start the API server in one terminal: make run"
	@echo "2. In another terminal, try these commands:"
	@echo "   $$ ./bin/client create \"Shopping\" \"Buy milk and bread\""
	@echo "   $$ ./bin-client list"
	@echo "   $$ ./bin/client get 1"
	@echo "   $$ ./bin/client update 1 \"Shopping list\" \"Buy milk, bread, eggs\""
	@echo "   $$ ./bin/client delete 1"