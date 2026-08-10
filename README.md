# Student CRUD REST API

A simple REST API built using **Go**, **Gin**, and **PostgreSQL** to perform CRUD (Create, Read, Update, Delete) operations on student records.

This project was built as part of the One2N SRE Bootcamp to learn REST API development, database integration, API best practices, Docker, Docker Compose, GNU Make, and the Twelve-Factor App methodology.

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
* One-command local development setup

---

## Tech Stack

* Go
* Gin
* PostgreSQL
* golang-migrate
* go-sqlmock
* Docker
* Docker Compose
* Make
* Postman

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

# Prerequisites

The following tools must already be installed:

* Go 1.25.7 or later
* Docker
* Docker Compose
* GNU Make
* golang-migrate CLI

PostgreSQL **does not need to be installed locally** because PostgreSQL runs as a Docker container.

---

# Environment Variables

For local development, create a `.env` file in the project root:

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=studentdb
DB_SSLMODE=disable
```

The application reads its configuration from environment variables.

The `.env` file is used for local development and is not copied into the Docker image.

---

# Install Go Dependencies

```bash
go mod tidy
```

---

# Step 3: One-Click Local Development Setup

The goal of this milestone is to allow a developer to start the complete application with minimal steps.

The application consists of:

```text
                 Docker Compose
                      │
             ┌────────┴────────┐
             │                 │
             ▼                 ▼
        student-api         postgres
        container           container
          :8080                :5432
             │                 │
             └──── Docker ─────┘
                  network
```

PostgreSQL data is persisted using a Docker named volume:

```text
postgres_data
```

---

# Start the Complete Environment

The recommended command is:

```bash
make setup
```

This performs the required steps in order:

```text
make setup
    │
    ├── Start PostgreSQL
    │
    ├── Wait for PostgreSQL to become ready
    │
    ├── Run database migrations
    │
    ├── Build REST API Docker image
    │
    └── Start REST API container
```

The individual Make targets can also be executed separately.

---

# Makefile Targets

## Start PostgreSQL

```bash
make db
```

This starts only the PostgreSQL service using Docker Compose:

```bash
docker compose up -d postgres
```

PostgreSQL is exposed on:

```text
localhost:5432
```

---

## Wait for PostgreSQL

```bash
make wait-for-db
```

This checks whether PostgreSQL is ready to accept connections before migrations are executed.

This prevents the following race condition:

```text
PostgreSQL container started
        ↓
PostgreSQL still initializing
        ↓
Migration starts too early
        ↓
Connection refused
```

The readiness check waits until PostgreSQL reports that it is ready.

---

## Run Database Migrations

```bash
make migrate-up
```

This applies the database migrations:

```bash
migrate \
    -path migrations \
    -database "postgres://postgres:postgres@localhost:5432/studentdb?sslmode=disable" \
    up
```

To roll back one migration:

```bash
make migrate-down
```

---

## Build REST API Docker Image

```bash
make docker-build
```

This builds the Docker image using the multi-stage Dockerfile.

The image is tagged using Semantic Versioning:

```text
student-api:1.0.1
```

The project intentionally avoids using the `latest` tag.

---

## Run REST API Container

```bash
make run-api
```

This starts the API using Docker Compose:

```bash
docker compose up -d student-api
```

The API is exposed on:

```text
http://localhost:8080
```

---

# Complete Setup Flow

The complete setup can be performed with:

```bash
make setup
```

Internally, this performs:

```text
make db
    ↓
PostgreSQL container
    ↓
make wait-for-db
    ↓
PostgreSQL ready
    ↓
make migrate-up
    ↓
students table created
    ↓
make docker-build
    ↓
student-api:1.0.1
    ↓
make run-api
    ↓
REST API running
```

This is the main objective of Step 3.

---

# Docker Compose

The Docker Compose configuration manages two services:

```text
postgres
student-api
```

The API connects to PostgreSQL using:

```env
DB_HOST=postgres
```

The value `postgres` is the Docker Compose service name.

Inside the Docker Compose network:

```text
student-api
      │
      │ postgres:5432
      ▼
postgres
```

The API should therefore **not** use:

```env
DB_HOST=localhost
```

when running inside Docker.

For local execution using `go run`, `localhost` is used because the Go process runs directly on the host.

---

# PostgreSQL Persistence

PostgreSQL uses a named Docker volume:

```yaml
volumes:
  - postgres_data:/var/lib/postgresql/data
```

This means that removing the PostgreSQL container does not automatically remove the database data.

To stop the environment:

```bash
make docker-compose-down
```

To stop the environment and remove the PostgreSQL volume:

```bash
docker compose down -v
```

**Warning:** removing the volume deletes the PostgreSQL data.

---

# Docker Compose Commands

Start the complete Compose environment manually:

```bash
make docker-compose-up
```

Stop the Compose environment:

```bash
make docker-compose-down
```

For normal development, prefer:

```bash
make setup
```

because it also performs the database startup, readiness check, migration, image build, and API startup sequence.

---

# Docker

The application uses a multi-stage Dockerfile.

```text
Stage 1: Builder
        │
        ├── Go compiler
        ├── Dependencies
        └── Application source
                 │
                 ▼
           Go application binary
                 │
                 ▼
Stage 2: Runtime
        │
        ├── Small runtime image
        └── Application binary
```

The Go compiler and build dependencies are not included in the final runtime image.

This reduces the final image size.

---

# Docker Image Versioning

Images use Semantic Versioning:

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

Current image:

```text
student-api:1.0.1
```

The `latest` tag is intentionally avoided so that the exact application version is known.

---

# Docker Image Optimization

The project uses several measures to reduce the Docker image footprint:

### Multi-stage build

The Go compiler and build dependencies are kept in the builder stage and are not included in the runtime image.

### Alpine runtime image

A lightweight Alpine-based image is used for the runtime stage.

### `.dockerignore`

Unnecessary files are excluded from the Docker build context.

For example:

```text
.git
.env
.gitignore
README.md
*.md
tmp/
bin/
```

The `.env` file is particularly important because application secrets and local configuration should not be baked into the Docker image.

---

# Local Go Development

The API can also be run directly without Docker:

```bash
make run
```

This executes:

```bash
go run ./cmd/server/main.go
```

For this mode, PostgreSQL must already be reachable on:

```text
localhost:5432
```

---

# Testing

Run all unit tests:

```bash
make test
```

The tests use `go-sqlmock`, so handler tests do not require a running PostgreSQL instance.

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

# Healthcheck

Test the application:

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

# Example API Request

## Create Student

```http
POST /api/v1/students
```

Request:

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

# Development Workflow

For the complete Docker-based development environment:

```bash
make setup
```

Run tests:

```bash
make test
```

Stop the environment:

```bash
make docker-compose-down
```

For local Go development:

```bash
make run
```

---

# Step 3 Completion Checklist

The Step 3 requirements are covered as follows:

| Requirement                                   | Implementation         |
| --------------------------------------------- | ---------------------- |
| API + dependent services using Docker Compose | `docker-compose.yml`   |
| Start DB container                            | `make db`              |
| Run DB migrations                             | `make migrate-up`      |
| Build REST API image                          | `make docker-build`    |
| Run REST API container                        | `make run-api`         |
| DB readiness check                            | `make wait-for-db`     |
| Complete startup workflow                     | `make setup`           |
| Makefile                                      | `Makefile`             |
| Docker Compose instructions                   | This README            |
| Make target execution order                   | `make setup` section   |
| Persistent PostgreSQL data                    | `postgres_data` volume |
| Semantic image version                        | `student-api:1.0.1`    |

---

## Step 3 Architecture

The final local development workflow is:

```text
                    make setup
                        │
                        ▼
                 Docker Compose
                        │
             ┌──────────┴──────────┐
             │                     │
             ▼                     │
          PostgreSQL               │
             │                     │
       healthcheck                 │
             │                     │
             ▼                     │
         migrations                │
             │                     │
             └──────────┐          │
                        ▼          │
                 Docker build      │
                        │          │
                        ▼          │
                student-api:1.0.1 │
                        │          │
                        ▼          ▼
                   student-api ↔ postgres
                       :8080       :5432
```

This completes the main **Docker Compose + GNU Make one-click development setup** requirements for Step 3.
