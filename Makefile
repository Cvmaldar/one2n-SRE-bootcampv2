APP_NAME=student-api
VERSION=1.0.1

.PHONY: run test build migrate-up migrate-down docker-build docker-run \
	docker-compose-up docker-compose-down db wait-for-db run-api

run:
	go run ./cmd/server/main.go

test:
	go test ./...

build:
	go build -o bin/$(APP_NAME) ./cmd/server

migrate-up:
	migrate \
		-path migrations \
		-database "postgres://postgres:postgres@localhost:5432/studentdb?sslmode=disable" \
		up

migrate-down:
	migrate \
		-path migrations \
		-database "postgres://postgres:postgres@localhost:5432/studentdb?sslmode=disable" \
		down 1

docker-build:
	docker build -t $(APP_NAME):$(VERSION) .

docker-run:
	docker run --env-file .env \
		--name $(APP_NAME) \
		-p 8080:8080 \
		$(APP_NAME):$(VERSION)

docker-compose-up:
	docker compose up

docker-compose-down:
	docker compose down

db:
	docker compose up -d postgres

wait-for-db:
	@echo "Waiting for PostgreSQL..."
	@until docker exec student-postgres pg_isready -U postgres -d studentdb > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "PostgreSQL is ready."

run-api:
	docker compose up -d student-api

setup:
	$(MAKE) db
	$(MAKE) wait-for-db
	$(MAKE) migrate-up
	$(MAKE) docker-build
	$(MAKE) run-api