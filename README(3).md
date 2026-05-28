# E-Commerce API

A production-ready RESTful E-Commerce API built with **Go**, featuring JWT authentication, role-based access control, rate limiting, and PostgreSQL database integration.

## 🚀 Live Demo

**Base URL:** `https://my-ecommerce-api.onrender.com`

**Health Check:** `GET /health`

---

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Local Setup](#local-setup)
  - [Environment Variables](#environment-variables)
- [Deployment](#deployment)
- [Project Structure](#project-structure)
- [Screenshots / Testing](#screenshots--testing)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🔐 **JWT Authentication** | Secure login/register with JSON Web Tokens |
| 👤 **Role-Based Access** | Customer and Admin roles with protected routes |
| 📦 **Product Management** | CRUD operations for products (Admin only for create/update/delete) |
| 🛒 **Order System** | Place orders with automatic stock deduction |
| ⏱️ **Rate Limiting** | IP-based request throttling to prevent abuse |
| 📝 **Request Validation** | Input validation using `validator/v10` |
| 🌐 **CORS Support** | Configured for cross-origin requests |
| 📊 **Structured Logging** | Production-grade logging with Zap |
| 🐳 **Docker Ready** | Dockerfile and docker-compose included |

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| **Language** | Go 1.22+ |
| **Router** | Chi (go-chi/chi) |
| **Database** | PostgreSQL |
| **ORM/Query Builder** | SQLC (type-safe SQL) |
| **Authentication** | JWT (golang-jwt/jwt) |
| **Password Hashing** | bcrypt |
| **Validation** | go-playground/validator |
| **Logging** | Uber Zap |
| **Rate Limiting** | golang.org/x/time/rate |
| **Deployment** | Render (Free Tier) |

---

## 🔌 API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/register` | Register new user |
| `POST` | `/api/v1/login` | Login and get JWT token |

### Protected Endpoints (Requires Bearer Token)

| Method | Endpoint | Description | Role |
|--------|----------|-------------|------|
| `GET` | `/api/v1/profile` | Get current user profile | Any |
| `GET` | `/api/v1/products` | List all products | Any |
| `GET` | `/api/v1/products/{id}` | Get single product | Any |
| `POST` | `/api/v1/orders` | Create new order | Any |
| `GET` | `/api/v1/orders` | Get my orders | Any |
| `GET` | `/api/v1/orders/{id}` | Get order details | Any |

### Admin Only Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/products` | Create product |
| `PUT` | `/api/v1/products/{id}` | Update product |
| `DELETE` | `/api/v1/products/{id}` | Delete product |
| `PATCH` | `/api/v1/orders/{id}/status` | Update order status |

---

## 🔑 Authentication

The API uses **Bearer Token** authentication. Include the token in the Authorization header:

```bash
curl https://api.example.com/api/v1/profile   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### Getting a Token

1. **Register** a new account:
```bash
curl -X POST https://api.example.com/api/v1/register   -H "Content-Type: application/json"   -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

2. **Login** to get your token:
```bash
curl -X POST https://api.example.com/api/v1/login   -H "Content-Type: application/json"   -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "customer"
    }
  }
}
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 14+ (or Docker)
- (Optional) Docker & Docker Compose

### Local Setup

```bash
# 1. Clone the repository
git clone https://github.com/nikhil-sharma-dotcom/my-ecommerce-api.git
cd my-ecommerce-api

# 2. Install dependencies
go mod tidy

# 3. Start PostgreSQL (using Docker)
docker run -d   --name postgres-ecom   -e POSTGRES_USER=postgres   -e POSTGRES_PASSWORD=password   -e POSTGRES_DB=ecommerce   -p 5432:5432   postgres:16-alpine

# 4. Run migrations
docker exec -i postgres-ecom psql -U postgres -d ecommerce < internal/adapters/postgresql/migrations/00001_create_products.sql
docker exec -i postgres-ecom psql -U postgres -d ecommerce < internal/adapters/postgresql/migrations/00002_create_orders.sql
docker exec -i postgres-ecom psql -U postgres -d ecommerce < internal/adapters/postgresql/migrations/00003_create_users.sql

# 5. Configure environment
cp .env.example .env
# Edit .env with your database credentials

# 6. Generate SQLC code
sqlc generate

# 7. Run the server
go run ./cmd/api
```

Server starts on `http://localhost:8080`

---

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `ecommerce` |
| `JWT_SECRET` | Secret key for JWT signing | *(required)* |

---

## 🌐 Deployment

### Render (Recommended - Free)

1. Push code to GitHub
2. Connect repo to [Render](https://render.com)
3. Create Web Service + PostgreSQL (Free tier)
4. Set environment variables in dashboard
5. Auto-deploys on every push

**Live URL:** `https://my-ecommerce-api.onrender.com`

### Docker

```bash
# Build image
docker build -t my-ecommerce-api .

# Run container
docker run -d   -p 8080:8080   -e PORT=8080   -e DB_HOST=host.docker.internal   -e JWT_SECRET=your-secret   my-ecommerce-api
```

### Docker Compose

```bash
docker-compose up -d
```

---

## 📁 Project Structure

```
my-ecommerce-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── adapters/
│   │   └── postgresql/
│   │       ├── migrations/      # SQL schema migrations
│   │       └── sqlc/            # Auto-generated type-safe SQL code
│   ├── auth/                    # JWT & password hashing
│   ├── config/                  # Environment configuration
│   ├── handlers/                # HTTP request handlers
│   ├── middleware/              # Auth, CORS, rate limiting, logging
│   ├── orders/                  # Order business logic
│   ├── products/                # Product business logic
│   ├── users/                   # User business logic
│   └── validator/               # Request validation
├── .env.example                 # Environment template
├── Dockerfile                   # Container build
├── docker-compose.yml           # Multi-container orchestration
├── go.mod                       # Go module definition
├── Makefile                     # Build automation
└── sqlc.yaml                    # SQLC configuration
```

---

## 🧪 Testing with cURL

### Register a User
```bash
curl -X POST https://my-ecommerce-api.onrender.com/api/v1/register   -H "Content-Type: application/json"   -d '{"email":"test@example.com","password":"password123","first_name":"Test","last_name":"User"}'
```

### Login
```bash
curl -X POST https://my-ecommerce-api.onrender.com/api/v1/login   -H "Content-Type: application/json"   -d '{"email":"test@example.com","password":"password123"}'
```

### Access Protected Route
```bash
curl https://my-ecommerce-api.onrender.com/api/v1/profile   -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### List Products
```bash
curl https://my-ecommerce-api.onrender.com/api/v1/products
```

---

## 🎓 Learning Outcomes

This project demonstrates:

- ✅ Building REST APIs with Go and Chi router
- ✅ Type-safe database operations with SQLC
- ✅ JWT-based authentication & authorization
- ✅ Middleware pattern for cross-cutting concerns
- ✅ Docker containerization
- ✅ Cloud deployment on Render
- ✅ Production-ready logging and error handling

---

## 📝 License

MIT License - feel free to use this project for learning or as a starter template.

---

## 👤 Author

**Nikhil Sharma**

- GitHub: [@nikhil-sharma-dotcom](https://github.com/nikhil-sharma-dotcom)
- Project: [my-ecommerce-api](https://github.com/nikhil-sharma-dotcom/my-ecommerce-api)

---

> **Note:** This project was built as part of a learning assignment and has been enhanced with production-grade features including authentication, rate limiting, and structured logging.
