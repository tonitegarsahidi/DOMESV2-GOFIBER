package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"domesv2/internal/service"
	"domesv2/pkg/response"
)

type UploadController struct {
	uploadService service.FileUploadService
}

func NewUploadController(uploadService service.FileUploadService) *UploadController {
	return &UploadController{
		uploadService: uploadService,
	}
}

func (ctrl *UploadController) UploadFile(c *fiber.Ctx) error {
	// Check multipart form file
	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, "No file uploaded", "INVALID_REQUEST_BODY")
	}

	// We pass type parameter (main, cover, additional)
	fileType := c.Query("type")
	if fileType == "" {
		fileType = c.FormValue("type")
	}

	url, err := ctrl.uploadService.UploadFile(file, fileType)
	if err != nil {
		return response.BadRequest(c, err.Error(), "VALIDATION_FAILED")
	}

	// Format size
	fileSizeStr := formatBytes(file.Size)

	// Return response matching the contract
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"data": fiber.Map{
			"url":           url,
			"filename":      filepath.Base(url),
			"original_name": file.Filename,
			"file_size":     fileSizeStr,
			"mime_type":     file.Header.Get("Content-Type"),
		},
	})
}

func (ctrl *UploadController) UploadAvatar(c *fiber.Ctx) error {
	file, err := c.FormFile("avatar")
	if err != nil {
		return response.BadRequest(c, "No avatar file uploaded", "INVALID_REQUEST_BODY")
	}

	url, err := ctrl.uploadService.UploadFile(file, "cover")
	if err != nil {
		return response.BadRequest(c, err.Error(), "VALIDATION_FAILED")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Avatar uploaded successfully",
		"data": fiber.Map{
			"avatar_url": url,
		},
	})
}

func (ctrl *UploadController) ValidateURL(c *fiber.Ctx) error {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	if req.URL == "" {
		return response.BadRequest(c, "URL is required", "VALIDATION_FAILED")
	}

	// Perform HTTP HEAD request with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(req.URL)
	if err != nil {
		// Try GET request in case HEAD is not supported by target server
		resp, err = client.Get(req.URL)
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"message": "URL is not accessible",
				"data": fiber.Map{
					"url":        req.URL,
					"accessible": false,
					"error":      err.Error(),
				},
			})
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode >= 400 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "URL is not accessible",
			"data": fiber.Map{
				"url":        req.URL,
				"accessible": false,
				"error":      fmt.Sprintf("HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode)),
			},
		})
	}

	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.ContentLength
	fileSizeStr := "Unknown"
	if contentLength > 0 {
		fileSizeStr = formatBytes(contentLength)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "URL is valid",
		"data": fiber.Map{
			"url":          req.URL,
			"accessible":   true,
			"content_type": contentType,
			"file_size":    fileSizeStr,
		},
	})
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
