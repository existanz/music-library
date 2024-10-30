package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"music-library/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

//	@title			Music library API
//	@version		1.0
//	@description	This is a music library API

//	@host		localhost:4001
//	@BasePath	/songs

func main() {
	setupLogger()
	server := server.NewServer()
	go gracefulShutdown(server)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}

func gracefulShutdown(apiServer *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
}

func setupLogger() {
	logLevel := new(slog.LevelVar)
	options := &slog.HandlerOptions{Level: logLevel}

	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel.Set(slog.LevelDebug)
		options.AddSource = true
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, options))
	slog.SetDefault(logger)
}
