package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/donovan-rincon/taller/internal/handlers"
	"github.com/donovan-rincon/taller/internal/repository"
	"github.com/jackc/pgx/v5"
)

type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		Port:            ":8080",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 30 * time.Second,
	}
}

type Server struct {
	httpServer *http.Server
	config     Config
}

// New creates a new Server instance with the given configuration and database connection.
func New(cfg Config, conn *pgx.Conn) *Server {
	eventRepo := repository.NewEventRepository(conn)
	eventHandler := handlers.NewEventHandler(eventRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/events", eventHandler.HandleEvents)
	mux.HandleFunc("/events/", eventHandler.HandleEventByID)

	httpServer := &http.Server{
		Addr:         cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		config:     cfg,
	}
}

// Start begins listening for HTTP requests on the configured port.
func (s *Server) Start() error {
	fmt.Printf("Server listening on %s\n", s.config.Port)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server within the configured timeout.
func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.config.ShutdownTimeout)
	defer cancel()

	fmt.Println("Shutting down server...")
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}
	fmt.Println("Server stopped gracefully")
	return nil
}
