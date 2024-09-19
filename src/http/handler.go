package http

import (
	"database/sql"
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

//go:embed resources/index.html
var indexHtml string

type Handler struct {
	db       *sql.DB
	redis    *redis.Client
	upgrader websocket.Upgrader
}

func newHandler(
	db *sql.DB,
	redis *redis.Client,
) *Handler {
	return &Handler{
		db:       db,
		redis:    redis,
		upgrader: websocket.Upgrader{},
	}
}

func (h *Handler) health(ctx echo.Context) error {
	if err := h.db.Ping(); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusServiceUnavailable,
			Message: err.Error(),
		}
	}
	return ctx.JSON(http.StatusOK, map[string]string{"database": "ok"})
}

const redisChannel = "messages"

func (h *Handler) index(c echo.Context) error {
	if !c.IsWebSocket() {
		return c.HTML(http.StatusOK, indexHtml)
	}

	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, "SELECT text FROM messages")
	if err != nil {
		return err
	}
	for rows.Next() {
		var message []byte
		if err := rows.Scan(&message); err != nil {
			return err
		}
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			return err
		}
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	sub := h.redis.Subscribe(ctx, redisChannel)
	defer sub.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-sub.Channel():
				if message == nil {
					continue
				}
				if err := ws.WriteMessage(websocket.TextMessage, []byte(message.Payload)); err != nil {
					slog.Error(err.Error())
				}
			}
		}
	}()

	for ctx.Err() == nil {
		_, data, err := ws.ReadMessage()
		if err != nil {
			slog.Error(err.Error())
			return nil
		}

		if _, err := h.db.ExecContext(ctx, "INSERT INTO messages (text) VALUES ($1)", string(data)); err != nil {
			return err
		}

		if err := h.redis.Publish(ctx, redisChannel, data).Err(); err != nil {
			return err
		}
	}
	return ctx.Err()
}
