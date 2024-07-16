package helpers

import (
	"log"
	"os"
)

func SafeGetEnv(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("The environment variable '%s' doesn't exist or is not set", key)
	}
	return os.Getenv(key)
}
