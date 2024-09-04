package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"svg-logos-uploader/internal/config"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context, cfg *config.Config) {
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

	c.Redirect(http.StatusSeeOther, "/?success=true")
}
