package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"PenbunAPI/config"
	"PenbunAPI/utils"
)

func JWTMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.UnauthorizedResponse(c)
		}

		tokenStr := parts[1]

		if config.IsTokenBlacklisted(tokenStr) {
			return utils.UnauthorizedResponse(c)
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return utils.UnauthorizedResponse(c)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.UnauthorizedResponse(c)
		}

		c.Locals("username", claims["username"])
		c.Locals("user_level", claims["user_level"])
		c.Locals("token", tokenStr)

		return c.Next()
	}
}
