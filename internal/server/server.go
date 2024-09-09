package server

import (
	"fmt"
	"net/http"
	"time"

	"log/slog"
	"svg-logos-uploader/internal/config"
)

type Server struct {
	port int
}

func NewServer(cfg *config.Config, log *slog.Logger) *http.Server {
	NewServer := &Server{
		port: cfg.Port,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(cfg, log),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
