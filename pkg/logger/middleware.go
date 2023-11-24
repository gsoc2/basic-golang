package logger

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type key int

const (
	levelHeader     = "log-level"
	echoIDKey       = "id"
	ctxKey      key = 0
)

// Middleware attaches a Logger instance with a request ID onto the context. It
// also logs every request along with metadata about the request.
func Middleware() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			level := c.Request().Header.Get(levelHeader)
			l := NewWithLevel(level)
			t1 := time.Now()
			id, err := uuid.NewV4()
			if err != nil {
				return errors.WithStack(err)
			}
			idStr := id.String()
			c.Set(echoIDKey, idStr)
			log := l.ID(idStr)
			// We set the logger on the underlying context.Context instead of
			// the echo.Context so that if we need to use the underlying
			// context.Context for anything during the request lifecycle, it
			// will also have the request ID on it.
			c.SetRequest(c.Request().WithContext(log.WithContext(c.Request().Context())))
			if err := next(c); err != nil {
				c.Error(err)
			}
			t2 := time.Now()

			data := Data{
				"status_code": c.Response().Status,
				"method":      c.Request().Method,
				"path":        c.Request().URL.Path,
				"route":       c.Path(),
				"duration":    t2.Sub(t1).Seconds() * 1000,
				"referer":     c.Request().Referer(),
				"user_agent":  c.Request().UserAgent(),
			}
			if accountID, ok := c.Get("accountID").(string); ok && accountID != "" {
				data["account_id"] = accountID
			}
			if userID, ok := c.Get("userID").(string); ok && userID != "" {
				data["user_id"] = userID
			}

			log.Root(data).Info("request handled")
			return nil
		}
	}
}

// IDFromEchoContext returns the request ID from the given echo.Context. If
// there is no request ID, then this will just return the empty string.
func IDFromEchoContext(c echo.Context) string {
	id, ok := c.Get(echoIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

// FromEchoContext returns a Logger from the given echo.Context. It fetches the
// logger on the underlying context.Context.
func FromEchoContext(c echo.Context) Logger {
	return FromContext(c.Request().Context())
}

// FromContext returns a Logger from the given context.Context. If there is no
// attached logger, then this will just return a new Logger instance.
func FromContext(ctx context.Context) Logger {
	var log Logger
	log, ok := ctx.Value(ctxKey).(Logger)
	if !ok {
		log = New()
	}
	return log
}
