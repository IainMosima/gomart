# GoMart

GoMart is a Go backend service designed for African e-commerce. It supports hierarchical product categories, order processing with automated notifications, and AWS Cognito authentication.

## Architecture

The project follows Domain-Driven Design (DDD) and clean architecture:

```
├── domains/                 # Business domains
│   ├── auth/               # Authentication & customer management
│   ├── category/           # Product categories (hierarchical)
│   ├── product/            # Product catalog
│   └── order/              # Order processing & notifications
├── infrastructures/        # External systems
│   ├── db/                 # Database layer (PostgreSQL + SQLC)
│   └── repository/         # Repository implementations
├── services/               # Business logic services
├── rest-server/            # HTTP layer (Gin framework)
│   ├── handlers/           # Request handlers
│   ├── routes/             # Route definitions
│   └── dtos/               # Data transfer objects
└── configs/                # Configuration files and environment variables
```

## Features

- AWS Cognito for authentication
- Unlimited-depth product categories
- Product catalog linked to categories
- Order processing with email and SMS notifications
- PostgreSQL database with SQLC for type-safe queries
- Email (SMTP) and SMS (Africa's Talking) notifications
- RESTful JSON API

## Getting Started

### Requirements

- Go 1.24+
- Docker & Docker Compose (recommended) or PostgreSQL 14+
- AWS Cognito account
- Africa's Talking account for SMS

### Setup

1. Clone the repository and copy environment variables:

```bash
git clone https://github.com/IainMosima/gomart.git
cd gomart
cp configs/app.env.example configs/app.env
```

2. Edit `configs/app.env` with your credentials and settings.

3. Start the database and apply migrations:

```bash
make postgres
make createdb
make migrateup
# Optional: seed test data
make seed
```

### Running the Application

- Using Docker:

```bash
make docker-up
make docker-logs
make docker-down
```

- Manually:

```bash
go run main.go
# Or run tests and then start
make test && go run main.go
```

Access the API at `http://localhost:8080`.

## API Examples

### Authentication

```bash
GET /auth/login

POST /auth/validate
{
  "access_token": "your-access-token"
}

POST /auth/refresh
{
  "refresh_token": "your-refresh-token"
}
```

### Categories

```bash
POST /categories
{
  "category_name": "Electronics"
}

POST /categories
{
  "category_name": "Television",
  "parent_id": "parent-category-id"
}

GET /categories

GET /categories?parent_id=parent-category-id

PUT /categories/category-id
{
  "category_name": "Fridges"
}
```

### Products

```bash
POST /products
{
  "product_name": "iPhone 16",
  "description": "Latest iPhone model",
  "price": 37000.00,
  "sku": "PHONE-IPHONE16-001",
  "stock_quantity": 50,
  "category_id": "category-id",
  "is_active": true
}

GET /products

GET /products/category/category-id

PUT /products/product-id
{
  "product_name": "iPhone 16 Pro",
  "price": 376000.00,
  "sku": "PHONE-IPHONE16PRO-001"
}
```

### Orders

```bash
POST /orders
Authorization: Bearer your-auth-token
{
  "customer_id": "customer-id",
  "items": [
    {
      "product_id": "product-id",
      "quantity": 1
    }
  ]
}

GET /orders/order-id/status
```

### Health Check

```bash
GET /health
```

## Testing

Run all tests:

```bash
make test
```

Run tests with coverage:

```bash
make test-coverage
```

Generate HTML coverage report:

```bash
make test-coverage-html
```

## Development

### Database Migrations

```bash
migrate create -ext sql -dir infrastructures/db/migration -seq migration_name
make migrateup
make migratedown
```

### Generate Mocks

```bash
make mockgen-all
make mockgen-order
make mockgen-auth
```

### Code Generation

```bash
make sqlc
```

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feature-name`)
3. Write tests for your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add feature'`)
6. Push to your branch (`git push origin feature-name`)
7. Open a Pull Request