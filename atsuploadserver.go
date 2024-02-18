package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

const saveDir = "/home/charliepi/Downloads/" // Directory to save uploaded files

func main() {
	e := echo.New()

	// Route for file upload
	e.POST("/upload", func(c echo.Context) error {
		// Get form data
		form, err := c.MultipartForm()
		if err != nil {
			return c.String(http.StatusBadRequest, "Failed to parse form data")
		}

		// Access filepicker field
		files := form.File["filepicker"]
		if len(files) == 0 {
			return c.String(http.StatusBadRequest, "No file uploaded")
		}

		file := files[0]
		src, err := file.Open()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to open file")
		}
		defer src.Close()

		// Generate unique filename
		filename := fmt.Sprintf("%d.%s", time.Now().UnixNano(), file.Filename)
		savePath := filepath.Join(saveDir, filename)

		// Create save directory if not exists
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create save directory")
		}

		// Save file to disk
		dst, err := os.Create(savePath)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to save file")
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to copy file")
		}

		return c.String(http.StatusOK, fmt.Sprintf("File uploaded successfully: %s", filename))
	})

	// Start the server
	e.Logger.Fatal(e.Start(":8000"))
}
