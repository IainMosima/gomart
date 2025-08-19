# GoMart

A domain-driven Go backend service for African e-commerce, featuring hierarchical product categories, order processing with automated notifications, and AWS Cognito authentication.

## Architecture

GoMart follows **Domain-Driven Design (DDD)** principles with clean architecture:

```
├── domains/                 # Business domains (DDD)
│   ├── auth/               # Customer authentication & management
│   ├── category/           # Hierarchical product categories
│   ├── product/            # Product catalog management
│   └── order/              # Order processing & notifications
├── infrastructures/        # External concerns
│   ├── db/                 # Database layer (PostgreSQL + SQLC)
│   └── repository/         # Repository implementations
├── services/               # Application services (business logic)
├── rest-server/            # HTTP layer (Gin framework)
│   ├── handlers/           # HTTP request handlers
│   ├── routes/             # Route definitions
│   └── dtos/               # Data transfer objects
└── configs/                # Configuration management
```

## Features

### Core Domains
- **Authentication**: AWS Cognito integration with customer management
- **Categories**: Unlimited-depth hierarchical product categorization
- **Products**: Complete product catalog with category relationships
- **Orders**: Order processing with automated email + SMS notifications

### Technical Features
- **Database**: PostgreSQL with SQLC for type-safe queries
- **Notifications**: Email (SMTP) + SMS (Africa's Talking) for African markets
- **Testing**: 95%+ test coverage with comprehensive mocking
- **Architecture**: Clean DDD with interface segregation
- **API**: RESTful endpoints with JSON responses

## Quick Start

### Prerequisites
- Go 1.24+
- PostgreSQL 14+
- AWS Cognito (for authentication)
- Africa's Talking account (for SMS)

### 1. Environment Setup

```bash
# Clone repository
git clone https://github.com/IainMosima/gomart.git
cd gomart

# Copy environment template
cp configs/app.env.example configs/app.env
```

### 2. Configure Environment

Copy the example and customize with your values:

```bash
cp configs/app.env.example configs/app.env
```

Edit `configs/app.env`:

```env
# Database
DB_SOURCE=postgresql://root:supersecret@127.0.0.1:5432/gomart_db?sslmode=disable
HTTP_SERVER_ADDRESS=:8080

# AWS Cognito
AWS_REGION=ap-northeast-1
COGNITO_CLIENT_ID=your-client-id
COGNITO_CLIENT_SECRET=your-client-secret
COGNITO_REDIRECT_URI=http://localhost:8080/cognito/callback
COGNITO_DOMAIN=your-cognito-domain.auth.region.amazoncognito.com
COGNITO_USER_POOL_ID=your-user-pool-id

# Africa's Talking SMS
atApiKeys=your-api-key
atUsername=sandbox
atShortCode=15548
atSandbox=true

# Email SMTP
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
EMAIL_FROM=your-email@gmail.com
```

### 3. Database Setup

```bash
# Start PostgreSQL
make postgres

# Create database
make createdb

# Run migrations
make migrateup

# (Optional) Seed test data
make seed
```

### 4. Run Application

```bash
# Start server
go run main.go

# Or use make command
make test && go run main.go
```

API will be available at `http://localhost:8080`

## API Overview

### Authentication
- `GET /auth/login` - Get Cognito authorization URL
- `GET /cognito/callback` - Handle Cognito callback
- `POST /auth/validate` - Validate access token
- `POST /auth/refresh` - Refresh access token

### Categories (Hierarchical)
- `GET /categories` - List categories (supports `parent_id` filter)
- `POST /categories` - Create category
- `GET /categories/:id` - Get category details
- `PUT /categories/:id` - Update category
- `DELETE /categories/:id` - Delete category

### Products
- `GET /products` - List all products
- `POST /products` - Create product
- `GET /products/:id` - Get product details
- `PUT /products/:id` - Update product
- `DELETE /products/:id` - Delete product
- `GET /products/category/:categoryId` - Get products by category

### Orders
- `POST /orders` - Create order (triggers email + SMS notifications)
- `GET /orders/:id/status` - Get order status


## Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

### Test Coverage
- **Service Layer**: 100% (business logic + error scenarios)
- **Handler Layer**: 95% (HTTP endpoints + validation)
- **Route Layer**: 100% (routing configuration)
- **Overall**: 95%+ comprehensive coverage

## Domain Architecture

### Auth Domain
- **Entity**: Customer (linked to AWS Cognito)
- **Repository**: User management operations
- **Service**: Cognito integration + token validation

### Category Domain
- **Entity**: Category (hierarchical with parent/child relationships)
- **Repository**: CRUD + hierarchy queries
- **Service**: Category management + validation

### Product Domain
- **Entity**: Product (linked to categories)
- **Repository**: Product CRUD + category filtering
- **Service**: Product management + category validation

### Order Domain
- **Entity**: Order + OrderItem
- **Repository**: Order CRUD + order item management
- **Service**: Order processing + notification orchestration
- **Notifications**: Email + SMS automated notifications

## Development

### Database Migrations
```bash
# Create new migration
migrate create -ext sql -dir infrastructures/db/migration -seq migration_name

# Apply migrations
make migrateup

# Rollback migrations
make migratedown
```

### Generate Mocks
```bash
# Generate all domain mocks
make mockgen-all

# Generate specific domain mocks
make mockgen-order
make mockgen-auth
```

### Code Generation
```bash
# Generate SQLC code
make sqlc
```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`make test`)
5. Commit changes (`git commit -m 'Add amazing feature'`)
6. Push to branch (`git push origin feature/amazing-feature`)
7. Open Pull Request