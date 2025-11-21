package http

import (
	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func JSON(c echo.Context, status int, data any) error {
	return c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c echo.Context, status int, msg string) error {
	return c.JSON(status, APIResponse{
		Success: false,
		Error:   msg,
	})
}
