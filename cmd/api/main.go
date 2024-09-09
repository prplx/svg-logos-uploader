package main

import (
	"fmt"
	"log/slog"
	"os"
	"svg-logos-uploader/internal/config"
	"svg-logos-uploader/internal/server"

	"github.com/gin-gonic/gin"
)

const (
	envProd = "production"
)

func main() {
	config := config.MustLoadConfig()
	logger := setupLogger(config.Env)

	if config.Env == envProd {
		gin.SetMode(gin.ReleaseMode)
	}

	server := server.NewServer(config, logger)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	if env == envProd {
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	} else {
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}
