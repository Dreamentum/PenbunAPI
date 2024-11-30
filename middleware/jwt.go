package middleware

import (
	"PenbunAPI/config" // สำหรับ Blacklist
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware เป็น middleware ที่ใช้ในการตรวจสอบ JWT Token
func JWTMiddleware(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing or malformed token")
		}

		// ตัดคำว่า "Bearer " ออกจาก Token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// ตรวจสอบว่า Token อยู่ใน Blacklist หรือไม่
		log.Println("[DEBUG] Token being validated:", tokenString)
		if config.IsBlacklisted(tokenString) {
			log.Println("[DEBUG] Token is blacklisted:", tokenString)
			return fiber.NewError(fiber.StatusUnauthorized, "Token is blacklisted")
		}
		
		// ตรวจสอบและ parse Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("user", claims)
			return c.Next()
		}

		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}
}

