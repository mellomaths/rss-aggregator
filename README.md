# RSS Aggregator

A Go-based REST API for aggregating and managing RSS feeds with user authentication and database persistence.

## Features

- **User Management**: Create and manage users with unique IDs
- **RESTful API**: Clean HTTP endpoints with JSON responses
- **Database Integration**: PostgreSQL with SQLC for type-safe database operations
- **CORS Support**: Cross-origin resource sharing enabled
- **Health Checks**: Built-in readiness endpoint for monitoring

## Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: Chi router
- **Database**: PostgreSQL
- **ORM**: SQLC for type-safe database queries
- **UUID**: Google UUID for unique identifiers
- **Configuration**: Environment variables with godotenv

## API Endpoints

### Health Check
- `GET /v1/healthz` - Server readiness check

### Users
- `POST /v1/users` - Create a new user
  - Request body: `{"name": "string"}`
  - Response: `201` with user object

## Getting Started

### Prerequisites

- Go 1.25.1 or later
- PostgreSQL database
- Environment variables configured

### Environment Variables

Create a `.env` file with the following variables:

```env
PORT=8080
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://username:password@localhost/dbname?sslmode=disable
```

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up your database schema:
   ```bash
   # Run SQL migrations from sql/schema/
   goose postgres postgres://postgres:postgres@localhost/dbname up
   ```

4. Generate database code:
   ```bash
   sqlc generate
   ```

5. Run the server:
   ```bash
   go run .
   ```

The server will start on the port specified in your `PORT` environment variable (default: 8080).

## Project Structure

```
├── internal/database/     # Database models and queries
├── sql/
│   ├── queries/          # SQLC query definitions
│   └── schema/           # Database migrations
├── handler_*.go          # HTTP request handlers
├── main.go              # Application entry point
└── settings.go          # Environment configuration
```

## Development

This project uses SQLC for database operations and Chi for HTTP routing. Database migrations are stored in the `sql/schema/` directory, and queries are defined in `sql/queries/`.

## License

[Add your license here]
