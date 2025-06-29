package middleware

import (
	"PenbunAPI/config"
	"time"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

func NewLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		user := "-"
		if token, ok := c.Locals("user").(*jwt.Token); ok {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if name, ok := claims["user_name"].(string); ok {
					user = name
				}
			}
		}

		err := c.Next()

		config.Logger.WithFields(map[string]interface{}{
			"time":    time.Now().Format(time.RFC3339),
			"user":    user,
			"method":  c.Method(),
			"path":    c.OriginalURL(),
			"status":  c.Response().StatusCode(),
			"latency": time.Since(start).String(),
		}).Info("API Request")

		return err
	}
}
