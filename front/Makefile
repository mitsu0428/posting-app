.PHONY: install start build test lint lint-fix generate-api clean

# Install dependencies
install:
	npm install

# Start development server
start:
	npm start

# Build for production
build:
	npm run build

# Run tests
test:
	npm test

# Run linter
lint:
	npm run lint

# Fix linting errors
lint-fix:
	npm run lint:fix

# Generate API client from OpenAPI schema
generate-api:
	npm run generate-api

# Clean build artifacts and node_modules
clean:
	rm -rf build node_modules styled-system src/generated

# Development setup
dev-setup: install generate-api

# Build for Cloud Run
build-cloudrun:
	docker build -f Dockerfile.cloudrun -t posting-app-frontend .

# Run development server with Docker
dev-docker:
	docker-compose up frontend