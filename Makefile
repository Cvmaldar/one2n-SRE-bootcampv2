run:
	go run ./cmd/server/main.go

test:
	go test ./internal/handlers/

build:
	go build -o bin/student-api ./cmd/server

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