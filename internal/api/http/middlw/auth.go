package middlw

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthReq(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing auth header"})
			}
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token format"})
			}

			token := parts[1]

			userID, err := ValidateJWT(token, secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			c.Set("user_id", userID)
			return next(c)
		}
	}
}
