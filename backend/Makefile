.PHONY: build run test clean deps migrate-up migrate-down docker-build docker-run

# Build the application
build:
	go build -o main .

# Run the application
run:
	go run .

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -f main posting-app

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Build Docker image
docker-build:
	docker build -t posting-app-backend .

# Build Docker image for Cloud Run
docker-build-cloudrun:
	docker build -f Dockerfile.cloudrun -t posting-app-backend-cloudrun .

# Run with Docker
docker-run:
	docker run -p 8080:8080 posting-app-backend

# Run development server with live reload (requires air)
dev:
	air

# Install development dependencies
dev-deps:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Database operations (requires migrate CLI)
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost/posting_app?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost/posting_app?sslmode=disable" down

# Run with docker-compose
docker-compose-up:
	docker-compose up backend

# Run full stack with docker-compose
stack-up:
	docker-compose up