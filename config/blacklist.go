package config

import (
	"log"
	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

var blacklist = make(map[string]bool)
var mu sync.Mutex

func AddToBlacklist(token string) {
	mu.Lock()
	defer mu.Unlock()
	if token != "" {
		blacklist[token] = true

		user := extractUsername(token)
		if user != "" {
			log.Printf("[DEBUG] Token added to blacklist for user: %s", user)
		} else {
			log.Println("[DEBUG] Token added to blacklist: [unknown user]")
		}
	}
}

func IsBlacklisted(token string) bool {
	mu.Lock()
	defer mu.Unlock()

	user := extractUsername(token)
	if user != "" {
		log.Printf("[DEBUG] Checking if token belongs to user: %s", user)
	} else {
		log.Println("[DEBUG] Checking if token is blacklisted: [unknown user]")
	}

	return blacklist[token]
}

// extractUsername ถอด username จาก JWT Token (ถ้ามี)
func extractUsername(tokenString string) string {
	if strings.TrimSpace(tokenString) == "" {
		return ""
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetEnv("JWT_SECRET")), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userName, ok := claims["user_name"].(string); ok {
			return userName
		}
	}
	return ""
}
