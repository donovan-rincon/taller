package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/donovan-rincon/taller/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EventRepository struct {
	conn *pgx.Conn
}

// NewEventRepository creates a new instance of EventRepository with the given database connection.
func NewEventRepository(conn *pgx.Conn) *EventRepository {
	return &EventRepository{conn: conn}
}

// Create inserts a new event into the database.
func (r *EventRepository) Create(ctx context.Context, event *models.Event) error {
	// SQL query to insert a new event, need to improve to avoid sql injection
	query := `
		INSERT INTO events (id, title, description, start_time, end_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.conn.Exec(ctx, query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

// GetAll retrieves all events from the database.
func (r *EventRepository) GetAll(ctx context.Context) ([]models.Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at
		FROM events
		ORDER BY start_time ASC
	`

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	events := make([]models.Event, 0)

	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// GetByID retrieves an event by its ID from the database.
func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at
		FROM events
		WHERE id = $1
	`

	var event models.Event
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrEventNotFound
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}
