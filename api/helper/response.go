package helper

import (
	"net/http"
	"weigo/api/database"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	// Secret .
	Secret = "12345"
)

// BuildErrorsResponse .
func BuildErrorsResponse(i interface{}) map[string]interface{} {
	return map[string]interface{}{
		"errors": i,
	}
}

// BuildErrorResponse .
func BuildErrorResponse(i interface{}) map[string]interface{} {
	return map[string]interface{}{
		"errors": []interface{}{i},
	}
}

// AuthMiddleware .
func AuthMiddleware() echo.MiddlewareFunc {
	return middleware.JWT([]byte(Secret))
}

// GetUserMiddleware .
func GetUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := claims["id"].(float64)
		p, err := database.ProfileGetByID(uint64(id))
		if err != nil {
			return c.String(http.StatusUnauthorized, "")
		} else if p == nil {
			return c.String(http.StatusNotFound, "")
		}
		c.Set("user", p)
		return next(c)
	}
}
