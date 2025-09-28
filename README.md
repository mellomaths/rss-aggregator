# RSS Aggregator

A Go-based REST API for aggregating and managing RSS feeds with user authentication, RSS validation, and database persistence.

## Features

- **User Management**: Create and manage users with unique IDs and API keys
- **RSS Feed Management**: Add, validate, and manage RSS feeds with real-time validation
- **Feed Following System**: Users can follow/unfollow specific RSS feeds
- **RSS Scraping**: Background service that automatically scrapes RSS feeds and stores posts
- **Post Management**: View and manage posts from followed feeds with pagination
- **RSS Validation**: Automatic validation of RSS feed URLs to ensure they point to valid RSS content
- **Pagination**: Built-in pagination support for feeds, posts, and feed follows
- **Concurrent Processing**: Multi-threaded RSS scraping with configurable concurrency
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
  - Response: `201` with user object including API key

- `GET /v1/users` - Get current user info (requires authentication)
  - Headers: `Authorization: ApiKey <api_key>`
  - Response: `200` with user object

### Feeds
- `POST /v1/feeds` - Create a new RSS feed (requires authentication)
  - Headers: `Authorization: ApiKey <api_key>`
  - Request body: `{"name": "string", "url": "string"}`
  - Response: `201` with feed object
  - Validation: URL must be a valid RSS feed (HTTP/HTTPS)
  - Examples:
      - `{"name": "Lane's Blog", "url": "https://www.wagslane.dev/index.xml"}`
      - `{"name": "Boot.dev Blog", "url": "https://blog.boot.dev/index.xml"}`

- `GET /v1/feeds` - Get all available feeds (public endpoint)
  - Response: `200` with feed list

### Feed Following
- `POST /v1/feeds/follows` - Follow a specific feed (requires authentication)
  - Headers: `Authorization: ApiKey <api_key>`
  - Request body: `{"feed_id": "uuid"}`
  - Response: `201` with feed follow object

- `GET /v1/feeds/follows` - Get feeds followed by user (requires authentication)
  - Headers: `Authorization: ApiKey <api_key>`
  - Query parameters: `limit` (int), `offset` (int)
  - Response: `200` with paginated feed follows list

- `DELETE /v1/feeds/follows/{feedFollowID}` - Unfollow a feed (requires authentication)
  - Headers: `Authorization: ApiKey <api_key>`
  - Response: `204` No Content

### Posts
- `GET /v1/posts` - Get posts from followed feeds (requires authentication)
  - Headers: `Authorization: ApiKey <api_key>`
  - Query parameters: `limit` (int), `offset` (int)
  - Response: `200` with paginated posts list

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
├── internal/
│   ├── auth/            # Authentication middleware
│   ├── database/        # Database models and queries
│   ├── infra/           # Infrastructure settings
│   └── models/          # Domain models (User, Feed, Paginated)
├── sql/
│   ├── queries/         # SQLC query definitions
│   └── schema/          # Database migrations
├── handler_*.go         # HTTP request handlers
├── middleware_*.go      # HTTP middleware
├── main.go             # Application entry point
└── settings.go         # Environment configuration
```

## RSS Scraping

The application includes a background RSS scraper that:

- **Automatic Scraping**: Continuously scrapes RSS feeds every 10 minutes
- **Concurrent Processing**: Processes up to 10 feeds simultaneously
- **Post Storage**: Automatically stores new posts from RSS feeds
- **Duplicate Prevention**: Prevents duplicate posts using unique URL constraints
- **Feed Tracking**: Tracks last fetch time for each feed to optimize scraping

## RSS Validation

The API includes robust RSS validation that:

- Validates URLs are accessible (HTTP/HTTPS)
- Checks content type for RSS/XML/Atom feeds
- Parses and validates RSS structure
- Provides detailed error messages for invalid feeds
- Uses a 200ms timeout for quick validation

## Data Models

- **Users**: User accounts with API keys for authentication
- **Feeds**: RSS feed sources with validation and metadata
- **Feed Follows**: User subscriptions to specific feeds
- **Posts**: Individual articles/posts from RSS feeds with metadata
- **Pagination**: Consistent pagination across all list endpoints

## Development

This project uses SQLC for database operations and Chi for HTTP routing. Database migrations are stored in the `sql/schema/` directory, and queries are defined in `sql/queries/`.

## License

[Add your license here]
