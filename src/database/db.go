package database

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/zerops-dev/warpcamp-showcase/src/env"
)

const (
	EnvHost     = "DB_HOST"
	EnvPort     = "DB_PORT"
	EnvUser     = "DB_USER"
	EnvPass     = "DB_PASS"
	EnvDatabase = "DB_DATABASE"
)

func NewConnection(ctx context.Context) (*sql.DB, error) {
	if err := env.Check(EnvHost, EnvUser, EnvPass, EnvDatabase); err != nil {
		return nil, err
	}

	port, err := strconv.ParseUint(os.Getenv(EnvPort), 10, 16)
	if err != nil {
		return nil, err
	}

	config, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}

	config.Host = os.Getenv(EnvHost)
	config.Port = uint16(port)
	config.Database = os.Getenv(EnvDatabase)
	config.User = os.Getenv(EnvUser)
	config.Password = os.Getenv(EnvPass)
	config.ConnectTimeout = time.Second * 5

	conn := stdlib.OpenDB(*config)
	return conn, conn.PingContext(ctx)
}
