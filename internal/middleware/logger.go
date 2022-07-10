package middleware

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetEchoLogger(e *echo.Echo, logOutput io.Writer) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logOutput,
	}))
}
