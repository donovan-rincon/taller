package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/donovan-rincon/taller/internal/db"
	"github.com/donovan-rincon/taller/internal/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Database configuration
	dbConfig := db.Config{
		// TODO: Move to environment variable or config file, left here for simplicity and testing purposes
		URL:            "postgresql://postgres.ioxmhlkgzcrrsryvrwcz:Ugi3LmYFqRCDApzK@aws-0-us-west-2.pooler.supabase.com:6543/postgres",
		ConnectTimeout: 10 * time.Second,
	}

	// Connect to database
	conn, err := db.Connect(ctx, dbConfig)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer func() {
		if err := db.Close(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "error closing database: %v\n", err)
		}
	}()

	fmt.Println("Connected to database successfully")

	// Server configuration
	serverConfig := server.DefaultConfig()
	srv := server.New(serverConfig, conn)

	// Graceful shutdown handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start()
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	case sig := <-shutdownChan:
		fmt.Printf("Received signal %v, initiating graceful shutdown...\n", sig)
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
	}

	return nil
}
