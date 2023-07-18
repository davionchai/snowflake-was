.PHONY: up check-logs down build

up:
	docker compose up -d

check-logs:
	docker compose logs -f -t

down:
	docker compose down --remove-orphans -v

build:
	docker build --tag snowflake-was:latest --file Dockerfile ..
