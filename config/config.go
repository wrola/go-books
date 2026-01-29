package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Security SecurityConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port string
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	APIKey          string
	AuthDisabled    bool
	RateLimitRPS    int
	RateLimitBurst  int
	CORSOrigins     string
	CORSMethods     string
	CORSHeaders     string
	CORSCredentials bool
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Database: loadDatabaseConfig(),
		Server:   loadServerConfig(),
		Security: loadSecurityConfig(),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	port := 5432
	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "require" // Secure default
	}

	return DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnvOrDefault("DB_USER", ""),
		Password: getEnvOrDefault("DB_PASSWORD", ""),
		Name:     getEnvOrDefault("DB_NAME", "books"),
		SSLMode:  sslMode,
	}
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port: getEnvOrDefault("SERVER_PORT", ":8080"),
	}
}

func loadSecurityConfig() SecurityConfig {
	rateRPS := 100
	if rps := os.Getenv("RATE_LIMIT_RPS"); rps != "" {
		if r, err := strconv.Atoi(rps); err == nil && r > 0 {
			rateRPS = r
		}
	}

	rateBurst := 200
	if burst := os.Getenv("RATE_LIMIT_BURST"); burst != "" {
		if b, err := strconv.Atoi(burst); err == nil && b > 0 {
			rateBurst = b
		}
	}

	return SecurityConfig{
		APIKey:          os.Getenv("API_KEY"),
		AuthDisabled:    os.Getenv("AUTH_DISABLED") == "true",
		RateLimitRPS:    rateRPS,
		RateLimitBurst:  rateBurst,
		CORSOrigins:     os.Getenv("CORS_ALLOWED_ORIGINS"),
		CORSMethods:     getEnvOrDefault("CORS_ALLOWED_METHODS", "GET, POST, PUT, DELETE, OPTIONS"),
		CORSHeaders:     getEnvOrDefault("CORS_ALLOWED_HEADERS", "Content-Type, Authorization, X-API-Key, X-Request-ID"),
		CORSCredentials: os.Getenv("CORS_ALLOW_CREDENTIALS") == "true",
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
