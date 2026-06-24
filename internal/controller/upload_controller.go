package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"domesv2/pkg/response"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	// Ensure uploads directory exists
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		fmt.Printf("Warning: Failed to create uploads directory: %v\n", err)
	}
	return &UploadController{}
}

func (ctrl *UploadController) UploadFile(c *fiber.Ctx) error {
	// Check multipart form file
	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, "No file uploaded", "INVALID_REQUEST_BODY")
	}

	// Limit to 50MB
	const maxFileSize = 50 * 1024 * 1024 // 50MB
	if file.Size > maxFileSize {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"success": false,
			"message": "File size exceeds the maximum limit of 50MB",
			"error":   "FILE_TOO_LARGE",
			"details": "File size exceeds the maximum limit of 50MB: FILE_TOO_LARGE",
		})
	}

	// Generate UUID filename
	ext := filepath.Ext(file.Filename)
	newFilename := uuid.New().String() + ext
	targetPath := filepath.Join("./uploads", newFilename)

	// Save file
	if err := c.SaveFile(file, targetPath); err != nil {
		return response.InternalServerError(c, "Failed to save file", "INTERNAL_ERROR")
	}

	// Format size
	fileSizeStr := formatBytes(file.Size)

	// Return response matching the contract
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"data": fiber.Map{
			"url":           "/uploads/" + newFilename,
			"filename":      newFilename,
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

	// Limit to 2MB
	const maxAvatarSize = 2 * 1024 * 1024 // 2MB
	if file.Size > maxAvatarSize {
		return response.BadRequest(c, "Avatar size exceeds the maximum limit of 2MB", "VALIDATION_FAILED")
	}

	// Validate format (jpg, png, webp)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return response.BadRequest(c, "Invalid file format. Only JPG, PNG, and WEBP are allowed", "VALIDATION_FAILED")
	}

	// Generate UUID filename
	newFilename := uuid.New().String() + ext
	targetPath := filepath.Join("./uploads", newFilename)

	// Save file
	if err := c.SaveFile(file, targetPath); err != nil {
		return response.InternalServerError(c, "Failed to save avatar", "INTERNAL_ERROR")
	}

	// Update user's avatar_url in database if authenticated user
	// Note: We'll retrieve user_id from token context if present
	if userIDVal := c.Locals("user_id"); userIDVal != nil {
		// Update user profile avatar URL
		// For simplicity, we can do this in the controller or ignore if it's just a general upload endpoint.
		// The contract says:
		// POST /api/upload/avatar -> Response contains avatar_url.
		// So we just return the URL, and front-end can pass it to profile update, or we can handle it here.
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Avatar uploaded successfully",
		"data": fiber.Map{
			"avatar_url": "/uploads/" + newFilename,
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
