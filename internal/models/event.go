package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTitle     = errors.New("title must be non-empty and at most 100 characters")
	ErrInvalidTimeRange = errors.New("start time must be before end time")
	ErrEventNotFound    = errors.New("event not found")
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

func (r *CreateEventRequest) Validate() error {
	if r.Title == "" || len(r.Title) > 100 {
		return ErrInvalidTitle
	}
	if !r.StartTime.Before(r.EndTime) {
		return ErrInvalidTimeRange
	}
	return nil
}
