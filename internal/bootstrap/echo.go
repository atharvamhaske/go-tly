package bootstrap

import (
	"github.com/atharvamhaske/go-tly/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEcho(l *logger.Logger) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Use(l.EchoLogger())

	return e
}
