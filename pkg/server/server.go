package server

import (
	"SpeechAnalytics/pkg/handlers"
	"SpeechAnalytics/pkg/logger"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	logger     *log.Logger
	HttpServer *http.Server
}

func New() *Server {
	port, ok := os.LookupEnv("TODO_PORT")
	if !ok || len(port) == 0 {
		port = "8080"
	}

	mux := handlers.Init()

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ErrorLog:     logger.Logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{logger: logger.Logger, HttpServer: server}
}
