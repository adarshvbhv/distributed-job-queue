package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort    string
	WorkerCount int
	DatabaseURL string
}

func Load() Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system environment variables")
	}
	port := os.Getenv("HTTP_PORT")
	worker := os.Getenv("WORKER_COUNT")
	database_url := os.Getenv("DATABASE_URL")

	if worker == "" {
		worker = "10"
	}

	workerCnt, _ := strconv.Atoi(worker)

	if port == "" {
		port = "8080"
	}

	return Config{
		HTTPPort:    port,
		WorkerCount: workerCnt,
		DatabaseURL: database_url,
	}
}
