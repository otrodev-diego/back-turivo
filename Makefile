env ?= .env
-include $(env)
export

compose = docker compose

.PHONY: up down logs db psql migrate-up migrate-down swagger

up:
	$(compose) up -d

down:
	$(compose) down

logs:
	$(compose) logs -f

db:
	docker exec -it turivo-postgres sh

psql:
	docker exec -it turivo-postgres psql -U postgres -d turivo

migrate-up:
	migrate -database "$(DB_DSN)" -path migrations up

migrate-down:
	migrate -database "$(DB_DSN)" -path migrations down

swagger:
	swagger generate spec -o swagger.json --scan-models
