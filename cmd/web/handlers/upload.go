package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"svg-logos-uploader/internal/config"
	"svg-logos-uploader/internal/github"

	b64 "encoding/base64"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context, cfg *config.Config) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	uploadsDir := "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err = os.Mkdir(uploadsDir, 0755)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/?error=true")
			return
		}
	}

	err := c.Request.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	// Get all files from the form field "files"
	files := c.Request.MultipartForm.File["files"]
	for _, file := range files {
		filename := filepath.Join(uploadsDir, file.Filename)
		err := c.SaveUploadedFile(file, filename)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/?error=true")
			return
		}
	}

	githubClient := github.NewGithubClient(cfg.GithubAccessToken)
	repoContent, err := githubClient.GetRepositories(context)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	sDec, _ := b64.StdEncoding.DecodeString(*repoContent.Content)
	fmt.Println(string(sDec))

	c.Redirect(http.StatusSeeOther, "/?success=true")
}
