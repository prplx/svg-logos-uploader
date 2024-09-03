package web

import (
	"net/http"

	"svg-logos-uploader/internal/config"
	"svg-logos-uploader/internal/lib/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context, cfg *config.Config) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.Redirect(http.StatusSeeOther, "/login?error=true")
		return
	}

	if username != cfg.AdminLogin || !checkPasswordHash(password, cfg.AdminPassword) {
		c.Redirect(http.StatusSeeOther, "/login?error=true")
		return
	}

	token, err := jwt.GenerateToken(username, cfg.JWTSecret)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login?error=true")
		return
	}
	c.SetCookie("session", token, 30*24*60*60, "/", "", cfg.Env == "production", true)
	c.Redirect(http.StatusSeeOther, "/")
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
