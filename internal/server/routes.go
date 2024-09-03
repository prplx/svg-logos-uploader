package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"svg-logos-uploader/cmd/web"

	"svg-logos-uploader/internal/config"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes(cfg *config.Config) http.Handler {
	r := gin.Default()

	r.Static("/assets", "./cmd/web/assets")

	r.POST("/hello", func(c *gin.Context) {
		web.HelloWebHandler(c.Writer, c.Request)
	})

	r.GET("/login", func(c *gin.Context) {
		loginError := c.Query("error") == "true"
		templ.Handler(web.LoginForm(loginError)).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/login", func(c *gin.Context) {
		web.LoginHandler(c, cfg)
	})

	return r
}
