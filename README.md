# Go Books API

A modern, clean-architecture book management API built with Go.

## Overview

This project demonstrates a complete implementation of a RESTful API for managing books using clean architecture, CQRS (Command Query Responsibility Segregation), and the command bus pattern.

The application follows hexagonal architecture principles to ensure separation of concerns, maintainability, and testability.

## Architecture

The application follows a layered architecture:

### Core Layer

- Contains all business logic
- Independent of external concerns like HTTP or databases
- Organized around commands and queries

### Ports Layer

- Adapts the core to different protocols (HTTP, gRPC)
- Translates external requests to commands and queries

### Infrastructure Layer

- Provides implementations for repositories
- Currently uses in-memory storage (can be replaced with a real database)

## Key Concepts

### Command Bus Pattern

The application uses a command bus to handle commands:

1. Commands represent intentions to change state (e.g., AddBookCommand)
2. Command handlers contain the logic for processing commands
3. A command bus dispatches commands to the appropriate handlers

This pattern provides:

- Separation of concerns
- Single responsibility principle
- Testability
- Decoupling

### Repository Pattern

The application uses repositories to abstract data access:

- `BookRepository` interface defines operations for book storage
- `InMemoryBookRepository` provides an in-memory implementation
- Repositories can be swapped without changing business logic

## Running the Application

```bash
# Run the server
go run main.go

# The server will start at http://localhost:8080
```

## Docker

This application can be run using Docker and Docker Compose.

### Prerequisites

- Docker
- Docker Compose

### Running with Docker Compose

1. Build and start the containers:

```bash
docker-compose up -d
```

2. The application will be available at http://localhost:8080

3. To stop the containers:

```bash
docker-compose down
```

4. To view logs:

```bash
docker-compose logs -f
```

### Building the Docker Image Manually

1. Build the image:

```bash
docker build -t books-service .
```

2. Run the container:

```bash
docker run -p 8080:8080 --name books-service books-service
```

Note: When running the container manually, you'll need to ensure the PostgreSQL database is accessible and the environment variables are set correctly.

## API Endpoints

### Books

- `POST /books` - Create a new book
- `GET /books` - Get all books
- `GET /books/:id` - Get a book by ID
- `GET /books/isbn/:isbn` - Get a book by ISBN
- `PUT /books/:id` - Update a book
- `DELETE /books/:id` - Delete a book

### Health Check

- `GET /health` - Check API health

## Testing

The application includes comprehensive unit tests for all commands:

```bash
# Run all tests
go test ./...

# Run specific tests
go test ./core/commands/...
```

## Project Structure

```
books/
├── core/                      # Core business logic
│   ├── commands/              # Command definitions and handlers
│   ├── models/                # Domain models
│   └── storage/               # Storage interfaces and implementations
│       ├── models/            # Storage-specific models
│       └── repositories/      # Repository implementations
├── infrastructure/            # Infrastructure components
│   ├── database.go            # Database connection and configuration
│   └── migrations.go          # Database migrations
├── ports/                     # Adapters for external interfaces
│   └── http-controlers/       # HTTP controllers and server
│       └── controllers/       # Individual HTTP controllers
├── Dockerfile                 # Docker configuration for the application
├── docker-compose.yml         # Docker Compose configuration
├── .dockerignore              # Files to exclude from Docker builds
└── main.go                    # Application entry point
```

## Design Decisions

### Why CQRS?

CQRS (Command Query Responsibility Segregation) separates read and write operations, which:

- Simplifies the model for complex domains
- Improves scalability and performance
- Makes the system more maintainable

### Why Command Bus?

The command bus pattern:

- Centralizes command handling
- Provides a clean way to handle cross-cutting concerns
- Makes it easy to add new commands
- Supports separation of concerns

## Future Improvements

- Implement a proper database repository
- Add user authentication and authorization
- Implement a Query Bus for CQRS completion
- Add pagination for collections
- Create API documentation using Swagger/OpenAPI

## License

MIT
