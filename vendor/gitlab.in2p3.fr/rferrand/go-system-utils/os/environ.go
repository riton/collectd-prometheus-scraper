package os

import (
	coreos "os"
)

// GetEnv returns environment variable
// if present. If not return fallback value
func GetEnv(key, fallback string) string {
	if value, ok := coreos.LookupEnv(key); ok {
		return value
	}
	return fallback
}
