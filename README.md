# TinyLynx

A lightweight link shortener built with Go, featuring basic concurrency and analytics tracking.

## Features

- ✅ URL shortening with unique short codes
- ✅ Database migrations with Goose
- ✅ Analytics tracking (clicks, visitors, devices, browsers)
- ✅ Worker pool for concurrent analytics processing
- ✅ Rate limiting
- ✅ Graceful shutdown
- ✅ Docker support

## Architecture

```
TinyLynx/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── handler/             # HTTP handlers
│   ├── model/               # Data models
│   ├── service/             # Business logic
│   ├── storage/             # Database operations
│   │   └── migrations/      # Database migrations
│   └── concurrency/         # Worker pool and rate limiting
├── go.mod
├── go.sum
├── Dockerfile
├── compose.yaml
└── .env.example
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose

### Local Development

1. Copy the environment file:
```bash
cp .env.example .env
```

2. Start the database:
```bash
docker-compose up -d
```

3. Run migrations:
```bash
go run cmd/main.go
```

4. Or build and run:
```bash
go build -o tinylynx cmd/main.go
./tinylynx
```

### API Endpoints

#### Create Short Link
```bash
POST /api/shorten
Content-Type: application/json

{
  "original_link": "https://example.com/very-long-url"
}
```

Response:
```json
{
  "short_code": "abc123",
  "original_link": "https://example.com/very-long-url",
  "created_at": "2026-03-18T10:30:00Z"
}
```

#### Redirect
```bash
GET /abc123
```

#### Get Link Stats
```bash
GET /api/stats?code=abc123
```

Response:
```json
{
  "link_id": 1,
  "total_clicks": 42,
  "unique_visitors": 35,
  "countries": 12,
  "devices": 8,
  "browsers": 6
}
```

#### Get Link Analytics
```bash
GET /api/analytics?code=abc123&limit=100
```

Response:
```json
[
  {
    "id": 1,
    "link_id": 1,
    "created_at": "2026-03-18T10:30:00Z",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
    "referrer": "https://google.com",
    "country": "US",
    "device": "desktop",
    "browser": "chrome",
    "platform": "windows"
  }
]
```

## Docker

### Build and Run with Docker

```bash
docker-compose up --build
```

### Using Docker Compose

The project includes a `compose.yaml` file for easy setup:

```bash
docker-compose up -d
```

## Database Schema

### links table
- `id`: Primary key
- `original_link`: The original URL
- `short_code`: Unique short code
- `created_at`: Timestamp

### link_analytics table
- `id`: Primary key
- `link_id`: Foreign key to links
- `created_at`: Timestamp
- `ip_address`: Client IP
- `user_agent`: Browser/user agent string
- `referrer`: Referring URL
- `country`: Country code
- `device`: Device type (mobile/tablet/desktop)
- `browser`: Browser name
- `platform`: Operating system

## Technology Stack

- **Go**: Programming language
- **PostgreSQL**: Database
- **pgx**: PostgreSQL driver
- **goose**: Database migrations
- **base62**: Short code generation

## Project Structure

### Configuration
- Loads environment variables from `.env`
- Configures database connection

### Storage Layer
- Database connection pooling
- Migration management
- SQL queries for links and analytics

### Service Layer
- Business logic for link operations
- Analytics recording and retrieval
- Short code generation using SHA256

### Handler Layer
- HTTP request handling
- Request validation
- Response formatting

### Concurrency
- Worker pool for analytics processing
- Rate limiting for API endpoints

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

## License

MIT