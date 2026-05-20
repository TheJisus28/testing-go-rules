// Package logger provides process-wide and request-scoped zerolog helpers.
package logger

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ctxKey struct{}

var global zerolog.Logger

// Init configures the global log level and console writer (call once from main).
func Init(level string) {
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	global = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).With().Timestamp().Logger()
}

func Get() *zerolog.Logger {
	return &global
}

func FromCtx(ctx context.Context) *zerolog.Logger {
	if log, ok := ctx.Value(ctxKey{}).(zerolog.Logger); ok {
		return &log
	}
	return &global
}

func WithLogger(ctx context.Context, log zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

// Middleware attaches a request-scoped logger and emits an access log line per request.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			log := global.With().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Logger()
			ctx := WithLogger(req.Context(), log)
			c.SetRequest(req.WithContext(ctx))
			err := next(c)
			if err != nil {
				log.Error().Err(err).Int("status", c.Response().Status).Msg("http: request failed")
			} else {
				log.Info().Int("status", c.Response().Status).Msg("http: request completed")
			}
			return err
		}
	}
}

func Writer() io.Writer {
	return os.Stdout
}
