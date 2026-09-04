package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	UploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{UploadDir: uploadDir}
}

func (h *UploadHandler) UploadImage(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "No image file provided")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		return echo.NewHTTPError(http.StatusBadRequest, "Only jpg, jpeg, png, webp files are allowed")
	}

	if err := os.MkdirAll(h.UploadDir, 0755); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create upload directory")
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(h.UploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save file")
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to write file")
	}

	url := fmt.Sprintf("/uploads/%s", filename)
	return c.JSON(http.StatusCreated, echo.Map{
		"code": 201, "status": "success", "url": url,
	})
}
