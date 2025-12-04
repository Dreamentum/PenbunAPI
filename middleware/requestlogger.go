package middleware

import (
	"PenbunAPI/config"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func NewLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		user := "-"
		if claims, ok := c.Locals("user").(jwt.MapClaims); ok {
			if name, ok := claims["user_name"].(string); ok {
				user = name
			}
		}

		config.Logger.WithFields(map[string]interface{}{
			"time":    time.Now().Format(time.RFC3339),
			"user":    user,
			"method":  c.Method(),
			"path":    c.OriginalURL(),
			"status":  c.Response().StatusCode(),
			"latency": time.Since(start).String(),
		}).Info("API Request")

		// Console Log (Simple Format)
		status := "success"
		statusCode := c.Response().StatusCode()
		if statusCode >= 400 && statusCode < 500 {
			status = "fail"
		} else if statusCode >= 500 {
			status = "error"
		}

		// Strip prefix for console log
		logPath := c.OriginalURL()
		if len(logPath) > 18 && logPath[:18] == "/api/v1/protected/" {
			logPath = logPath[18:]
		}

		log.Printf("[User: %s] [Method: %s] [Path: %s] [Status: %s]", user, c.Method(), logPath, status)

		return err
	}
}
