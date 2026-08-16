# Student CRUD REST API — One2N SRE Bootcamp

A production-oriented Student CRUD REST API built using **Go, Gin, and PostgreSQL**, extended with Docker, Kubernetes, Helm, External Secrets, HashiCorp Vault, GitHub Actions, and Argo CD.

This project was built as part of the **One2N SRE Bootcamp** to progressively learn backend development, containerization, Kubernetes, CI/CD, GitOps, secrets management, and observability.

---

# Table of Contents

* [Overview](#overview)
* [Features](#features)
* [Architecture](#architecture)
* [Tech Stack](#tech-stack)
* [Project Structure](#project-structure)
* [Prerequisites](#prerequisites)
* [Local Development](#local-development)
* [Docker](#docker)
* [Kubernetes](#kubernetes)
* [Helm](#helm)
* [Secrets Management](#secrets-management)
* [CI Pipeline](#ci-pipeline)
* [GitOps with Argo CD](#gitops-with-argo-cd)
* [Deployment Flow](#deployment-flow)
* [API Endpoints](#api-endpoints)
* [Health Check](#health-check)
* [Testing](#testing)
* [Verification](#verification)
* [Troubleshooting](#troubleshooting)

---

# Overview

The project started as a simple Go REST API and was progressively evolved into a Kubernetes-based application with automated deployments.

The application provides CRUD operations for student records and uses PostgreSQL as its database.

The infrastructure evolved through the following stages:

```text
Go REST API
    ↓
PostgreSQL
    ↓
Docker
    ↓
Docker Compose
    ↓
Kubernetes
    ↓
Helm
    ↓
Vault + External Secrets
    ↓
GitHub Actions
    ↓
Argo CD / GitOps
```

The final deployment model uses Git as the source of truth for Kubernetes configuration.

---

# Features

## Application

* Create a student
* Get all students
* Get a student by ID
* Update student details
* Delete a student
* API versioning using `/api/v1`
* PostgreSQL integration
* Database schema migrations
* Health check endpoint
* Environment-variable based configuration
* Unit tests using `go-sqlmock`
* Request logging

## Containerization

* Multi-stage Docker build
* Docker Compose
* PostgreSQL Docker volume
* Semantic Docker image versioning
* `.dockerignore`
* One-command local development setup

## Kubernetes

* Kubernetes Deployment
* Kubernetes Service
* PostgreSQL deployment
* Namespace-based application isolation
* Node selectors
* Database migrations
* Helm-based deployment

## Secrets Management

* HashiCorp Vault
* External Secrets Operator
* `ClusterSecretStore`
* Secrets synchronized from Vault into Kubernetes

## CI/CD and GitOps

* GitHub Actions
* Self-hosted GitHub Actions runner
* Automated build
* Automated tests
* Automated linting
* Docker image build and push
* Automated Helm image-tag update
* Automated Git commit
* Argo CD
* App-of-Apps pattern
* Automated synchronization
* Self-healing
* Pruning

---

# Architecture

The final architecture is:

```text
                         GitHub
                           │
                           │ source code
                           ▼
                  GitHub Actions CI
                           │
              ┌────────────┼────────────┐
              │            │            │
            Build        Test         Lint
              │
              ▼
        Docker Build
              │
              ▼
        Docker Hub
              │
              │ image
              ▼
       Update Helm values
              │
              ▼
        Git commit/push
              │
              ▼
          Argo CD
              │
          Auto Sync
              │
              ▼
        Kubernetes
              │
       ┌──────┼──────────────┐
       │      │              │
       ▼      ▼              ▼
  student-api PostgreSQL External Secrets
                              │
                              ▼
                            Vault
```

---

# Application Architecture

The Go application is structured using a separation between handlers, models, and database access.

```text
Client
   │
   ▼
Gin HTTP Server
   │
   ▼
Handlers
   │
   ▼
Database Layer
   │
   ▼
PostgreSQL
```

---

# Tech Stack

## Application

* Go
* Gin
* PostgreSQL
* golang-migrate
* go-sqlmock

## Containerization

* Docker
* Docker Compose

## Kubernetes

* Kubernetes
* Minikube
* Helm

## Secrets

* HashiCorp Vault
* External Secrets Operator

## CI/CD

* GitHub Actions
* Self-hosted GitHub Actions runner
* Docker Hub

## GitOps

* Argo CD

## Tooling

* GNU Make
* Postman
* Git

---

# Project Structure

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
├── helm
│   ├── student-api
│   ├── database
│   └── external-secrets
│
├── argocd
│   ├── install-values.yaml
│   ├── root-app.yaml
│   ├── repo-secret.yaml
│   └── applications
│       ├── student-api.yaml
│       ├── database.yaml
│       └── external-secrets-config.yaml
│
├── .github
│   └── workflows
│       └── ci.yml
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

The following tools should be installed:

* Go 1.25.7 or later
* Docker
* Docker Compose
* GNU Make
* golang-migrate CLI
* kubectl
* Minikube
* Helm

PostgreSQL does not need to be installed locally because it can run as a Docker container.

---

# Local Development

## Environment Variables

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

The `.env` file is intended for local development and should not be included in the Docker image.

---

# Install Go Dependencies

```bash
go mod tidy
```

---

# One-Click Local Development

The complete local environment can be started using:

```bash
make setup
```

The setup performs:

```text
make setup
    │
    ├── Start PostgreSQL
    │
    ├── Wait for PostgreSQL
    │
    ├── Run migrations
    │
    ├── Build Docker image
    │
    └── Start REST API
```

---

# Makefile Targets

## Start PostgreSQL

```bash
make db
```

This starts PostgreSQL using Docker Compose.

PostgreSQL is exposed on:

```text
localhost:5432
```

## Wait for PostgreSQL

```bash
make wait-for-db
```

This waits until PostgreSQL is ready before running migrations.

This avoids:

```text
PostgreSQL container starts
        ↓
PostgreSQL still initializing
        ↓
Migration starts
        ↓
Connection refused
```

## Run migrations

```bash
make migrate-up
```

Rollback the latest migration:

```bash
make migrate-down
```

## Build Docker image

```bash
make docker-build
```

## Run API

```bash
make run-api
```

## Run tests

```bash
make test
```

## Start complete environment

```bash
make setup
```

---

# Docker

The application uses a multi-stage Docker build.

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
    ├── Lightweight runtime
    └── Application binary
```

The Go compiler and build dependencies are excluded from the final runtime image.

---

# Docker Image Versioning

Images use semantic versioning:

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

The project avoids the `latest` tag so that deployments always reference an explicit version.

Example:

```text
cvmaldar234/student-api:1.0.8
```

---

# Docker Compose

Docker Compose manages:

```text
postgres
student-api
```

The API connects to PostgreSQL using the Compose service name:

```env
DB_HOST=postgres
```

Inside the Compose network:

```text
student-api
      │
      │ postgres:5432
      ▼
   postgres
```

When the Go application runs directly on the host, PostgreSQL is accessed through:

```text
localhost:5432
```

---

# PostgreSQL Persistence

PostgreSQL uses a named Docker volume:

```yaml
volumes:
  - postgres_data:/var/lib/postgresql/data
```

Stopping the PostgreSQL container does not automatically remove the data.

To stop the environment:

```bash
make docker-compose-down
```

To remove the PostgreSQL volume:

```bash
docker compose down -v
```

> Warning: removing the volume deletes PostgreSQL data.

---

# Kubernetes

The application is deployed to Kubernetes using Helm.

The project uses Minikube with multiple nodes.

Nodes are assigned different responsibilities using labels.

Example:

```text
application
dependent_services
database
```

Node selectors are used to place workloads on the appropriate nodes.

Check nodes:

```bash
kubectl get nodes --show-labels
```

---

# Kubernetes Namespaces

The main application namespace is:

```text
student-api
```

Argo CD runs in:

```text
argocd
```

External Secrets Operator runs in:

```text
external-secrets-system
```

Vault runs in:

```text
vault
```

---

# Helm

Helm is used to package and deploy the Kubernetes resources.

The project contains separate charts for:

```text
helm/
├── student-api/
├── database/
└── external-secrets/
```

Example application values:

```yaml
replicaCount: 2

image:
  repository: cvmaldar234/student-api
  tag: "1.0.8"
  pullPolicy: Always

service:
  type: ClusterIP
  port: 8080
  targetPort: 8080

nodeSelector:
  type: application
```

The Helm `values.yaml` file is the desired configuration used by Argo CD.

---

# Database

The database is deployed separately from the API.

The application connects to PostgreSQL using:

```text
DB_HOST=student-api-db
DB_PORT=5432
```

The database deployment is managed through the database Helm chart.

Verify:

```bash
kubectl get pods -n student-api
```

---

# Secrets Management

The project uses:

```text
HashiCorp Vault
       ↓
External Secrets Operator
       ↓
Kubernetes Secret
       ↓
student-api
```

The External Secrets Operator installation provides the Kubernetes CRDs and controllers required for External Secrets resources.

The configuration chart creates resources such as:

```text
ClusterSecretStore
ExternalSecret
```

The application does not need to store database credentials directly in Git.

---

# Vault and External Secrets Flow

```text
                 Vault
                   │
                   │ secret
                   ▼
          ClusterSecretStore
                   │
                   ▼
           ExternalSecret
                   │
                   ▼
          Kubernetes Secret
                   │
                   ▼
             student-api
```

Verify External Secrets components:

```bash
kubectl get pods -n external-secrets-system
```

Verify External Secrets resources:

```bash
kubectl get externalsecrets -A
kubectl get clustersecretstores
```

---

# CI Pipeline

GitHub Actions is used for continuous integration.

The workflow runs on a self-hosted GitHub Actions runner.

The CI pipeline performs:

```text
Checkout
   ↓
Build
   ↓
Test
   ↓
Lint
   ↓
Docker Login
   ↓
Docker Build
   ↓
Docker Push
   ↓
Update Helm values.yaml
   ↓
Commit
   ↓
Push to GitHub
```

---

# CI Image Tagging

The image tag is generated from the GitHub Actions run number:

```bash
IMAGE_TAG=1.0.${{ github.run_number }}
```

For example:

```text
GitHub Actions run #8
        ↓
IMAGE_TAG=1.0.8
```

The Docker image is then:

```text
cvmaldar234/student-api:1.0.8
```

The same version is written to:

```text
helm/student-api/values.yaml
```

Example:

```yaml
image:
  repository: cvmaldar234/student-api
  tag: "1.0.8"
```

This keeps the container registry and Kubernetes desired state consistent.

---

# Why `contents: write` is Required

The GitHub Actions workflow updates:

```text
helm/student-api/values.yaml
```

and commits the change back to the repository.

Therefore the workflow needs:

```yaml
permissions:
  contents: write
```

Without write permission, the workflow can read the repository but cannot push the updated Helm values.

---

# GitOps with Argo CD

Argo CD is used for continuous delivery using GitOps.

The fundamental principle is:

> Git is the source of truth for the desired Kubernetes state.

The deployment flow is:

```text
Git
 ↓
Argo CD
 ↓
Helm
 ↓
Kubernetes
```

Argo CD continuously compares the desired state stored in Git with the live state of the Kubernetes cluster.

---

# Argo CD Installation

Create the namespace:

```bash
kubectl create namespace argocd
```

Add the Argo Helm repository:

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update
```

Install Argo CD:

```bash
helm install argocd argo/argo-cd \
  -n argocd \
  -f argocd/install-values.yaml
```

Verify:

```bash
kubectl get pods -n argocd
```

Check where the components are scheduled:

```bash
kubectl get pods -n argocd -o wide
```

The Argo CD components are configured to run on the node labelled:

```text
type=dependent_services
```

---

# Declarative Argo CD Configuration

Argo CD configuration is stored in Git.

```text
argocd/
├── install-values.yaml
├── root-app.yaml
├── repo-secret.yaml
└── applications/
    ├── student-api.yaml
    ├── database.yaml
    └── external-secrets-config.yaml
```

The repository secret allows Argo CD to access the private GitHub repository.

Credentials should not be committed as plaintext secrets.

---

# App-of-Apps Pattern

The project uses the Argo CD App-of-Apps pattern.

The root application manages the child applications.

```text
root-app
   │
   ├── student-api
   ├── database
   └── external-secrets-config
```

The root application points to:

```text
argocd/applications
```

Example:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: root-app
  namespace: argocd

spec:
  project: default

  source:
    repoURL: https://github.com/Cvmaldar/one2n-SRE-bootcamp.git
    targetRevision: main
    path: argocd/applications

  destination:
    server: https://kubernetes.default.svc
    namespace: argocd

  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

Apply it:

```bash
kubectl apply -f argocd/root-app.yaml
```

---

# Argo CD Applications

Check all applications:

```bash
kubectl get applications -n argocd
```

Expected applications:

```text
database
external-secrets-config
root-app
student-api
```

Expected state:

```text
SYNC STATUS: Synced
HEALTH STATUS: Healthy
```

---

# Automated Sync

Argo CD is configured with:

```yaml
syncPolicy:
  automated:
    prune: true
    selfHeal: true
```

## Automated Sync

When Git changes, Argo CD automatically synchronizes the Kubernetes cluster.

## Prune

Resources removed from the Git desired state can be removed from the cluster.

## Self Heal

If the live Kubernetes configuration is changed manually and differs from Git, Argo CD can reconcile it back to the desired Git state.

---

# Helm as the Source of Truth

Argo CD uses the Helm charts and their `values.yaml` files stored in Git.

For the Student API:

```text
helm/student-api/
```

The image is configured in:

```text
helm/student-api/values.yaml
```

For example:

```yaml
image:
  repository: cvmaldar234/student-api
  tag: "1.0.8"
```

Argo CD renders the Helm chart using these values and deploys the resulting Kubernetes resources.

---

# Complete GitOps Deployment Flow

A normal application deployment follows:

```text
Developer
    │
    │ git push
    ▼
GitHub
    │
    ▼
GitHub Actions
    │
    ├── Build
    ├── Test
    ├── Lint
    ├── Docker Build
    └── Docker Push
            │
            ▼
       Docker Hub
            │
            ▼
    Update values.yaml
            │
            ▼
       Git commit
            │
            ▼
       Git push
            │
            ▼
         Argo CD
            │
            │ detects Git change
            ▼
        Helm rendering
            │
            ▼
        Auto Sync
            │
            ▼
        Kubernetes
            │
            ▼
      New application pods
```

No manual `helm upgrade` or `kubectl apply` is required for normal application deployments after Argo CD is configured.

---

# Verify GitOps Deployment

Check Argo CD:

```bash
kubectl get applications -n argocd
```

Monitor synchronization:

```bash
kubectl get applications -n argocd -w
```

Monitor the application:

```bash
kubectl get pods -n student-api -w
```

Check the deployed image:

```bash
kubectl get deployment student-api \
  -n student-api \
  -o jsonpath='{.spec.template.spec.containers[0].image}'
```

Expected:

```text
cvmaldar234/student-api:<image-tag>
```

The image tag should match:

```text
helm/student-api/values.yaml
```

---

# One-Click Deployment

After the GitOps setup is complete, deployment becomes:

```text
Code change
     ↓
git push
     ↓
GitHub Actions
     ↓
CI
     ↓
Docker image
     ↓
Helm values update
     ↓
Git commit
     ↓
Argo CD
     ↓
Auto Sync
     ↓
Kubernetes
```

The developer does not need to manually execute deployment commands.

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

# Health Check

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

# Testing

Run all unit tests:

```bash
make test
```

The tests use `go-sqlmock`, so handler tests do not require a running PostgreSQL instance.

Run linting:

```bash
make lint
```

Build the application:

```bash
make build
```

---

# Kubernetes Verification

Check all namespaces:

```bash
kubectl get ns
```

Check application resources:

```bash
kubectl get all -n student-api
```

Check database:

```bash
kubectl get pods -n student-api
```

Check External Secrets:

```bash
kubectl get pods -n external-secrets-system
```

Check Vault:

```bash
kubectl get pods -n vault
```

Check Argo CD:

```bash
kubectl get pods -n argocd
```

Check Argo applications:

```bash
kubectl get applications -n argocd
```

---

# Troubleshooting

## Check Helm releases

```bash
helm list -A
```

## Check Argo CD applications

```bash
kubectl get applications -n argocd
```

## Describe an Argo CD application

```bash
kubectl describe application student-api -n argocd
```

## Check application events

```bash
kubectl get events -n student-api --sort-by=.lastTimestamp
```

## Check application logs

```bash
kubectl logs -n student-api deployment/student-api
```

## Check Helm values

```bash
helm get values student-api -n student-api
```

## Check deployed image

```bash
kubectl get deployment student-api \
  -n student-api \
  -o jsonpath='{.spec.template.spec.containers[0].image}'
```

---

# Development Workflow

For local development:

```bash
make setup
```

Run tests:

```bash
make test
```

Build:

```bash
make build
```

Run the API:

```bash
make run
```

Stop Docker Compose:

```bash
make docker-compose-down
```

---

# Final Architecture

The completed project demonstrates the progression from application development to an automated Kubernetes GitOps deployment platform:

```text
                        GitHub
                           │
                    Source Repository
                           │
                           ▼
                  GitHub Actions CI
                           │
              ┌────────────┼────────────┐
              │            │            │
            Build        Test         Lint
              │
              ▼
         Docker Build
              │
              ▼
          Docker Hub
              │
              ▼
      Helm values.yaml update
              │
              ▼
          Git commit
              │
              ▼
           Argo CD
              │
        GitOps Reconciliation
              │
              ▼
         Kubernetes
              │
       ┌──────┼──────────────┐
       │      │              │
       ▼      ▼              ▼
    API     Database    External Secrets
                             │
                             ▼
                           Vault
```

The project demonstrates:

```text
Go
 ↓
REST API
 ↓
PostgreSQL
 ↓
Docker
 ↓
Kubernetes
 ↓
Helm
 ↓
Vault
 ↓
External Secrets
 ↓
GitHub Actions
 ↓
GitOps
 ↓
Argo CD
```

This completes the application, containerization, Kubernetes, secrets-management, CI, and GitOps deployment workflow.
