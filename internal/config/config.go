package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Worker   WorkerConfig
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type WorkerConfig struct {
	Count     int
	QueueSize int
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "exams_db"),
		},
		Worker: WorkerConfig{
			Count:     getIntEnv("WORKER_COUNT", 5),
			QueueSize: getIntEnv("QUEUE_SIZE", 100),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
