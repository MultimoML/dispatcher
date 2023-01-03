package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnection         string
	DBName               string
	Port                 string
	ConfigServer         string
	ConfigServerPort     string
	ConfigServerGrpcPort string
}

func LoadConfig() *Config {
	env := os.Getenv("ACTIVE_ENV")
	if env == "" {
		env = "dev"
	}

	switch env {
	case "dev":
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	case "prod":
	default:
		log.Fatal("Unknown environment")
	}

	log.Println("Loaded environment variables for", env)

	port := "6001"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return &Config{
		DBConnection:         os.Getenv("DB_CONNECTION_STRING"),
		DBName:               os.Getenv("DB_NAME"),
		ConfigServer:         os.Getenv("CONFIG_SERVER"),
		ConfigServerPort:     os.Getenv("CONFIG_SERVER_PORT"),
		ConfigServerGrpcPort: os.Getenv("CONFIG_SERVER_PORT_GRPC"),
		Port:                 port,
	}
}
