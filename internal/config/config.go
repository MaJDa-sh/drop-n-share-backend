package config

import (
	"log"
	"os"
	"strconv"
)

var WebSocketMaxConnections int

func LoadConfig() {

	WebSocketMaxConnections = getEnv("WEBSOCKET_MAX_CONNECTIONS", 100)
}

func getEnv(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {

		log.Printf("Invalid value for %s, using default: %s\n", key, err)
		return defaultValue
	}
	return intValue
}
