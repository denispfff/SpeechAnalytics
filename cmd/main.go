package main

import (
	"SpeechAnalytics/pkg/database"
	"SpeechAnalytics/pkg/logger"
	"SpeechAnalytics/pkg/repositories"
	"SpeechAnalytics/pkg/services"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Println("Ошибка загрузки .env:", envErr)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Moscow",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	auth := os.Getenv("API_KEY")
	modelUri := os.Getenv("MODEL_URI")

	logger.Init()
	database.ConnectDB(dsn)
	repositories.InitFilePaths()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	services.StartProcessing(ctx, auth, modelUri)
}
