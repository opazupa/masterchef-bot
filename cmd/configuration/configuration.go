package configuration

import (
	"os"
	"strconv"
)

// Configuration type
type Configuration struct {
	APIKey    string
	DebugMode bool
}

// Get struct
func Get() *Configuration {
	return &Configuration{
		APIKey:    getEnv("API_KEY", ""),
		DebugMode: getEnvAsBool("DEBUG_MODE", true),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
