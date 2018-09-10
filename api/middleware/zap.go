package middleware

import (
	"time"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// ZapLogger ...
func ZapLogger(log *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			fields := []interface{}{
				"id", id,
				"status", res.Status,
				"latency", time.Since(start).String(),
				"method", req.Method,
				"uri", req.RequestURI,
				"host", req.Host,
				"remote_ip", c.RealIP(),
			}

			n := res.Status
			switch {
			case n >= 500:
				log.Errorw("Server error", fields...)
			case n >= 400:
				log.Warnw("Client error", fields...)
			case n >= 300:
				log.Infow("Redirection", fields...)
			default:
				log.Infow("Success", fields...)
			}

			return nil
		}
	}
}
