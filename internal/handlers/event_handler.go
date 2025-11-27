package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/donovan-rincon/taller/internal/models"
	"github.com/donovan-rincon/taller/internal/repository"
	"github.com/google/uuid"
)

const requestTimeout = 10 * time.Second

type EventHandler struct {
	repo *repository.EventRepository
}

// NewEventHandler creates a new EventHandler with the given EventRepository.
func NewEventHandler(repo *repository.EventRepository) *EventHandler {
	return &EventHandler{repo: repo}
}

// HandleEvents handles requests for the /events endpoint.
func (h *EventHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

// HandleEventByID handles requests for a specific event identified by its ID.
func (h *EventHandler) HandleEventByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractEventID(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

// create handles the creation of a new event.
func (h *EventHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_payload", "Invalid request payload")
		return
	}

	event, err := newEvent(req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.repo.Create(ctx, event); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "timeout", "Request timed out")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to create event")
		return
	}

	respondJSON(w, http.StatusCreated, event)
}

// newEvent creates a new Event instance from the given CreateEventRequest.
func newEvent(req models.CreateEventRequest) (*models.Event, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &models.Event{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

// getAll handles the retrieval of all events.
func (h *EventHandler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	events, err := h.repo.GetAll(ctx)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "timeout", "Request timed out")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to fetch events")
		return
	}

	respondJSON(w, http.StatusOK, events)
}

// getByID handles the retrieval of a single event by its ID.
func (h *EventHandler) getByID(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	event, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			respondError(w, http.StatusNotFound, "not_found", "Event not found")
			return
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "timeout", "Request timed out")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to fetch event")
		return
	}

	respondJSON(w, http.StatusOK, event)
}

// extractEventID extracts the event ID from the URL path.
func extractEventID(path string) (uuid.UUID, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 || parts[1] == "" {
		return uuid.Nil, errors.New("missing event ID")
	}

	id, err := uuid.Parse(parts[1])
	if err != nil {
		return uuid.Nil, errors.New("invalid event ID format")
	}

	return id, nil
}
