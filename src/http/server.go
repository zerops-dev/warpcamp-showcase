package http

import (
	"net/http"
	"os"
	"time"

	"github.com/zerops-dev/warpcamp-showcase/src/env"
)

const (
	EnvAddress = "HTTP_ADDRESS"
)

func newServer() (*http.Server, error) {
	if err := env.Check(EnvAddress); err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:              os.Getenv(EnvAddress),
		ReadTimeout:       time.Second * 3,
		ReadHeaderTimeout: time.Second * 3,
		WriteTimeout:      time.Second * 3,
		IdleTimeout:       time.Second * 3,
	}, nil
}
