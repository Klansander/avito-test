
.PHONY: help build up down clean

all:  swag build up

# Описание целей
help: ## Отображает список доступных команд
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Собирает приложения с помощью Docker Compose
	docker-compose build

up: ## Запускает приложения с помощью Docker Compose
	docker-compose up -d

swag: ## Генерация документации
	swag init -g ./app/cmd/app/main.go -o ./app/docs

down: ## Останавливает приложения
	docker-compose down

clean: ## Останавливает и удаляет все контейнеры и тома
	docker-compose down -v --remove-orphans

test.integration: ## Запуск тестов
	docker-compose   -f app/tests/docker-compose-test.yml down
	docker-compose -f app/tests/docker-compose-test.yml build
	docker-compose   -f app/tests/docker-compose-test.yml up -d
	GIN_MODE=release go test  -v ./app/tests/
	docker-compose   -f app/tests/docker-compose-test.yml down
lint: ## Запуск линтера
	golangci-lint -c .golangci.yml run ./app/...

