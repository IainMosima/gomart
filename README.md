# E-Commerce Backend Service

A comprehensive e-commerce backend service built with Go, featuring hierarchical product categories, customer management, order processing, OpenID Connect authentication, and automated notifications.

## 🏗️ Architecture Overview

This project follows Domain-Driven Design (DDD) principles with a clean architecture approach:

```
├── domains/                 # Business domains
│   └── ecommerce/
│       ├── entity/         # Domain entities
│       ├── repository/     # Repository interfaces
│       ├── service/        # Business logic interfaces
│       └── schema/         # Request/Response DTOs
├── infrastructures/        # External concerns
│   ├── db/                # Database layer
│   └── repository/        # Repository implementations
├── services/              # Application services
│   ├── api-gateway/       # REST API layer
│   └── ecommerce/         # Business logic implementation
└── shared/               # Shared utilities
    ├── auth/             # Authentication
    └── utils/            # Common utilities
```

## 🚀 Features

### Core Features
- **Hierarchical Categories**: Unlimited depth category trees for products
- **Product Management**: CRUD operations with bulk upload capabilities
- **Customer Management**: Secure customer profiles with OpenID Connect
- **Order Processing**: Complete order lifecycle with inventory management
- **Real-time Notifications**: SMS via Africa's Talking + Email notifications
- **Authentication & Authorization**: OpenID Connect with role-based access

### Technical Features
- **Database**: PostgreSQL with migrations
- **Caching**: Redis for performance optimization
- **API Documentation**: Auto-generated Swagger docs
- **Testing**: Unit, integration, and benchmark tests with >80% coverage
- **CI/CD**: GitHub Actions with automated testing and deployment
- **Containerization**: Docker and Docker Compose support
- **Kubernetes**: Production-ready K8s manifests
- **Security**: Vulnerability scanning and secure coding practices

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)
- Kubernetes cluster (for production deployment)

## 🛠️ Quick Start

### 1. Environment Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/ecommerce-backend.git
cd ecommerce-backend

# Copy environment variables
cp configs/app.env.example configs/app.env
```

### 2. Configure Environment Variables

Edit `configs/app.env`:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ecommerce
DB_USER=postgres
DB_PASSWORD=your-password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# OpenID Connect
OIDC_ISSUER_URL=https://your-oidc-provider.com
OIDC_CLIENT_ID=your-client-id
OIDC_CLIENT_SECRET=your-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/api/v1/auth/callback

# Africa's Talking SMS
AFRICAS_TALKING_API_KEY=your-api-key
AFRICAS_TALKING_USERNAME=sandbox

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
ADMIN_EMAIL=admin@yourdomain.com
```

### 3. Database Setup

```bash
# Start PostgreSQL and Redis
docker-compose up postgres redis -d

# Run migrations
make migrate-up
```

### 4. Run the Application

```bash
# Development mode
make dev

# Or with Docker
make docker-run
```

The API will be available at `http://localhost:8080`

## 📚 API Documentation

### Authentication Endpoints
- `GET /api/v1/auth/login` - Get OpenID Connect authorization URL
- `POST /api/v1/auth/callback` - Handle OIDC callback
- `POST /api/v1/auth/logout` - Logout (client-side token removal)

### Category Endpoints
- `GET /api/v1/categories` - List categories (supports parent_id filter)
- `POST /api/v1/categories` - Create category (admin only)
- `GET /api/v1/categories/{id}` - Get category by ID
- `GET /api/v1/categories/{id}/average-price` - Get average price for category
- `PUT /api/v1/categories/{id}` - Update category (admin only)
- `DELETE /api/v1/categories/{id}` - Delete category (admin only)

### Product Endpoints
- `GET /api/v1/products` - List products (supports category_id filter)
- `POST /api/v1/products` - Create product (admin only)
- `POST /api/v1/products/bulk` - Bulk upload products (admin only)
- `GET /api/v1/products/{id}` - Get product by ID
- `PUT /api/v1/products/{id}` - Update product (admin only)
- `DELETE /api/v1/products/{id}` - Delete product (admin only)

### Customer Endpoints
- `POST /api/v1/customers` - Create customer (self-registration)
- `GET /api/v1/customers/me` - Get current customer info
- `PUT /api/v1/customers/me` - Update current customer
- `GET /api/v1/customers` - List all customers (admin only)
- `GET /api/v1/customers/{id}` - Get customer by ID (admin only)

### Order Endpoints
- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders/my-orders` - Get current customer's orders
- `GET /api/v1/orders/{id}` - Get order details
- `PUT /api/v1/orders/{id}/cancel` - Cancel order
- `GET /api/v1/orders` - List all orders (admin only)
- `PUT /api/v1/orders/{id}/status` - Update order status (admin only)

## 🏗️ Database Schema

### Categories
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES categories(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Products
```sql
CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    sku VARCHAR(100) UNIQUE NOT NULL,
    stock_quantity INTEGER NOT NULL,
    category_id UUID REFERENCES categories(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Customers
```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    openid_sub VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Orders
```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    customer_id UUID REFERENCES customers(id),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    status order_status DEFAULT 'pending',
    total_amount DECIMAL(10,2) NOT NULL,
    shipping_address TEXT NOT NULL,
    billing_address TEXT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

## 🧪 Testing

### Run Tests
```bash
# Unit tests
make test

# Integration tests
make test-integration

# Linting
make lint
```

### Test Coverage
The project maintains >80% test coverage. Coverage reports are generated automatically:
- `coverage.out` - Coverage data
- `coverage.html` - HTML coverage report

### Testing Strategy
- **Unit Tests**: Test business logic in isolation using mocks
- **Integration Tests**: Test database interactions with real database
- **Benchmark Tests**: Performance testing for critical paths
- **End-to-End Tests**: Complete API workflow testing

## 🚀 Deployment

### Docker Deployment
```bash
# Build and run with Docker Compose
make docker-run

# Stop services
make docker-down
```

### Kubernetes Deployment
```bash
# Apply secrets and config
make k8s-apply-secrets

# Deploy to Kubernetes
make k8s-deploy

# Port forward for local access
make k8s-port-forward
```

### Production Considerations
- Use managed databases (RDS, Cloud SQL)
- Implement proper logging and monitoring
- Set up SSL/TLS certificates
- Configure auto-scaling policies
- Implement backup strategies
- Set up alerting for critical metrics

## 🔒 Security

### Authentication Flow
1. User initiates login via `/api/v1/auth/login`
2. System returns OpenID Connect authorization URL
3. User authenticates with OIDC provider
4. Provider redirects to `/api/v1/auth/callback` with code
5. System exchanges code for tokens and validates ID token
6. System creates/updates customer record and issues JWT
7. Client uses JWT for subsequent API requests

### Security Features
- OpenID Connect integration
- JWT-based authentication
- Role-based authorization (customer/admin)
- SQL injection prevention with parameterized queries
- Input validation and sanitization
- Rate limiting (recommended for production)
- HTTPS enforcement in production
- Vulnerability scanning in CI/CD

## 📊 Monitoring & Observability

### Health Checks
- `GET /health` - Application health status
- Database connectivity check
- Redis connectivity check

### Metrics (Recommended)
- Request/response metrics
- Database query performance
- Cache hit/miss rates
- Order processing metrics
- Authentication success/failure rates

### Logging
- Structured logging with contextual information
- Request/response logging
- Error logging with stack traces
- Audit logging for sensitive operations