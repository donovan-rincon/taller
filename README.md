# Taller - Events API

A simple REST API for managing events, built with Go and PostgreSQL.

## Project Structure

```
taller/
├── main.go                              # Entry point with graceful shutdown
├── go.mod
└── internal/
    ├── models/
    │   └── event.go                     # Event model and validation
    ├── db/
    │   └── db.go                        # Database connection management
    ├── repository/
    │   └── event_repository.go          # Database operations
    ├── handlers/
    │   ├── response.go                  # JSON response helpers
    │   └── event_handler.go             # HTTP handlers
    └── server/
        └── server.go                    # HTTP server configuration
```

## Requirements

- Go 1.21+
- PostgreSQL database (cloud stored supabase used, connection string set in code)

## Setup

1. Clone the repository
2. Update the database connection string in `main.go` (already set the properone, not safe but for testing purpuses is fine, should be in an env file)
3. Run the application:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Get All Events

```
GET /events
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "description": "Weekly sync",
    "start_time": "2025-11-27T10:00:00Z",
    "end_time": "2025-11-27T11:00:00Z",
    "created_at": "2025-11-26T15:30:00Z"
  }
]
```

### Get Event by ID

```
GET /events/{id}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Team Meeting",
  "description": "Weekly sync",
  "start_time": "2025-11-27T10:00:00Z",
  "end_time": "2025-11-27T11:00:00Z",
  "created_at": "2025-11-26T15:30:00Z"
}
```

**Error Response (404):**
```json
{
  "error": "not_found",
  "message": "Event not found"
}
```

### Create Event

```
POST /events
```

**Request Body:**
```json
{
  "title": "Team Meeting",
  "description": "Weekly sync",
  "start_time": "2025-11-27T10:00:00Z",
  "end_time": "2025-11-27T11:00:00Z"
}
```

**Validation Rules:**
- `title`: Required, max 100 characters
- `start_time`: Must be before `end_time`
- `description`: Optional

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Team Meeting",
  "description": "Weekly sync",
  "start_time": "2025-11-27T10:00:00Z",
  "end_time": "2025-11-27T11:00:00Z",
  "created_at": "2025-11-26T15:30:00Z"
}
```

**Error Response (400):**
```json
{
  "error": "validation_error",
  "message": "title must be non-empty and at most 100 characters"
}
```

## Database Schema

```sql
CREATE TABLE events (
  id UUID PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  description TEXT,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL
);
```

## Dependencies

- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [google/uuid](https://github.com/google/uuid) - UUID generation