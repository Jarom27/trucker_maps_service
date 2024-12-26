package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	Port        string
	RabbitHost  string
	RabbitPort  string
	RabbitUser  string
	RabbitPass  string
	RabbitQueue string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // Cargar .env (si existe)

	return &Config{
		Host:        getEnv("TRUCKER_HOST", "localhost"),
		Port:        getEnv("TRUCKER_PORT", "6000"),
		RabbitHost:  getEnv("RABBIT_HOST", "localhost"),
		RabbitPort:  getEnv("RABBIT_PORT", "5672"),
		RabbitUser:  getEnv("RABBIT_DEFAULT_USER", ""),
		RabbitPass:  getEnv("RABBIT_DEFAULT_PASS", ""),
		RabbitQueue: getEnv("RABBIT_QUEUE_NAME", ""),
	}
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
