# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Backend (Go)
- `cd backend && go run main.go` - Run backend server
- `cd backend && go test ./...` - Run all tests
- `cd backend && go mod tidy` - Clean up dependencies

### Frontend (React + TypeScript)
- `cd front && npm install` - Install dependencies
- `cd front && npm start` - Start development server
- `cd front && npm run build` - Build for production
- `cd front && npm test` - Run tests
- `cd front && npm run lint` - Run linter
- `cd front && npm run lint:fix` - Fix linting errors
- `cd front && npm run generate-api` - Generate TypeScript API client from OpenAPI spec

### Full Stack Development
- `docker-compose up` - Start all services (database, backend, frontend)
- `docker-compose up postgres` - Start only PostgreSQL database
- `docker-compose down` - Stop all services

### API Code Generation
The frontend uses Orval to generate TypeScript API clients from the OpenAPI specification at `api/schema.yaml`. Always run `npm run generate-api` in the front directory after making changes to the API schema.

## Architecture Overview

This is a subscription-based posting platform with Clean Architecture principles:

### Backend Structure (Go)
- **Clean Architecture layers**:
  - `domain/` - Business entities (User, Post, Subscription)
  - `usecase/` - Business logic and application services
  - `repository/` - Data access layer interfaces and implementations
  - `handler/` - HTTP handlers and API endpoints
  - `infrastructure/` - External service implementations (database, JWT, email)
  - `di/` - Dependency injection container

- **Key dependencies**:
  - Chi router for HTTP routing
  - PostgreSQL with database/sql
  - JWT authentication with golang-jwt/jwt
  - Stripe SDK for payment processing
  - bcrypt for password hashing

### Frontend Structure (React + TypeScript)
- **Architecture**:
  - `src/pages/` - Route components and page layouts
  - `src/components/` - Reusable UI components
  - `src/context/` - React Context providers (AuthContext)
  - `src/generated/` - Auto-generated API client code (do not edit manually)
  - `src/utils/` - Utility functions and API configuration

- **Key technologies**:
  - React Query for server state management
  - Material-UI for admin interface components
  - PandaCSS for styling
  - react-hook-form + zod for form validation
  - Stripe React components for payment UI

### Database
PostgreSQL with migrations in `backend/migrations/`:
- `001_initial_schema.sql` - Core tables (users, posts, replies, subscriptions)
- `002_seed_data.sql` - Default admin user and test data

## Key Features and Business Logic

### Authentication & Authorization
- JWT-based auth with access (15m) and refresh (720h) tokens
- Role-based access control (admin vs regular users)
- Password reset flow with secure tokens

### Content Management
- Posts require subscription for creation
- Admin approval workflow for posts
- Anonymous and named replies
- Image upload with compression support

### Subscription System
- Stripe integration for recurring payments
- Webhook handling for real-time subscription updates
- Batch processing for webhook failure recovery

## Environment Configuration

Required environment variables (see docker-compose.yml for examples):
- `DB_*` - PostgreSQL connection details
- `JWT_SECRET` - Secret for signing JWT tokens
- `STRIPE_*` - Stripe API keys and webhook secrets
- `BASE_URL` - Frontend URL for redirects

## Common Development Patterns

### Backend
- All handlers follow dependency injection pattern via the DI container
- Repository pattern for data access with interface abstractions
- Use slog for structured logging
- Environment configuration via envconfig package
- Database transactions for multi-step operations

### Frontend
- All API calls use the generated client from `src/generated/api.ts`
- Global auth state managed via AuthContext
- React Query for caching and server state
- Form validation with react-hook-form + zod schemas

### Testing Strategy
- Backend: Unit tests for usecases and repositories
- Frontend: Component tests with React Testing Library
- Integration tests via Docker Compose setup

## Known Issues and Workarounds

### Backend Migration Issues
If you see "failed to run migrations: first .: file does not exist", the backend is trying to run database migrations but can't find the migration files. Ensure migrations are properly mounted in the Docker container.

### API Client Generation
After changing the OpenAPI schema (`api/schema.yaml`), you must regenerate the frontend API client by running `npm run generate-api` in the front directory before the frontend will work with the new API changes.

### Development Setup
For the full application to work properly:
1. Start PostgreSQL first: `docker-compose up postgres`
2. Run migrations by starting the backend: `docker-compose up backend`
3. Start frontend: `docker-compose up frontend` or `cd front && npm start`

## File Upload System
- Images are uploaded to `/app/uploads` in the backend container
- Frontend uses react-dropzone with compressorjs for client-side compression
- File validation is handled on both client and server side