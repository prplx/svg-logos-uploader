package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"svg-logos-uploader/cmd/web"
	"svg-logos-uploader/cmd/web/handlers"

	"log/slog"
	"svg-logos-uploader/internal/config"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes(cfg *config.Config, log *slog.Logger) http.Handler {
	r := gin.Default()

	r.Static("/assets", "./cmd/web/assets")

	r.GET("/login", func(c *gin.Context) {
		loginError := c.Query("error") == "true"
		templ.Handler(web.LoginForm(loginError)).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/login", func(c *gin.Context) {
		handlers.LoginHandler(c, cfg, log)
	})

	r.GET("/", authMiddleware(cfg), func(c *gin.Context) {
		uploadError := c.Query("error") == "true"
		uploadSuccess := c.Query("success") == "true"
		templ.Handler(web.UploadForm(uploadError, uploadSuccess)).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/upload", authMiddleware(cfg), func(c *gin.Context) {
		handlers.UploadHandler(c, cfg, log)
	})

	return r
}
