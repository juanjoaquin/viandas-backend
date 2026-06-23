package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	FormatJSON = "json"
	FormatText = "text"
)

func New(format string) *slog.Logger {
	if strings.EqualFold(format, FormatJSON) {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func RequestLogger(format string) echo.MiddlewareFunc {
	if strings.EqualFold(format, FormatJSON) {
		return middleware.RequestLogger()
	}

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogRequestID: true,
		HandleError:  true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			latency := v.Latency.Round(time.Microsecond)
			requestID := v.RequestID
			if requestID == "" {
				requestID = "-"
			}

			if v.Error != nil {
				c.Logger().Error(fmt.Sprintf(
					"%s %s → %d (%s) [%s] from %s — %v",
					v.Method, v.URI, v.Status, latency, requestID, v.RemoteIP, v.Error,
				))
				return nil
			}

			c.Logger().Info(fmt.Sprintf(
				"%s %s → %d (%s) [%s] from %s",
				v.Method, v.URI, v.Status, latency, requestID, v.RemoteIP,
			))
			return nil
		},
	})
}

func SetDefault(l *slog.Logger) {
	slog.SetDefault(l)
}
