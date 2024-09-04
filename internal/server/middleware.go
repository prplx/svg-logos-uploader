package server

import (
	"net/http"
	"svg-logos-uploader/internal/config"
	"svg-logos-uploader/internal/lib/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

func authMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(cfg.SessionCookieName)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		if tokenString == "" {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		token, err := jwt.ParseToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		if token["exp"].(float64) < float64(time.Now().Unix()) {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.Next()
	}
}
