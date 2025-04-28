# Task Manager API 🚀

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/ahmednurovic/task-manager-api/actions/workflows/go.yml/badge.svg)](https://github.com/yourusername/task-manager-api/actions)

A modern task management REST API built with Golang, featuring JWT authentication, PostgreSQL storage, and Docker support. Perfect for learning Golang backend development or as a starter for your next project!

## Features ✨

- 🔐 JWT Authentication
- 📝 Full CRUD operations for tasks
- 🐘 PostgreSQL database with migrations
- 📚 Swagger API documentation
- 🐳 Dockerized development environment
- 📈 Structured logging with Zap

## Tech Stack 💻

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Logging**: Zap
- **Containerization**: Docker
- **Documentation**: Swagger/OpenAPI

## Getting Started 🚦

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/task-manager-api.git
cd task-manager-api

# Copy environment file
cp .env.example .env

# Start services
docker-compose up -d

# Run database migrations
goose -dir migrations postgres "user=user password=password dbname=taskdb sslmode=disable" up

# Start the server
go run cmd/main.go