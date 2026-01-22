package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Test with default values
	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Verify server defaults
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, 8080, config.Server.Port)

	// Verify database defaults
	assert.Equal(t, "localhost", config.Database.Primary.Host)
	assert.Equal(t, 5432, config.Database.Primary.Port)
	assert.Equal(t, "userdb", config.Database.Primary.Database)

	// Verify replicas
	assert.Len(t, config.Database.Replicas, 2)
	assert.Equal(t, 5433, config.Database.Replicas[0].Port)
	assert.Equal(t, 5434, config.Database.Replicas[1].Port)

	// Verify pool config
	assert.Equal(t, 10, config.Database.Pool.MinConnections)
	assert.Equal(t, 100, config.Database.Pool.MaxConnections)
}

func TestLoadWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_PRIMARY_HOST", "db-primary")
	os.Setenv("DB_PRIMARY_PORT", "5555")
	os.Setenv("DB_POOL_MAX_CONNECTIONS", "200")

	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_PRIMARY_HOST")
		os.Unsetenv("DB_PRIMARY_PORT")
		os.Unsetenv("DB_POOL_MAX_CONNECTIONS")
	}()

	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Verify environment variables are used
	assert.Equal(t, "127.0.0.1", config.Server.Host)
	assert.Equal(t, 9090, config.Server.Port)
	assert.Equal(t, "db-primary", config.Database.Primary.Host)
	assert.Equal(t, 5555, config.Database.Primary.Port)
	assert.Equal(t, 200, config.Database.Pool.MaxConnections)
}

func TestDatabaseInstanceConfigDSN(t *testing.T) {
	config := DatabaseInstanceConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		User:     "testuser",
		Password: "testpass",
		SSLMode:  "disable",
	}

	expectedDSN := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	assert.Equal(t, expectedDSN, config.DSN())
}

func TestGetEnvAsDuration(t *testing.T) {
	os.Setenv("TEST_DURATION", "45m")
	defer os.Unsetenv("TEST_DURATION")

	duration := getEnvAsDuration("TEST_DURATION", 30*time.Minute)
	assert.Equal(t, 45*time.Minute, duration)

	// Test with invalid duration (should return default)
	os.Setenv("TEST_DURATION", "invalid")
	duration = getEnvAsDuration("TEST_DURATION", 30*time.Minute)
	assert.Equal(t, 30*time.Minute, duration)
}
