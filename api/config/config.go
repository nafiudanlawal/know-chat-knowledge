package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	// Database configuration
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		URL      string
	}

	// File storage configuration
	FileStorage struct {
		Path string
	}
}

// New creates a new configuration from environment variables
func New() *Config {
	config := &Config{}

	// Setup database config
	config.DB.Host = getEnv("PGHOST", "localhost")
	config.DB.Port = getEnv("PGPORT", "5432")
	config.DB.User = getEnv("PGUSER", "postgres")
	config.DB.Password = getEnv("PGPASSWORD", "postgres")
	config.DB.Name = getEnv("PGDATABASE", "knowledge")
	config.DB.URL = getEnv("DATABASE_URL", "")

	// Setup file storage
	config.FileStorage.Path = getEnv("FILE_STORAGE_PATH", "./uploads")

	// Create upload directory if it doesn't exist
	if _, err := os.Stat(config.FileStorage.Path); os.IsNotExist(err) {
		err := os.MkdirAll(config.FileStorage.Path, 0755)
		if err != nil {
			fmt.Printf("Error creating upload directory: %v\n", err)
		}
	}

	return config
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", // TODO add prod sslmode conditions
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name)
}

// getEnv retrieves the value of the environment variable named by the key
// If the variable is not present, returns the default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
