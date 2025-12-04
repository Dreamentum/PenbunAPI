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
		if config.IsBlacklisted(tokenString) {
			log.Println("[DEBUG] Token is blacklisted")
			return fiber.NewError(fiber.StatusUnauthorized, "Token is blacklisted")
		}

		// ตรวจสอบและ parse Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			log.Println("[DEBUG] Invalid token format or signature")
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Println("[DEBUG] Invalid claims structure")
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		// แสดงชื่อผู้ใช้งานแทน token
		if userName, ok := claims["user_name"]; ok {
			log.Printf("[DEBUG] Token validated for user: %v", userName)
		} else {
			log.Println("[DEBUG] Token valid but missing user_name in claims")
		}

		// เก็บข้อมูลผู้ใช้งานใน context
		c.Locals("user", claims)
		return c.Next()
	}
}
