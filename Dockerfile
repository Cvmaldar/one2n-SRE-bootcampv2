# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o student-api ./cmd/server


# Stage 2: Run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/student-api .

EXPOSE 8080

CMD ["./student-api"]