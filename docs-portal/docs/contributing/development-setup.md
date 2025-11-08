<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 3
---

# Development Environment Setup

This guide will help you set up a complete local development environment for contributing to DictaMesh.

## Prerequisites

### Required Software

Before you begin, ensure you have the following installed:

#### 1. Go

DictaMesh is written in Go. Install Go 1.21 or later:

```bash
# On macOS with Homebrew
brew install go

# On Ubuntu/Debian
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Verify installation
go version
```

**Configure Go environment:**

```bash
# Add to ~/.bashrc or ~/.zshrc
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export GO111MODULE=on
```

#### 2. Docker and Docker Compose

Required for running infrastructure services:

```bash
# On macOS
brew install docker docker-compose

# On Ubuntu/Debian
sudo apt-get update
sudo apt-get install docker.io docker-compose

# Verify installation
docker --version
docker-compose --version

# Add your user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker
```

#### 3. Git

```bash
# On macOS
brew install git

# On Ubuntu/Debian
sudo apt-get install git

# Configure Git
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

#### 4. Make

```bash
# On macOS
xcode-select --install

# On Ubuntu/Debian
sudo apt-get install build-essential
```

### Recommended Tools

These tools will improve your development experience:

#### 1. golangci-lint

Comprehensive linter for Go:

```bash
# Install using go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Verify installation
golangci-lint --version
```

#### 2. gopls

Go language server for IDE integration:

```bash
go install golang.org/x/tools/gopls@latest
```

#### 3. Air (Live Reload)

Automatic reloading during development:

```bash
go install github.com/cosmtrek/air@latest
```

#### 4. PostgreSQL Client

For database management:

```bash
# On macOS
brew install postgresql

# On Ubuntu/Debian
sudo apt-get install postgresql-client
```

#### 5. kubectl (Optional)

For Kubernetes development:

```bash
# On macOS
brew install kubectl

# On Ubuntu/Debian
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
```

## Project Setup

### 1. Clone the Repository

```bash
# Clone your fork
git clone https://github.com/YOUR-USERNAME/dictamesh.git
cd dictamesh

# Add upstream remote
git remote add upstream https://github.com/Click2-Run/dictamesh.git

# Verify remotes
git remote -v
```

### 2. Install Dependencies

```bash
# Download Go modules
go mod download

# Verify dependencies
go mod verify

# Install development tools
make install-tools
```

### 3. Set Up Infrastructure

Start the required infrastructure services:

```bash
# Start PostgreSQL, Kafka, Redis, and other services
docker-compose up -d

# Verify services are running
docker-compose ps

# Check logs if needed
docker-compose logs -f
```

**Default service ports:**
- PostgreSQL: `localhost:5432`
- Kafka: `localhost:9092`
- Redis: `localhost:6379`
- Jaeger UI: `http://localhost:16686`
- Prometheus: `http://localhost:9090`

### 4. Configure Environment

Create a `.env` file for local development:

```bash
# Copy example environment file
cp .env.example .env

# Edit with your preferred editor
nano .env
```

**Example `.env` configuration:**

```bash
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=dictamesh
DATABASE_USER=dictamesh
DATABASE_PASSWORD=dictamesh_dev
DATABASE_SSL_MODE=disable

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=dictamesh-dev

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Observability
JAEGER_ENDPOINT=http://localhost:14268/api/traces
METRICS_PORT=9090

# Development
LOG_LEVEL=debug
ENV=development
```

### 5. Initialize Database

Run database migrations:

```bash
# Create database and run migrations
make db-setup

# Or manually
psql -h localhost -U dictamesh -c "CREATE DATABASE dictamesh;"
make migrate-up

# Verify migrations
make migrate-status
```

## Development Workflow

### Running the Application

#### Start All Services

```bash
# Run with hot reload (recommended for development)
make dev

# Or run specific services
make run-gateway
make run-adapter-example
```

#### Run Without Hot Reload

```bash
# Build and run
make build
./bin/dictamesh-gateway

# Or use go run
go run cmd/gateway/main.go
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./pkg/connector/...

# Run with verbose output
go test -v ./...

# Run integration tests
make test-integration

# Run specific test
go test -v -run TestMongoDBConnector_Connect ./pkg/connector
```

**Test Coverage Reports:**

```bash
# Generate HTML coverage report
make test-coverage-html

# Open in browser (macOS)
open coverage.html

# Open in browser (Linux)
xdg-open coverage.html
```

### Linting and Formatting

```bash
# Run all linters
make lint

# Fix auto-fixable issues
make lint-fix

# Format code
make fmt

# Run goimports
make imports
```

### Building

```bash
# Build all binaries
make build

# Build specific binary
go build -o bin/gateway cmd/gateway/main.go

# Build with race detector (for testing)
go build -race -o bin/gateway cmd/gateway/main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 make build
```

### Database Operations

```bash
# Create new migration
make migrate-create NAME=add_users_table

# Apply migrations
make migrate-up

# Rollback one migration
make migrate-down

# Rollback all migrations
make migrate-reset

# Check migration status
make migrate-status

# Force migration version (use with caution)
make migrate-force VERSION=1
```

### Working with Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f [service-name]

# Stop all services
docker-compose down

# Rebuild services
docker-compose up -d --build

# Remove volumes (deletes data)
docker-compose down -v
```

## IDE Setup

### Visual Studio Code

**Recommended Extensions:**

```json
{
  "recommendations": [
    "golang.go",
    "ms-azuretools.vscode-docker",
    "eamodio.gitlens",
    "esbenp.prettier-vscode",
    "dbaeumer.vscode-eslint"
  ]
}
```

**Settings (`.vscode/settings.json`):**

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.formatTool": "goimports",
  "go.testFlags": ["-v", "-race"],
  "editor.formatOnSave": true,
  "[go]": {
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  },
  "go.testEnvVars": {
    "DATABASE_HOST": "localhost",
    "KAFKA_BROKERS": "localhost:9092"
  }
}
```

**Debug Configuration (`.vscode/launch.json`):**

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Gateway",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/gateway/main.go",
      "envFile": "${workspaceFolder}/.env",
      "args": []
    },
    {
      "name": "Launch Adapter",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/adapter/main.go",
      "envFile": "${workspaceFolder}/.env"
    },
    {
      "name": "Debug Test",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}",
      "args": [
        "-test.v",
        "-test.run",
        "TestName"
      ]
    }
  ]
}
```

### GoLand / IntelliJ IDEA

**Configuration:**

1. Open the project in GoLand
2. Go to **Preferences → Go → GOROOT** and select Go 1.21+
3. Go to **Preferences → Go → GOPATH** and ensure it's set correctly
4. Enable **Go Modules** integration
5. Install **golangci-lint** plugin
6. Configure **File Watchers** for auto-formatting

**Run Configurations:**

Create run configurations for:
- Gateway server
- Example adapter
- Tests
- Docker Compose

### Vim/Neovim

**Recommended Plugins:**

- [vim-go](https://github.com/fatih/vim-go)
- [coc.nvim](https://github.com/neoclide/coc.nvim) with coc-go

**vim-go Configuration:**

```vim
" .vimrc or init.vim
let g:go_fmt_command = "goimports"
let g:go_auto_type_info = 1
let g:go_metalinter_autosave = 1
let g:go_metalinter_command = "golangci-lint"
```

## Debugging

### Using Delve

Install and use Delve for debugging:

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Start debug session
dlv debug cmd/gateway/main.go

# Common Delve commands
(dlv) break main.main      # Set breakpoint
(dlv) continue             # Continue execution
(dlv) next                 # Step over
(dlv) step                 # Step into
(dlv) print variableName   # Print variable
(dlv) locals               # Show local variables
```

### Debugging Tests

```bash
# Debug specific test
dlv test ./pkg/connector -- -test.run TestMongoDBConnector_Connect

# Debug with IDE integration (VS Code, GoLand)
# Use the debug test configuration
```

### Debugging with Docker

```bash
# Attach to running container
docker exec -it dictamesh-gateway /bin/sh

# View logs
docker logs -f dictamesh-gateway

# Run with debug port exposed
docker run -p 2345:2345 dictamesh-gateway
dlv attach $(pgrep dictamesh-gateway)
```

### Common Debugging Tips

1. **Enable Verbose Logging**: Set `LOG_LEVEL=debug` in `.env`
2. **Use Printf Debugging**: Add `fmt.Printf()` statements temporarily
3. **Check Service Health**: Verify all infrastructure services are running
4. **Review Logs**: Check both application and service logs
5. **Use Tracing**: View requests in Jaeger UI at `http://localhost:16686`

## Troubleshooting

### Common Issues

#### Go Modules Issues

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
rm go.sum
go mod tidy
go mod download
```

#### Docker Permission Issues (Linux)

```bash
# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Or use sudo with docker commands
sudo docker-compose up -d
```

#### Port Already in Use

```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in configuration
```

#### Database Connection Issues

```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Test connection
psql -h localhost -U dictamesh -d dictamesh

# Check logs
docker-compose logs postgres

# Reset database
docker-compose down -v
docker-compose up -d postgres
make db-setup
```

#### Kafka Connection Issues

```bash
# Check Kafka is running
docker-compose ps kafka

# View Kafka logs
docker-compose logs kafka

# Test Kafka connection
docker-compose exec kafka kafka-topics.sh --list --bootstrap-server localhost:9092
```

### Getting Help

If you encounter issues:

1. Check the [Troubleshooting Guide](../operations/troubleshooting.md)
2. Search [GitHub Issues](https://github.com/Click2-Run/dictamesh/issues)
3. Ask in [GitHub Discussions](https://github.com/Click2-Run/dictamesh/discussions)
4. Review service logs: `docker-compose logs -f`

## Performance Tips

### Optimize Build Times

```bash
# Use build cache
export GOCACHE=$(go env GOCACHE)

# Parallel compilation
go build -p 8

# Skip tests during development
go build -tags dev
```

### Optimize Test Runs

```bash
# Run tests in parallel
go test -parallel 8 ./...

# Cache test results
go test -count=1 ./...

# Run only short tests
go test -short ./...
```

### Development Containers

For consistent environments across teams:

```bash
# Build development container
docker build -f Dockerfile.dev -t dictamesh-dev .

# Run development container
docker run -it -v $(pwd):/app dictamesh-dev bash
```

## Next Steps

Now that your environment is set up:

1. Read the [Contributing Guide](./contributing.md)
2. Review the [Architecture Documentation](../architecture/overview.md)
3. Try building an [Example Adapter](../guides/building-adapters.md)
4. Explore the [API Reference](../api-reference/rest-api.md)

## Makefile Commands Reference

Here's a quick reference of available Make commands:

```bash
# Development
make dev                 # Run with hot reload
make build              # Build all binaries
make clean              # Clean build artifacts

# Testing
make test               # Run all tests
make test-coverage      # Run tests with coverage
make test-integration   # Run integration tests

# Code Quality
make lint               # Run linters
make lint-fix           # Fix auto-fixable issues
make fmt                # Format code

# Database
make db-setup           # Create database and run migrations
make migrate-up         # Apply migrations
make migrate-down       # Rollback one migration
make migrate-create     # Create new migration

# Docker
make docker-build       # Build Docker image
make docker-up          # Start services
make docker-down        # Stop services

# Documentation
make docs               # Build documentation
make docs-serve         # Serve documentation locally

# Help
make help               # Show all available commands
```

---

**Previous**: [← Code of Conduct](./code-of-conduct.md) | **Next**: [Contributing Guide →](./contributing.md)
