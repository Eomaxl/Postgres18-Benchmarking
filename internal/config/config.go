package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig holds database specific configuration
type DatabaseConfig struct {
	Primary  DatabaseInstanceConfig
	Replicas []DatabaseInstanceConfig
	Pool     PoolConfig
}

// DatabaseInstanceConfig holds configuration for a single database instance
type DatabaseInstanceConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string
}

// PoolConfig holds configuration pool configuration
type PoolConfig struct {
	MinConnections int
	MaxConnections int
	MaxIdleTime    time.Duration
	MaxLifetime    time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Primary: DatabaseInstanceConfig{
				Host:     getEnv("DB_PRIMARY_HOST", "localhost"),
				Port:     getEnvAsInt("DB_PRIMARY_PORT", 5432),
				Database: getEnv("DB_NAME", "userdb"),
				User:     getEnv("DB_USER", "postgres"),
				Password: getEnv("DB_PASSWORD", "postgres"),
				SSLMode:  getEnv("DB_SSLMODE", "disable"),
			},
			Replicas: []DatabaseInstanceConfig{
				{
					Host:     getEnv("DB_REPLICA1_HOST", "localhost"),
					Port:     getEnvAsInt("DB_REPLICA1_PORT", 5433),
					Database: getEnv("DB_NAME", "userdb"),
					User:     getEnv("DB_USER", "postgres"),
					Password: getEnv("DB_PASSWORD", "postgres"),
					SSLMode:  getEnv("DB_SSLMODE", "disable"),
				},
				{
					Host:     getEnv("DB_REPLICA2_HOST", "localhost"),
					Port:     getEnvAsInt("DB_REPLICA2_PORT", 5434),
					Database: getEnv("DB_NAME", "userdb"),
					User:     getEnv("DB_USER", "postgres"),
					Password: getEnv("DB_PASSWORD", "postgres"),
					SSLMode:  getEnv("DB_SSLMODE", "disable"),
				},
			},
			Pool: PoolConfig{
				MinConnections: getEnvAsInt("DB_POOL_MIN_CONNECTIONS", 10),
				MaxConnections: getEnvAsInt("DB_POOL_MAX_CONNECTIONS", 100),
				MaxIdleTime:    getEnvAsDuration("DB_POOL_MAX_IDLE_TIME", 30*time.Minute),
				MaxLifetime:    getEnvAsDuration("DB_POOL_MAX_LIFETIME", 1*time.Hour),
			},
		},
	}

	return config, nil
}

// DSN returns the postgreSQL connection string for a database instance
func (d *DatabaseInstanceConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode,
	)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or return a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration retrieves an environment variable as a duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
