# Student CRUD REST API

A simple REST API built using **Go**, **Gin**, and **PostgreSQL** to perform CRUD (Create, Read, Update, Delete) operations on student records.

This project was built as part of the One2N SRE Bootcamp to learn REST API development, database integration, API best practices, and the Twelve-Factor App methodology.

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

---

## Tech Stack

* Go
* Gin
* PostgreSQL
* golang-migrate
* go-sqlmock
* Docker (PostgreSQL)
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
├── .env
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## Prerequisites

Before running the application, ensure the following are installed:

* Go 1.24 or later
* Docker
* PostgreSQL (or PostgreSQL running in Docker)
* golang-migrate CLI
* Make

---

## Running PostgreSQL with Docker

```bash
docker run --name postgres \
-e POSTGRES_USER=postgres \
-e POSTGRES_PASSWORD=postgres \
-e POSTGRES_DB=studentdb \
-p 5432:5432 \
-d postgres
```

---

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

---

## Install Dependencies

```bash
go mod tidy
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

---

## API Endpoints

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

Request Body

```json
{
  "name": "Chinmay",
  "age": 24,
  "email": "chinmay@example.com"
}
```

Response

```json
{
  "id": 1,
  "name": "Chinmay",
  "age": 24,
  "email": "chinmay@example.com"
}
```

---

## Testing

The project includes unit tests for the API handlers using `go-sqlmock`, allowing database interactions to be tested without requiring a running PostgreSQL instance.

Run all tests:

```bash
go test ./...
```

---
