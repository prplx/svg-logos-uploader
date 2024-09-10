package handlers

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"svg-logos-uploader/internal/config"
	"svg-logos-uploader/internal/github"
	"svg-logos-uploader/internal/lib/sl"
	"svg-logos-uploader/internal/markdown"

	"github.com/gin-gonic/gin"
)

const (
	uploadsDir = "uploads"
	owner      = "prplx"
	repo       = "svg-logos"
	mainBranch = "master"
)

func UploadHandler(c *gin.Context, cfg *config.Config, logger *slog.Logger) {
	log := logger.With("op", "cmd/web/handlers/upload")
	context, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.UploadTimeout)*time.Second)

	defer cancel()

	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err = os.Mkdir(uploadsDir, 0755)
		if err != nil {
			log.Error("cannot create uploads directory: ", sl.Err(err))
			c.Redirect(http.StatusSeeOther, "/?error=true")
			return
		}
	}

	err := c.Request.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		log.Error("cannot parse multipart form: ", sl.Err(err))
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	files := c.Request.MultipartForm.File["files"]
	fileNames := []string{}
	for _, file := range files {
		filename := filepath.Join(uploadsDir, file.Filename)
		err := c.SaveUploadedFile(file, filename)
		if err != nil {
			log.Error("cannot save uploaded file: ", sl.Err(err))
			c.Redirect(http.StatusSeeOther, "/?error=true")
			return
		}
		fileNames = append(fileNames, filename)
	}

	githubClient := github.NewGithubClient(cfg.GithubAccessToken)
	_, dirContent, err := githubClient.GetRepositoryContent(context, owner, repo, "svg")
	if err != nil {
		log.Error("cannot get repository content: ", sl.Err(err))
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	repoFiles := []string{}
	for _, file := range dirContent {
		repoFiles = append(repoFiles, *file.Name)
	}

	newMarkdown, err := markdown.AddFilesToMarkdown(append(repoFiles, fileNames...))
	if err != nil {
		log.Error("cannot add files to markdown: ", sl.Err(err))
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, newMarkdown)
	if err != nil {
		log.Error("cannot copy markdown to buffer: ", sl.Err(err))
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	err = os.WriteFile(uploadsDir+"/README.md", []byte(buf.String()), 0644)
	if err != nil {
		log.Error("cannot write markdown to file: ", sl.Err(err))
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	branchName := github.GenerateBranchNameFromUploadedFiles(fileNames)

	err = githubClient.CreateBranch(context, owner, repo, mainBranch, branchName)
	if err != nil {
		log.Error("cannot create branch: ", sl.Err(err))
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	err = githubClient.CreateTree(context, owner, repo, branchName, append(fileNames, uploadsDir+"/README.md"))
	if err != nil {
		log.Error("cannot create tree: ", sl.Err(err))
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	err = githubClient.CreatePullRequest(context, owner, repo, branchName, mainBranch, github.GenerateCommitMessageFromUploadedFiles(append(fileNames, uploadsDir+"/README.md")))
	if err != nil {
		log.Error("cannot create pull request: ", sl.Err(err))
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	err = clearDirectory(uploadsDir)
	if err != nil {
		log.Error("cannot clear uploads directory: ", sl.Err(err))
		c.Redirect(http.StatusSeeOther, "/?error=true")
		return
	}

	c.Redirect(http.StatusSeeOther, "/?success=true")
}

func clearDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return nil
}
