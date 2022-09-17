package middleware

import (
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func LogMiddleware(e *echo.Echo) {
	var logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger().Level(zerolog.GlobalLevel())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURIPath: true,
		LogStatus:  true,
		LogLatency: true,
		LogMethod:  true,
		LogHost:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Int("status", v.Status).
				Str("Method", v.Method).
				Str("Host", v.Host).
				Str("Path", v.URIPath).
				Str("Latency", v.Latency.String()).
				Msg("request")
			return nil
		},
	}))
}
