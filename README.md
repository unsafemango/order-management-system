# Order Management System

A microservices-based order management system built with Go, featuring separate services for products and orders with PostgreSQL as the database.

## Architecture

This system consists of two main microservices:

- **Products Service** (`products-service`): Manages product catalog, inventory, and stock reservations
- **Orders Service** (`orders-service`): Handles order creation, status updates, and order management

Both services communicate via HTTP APIs and share a PostgreSQL database.

## Services

### Products Service (Port 8081)

Manages product inventory and stock reservations.

**Endpoints:**
- `GET /v1/health` - Health check
- `POST /v1/products` - Create a new product
- `GET /v1/products/{id}` - Get product by ID
- `POST /v1/products/{id}/reserve` - Reserve stock for a product
- `POST /v1/products/{id}/release` - Release reserved stock

**Product Model:**
```json
{
  "id": 1,
  "sku": "PROD-001",
  "name": "Product Name",
  "price": 2999,
  "stock": 100,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Orders Service (Port 8080)

Handles order processing and status management.

**Endpoints:**
- `GET /v1/health` - Health check
- `POST /v1/orders` - Create a new order
- `GET /v1/orders/{id}` - Get order by ID
- `PATCH /v1/orders/{id}/status` - Update order status

**Order Model:**
```json
{
  "id": 1,
  "customer_id": "customer-123",
  "status": "NEW",
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "price": 2999
    }
  ],
  "total": 5998,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Order Statuses:**
- `NEW` - Order created
- `PAID` - Payment processed
- `CANCELLED` - Order cancelled
- `FULFILLED` - Order completed

## Technology Stack

- **Language**: Go 1.23.4
- **Database**: PostgreSQL 16
- **HTTP Router**: Chi v5
- **Database Driver**: pgx/v5
- **Containerization**: Docker & Docker Compose

## Prerequisites

- Docker and Docker Compose
- Go 1.23.4+ (for local development)

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd order-management-system
   ```

2. **Start all services:**
   ```bash
   make up
   ```
   or
   ```bash
   docker compose up --build -d
   ```

3. **View logs:**
   ```bash
   make logs
   ```
   or
   ```bash
   docker compose logs -f --tail=200
   ```

4. **Stop services:**
   ```bash
   make down
   ```
   or
   ```bash
   docker compose down -v
   ```

### Service URLs

Once started, the services will be available at:
- **Products Service**: http://localhost:8081
- **Orders Service**: http://localhost:8080
- **PostgreSQL**: localhost:5432

## API Usage Examples

### Create a Product

```bash
curl -X POST http://localhost:8081/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "PROD-001",
    "name": "Sample Product",
    "price": 2999,
    "stock": 100
  }'
```

### Create an Order

```bash
curl -X POST http://localhost:8080/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "customer-123",
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'
```

### Update Order Status

```bash
curl -X PATCH http://localhost:8080/v1/orders/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "PAID"
  }'
```

## Database

The system uses PostgreSQL with the following configuration:
- **Database**: retail
- **User**: app
- **Password**: apppass
- **Port**: 5432

Database migrations are automatically applied when services start.

## Development

### Project Structure

```
order-management-system/
├── orders-service/          # Orders microservice
│   ├── cmd/orders/         # Main application entry point
│   ├── internal/
│   │   ├── api/           # HTTP handlers and routing
│   │   ├── models/        # Data models
│   │   ├── repo/          # Database repository layer
│   │   └── env/           # Environment configuration
│   └── Dockerfile
├── products-service/        # Products microservice
│   ├── cmd/products/      # Main application entry point
│   ├── internal/
│   │   ├── api/           # HTTP handlers and routing
│   │   ├── models/        # Data models
│   │   ├── repo/          # Database repository layer
│   │   └── env/           # Environment configuration
│   └── Dockerfile
├── docker-compose.yml      # Service orchestration
├── Makefile               # Development commands
└── go.work               # Go workspace configuration
```

### Local Development

1. **Start PostgreSQL:**
   ```bash
   docker compose up pg -d
   ```

2. **Run Products Service:**
   ```bash
   cd products-service
   go run cmd/products/main.go
   ```

3. **Run Orders Service:**
   ```bash
   cd orders-service
   go run cmd/orders/main.go
   ```

## Environment Variables

### Products Service
- `DB_DSN` - PostgreSQL connection string
- `HTTP_ADDR` - HTTP server address (default: :8081)

### Orders Service
- `DB_DSN` - PostgreSQL connection string
- `HTTP_ADDR` - HTTP server address (default: :8080)
- `PRODUCTS_URL` - Products service URL for inter-service communication

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
