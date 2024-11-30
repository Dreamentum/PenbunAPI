package config

import (
	"log"
	"sync"
)

var blacklist = make(map[string]bool)
var mu sync.Mutex

func AddToBlacklist(token string) {
	mu.Lock()
	defer mu.Unlock()
	if token != "" {
		blacklist[token] = true
		log.Println("[DEBUG] Token added to blacklist:", token)
	}
}

func IsBlacklisted(token string) bool {
	mu.Lock()
	defer mu.Unlock()
	log.Println("[DEBUG] Checking if token is blacklisted:", token)
	return blacklist[token]
}
