package handlers

import (
	"net/http"

	"svg-logos-uploader/internal/config"
	"svg-logos-uploader/internal/lib/jwt"
	"svg-logos-uploader/internal/lib/sl"

	"log/slog"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context, cfg *config.Config, logger *slog.Logger) {
	log := logger.With("op", "cmd/web/handlers/login")

	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		log.Error("empty username or password")
		c.Redirect(http.StatusSeeOther, "/login?error=true")
		return
	}

	if username != cfg.AdminLogin || !checkPasswordHash(password, cfg.AdminPassword) {
		log.Error("invalid admin's user name or password")
		c.Redirect(http.StatusSeeOther, "/login?error=true")
		return
	}

	token, err := jwt.GenerateToken(username, cfg.JWTSecret)
	if err != nil {
		log.Error("cannot generate token: ", sl.Err(err))
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
