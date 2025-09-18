# Claude Code Project Configuration

## Project Structure

This is a Go web application template using the Gin framework with the following structure:

```
real-estate-tracker/
├── api/                    # API specifications
│   └── openapi.yaml       # OpenAPI/Swagger documentation
├── cmd/                   # Application entry points
│   └── server/
│       └── main.go       # Main server application
├── internal/             # Private application code
│   ├── handlers/         # HTTP request handlers
│   │   └── static_handlers.go
│   ├── middlewares/      # HTTP middlewares
│   ├── models/          # Data models
│   ├── routes/          # Route definitions
│   │   └── static_routes.go
│   ├── services/        # Business logic services
│   └── utils/           # Utility functions
├── prisma/              # Database schema and migrations
│   └── schema.prisma    # Prisma database schema
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── package.json        # Node.js dependencies (for Prisma)
└── yarn.lock          # Yarn lock file
```

## Development Guidelines

### File Organization
- Keep all Go source code in appropriate packages under `internal/`
- Place reusable handlers in `internal/handlers/`
- Define routes in `internal/routes/`
- Business logic goes in `internal/services/`
- Database models in `internal/models/`
- Utility functions in `internal/utils/`

### Naming Conventions
- Use snake_case for file names (e.g., `static_handlers.go`)
- Use PascalCase for Go types and exported functions
- Use camelCase for unexported functions and variables

### Dependencies
- Go modules for Go dependencies
- Yarn for Node.js dependencies (Prisma tooling)
- Prisma for database management

## Commands

### Build and Run
```bash
go run cmd/server/main.go
```

### Database Operations
```bash
npx prisma generate    # Generate Prisma client
npx prisma migrate dev # Run database migrations
npx prisma studio     # Open Prisma Studio
```

### Testing
```bash
go test ./...         # Run all tests
```
