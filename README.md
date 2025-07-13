# Posting App

A subscription-based posting platform with admin moderation, built with Go backend and React frontend.

## Features

### Authentication & User Management
- User registration with email verification
- JWT-based authentication (access + refresh tokens)
- Password reset functionality
- Admin and regular user roles
- Account deactivation and user banning

### Content Management
- Create, edit, and delete posts (with approval workflow)
- Reply to posts (anonymous or with username)
- Image upload support (thumbnails)
- Content moderation by admins
- Rich text content support

### Subscription System
- Stripe integration for subscription management
- Content creation restricted to active subscribers
- Subscription status tracking
- Webhook support for real-time updates
- Batch sync for webhook failure recovery

### Admin Features
- Post approval/rejection workflow
- User management and banning
- Content moderation dashboard
- Admin-only access controls

### Security & Performance
- Rate limiting
- XSS/CSRF protection
- Secure file uploads
- Optimized database queries
- Comprehensive logging

## Tech Stack

### Backend
- **Language**: Go 1.21
- **Framework**: Chi router
- **Database**: PostgreSQL with SQL migrations
- **Authentication**: JWT with bcrypt
- **Payment**: Stripe API
- **Email**: SendGrid
- **Validation**: go-playground/validator
- **Logging**: slog + zerolog (JSON)

### Frontend
- **Framework**: React 18 + TypeScript
- **Styling**: PandaCSS + Material-UI (admin)
- **State Management**: React Query + Context API
- **Forms**: react-hook-form + zod validation
- **File Upload**: react-dropzone + compressorjs
- **Payments**: Stripe React components
- **Testing**: React Testing Library

### Infrastructure
- **API**: OpenAPI 3.0 specification
- **Containers**: Docker + Docker Compose
- **Database**: PostgreSQL 15
- **File Storage**: Local filesystem (configurable)
- **Deployment**: Cloud Run ready

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for local frontend development)
- Go 1.21+ (for local backend development)

### Environment Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd posting-app
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Configure your environment variables in `.env`:
```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=posting_app

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_DURATION=15m
JWT_REFRESH_DURATION=720h

# Stripe
STRIPE_API_KEY=sk_test_your_stripe_secret_key
STRIPE_PRICE_ID=price_your_stripe_price_id
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# SendGrid
SENDGRID_API_KEY=your_sendgrid_api_key

# Application
BASE_URL=http://localhost:3000
PORT=8080
```

### Development with Docker

1. Start all services:
```bash
docker-compose up
```

2. The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: localhost:5432

### Local Development

#### Backend
```bash
cd backend
make deps
make run
```

#### Frontend
```bash
cd front
npm install
npm run generate-api  # Generate API client from OpenAPI spec
npm start
```

## Default Admin Account

The application seeds a default admin account:
- **Email**: admin@example.com
- **Password**: admin123

Access the admin dashboard at: http://localhost:3000/admin

## API Documentation

The API is documented using OpenAPI 3.0. The specification is located at `api/schema.yaml`.

Key endpoints:
- `POST /auth/login` - User authentication
- `POST /auth/register` - User registration
- `GET /posts` - List approved posts
- `POST /posts` - Create new post (requires subscription)
- `POST /posts/{id}/replies` - Add reply to post
- `GET /admin/posts` - Admin post management
- `POST /subscription/create-checkout-session` - Create Stripe checkout

## Database Schema

The application uses PostgreSQL with the following main tables:
- `users` - User accounts and subscription status
- `posts` - User-created content with approval status
- `replies` - Comments on posts (can be anonymous)
- `subscriptions` - Stripe subscription tracking
- `password_resets` - Password reset tokens

## Deployment

### Cloud Run Deployment

1. Build production images:
```bash
# Backend
cd backend
make docker-build-cloudrun

# Frontend
cd front
make build-cloudrun
```

2. Deploy to Cloud Run with appropriate environment variables

### Environment Variables for Production

Ensure the following environment variables are configured:
- All database connection details
- JWT secret (generate a secure random string)
- Stripe API keys and webhook secrets
- SendGrid API key for email functionality
- BASE_URL pointing to your frontend domain

## Testing

### Backend Tests
```bash
cd backend
make test
make test-coverage
```

### Frontend Tests
```bash
cd front
npm test
```

### Linting
```bash
# Backend
cd backend
make lint

# Frontend
cd front
npm run lint
```

## Security Considerations

1. **Environment Variables**: Never commit secrets to version control
2. **JWT Secrets**: Use strong, randomly generated secrets in production
3. **HTTPS**: Always use HTTPS in production
4. **Database**: Use strong passwords and restrict network access
5. **File Uploads**: Validate file types and sizes
6. **Rate Limiting**: Configure appropriate rate limits for your use case

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Run the test suite: `make test` (backend) and `npm test` (frontend)
5. Commit your changes: `git commit -am 'Add feature'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
1. Check the documentation in this README
2. Review the API specification in `api/schema.yaml`
3. Open an issue on GitHub
4. Contact the development team

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React App     │    │   Go Backend    │    │   PostgreSQL    │
│                 │    │                 │    │                 │
│ • Auth Context  │◄──►│ • JWT Auth      │◄──►│ • Users         │
│ • Post Mgmt     │    │ • Post API      │    │ • Posts         │
│ • Admin UI      │    │ • Admin API     │    │ • Replies       │
│ • Stripe UI     │    │ • Stripe API    │    │ • Subscriptions │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │
         │              ┌─────────────────┐
         │              │   Stripe API    │
         └──────────────►│                 │
                        │ • Subscriptions │
                        │ • Webhooks      │
                        └─────────────────┘
```

The application follows Clean Architecture principles with clear separation between:
- **Presentation Layer**: React components and pages
- **Business Logic**: Go usecases and domain models
- **Data Layer**: PostgreSQL repositories and external APIs