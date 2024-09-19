package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/zerops-dev/warpcamp-showcase/src/database"
	"github.com/zerops-dev/warpcamp-showcase/src/redis"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "starts http server",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.NewConnection(cmd.Context())
			if err != nil {
				return err
			}
			defer db.Close()

			redisClient, err := redis.NewRedis()
			if err != nil {
				return err
			}

			e := echo.New()
			e.Logger.SetOutput(os.Stdout)
			e.Logger.SetLevel(log.DEBUG)
			e.Use(logMiddleware)
			router := e.Router()
			e.HTTPErrorHandler = func(err error, c echo.Context) {
				slog.Error(err.Error())
				e.DefaultHTTPErrorHandler(err, c)
			}

			handler := newHandler(db, redisClient)

			router.Add(http.MethodGet, "/", handler.index)
			router.Add(http.MethodGet, "/health", handler.health)

			server, err := newServer()
			if err != nil {
				return err
			}

			var shutdownErr error
			go func() {
				<-cmd.Context().Done()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				if err := server.Shutdown(ctx); err != nil {
					shutdownErr = err
				}
			}()

			if err := e.StartServer(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}

			return shutdownErr
		},
	}

	return cmd
}

func logMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := c.Request()
		slog.
			With("method", request.Method).
			With("url", request.URL).
			With("host", request.Host).
			With("contentLen", request.ContentLength).
			Info("request")
		return next(c)
	}
}
