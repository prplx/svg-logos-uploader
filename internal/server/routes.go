package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"svg-logos-uploader/cmd/web"
	"svg-logos-uploader/cmd/web/handlers"

	"log/slog"
	"svg-logos-uploader/internal/config"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes(cfg *config.Config, log *slog.Logger) http.Handler {
	assetsFS, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.StaticFS("/assets", http.FS(assetsFS))

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

	r.GET("/debug/files", func(c *gin.Context) {
		var files []string
		err := fs.WalkDir(web.Files, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, files)
	})

	r.Use(func(c *gin.Context) {
		fmt.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	return r
}
