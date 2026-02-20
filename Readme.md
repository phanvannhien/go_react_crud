# Full REST App

Go + React full-stack application with User, Auth & Product modules.

## Tech Stack

**Backend:** Go 1.24, Echo, PostgreSQL (pgx), sqlc, dbmate  
**Frontend:** React 19, Vite, TypeScript, TailwindCSS, Zod

## How to Run

### Backend

```bash
# 1. Set up environment
cp .env.example .env   # Edit DATABASE_URL, JWT_SECRET

# 2. Install dbmate (if not installed)
brew install dbmate

# 3. Run migrations
dbmate up

# 4. Start server (port 8080)
go run main.go
```

### Frontend

```bash
cd ui
npm install
npm run dev            # Starts on http://localhost:5173
```

### Seed Admin User

```bash
# Default admin (admin@example.com / admin123)
go run cmd/seed/main.go

# Custom credentials
go run cmd/seed/main.go myemail@example.com mypassword
```

## API Docs (OpenAPI)

OpenAPI specs are located in the `docs/openapi/` directory:

| Module   | File                                           |
|----------|------------------------------------------------|
| Auth     | [`docs/openapi/auth.yaml`](docs/openapi/auth.yaml)       |
| User     | [`docs/openapi/user.yaml`](docs/openapi/user.yaml)       |
| Product  | [`docs/openapi/product.yaml`](docs/openapi/product.yaml) |

You can view them with any OpenAPI viewer:

```bash
# Option 1: Swagger UI via Docker
docker run -p 8081:8080 -e SWAGGER_JSON=/spec/auth.yaml -v $(pwd)/docs/openapi:/spec swaggerapi/swagger-ui

# Option 2: Open in https://editor.swagger.io
#   → File → Import file → select any .yaml from docs/openapi/

# Option 3: VS Code extension "OpenAPI (Swagger) Editor"
```

## API Endpoints

### Auth (Public)
| Method | Endpoint              | Description     |
|--------|-----------------------|-----------------|
| POST   | `/api/auth/register`  | Register user   |
| POST   | `/api/auth/login`     | Login (get JWT) |

### Users (Auth Required)
| Method | Endpoint          | Description                          |
|--------|-------------------|--------------------------------------|
| GET    | `/api/users`      | List users (offset pagination)       |
| GET    | `/api/users/:id`  | Get user by ID                       |
| PUT    | `/api/users/:id`  | Update user                          |
| DELETE | `/api/users/:id`  | Delete user                          |

### Products (Auth Required)
| Method | Endpoint             | Description                         |
|--------|----------------------|-------------------------------------|
| GET    | `/api/products`      | List products (cursor pagination)   |
| GET    | `/api/products/:id`  | Get product by ID                   |
| POST   | `/api/products`      | Create product                      |
| PUT    | `/api/products/:id`  | Update product                      |
| DELETE | `/api/products/:id`  | Delete product                      |

## Project Structure

```
├── main.go                    # Entry point
├── db/
│   ├── migrations/            # dbmate SQL migrations
│   ├── query/                 # sqlc query files
│   └── sqlc/                  # sqlc generated Go code
├── internal/
│   ├── handler/               # HTTP handlers (DTOs + validation)
│   ├── service/               # Business logic
│   ├── routes/                # Route registration
│   ├── middleware/            # JWT auth middleware
│   ├── response/              # JSON response helpers
│   └── test/                  # Integration tests
├── docs/openapi/              # OpenAPI specs
└── ui/                        # React frontend (Vite + TypeScript)
    └── src/
        ├── lib/               # API client, auth context
        └── features/          # auth, user, product modules
```
# 📊 Kiến trúc hoàn chỉnh
```

User Prompt  
    ↓  
agent_runtime  
    ↓  
risk_scoring  
    ↓  
skill_router  
    ↓  
skills activated  
    ↓  
execution_policy  
    ↓  
safe output
```