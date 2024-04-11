
.PHONY: help build up down clean

all: clean build up

# Описание целей
help: ## Отображает список доступных команд
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Собирает приложения с помощью Docker Compose
	docker-compose build

up: ## Запускает приложения с помощью Docker Compose
	docker-compose up -d

down: ## Останавливает приложения
	docker-compose down

clean: ## Останавливает и удаляет все контейнеры и тома
	docker-compose down -v --remove-orphans
test.integration:
	docker-compose   -f app/tests/docker-compose-test.yml down
	docker-compose -f app/tests/docker-compose-test.yml build
	docker-compose   -f app/tests/docker-compose-test.yml up -d
	GIN_MODE=release go test  -v ./app/tests/




#	GIN_MODE=release go test -v ./tests/
#	docker-compose down -v --remove-orphans