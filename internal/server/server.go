package server

import (
	"fmt"
	"net/http"
	"time"

	"svg-logos-uploader/internal/config"
)

type Server struct {
	port int
}

func NewServer(cfg *config.Config) *http.Server {
	NewServer := &Server{
		port: cfg.Port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(cfg),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
