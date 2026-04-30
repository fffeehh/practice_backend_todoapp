include .env
export

export CURRENT_UID := $(shell id -u)
export CURRENT_GID := $(shell id -g)

# динамически задаем полный путь до корня проекта. С помощью $(shell команда) задаем выполнение shell команды
export PROJECT_ROOT=$(CURDIR)

env-up:
	@mkdir -p out/logs
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

# таргет для полного перезапуска бд (с удалением volume файлов)
env-cleanup:
	# предусматриваем подтверждение пользователя
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/n]: " ans; \
	if [ "$$ans" = "y" ] ; then \
		docker compose down todoapp-postgres port-forwarder &&\
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутсвует необходимый параметр seq" \
		exit 1; \
		fi
	@mkdir -p migrations
	@docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up	

migrate-down:
	@make migrate-action action=down


migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action" \
		exit 1; \
		fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/todoapp/main.go
