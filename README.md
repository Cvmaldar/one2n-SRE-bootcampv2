# Student CRUD REST API

A simple REST API built using **Go**, **Gin**, and **PostgreSQL** to perform CRUD (Create, Read, Update, Delete) operations on student records.

This project was built as part of the One2N SRE Bootcamp to learn REST API development, database integration, API best practices, Docker, and the Twelve-Factor App methodology.

---

## Features

* Create a student
* Get all students
* Get a student by ID
* Update student details
* Delete a student
* API versioning (`/api/v1`)
* PostgreSQL integration
* Database schema migrations
* Healthcheck endpoint
* Environment variable based configuration
* Unit tests using `go-sqlmock`
* Request logging
* Multi-stage Docker build
* Docker Compose
* PostgreSQL Docker volume
* Semantic versioned Docker images

---

## Tech Stack

* Go
* Gin
* PostgreSQL
* golang-migrate
* go-sqlmock
* Docker
* Docker Compose
* Postman
* Make

---

## Project Structure

```text
.
├── cmd
│   └── server
│       └── main.go
│
├── internal
│   ├── db
│   │   └── db.go
│   ├── handlers
│   │   ├── student.go
│   │   └── student_test.go
│   └── models
│       └── student.go
│
├── migrations
│   ├── 000001_create_students.up.sql
│   └── 000001_create_students.down.sql
│
├── postman
│   └── Student_API.postman_collection.json
│
├── .dockerignore
├── .env
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## Prerequisites

Before running the application, ensure the following are installed:

* Go 1.25.7 or later
* Docker
* Docker Compose
* golang-migrate CLI
* Make

PostgreSQL can either run locally or through Docker Compose.

---

# Local Development

## Environment Variables

Create a `.env` file in the project root.

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=studentdb
DB_SSLMODE=disable
```

The application reads configuration from environment variables.

The `.env` file is used for local development and is not copied into the Docker image.

---

## Install Dependencies

```bash
go mod tidy
```

---

## Running PostgreSQL with Docker

If PostgreSQL is not installed locally, it can be started using Docker:

```bash
docker run --name postgres \
-e POSTGRES_USER=postgres \
-e POSTGRES_PASSWORD=postgres \
-e POSTGRES_DB=studentdb \
-p 5432:5432 \
-d postgres
```

Check that PostgreSQL is running:

```bash
docker ps
```

---

## Run Database Migrations

Apply migrations:

```bash
make migrate-up
```

Rollback one migration:

```bash
make migrate-down
```

---

## Run the Application

```bash
make run
```

The server starts on:

```text
http://localhost:8080
```

---

## Run Tests

Execute all unit tests:

```bash
make test
```

The tests use `go-sqlmock`, so the handler tests do not require a running PostgreSQL instance.

---

# API Endpoints

| Method | Endpoint                | Description       |
| ------ | ----------------------- | ----------------- |
| GET    | `/healthcheck`          | Health check      |
| POST   | `/api/v1/students`      | Create student    |
| GET    | `/api/v1/students`      | Get all students  |
| GET    | `/api/v1/students/{id}` | Get student by ID |
| PUT    | `/api/v1/students/{id}` | Update student    |
| DELETE | `/api/v1/students/{id}` | Delete student    |

---

## Example Request

### Create Student

```http
POST /api/v1/students
```

Request body:

```json
{
  "name": "Chinmay",
  "age": 24,
  "email": "chinmay@example.com"
}
```

Response:

```json
{
  "id": 1,
  "name": "Chinmay",
  "age": 24,
  "email": "chinmay@example.com"
}
```

---

# Docker

The application uses a **multi-stage Dockerfile**.

The Dockerfile contains two stages:

```text
Stage 1: Builder
        ↓
Go compiler + dependencies
        ↓
Compile application
        ↓
student-api binary

Stage 2: Runtime
        ↓
Small Alpine image
        ↓
student-api binary
```

The Go compiler and build dependencies are not included in the final runtime image.

This helps reduce the final Docker image size.

---

## Build Docker Image

The current image version is:

```text
student-api:1.0.1
```

Build the image using:

```bash
make docker-build
```

This is equivalent to:

```bash
docker build -t student-api:1.0.1 .
```

The project intentionally uses a semantic version instead of the `latest` tag.

---

## Run Docker Image

Environment variables are injected into the container at runtime.

Run:

```bash
make docker-run
```

This executes:

```bash
docker run --env-file .env \
	--name student-api \
	-p 8080:8080 \
	student-api:1.0.1
```

The `.env` file is read by Docker on the host.

It is **not copied into the Docker image**.

The application therefore receives configuration through runtime environment variables.

---

# Docker Compose

Docker Compose runs the application and PostgreSQL together.

Start the complete environment:

```bash
make docker-compose-up
```

This creates:

```text
                  Docker Compose
                       │
             ┌─────────┴─────────┐
             │                   │
             ▼                   ▼
       student-api           postgres
       container             container
          :8080                 :5432
             │                   │
             └──── Docker ───────┘
                  network
                       │
                       ▼
                postgres_data
                   volume
```

The API connects to PostgreSQL using:

```env
DB_HOST=postgres
```

`postgres` is the PostgreSQL service name defined in Docker Compose.

Inside the Docker network, Docker's internal DNS resolves `postgres` to the PostgreSQL container.

This is different from local development, where:

```env
DB_HOST=localhost
```

is used.

---

## Stop Docker Compose

```bash
make docker-compose-down
```

This stops and removes the containers.

The PostgreSQL volume remains.

To remove the PostgreSQL volume as well:

```bash
docker compose down -v
```

**Warning:** removing the volume deletes the PostgreSQL data stored in it.

---

# Makefile Commands

| Command                    | Purpose                   |
| -------------------------- | ------------------------- |
| `make run`                 | Run the API locally       |
| `make test`                | Run all unit tests        |
| `make build`               | Build the Go binary       |
| `make migrate-up`          | Apply database migrations |
| `make migrate-down`        | Roll back one migration   |
| `make docker-build`        | Build the Docker image    |
| `make docker-run`          | Run the Docker image      |
| `make docker-compose-up`   | Start API and PostgreSQL  |
| `make docker-compose-down` | Stop API and PostgreSQL   |

---

# Docker Image Versioning

Docker images use Semantic Versioning:

```text
MAJOR.MINOR.PATCH
```

Examples:

```text
1.0.0
1.0.1
1.1.0
2.0.0
```

The current image is:

```text
student-api:1.0.1
```

The `latest` tag is intentionally avoided so that the exact application version being deployed is known.

---

# Docker Image Optimization

The project uses several measures to reduce the Docker image footprint:

### Multi-stage build

The Go compiler and build dependencies are kept in the builder stage and are not included in the runtime image.

### Alpine runtime image

The runtime uses a lightweight Alpine-based image.

### `.dockerignore`

The following files and directories are excluded from the Docker build context:

```text
.git
.env
.gitignore
README.md
*.md
tmp/
bin/
```

This prevents unnecessary files and sensitive configuration from being sent to Docker during the build.

---

# Healthcheck

The application provides:

```text
GET /healthcheck
```

Test it with:

```bash
curl http://localhost:8080/healthcheck
```

Expected response:

```json
{
  "status": "ok"
}
```

---

# Testing

The project includes unit tests for the API handlers using `go-sqlmock`.

This allows database interactions to be tested without requiring a running PostgreSQL instance.

Run:

```bash
go test ./...
```

or:

```bash
make test
```

---

# Development Workflow

### Local development

```bash
make migrate-up
make run
```

### Run tests

```bash
make test
```

### Build Docker image

```bash
make docker-build
```

### Run Docker image

```bash
make docker-run
```

### Run complete Docker environment

```bash
make docker-compose-up
```

### Stop Docker environment

```bash
make docker-compose-down
```
