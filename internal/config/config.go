package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPort              string
	AppPort             string
	Env                 string
	AppName             string
	DBHost              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
	JWTSecret           string
	LogLevel            string
	AllowedOrigins      []string
	JWTExpiry           time.Duration
	RefreshTokenExpiry  time.Duration
	MidtransServerKey   string
	MidtransClientKey   string
	MidtransEnvironment string
}

func Load() (*Config, error) {
	_ = godotenv.Load() //nolint:errcheck

	cfg := &Config{
		AppPort:             getEnv("PORT", "8080"),
		Env:                 getEnv("ENV", "development"),
		AppName:             getEnv("APP_NAME", "Matchaciee API"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", ""),
		DBName:              getEnv("DB_NAME", "matchaciee_dev"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
		JWTSecret:           getEnv("JWT_SECRET", "rahasiamatcha"),
		JWTExpiry:           getEnvAsDuration("JWT_EXPIRY", 1*time.Hour),
		RefreshTokenExpiry:  getEnvAsDuration("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
		AllowedOrigins:      getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		MidtransServerKey:   getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey:   getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransEnvironment: getEnv("MIDTRANS_ENVIRONMENT", "sandbox"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	if c.DBPassword == "" && c.Env == "production" {
		return fmt.Errorf("DB_PASSWORD is required in production")
	}

	// Validate Midtrans configuration in production
	if c.Env == "production" {
		if c.MidtransServerKey == "" {
			return fmt.Errorf("MIDTRANS_SERVER_KEY is required in production")
		}
		if c.MidtransClientKey == "" {
			return fmt.Errorf("MIDTRANS_CLIENT_KEY is required in production")
		}
	}

	// Validate Midtrans environment value
	if c.MidtransEnvironment != "sandbox" && c.MidtransEnvironment != "production" {
		return fmt.Errorf("MIDTRANS_ENVIRONMENT must be either 'sandbox' or 'production'")
	}

	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
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
