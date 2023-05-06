package utils

import "os"

func GetEnv(key, fallback string) string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return port
	}
	return fallback
}
