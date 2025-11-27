package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	conn     *pgx.Conn
	connOnce sync.Once
	connErr  error
)

type Config struct {
	URL            string
	ConnectTimeout time.Duration
}

func Connect(ctx context.Context, cfg Config) (*pgx.Conn, error) {
	connOnce.Do(func() {
		if cfg.ConnectTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, cfg.ConnectTimeout)
			defer cancel()
		}

		conn, connErr = pgx.Connect(ctx, cfg.URL)
		if connErr != nil {
			connErr = fmt.Errorf("failed to connect to database: %w", connErr)
			return
		}

		if err := conn.Ping(ctx); err != nil {
			connErr = fmt.Errorf("failed to ping database: %w", err)
			conn.Close(ctx)
			conn = nil
			return
		}
	})

	return conn, connErr
}

func Close(ctx context.Context) error {
	if conn != nil {
		return conn.Close(ctx)
	}
	return nil
}

func GetConn() *pgx.Conn {
	return conn
}
