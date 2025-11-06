package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	MongoDB   MongoDBConfig
	JWT       JWTConfig
	CSRF      CSRFConfig
	Cookie    CookieConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type CSRFConfig struct {
	Secret string
}

type CookieConfig struct {
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite string
}

type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "mileapp_db"),
			Timeout:  time.Duration(getEnvAsInt("MONGODB_TIMEOUT", 10)) * time.Second,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
			Expiry: time.Duration(getEnvAsInt("JWT_EXPIRY_MINUTES", 15)) * time.Minute,
		},
		CSRF: CSRFConfig{
			Secret: getEnv("CSRF_SECRET", ""),
		},
		Cookie: CookieConfig{
			Domain:   getEnv("COOKIE_DOMAIN", "localhost"),
			Secure:   getEnvAsBool("COOKIE_SECURE", false),
			HTTPOnly: getEnvAsBool("COOKIE_HTTP_ONLY", true),
			SameSite: getEnv("COOKIE_SAME_SITE", "Strict"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173"}),
			AllowedMethods:   getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders:   getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization", "X-CSRF-Token"}),
			ExposeHeaders:    getEnvAsSlice("CORS_EXPOSE_HEADERS", []string{}),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvAsInt("CORS_MAX_AGE", 3600),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if c.CSRF.Secret == "" {
		return fmt.Errorf("CSRF_SECRET is required")
	}

	if c.MongoDB.URI == "" {
		return fmt.Errorf("MONGODB_URI is required")
	}

	if c.MongoDB.Database == "" {
		return fmt.Errorf("MONGODB_DATABASE is required")
	}

	return nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	values := strings.Split(valueStr, ",")
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}

	return values
}
